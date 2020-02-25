package delegates

import (
	"math/big"
	"strconv"
	"strings"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/pkg/errors"
)

// DelegationEarnings -
type DelegationEarnings struct {
	Fee          gotezos.BigInt
	GrossRewards gotezos.BigInt
	NetRewards   gotezos.BigInt
	Share        float64
}

// Baker -
type Baker struct {
	gt        *gotezos.GoTezos
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

	sbInt, err := strconv.Atoi(*stakingBalance)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get delegation earnings for cycle %d", cycle)
	}

	sbBigInt := gotezos.BigInt{*big.NewInt(int64(sbInt))}

	var delegationEarnings []DelegationEarnings
	for _, delegation := range *delegations {
		d, err := b.processDelegation(&processDelegationInput{
			delegation:           delegation,
			stakingBalance:       sbBigInt,
			frozenBalanceRewards: frozenBalanceRewards,
			blockHash:            networkCycle.BlockHash,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get delegation earnings for cycle %d", cycle)
		}

		delegationEarnings = append(delegationEarnings, *d)
	}

	return &delegationEarnings, nil
}

type processDelegationInput struct {
	delegation           string
	stakingBalance       gotezos.BigInt
	frozenBalanceRewards *gotezos.FrozenBalance
	blockHash            string
}

func (b *Baker) processDelegation(input *processDelegationInput) (*DelegationEarnings, error) {
	delegationEarnings := &DelegationEarnings{}
	balance, err := b.gt.Balance(input.blockHash, input.delegation)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to process delegation earnings for delegation %s", input.delegation)
	}

	bInt, err := strconv.Atoi(*balance)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to process delegation earnings for delegation %s", input.delegation)
	}

	bBigInt := gotezos.BigInt{*big.NewInt(int64(bInt))}
	delegationEarnings.Share = float64(bBigInt.Int64()) / float64(input.stakingBalance.Int64())

	rewards, err := strconv.Atoi(input.frozenBalanceRewards.Rewards)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to process delegation earnings for delegation %s", input.delegation)
	}

	grossRewardsFloat := delegationEarnings.Share * float64(rewards)
	feeFloat := grossRewardsFloat * b.fee

	delegationEarnings.GrossRewards = gotezos.BigInt{*big.NewInt(int64(grossRewardsFloat))}
	delegationEarnings.Fee = gotezos.BigInt{*big.NewInt(int64(feeFloat))}
	delegationEarnings.NetRewards = gotezos.BigInt{*delegationEarnings.GrossRewards.Sub(big.NewInt(delegationEarnings.GrossRewards.Int64()), big.NewInt(delegationEarnings.Fee.Int64()))}

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

func PrintDelegationEarnings(table bool, delegationEarnings *[]DelegationEarnings) {

}
