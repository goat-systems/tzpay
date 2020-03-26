package test

import (
	"errors"
	"math/big"

	gotezos "github.com/goat-systems/go-tezos/v2"
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
func (g *GoTezosMock) Counter(blockhash, pkh string) (*int, error) {
	counter := 0
	if g.CounterErr {
		return &counter, errors.New("failed to get block")
	}
	counter = 100
	return &counter, nil
}

// Balance -
func (g *GoTezosMock) Balance(blockhash, address string) (*big.Int, error) {
	if g.BalanceErr {
		return big.NewInt(0), errors.New("failed to get balance")
	}
	return big.NewInt(10000000000), nil
}

// FrozenBalance -
func (g *GoTezosMock) FrozenBalance(cycle int, delegate string) (*gotezos.FrozenBalance, error) {
	if g.FrozenBalanceErr {
		return &gotezos.FrozenBalance{}, errors.New("failed to get frozen balance")
	}
	return &gotezos.FrozenBalance{
		Deposits: gotezos.Int{Big: big.NewInt(10000000000)},
		Fees:     gotezos.Int{Big: big.NewInt(3000)},
		Rewards:  gotezos.Int{Big: big.NewInt(70000000)},
	}, nil
}

// DelegatedContractsAtCycle -
func (g *GoTezosMock) DelegatedContractsAtCycle(cycle int, delegate string) (*[]string, error) {
	if g.DelegatedContractsErr {
		return &[]string{}, errors.New("failed to get delegated contracts at cycle")
	}
	return &[]string{
		"tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
		"tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
		"tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
	}, nil
}

// Cycle -
func (g *GoTezosMock) Cycle(cycle int) (*gotezos.Cycle, error) {
	if g.CycleErr {
		return &gotezos.Cycle{}, errors.New("failed to get cycle")
	}
	return &gotezos.Cycle{
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
func (g *GoTezosMock) InjectionOperation(input *gotezos.InjectionOperationInput) (*[]byte, error) {
	if g.StakingBalanceErr {
		return nil, errors.New("failed to inject operation")
	}
	resp := []byte("ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2rfb1M")
	return &resp, nil
}
