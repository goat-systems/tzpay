package baker

import (
	"context"
	"math/big"
	"unicode"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v2/cli/internal/enviroment"
	"github.com/pkg/errors"
)

// DelegationEarning -
type DelegationEarning struct {
	Delegation   string
	Fee          *big.Int
	GrossRewards *big.Int
	NetRewards   *big.Int
	Share        float64
}

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

type processDelegationsOutput struct {
	delegationEarning DelegationEarning
	err               error
}

type processDelegationInput struct {
	delegation           string
	stakingBalance       *big.Int
	frozenBalanceRewards *gotezos.FrozenBalance
	blockHash            string
}

// Payout contains all needed information for a payout
type Payout struct {
	DelegationEarnings DelegationEarnings `json:"delegaions"`
	Cycle              int                `json:"cycle"`
	FrozenBalance      *big.Int           `json:"rewards"`
	StakingBalance     *big.Int           `json:"staking_balance"`
	Delegate           string             `json:"delegate"`
}

// DelegationEarnings contains list of DelegationEarning and implements sort.
type DelegationEarnings []DelegationEarning

func (d DelegationEarnings) Len() int { return len(d) }
func (d DelegationEarnings) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d DelegationEarnings) Less(i, j int) bool {
	iRunes := []rune(d[i].Delegation)
	jRunes := []rune(d[j].Delegation)

	max := len(iRunes)
	if max > len(jRunes) {
		max = len(jRunes)
	}

	for idx := 0; idx < max; idx++ {
		ir := iRunes[idx]
		jr := jRunes[idx]

		lir := unicode.ToLower(ir)
		ljr := unicode.ToLower(jr)

		if lir != ljr {
			return lir < ljr
		}

		if ir != jr {
			return ir < jr
		}
	}

	return false
}

// NewBaker returns a pointer to a new Baker
func NewBaker(gt gotezos.IFace) *Baker {
	return &Baker{gt: gt}
}

// Payouts returns all payouts for a cycle
func (b *Baker) Payouts(ctx context.Context, cycle int) (*Payout, error) {
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

	payouts := Payout{
		Delegate:       params.Delegate,
		StakingBalance: stakingBalance,
		Cycle:          cycle,
		FrozenBalance:  frozenBalanceRewards.Rewards.Big,
	}
	for _, delegation := range out {
		if delegation.err != nil {
			err = errors.Wrapf(delegation.err, "failed to get payout for delegation %s", delegation.delegationEarning.Delegation)
		} else {
			payouts.DelegationEarnings = append(payouts.DelegationEarnings, delegation.delegationEarning)
		}
	}

	return &payouts, err
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

func (b *Baker) processDelegation(ctx context.Context, input *processDelegationInput) (*DelegationEarning, error) {
	params := enviroment.GetEnviromentFromContext(ctx)
	delegationEarning := &DelegationEarning{Delegation: input.delegation}
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
