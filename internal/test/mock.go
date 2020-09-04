package test

import (
	"errors"
	"testing"

	gotezos "github.com/goat-systems/go-tezos/v3"
	"github.com/goat-systems/tzpay/v2/internal/tzkt"
	"github.com/stretchr/testify/assert"
)

// TzktMock is a test helper mocking the tzkt package
type TzktMock struct {
	tzkt.IFace
	TransactionsErr bool
}

func (t *TzktMock) GetTransactions(options ...tzkt.URLParameters) ([]tzkt.Transaction, error) {
	if t.TransactionsErr {
		return []tzkt.Transaction{}, errors.New("failed to get counter")
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

// GoTezosMock is a test helper mocking the GoTezos lib
type GoTezosMock struct {
	gotezos.IFace
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
}

// Head -
func (g *GoTezosMock) Head() (*gotezos.Block, error) {
	if g.HeadErr {
		return &gotezos.Block{}, errors.New("failed to get block")
	}
	return &gotezos.Block{
		Hash: "BLfEWKVudXH15N8nwHZehyLNjRuNLoJavJDjSZ7nq8ggfzbZ18p",
	}, nil
}

// Counter -
func (g *GoTezosMock) Counter(blockhash, pkh string) (int, error) {
	counter := 0
	if g.CounterErr {
		return counter, errors.New("failed to get counter")
	}
	counter = 100
	return counter, nil
}

// Balance -
func (g *GoTezosMock) Balance(blockhash, address string) (int, error) {
	if g.BalanceErr {
		return 0, errors.New("failed to get balance")
	}
	return 5000000, nil
}

// FrozenBalance -
func (g *GoTezosMock) FrozenBalance(cycle int, delegate string) (gotezos.FrozenBalance, error) {
	if g.FrozenBalanceErr {
		return gotezos.FrozenBalance{}, errors.New("failed to get frozen balance")
	}
	return gotezos.FrozenBalance{
		Deposits: 10000000000,
		Fees:     3000,
		Rewards:  70000000,
	}, nil
}

// DelegatedContracts -
func (g *GoTezosMock) DelegatedContracts(input gotezos.DelegatedContractsInput) ([]string, error) {
	if g.DelegatedContractsErr {
		return []string{}, errors.New("failed to get delegated contracts at cycle")
	}

	return []string{
		"KT1LinsZAnyxajEv4eNFWtwHMdyhbJsGfvp3",
		"KT1K4xei3yozp7UP5rHV5wuoDzWwBXqCGRBt",
		"KT1GcSsQaTtMB2HvUKU9b6WRFUnGpGx9JwGk",
	}, nil
}

// Cycle -
func (g *GoTezosMock) Cycle(cycle int) (gotezos.Cycle, error) {
	if g.CycleErr {
		return gotezos.Cycle{}, errors.New("failed to get cycle")
	}
	return gotezos.Cycle{
		RandomSeed:   "some_seed",
		RollSnapshot: 10,
		BlockHash:    "some_hash",
	}, nil
}

// StakingBalance -
func (g *GoTezosMock) StakingBalance(input gotezos.StakingBalanceInput) (int, error) {
	if g.StakingBalanceErr {
		return 0, errors.New("failed to get staking balance")
	}
	return 10000000000, nil
}

// InjectionOperation -
func (g *GoTezosMock) InjectionOperation(input gotezos.InjectionOperationInput) (string, error) {
	if g.InjectionOperationErr {
		return "", errors.New("failed to inject operation")
	}
	return "ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2rfb1M", nil
}

// OperationHashes -
func (g *GoTezosMock) OperationHashes(blockhash string) ([][]string, error) {
	if g.OperationHashesErr {
		return nil, errors.New("failed to get operation hashes")
	}

	return [][]string{
		{
			"ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2rfb1M",
			"ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2rfFGD",
		},
	}, nil
}

func (g *GoTezosMock) ContractStorage(blockhash string, KT1 string) ([]byte, error) {
	if g.ContractStorageErr {
		return nil, errors.New("failed to get contract storage")
	}

	return []byte(`{"prim":"Pair","args":[{"int":"16033"},{"prim":"Pair","args":[{"prim":"Pair","args":[{"prim":"False"},{"prim":"Pair","args":[{"prim":"False"},{"int":"23567891"}]}]},{"prim":"Pair","args":[{"prim":"Pair","args":[{"string":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV"},{"string":"KT1GQcLae1ve1ZEPNfD9z1dyv5ev9ki39SNW"}]},{"prim":"Pair","args":[{"int":"123456"},{"int":"23567891"}]}]}]}]}`), nil
}

func (g *GoTezosMock) BigMap(input gotezos.BigMapInput) ([]byte, error) {
	if g.BigMapErr {
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
