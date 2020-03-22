package baker

import (
	"context"
	"math/big"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v2/cli/internal/db/model"
	"github.com/goat-systems/tzpay/v2/cli/internal/enviroment"
	"github.com/pkg/errors"
)

// Baker is a tezos baker that can get payouts and execute them.
type Baker struct {
	gt gotezos.IFace
}

type processDelegationsInput struct {
	delegations          *[]string
	stakingBalance       *big.Int
	frozenBalanceRewards *gotezos.FrozenBalance
	blockHash            string
}

type processDelegateInput struct {
	delegate             string
	delegations          []model.DelegationEarning
	stakingBalance       *big.Int
	frozenBalanceRewards *gotezos.FrozenBalance
	blockHash            string
}

type processDelegationsOutput struct {
	delegationEarning model.DelegationEarning
	err               error
}

type processDelegationInput struct {
	delegation           string
	stakingBalance       *big.Int
	frozenBalanceRewards *gotezos.FrozenBalance
	blockHash            string
}

// NewBaker returns a pointer to a new Baker
func NewBaker(gt gotezos.IFace) *Baker {
	return &Baker{gt: gt}
}

// Payouts returns all payouts for a cycle
func (b *Baker) Payouts(ctx context.Context, cycle int) (*model.Payout, error) {
	params := enviroment.GetEnviromentFromContext(ctx)
	frozenBalanceRewards, err := b.gt.FrozenBalance(cycle, params.Delegate)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get delegation earnings for cycle %d", cycle)
	}

	delegations, err := b.gt.DelegatedContractsAtCycle(cycle, params.Delegate)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get delegation earnings for cycle %d", cycle)
	}

	networkCycle, err := b.gt.Cycle(cycle)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get delegation earnings for cycle %d", cycle)
	}

	stakingBalance, err := b.gt.StakingBalance(networkCycle.BlockHash, params.Delegate)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get delegation earnings for cycle %d", cycle)
	}

	out := b.proccessDelegations(ctx, &processDelegationsInput{
		delegations:          delegations,
		stakingBalance:       stakingBalance,
		frozenBalanceRewards: frozenBalanceRewards,
		blockHash:            networkCycle.BlockHash,
	})

	payouts := model.Payout{
		StakingBalance: stakingBalance,
		CycleHash:      networkCycle.BlockHash,
		Cycle:          cycle,
		FrozenBalance:  frozenBalanceRewards.Rewards.Big,
	}

	for _, delegation := range out {
		if delegation.err != nil {
			err = errors.Wrapf(delegation.err, "failed to get payout for delegation %s", delegation.delegationEarning.Address)
		} else {
			payouts.DelegationEarnings = append(payouts.DelegationEarnings, delegation.delegationEarning)
		}
	}

	payouts.DelegateEarnings, err = b.processDelegate(ctx, &processDelegateInput{
		delegate:             params.Delegate,
		delegations:          payouts.DelegationEarnings,
		stakingBalance:       stakingBalance,
		frozenBalanceRewards: frozenBalanceRewards,
		blockHash:            networkCycle.BlockHash,
	})
	if err != nil {
		err = errors.Wrap(err, "failed to get contruct payout info for delegate")
	}

	return &payouts, err
}

func (b *Baker) processDelegate(ctx context.Context, input *processDelegateInput) (model.DelegateEarnings, error) {
	delegateEarning := model.DelegateEarnings{
		Address: input.delegate,
		Net:     big.NewInt(0),
	}
	balance, err := b.gt.Balance(input.blockHash, input.delegate)
	if err != nil {
		return delegateEarning, errors.Wrapf(err, "failed to process delegate earnings for %s", input.delegate)
	}

	delegateEarning.Share = float64(balance.Int64()) / float64(input.stakingBalance.Int64())
	rewardsFloat := delegateEarning.Share * float64(input.frozenBalanceRewards.Rewards.Big.Int64())
	delegateEarning.Rewards = big.NewInt(int64(rewardsFloat))

	fees := big.NewInt(0)
	for _, delegation := range input.delegations {
		fees.Add(fees, delegation.Fee)
	}

	delegateEarning.Fees = fees
	delegateEarning.Net.Add(delegateEarning.Fees, delegateEarning.Rewards)

	return delegateEarning, nil
}

func (b *Baker) proccessDelegations(ctx context.Context, input *processDelegationsInput) []processDelegationsOutput {
	numJobs := len(*input.delegations)
	jobs := make(chan processDelegationInput, numJobs)
	results := make(chan processDelegationsOutput, numJobs)

	for i := 0; i < 50; i++ {
		go b.proccessDelegationWorker(ctx, jobs, results)
	}

	for _, pd := range *input.delegations {
		jobs <- processDelegationInput{
			delegation:           pd,
			stakingBalance:       input.stakingBalance,
			frozenBalanceRewards: input.frozenBalanceRewards,
			blockHash:            input.blockHash,
		}
	}
	close(jobs)

	var out []processDelegationsOutput
	for i := 1; i <= numJobs; i++ {
		out = append(out, <-results)
	}
	close(results)

	return out
}

func (b *Baker) proccessDelegationWorker(ctx context.Context, jobs <-chan processDelegationInput, results chan<- processDelegationsOutput) {
	for j := range jobs {
		d, err := b.processDelegation(ctx, &j)
		if err != nil {
			results <- processDelegationsOutput{
				err: err,
			}
		} else {
			results <- processDelegationsOutput{
				delegationEarning: *d,
			}
		}
	}
}

func (b *Baker) processDelegation(ctx context.Context, input *processDelegationInput) (*model.DelegationEarning, error) {
	params := enviroment.GetEnviromentFromContext(ctx)
	delegationEarning := &model.DelegationEarning{Address: input.delegation}
	balance, err := b.gt.Balance(input.blockHash, input.delegation)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to process delegation earnings for delegation %s", input.delegation)
	}

	delegationEarning.Share = float64(balance.Int64()) / float64(input.stakingBalance.Int64())
	grossRewardsFloat := delegationEarning.Share * float64(input.frozenBalanceRewards.Rewards.Big.Int64())
	feeFloat := grossRewardsFloat * params.BakersFee

	delegationEarning.GrossRewards = big.NewInt(int64(grossRewardsFloat))
	delegationEarning.Fee = big.NewInt(int64(feeFloat))
	delegationEarning.NetRewards = big.NewInt(0).Sub(delegationEarning.GrossRewards, delegationEarning.Fee)

	return delegationEarning, nil
}

// ForgePayout converts Payout into operation contents and forges them locally
func (b *Baker) ForgePayout(ctx context.Context, payout model.Payout) (string, error) {
	base := enviroment.GetEnviromentFromContext(ctx)
	head, err := b.gt.Head()
	if err != nil {
		return "", errors.Wrap(err, "failed to forge payout")
	}

	counter, err := b.gt.Counter(head.Hash, base.Wallet.Address)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge payout")
	}

	transactions := b.constructPayoutContents(ctx, *counter, payout)

	forge, err := gotezos.ForgeTransactionOperation(head.Hash, transactions...)
	if err != nil {
		return "", errors.Wrap(err, "failed to forge payout")
	}

	return *forge, nil
}

func (b *Baker) constructPayoutContents(ctx context.Context, counter int, payout model.Payout) []gotezos.ForgeTransactionOperationInput {
	base := enviroment.GetEnviromentFromContext(ctx)
	var contents []gotezos.ForgeTransactionOperationInput
	for _, delegation := range payout.DelegationEarnings {
		if delegation.NetRewards.Int64() >= int64(base.MinimumPayment) {
			counter++
			contents = append(contents, gotezos.ForgeTransactionOperationInput{
				Source:       base.Wallet.Address,
				Destination:  delegation.Address,
				Amount:       gotezos.Int{Big: delegation.NetRewards},
				Fee:          gotezos.Int{Big: big.NewInt(int64(base.NetworkFee))}, //TODO: expose NewInt function in GoTezos
				GasLimit:     gotezos.Int{Big: big.NewInt(int64(base.GasLimit))},
				Counter:      counter,
				StorageLimit: gotezos.Int{Big: big.NewInt(int64(0))},
			})
		}
	}
	return contents
}
