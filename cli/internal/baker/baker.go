package baker

import (
	"context"
	"math/big"
	"unicode"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/payman/v2/cmd/internal/enviroment"
	"github.com/pkg/errors"
)

// DelegationEarnings -
type DelegationEarnings struct {
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
	delegationEarnings DelegationEarnings
	err                error
}

type processDelegationInput struct {
	delegation           string
	stakingBalance       *big.Int
	frozenBalanceRewards *gotezos.FrozenBalance
	blockHash            string
}

// Payouts contains a alphabetically case sensitive sorted list of DelegationEarnings
type Payouts []DelegationEarnings

func (p Payouts) Len() int { return len(p) }
func (p Payouts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Payouts) Less(i, j int) bool {
	iRunes := []rune(p[i].Delegation)
	jRunes := []rune(p[j].Delegation)

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
func (b *Baker) Payouts(ctx context.Context, cycle int) (*Payouts, error) {
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

	payouts := Payouts{}
	for _, delegation := range out {
		if delegation.err != nil {
			err = errors.Wrapf(delegation.err, "failed to get payout for delegation %s", delegation.delegationEarnings.Delegation)
		} else {
			payouts = append(payouts, delegation.delegationEarnings)
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
				delegationEarnings: *d,
			}
		}
	}
}

func (b *Baker) processDelegation(ctx context.Context, input *processDelegationInput) (*DelegationEarnings, error) {
	params := enviroment.GetEnviromentFromContext(ctx)
	delegationEarnings := &DelegationEarnings{Delegation: input.delegation}
	balance, err := b.gt.Balance(input.blockHash, input.delegation)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to process delegation earnings for delegation %s", input.delegation)
	}

	delegationEarnings.Share = float64(balance.Int64()) / float64(input.stakingBalance.Int64())
	grossRewardsFloat := delegationEarnings.Share * float64(input.frozenBalanceRewards.Rewards.Big.Int64())
	feeFloat := grossRewardsFloat * params.BakersFee

	delegationEarnings.GrossRewards = big.NewInt(int64(grossRewardsFloat))
	delegationEarnings.Fee = big.NewInt(int64(feeFloat))
	delegationEarnings.NetRewards = big.NewInt(0).Sub(delegationEarnings.GrossRewards, delegationEarnings.Fee)

	return delegationEarnings, nil
}
