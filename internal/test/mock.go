package test

import (
	"errors"
	"math/big"
	"testing"

	gotezos "github.com/goat-systems/go-tezos/v3"
	"github.com/stretchr/testify/assert"
)

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
func (g *GoTezosMock) Balance(blockhash, address string) (*big.Int, error) {
	if g.BalanceErr {
		return big.NewInt(0), errors.New("failed to get balance")
	}
	return big.NewInt(5000000), nil
}

// FrozenBalance -
func (g *GoTezosMock) FrozenBalance(cycle int, delegate string) (gotezos.FrozenBalance, error) {
	if g.FrozenBalanceErr {
		return gotezos.FrozenBalance{}, errors.New("failed to get frozen balance")
	}
	return gotezos.FrozenBalance{
		Deposits: gotezos.NewInt(10000000000),
		Fees:     gotezos.NewInt(3000),
		Rewards:  gotezos.NewInt(70000000),
	}, nil
}

// DelegatedContractsAtCycle -
func (g *GoTezosMock) DelegatedContractsAtCycle(cycle int, delegate string) ([]*string, error) {
	if g.DelegatedContractsErr {
		return []*string{}, errors.New("failed to get delegated contracts at cycle")
	}
	strs := []string{
		"KT1LinsZAnyxajEv4eNFWtwHMdyhbJsGfvp3",
		"KT1K4xei3yozp7UP5rHV5wuoDzWwBXqCGRBt",
		"KT1GcSsQaTtMB2HvUKU9b6WRFUnGpGx9JwGk",
	}

	var rtnstrs []*string
	for i := range strs {
		rtnstrs = append(rtnstrs, &strs[i])
	}

	return rtnstrs, nil
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
func (g *GoTezosMock) StakingBalance(blockhash, delegate string) (*big.Int, error) {
	if g.StakingBalanceErr {
		return big.NewInt(0), errors.New("failed to get staking balance")
	}
	return big.NewInt(10000000000), nil
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

// CheckErr -
func CheckErr(t *testing.T, wantErr bool, errContains string, err error) {
	if wantErr {
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errContains)
	} else {
		assert.Nil(t, err)
	}
}
