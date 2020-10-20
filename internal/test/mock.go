package test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/goat-systems/go-tezos/v3/rpc"
	"github.com/goat-systems/tzpay/v3/internal/tzkt"
	"github.com/stretchr/testify/assert"
)

// TzktMock is a test helper mocking the tzkt package
type TzktMock struct {
	tzkt.IFace
	TransactionsErr bool
	RewardsSplitErr bool
}

var _ rpc.IFace = &RPCMock{}

func (t *TzktMock) GetTransactions(options ...tzkt.URLParameters) ([]tzkt.Transaction, error) {
	if t.TransactionsErr {
		return []tzkt.Transaction{}, errors.New("failed to get transaction")
	}

	return []tzkt.Transaction{
		{
			Sender: struct {
				Name    string "json:\"name\""
				Address string "json:\"address\""
			}{
				Address: "tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV",
			}, // sanity
		},
	}, nil

}

func (t *TzktMock) GetRewardsSplit(delegate string, cycle int, options ...tzkt.URLParameters) (tzkt.RewardsSplit, error) {
	if t.RewardsSplitErr {
		return tzkt.RewardsSplit{}, errors.New("failed to get rewards split")
	}

	var rewardsSplit tzkt.RewardsSplit
	v := []byte(`{"cycle":270,"stakingBalance":740613513605,"delegatedBalance":555430526884,"numDelegators":107,"expectedBlocks":4.43,"expectedEndorsements":141.71,"futureBlocks":0,"futureBlockRewards":0,"futureBlockDeposits":0,"ownBlocks":5,"ownBlockRewards":191250000,"extraBlocks":0,"extraBlockRewards":0,"missedOwnBlocks":2,"missedOwnBlockRewards":77500000,"missedExtraBlocks":0,"missedExtraBlockRewards":0,"uncoveredOwnBlocks":0,"uncoveredOwnBlockRewards":0,"uncoveredExtraBlocks":0,"uncoveredExtraBlockRewards":0,"blockDeposits":2560000000,"futureEndorsements":0,"futureEndorsementRewards":0,"futureEndorsementDeposits":0,"endorsements":126,"endorsementRewards":157500000,"missedEndorsements":16,"missedEndorsementRewards":20000000,"uncoveredEndorsements":0,"uncoveredEndorsementRewards":0,"endorsementDeposits":8064000000,"ownBlockFees":47180,"extraBlockFees":0,"missedOwnBlockFees":54607,"missedExtraBlockFees":0,"uncoveredOwnBlockFees":0,"uncoveredExtraBlockFees":0,"doubleBakingRewards":0,"doubleBakingLostDeposits":0,"doubleBakingLostRewards":0,"doubleBakingLostFees":0,"doubleEndorsingRewards":0,"doubleEndorsingLostDeposits":0,"doubleEndorsingLostRewards":0,"doubleEndorsingLostFees":0,"revelationRewards":0,"revelationLostRewards":0,"revelationLostFees":0,"delegators":[{"address":"tz1icdoLr8vof5oXiEKCFSyrVoouGiKDQ3Gd","balance":60545965782,"currentBalance":60739073316,"emptied":false},{"address":"KT1FPyY6mAhnzyVGP8ApGvuRyF7SKcT9TDWy","balance":60075572992,"currentBalance":60267312348,"emptied":false},{"address":"KT1LgkGigaMrnim3TonQWfwDHnM3fHkF1jMv","balance":57461165021,"currentBalance":57644560137,"emptied":false},{"address":"KT1C8S2vLYbzgQHhdC8MBehunhcp1Q9hj6MC","balance":55305195039,"currentBalance":176566401,"emptied":false}]}`)
	json.Unmarshal(v, &rewardsSplit)

	return rewardsSplit, nil
}

// RPCMock is a test helper mocking the go-tezos/rpc lib
type RPCMock struct {
	rpc.IFace
	HeadErr               bool
	CounterErr            bool
	BalanceErr            bool
	FrozenBalanceErr      bool
	DelegatedContractsErr bool
	CycleErr              bool
	StakingBalanceErr     bool
	InjectionOperationErr bool
	OperationHashesErr    bool
	ForgeOperationErr     bool
	ContractStorageErr    bool
	BigMapErr             bool
	BakingRightsErr       bool
	EndorsingRightsErr    bool
}

// EndorsingRights -
func (r *RPCMock) EndorsingRights(input rpc.EndorsingRightsInput) (*rpc.EndorsingRights, error) {
	if r.EndorsingRightsErr {
		return &rpc.EndorsingRights{}, errors.New("failed to get endorsing rights")
	}

	return &rpc.EndorsingRights{
		{
			Level:    100,
			Delegate: "some_delegate",
		},
	}, nil
}

// BakingRights -
func (r *RPCMock) BakingRights(input rpc.BakingRightsInput) (*rpc.BakingRights, error) {
	if r.BakingRightsErr {
		return &rpc.BakingRights{}, errors.New("failed to get baking rights")
	}

	return &rpc.BakingRights{
		{
			Level:    100,
			Delegate: "some_delegate",
		},
	}, nil
}

// Head -
func (r *RPCMock) Head() (*rpc.Block, error) {
	if r.HeadErr {
		return &rpc.Block{}, errors.New("failed to get block")
	}
	return &rpc.Block{
		Hash: "BLfEWKVudXH15N8nwHZehyLNjRuNLoJavJDjSZ7nq8ggfzbZ18p",
	}, nil
}

// Counter -
func (r *RPCMock) Counter(blockhash, pkh string) (int, error) {
	counter := 0
	if r.CounterErr {
		return counter, errors.New("failed to get counter")
	}
	counter = 100
	return counter, nil
}

// Balance -
func (r *RPCMock) Balance(input rpc.BalanceInput) (int, error) {
	if r.BalanceErr {
		return 0, errors.New("failed to get balance")
	}
	return 5000000, nil
}

// FrozenBalance -
func (r *RPCMock) FrozenBalance(cycle int, delegate string) (rpc.FrozenBalance, error) {
	if r.FrozenBalanceErr {
		return rpc.FrozenBalance{}, errors.New("failed to get frozen balance")
	}
	return rpc.FrozenBalance{
		Deposits: 10000000000,
		Fees:     3000,
		Rewards:  70000000,
	}, nil
}

// DelegatedContracts -
func (r *RPCMock) DelegatedContracts(input rpc.DelegatedContractsInput) ([]string, error) {
	if r.DelegatedContractsErr {
		return []string{}, errors.New("failed to get delegated contracts at cycle")
	}

	return []string{
		"KT1LinsZAnyxajEv4eNFWtwHMdyhbJsGfvp3",
		"KT1K4xei3yozp7UP5rHV5wuoDzWwBXqCGRBt",
		"KT1GcSsQaTtMB2HvUKU9b6WRFUnGpGx9JwGk",
	}, nil
}

// Cycle -
func (r *RPCMock) Cycle(cycle int) (rpc.Cycle, error) {
	if r.CycleErr {
		return rpc.Cycle{}, errors.New("failed to get cycle")
	}
	return rpc.Cycle{
		RandomSeed:   "some_seed",
		RollSnapshot: 10,
		BlockHash:    "some_hash",
	}, nil
}

// StakingBalance -
func (r *RPCMock) StakingBalance(input rpc.StakingBalanceInput) (int, error) {
	if r.StakingBalanceErr {
		return 0, errors.New("failed to get staking balance")
	}
	return 10000000000, nil
}

// InjectionOperation -
func (r *RPCMock) InjectionOperation(input rpc.InjectionOperationInput) (string, error) {
	if r.InjectionOperationErr {
		return "", errors.New("failed to inject operation")
	}
	return "ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2rfb1M", nil
}

// OperationHashes -
func (r *RPCMock) OperationHashes(blockhash string) ([][]string, error) {
	if r.OperationHashesErr {
		return nil, errors.New("failed to get operation hashes")
	}

	return [][]string{
		{
			"ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2rfb1M",
			"ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2rfFGD",
		},
	}, nil
}

func (r *RPCMock) ContractStorage(blockhash string, KT1 string) ([]byte, error) {
	if r.ContractStorageErr {
		return nil, errors.New("failed to get contract storage")
	}

	return []byte(`{"prim":"Pair","args":[{"int":"16033"},{"prim":"Pair","args":[{"prim":"Pair","args":[{"prim":"False"},{"prim":"Pair","args":[{"prim":"False"},{"int":"23567891"}]}]},{"prim":"Pair","args":[{"prim":"Pair","args":[{"string":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV"},{"string":"KT1GQcLae1ve1ZEPNfD9z1dyv5ev9ki39SNW"}]},{"prim":"Pair","args":[{"int":"123456"},{"int":"23567891"}]}]}]}]}`), nil
}

func (r *RPCMock) BigMap(input rpc.BigMapInput) ([]byte, error) {
	if r.BigMapErr {
		return nil, errors.New("failed to get contract storage")
	}

	return []byte(`{"prim":"Pair","args":[{"int":"23567891"},[]]}`), nil
}

// CheckErr -
func CheckErr(t *testing.T, wantErr bool, errContains string, err error) {
	if wantErr {
		assert.Error(t, err)
		if err != nil {
			assert.Contains(t, err.Error(), errContains)
		}
	} else {
		assert.Nil(t, err)
	}
}
