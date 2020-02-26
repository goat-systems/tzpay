package delegates

import (
	"math/big"
	"strings"
	"sync"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/pkg/errors"
)

// DelegationEarnings -
type DelegationEarnings struct {
	Fee          *big.Int
	GrossRewards *big.Int
	NetRewards   *big.Int
	Share        float64
}

// Baker -
type Baker struct {
	gt        gotezos.IFace
	address   string
	fee       float64
	blackList []string
}

// NewBakerInput -
type NewBakerInput struct {
	GoTezos   *gotezos.GoTezos
	Address   string
	Fee       float64
	BlackList []string
}

// NewBaker -
func NewBaker(input *NewBakerInput) *Baker {
	return &Baker{
		gt:        input.GoTezos,
		address:   input.Address,
		fee:       input.Fee,
		blackList: input.BlackList,
	}
}

// GetDelegationEarnings -
func (b *Baker) GetDelegationEarnings(cycle int) (*[]DelegationEarnings, error) {
	frozenBalanceRewards, err := b.gt.FrozenBalance(cycle, b.address)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get delegation earnings for cycle %d", cycle)
	}

	delegations, err := b.gt.DelegatedContractsAtCycle(cycle, b.address)

	networkCycle, err := b.gt.Cycle(cycle)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get delegation earnings for cycle %d", cycle)
	}

	stakingBalance, err := b.gt.StakingBalance(networkCycle.BlockHash, b.address)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get delegation earnings for cycle %d", cycle)
	}

	outChan, wg := b.proccessDelegations(&processDelegationsInput{
		delegations:          delegations,
		stakingBalance:       stakingBalance,
		frozenBalanceRewards: frozenBalanceRewards,
		blockHash:            networkCycle.BlockHash,
	})

	var errs []error
	var delegationsEarnings []DelegationEarnings

	for {
		select {
		case out := <-outChan:
		case wg.Wait():

		}
	}
	for out := range outChan {
		if out.err != nil {
			errs = append(errs, err)
		} else {
			delegationsEarnings = append(delegationsEarnings, out.delegationEarnings)
		}

	}

	if len(errs) > 0 {
		return &delegationsEarnings, multierror(errs)
	}

	return &delegationsEarnings, nil
}

func multierror(errs []error) error {
	err := errors.New("")
	for _, e := range errs {
		err = errors.Wrap(err, e.Error())
	}
	return err
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

func (b *Baker) proccessDelegations(input *processDelegationsInput) (chan processDelegationsOutput, *sync.WaitGroup) {
	outchan := make(chan processDelegationsOutput, len(*input.delegations))
	wg := &sync.WaitGroup{}

	for _, delegation := range *input.delegations {
		wg.Add(1)
		go func(del string) {
			d, err := b.processDelegation(&processDelegationInput{
				delegation:           del,
				stakingBalance:       input.stakingBalance,
				frozenBalanceRewards: input.frozenBalanceRewards,
				blockHash:            input.blockHash,
			})
			if err != nil {
				outchan <- processDelegationsOutput{err: err}
				wg.Done()
			} else {
				outchan <- processDelegationsOutput{delegationEarnings: *d}
				wg.Done()
			}
		}(delegation)
	}

	return outchan, wg
}

type processDelegationInput struct {
	delegation           string
	stakingBalance       *big.Int
	frozenBalanceRewards *gotezos.FrozenBalance
	blockHash            string
}

func (b *Baker) processDelegation(input *processDelegationInput) (*DelegationEarnings, error) {
	delegationEarnings := &DelegationEarnings{}
	balance, err := b.gt.Balance(input.blockHash, input.delegation)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to process delegation earnings for delegation %s", input.delegation)
	}

	delegationEarnings.Share = float64(balance.Int64()) / float64(input.stakingBalance.Int64())
	grossRewardsFloat := delegationEarnings.Share * float64(input.frozenBalanceRewards.Rewards.Big.Int64())
	feeFloat := grossRewardsFloat * b.fee

	delegationEarnings.GrossRewards = big.NewInt(int64(grossRewardsFloat))
	delegationEarnings.Fee = big.NewInt(int64(feeFloat))
	delegationEarnings.NetRewards = big.NewInt(0).Sub(delegationEarnings.GrossRewards, delegationEarnings.Fee)

	return delegationEarnings, nil
}

// ParseBlackList -
func ParseBlackList(list string) []string {
	blacklist := strings.Split(list, ",")
	for i := range blacklist {
		blacklist[i] = strings.Trim(blacklist[i], " ")
	}

	return blacklist
}
