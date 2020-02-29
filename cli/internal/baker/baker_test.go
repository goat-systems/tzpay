package baker

import (
	"context"
	"errors"
	"math/big"
	"sort"
	"testing"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v2/cli/internal/enviroment"
	"github.com/stretchr/testify/assert"
)

func Test_processDelegation(t *testing.T) {
	type want struct {
		err                bool
		errContains        string
		delegationEarnings *DelegationEarning
	}

	cases := []struct {
		name  string
		input *processDelegationInput
		want  want
	}{
		{
			"is successful",
			&processDelegationInput{
				stakingBalance: big.NewInt(100000000000),
				frozenBalanceRewards: &gotezos.FrozenBalance{
					Rewards: gotezos.Int{Big: big.NewInt(800000000)},
				},
			},
			want{
				false,
				"",
				&DelegationEarning{Fee: big.NewInt(4000000), GrossRewards: big.NewInt(80000000), NetRewards: big.NewInt(76000000), Share: 0.1},
			},
		},
		{
			"handles failure",
			&processDelegationInput{
				stakingBalance: big.NewInt(100000000000),
				frozenBalanceRewards: &gotezos.FrozenBalance{
					Rewards: gotezos.Int{Big: big.NewInt(800000000)},
				},
			},
			want{
				true,
				"failed to get balance",
				nil,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			baker := &Baker{
				gt: &gotezosMock{
					balanceErr: tt.want.err,
				},
			}
			out, err := baker.processDelegation(goldenContext, tt.input)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.delegationEarnings, out)
		})
	}
}

func Test_processDelegations(t *testing.T) {
	type want struct {
		err        bool
		errcount   int
		successful int
	}

	cases := []struct {
		name  string
		input *processDelegationsInput
		want  want
	}{
		{
			"is successful",
			&processDelegationsInput{
				delegations: &[]string{
					"some_delegation",
					"some_delegation1",
					"some_delegation2",
				},
				stakingBalance: big.NewInt(100000000000),
				frozenBalanceRewards: &gotezos.FrozenBalance{
					Rewards: gotezos.Int{Big: big.NewInt(800000000)},
				},
			},
			want{
				false,
				0,
				3,
			},
		},
		{
			"handles failure",
			&processDelegationsInput{
				delegations: &[]string{
					"some_delegation",
					"some_delegation1",
					"some_delegation2",
				},
				stakingBalance: big.NewInt(100000000000),
				frozenBalanceRewards: &gotezos.FrozenBalance{
					Rewards: gotezos.Int{Big: big.NewInt(800000000)},
				},
			},
			want{
				true,
				3,
				0,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			baker := &Baker{
				gt: &gotezosMock{
					balanceErr: tt.want.err,
				},
			}

			out := baker.proccessDelegations(goldenContext, tt.input)
			successful := 0
			errcount := 0

			for _, o := range out {
				if o.err != nil {
					errcount++
				} else {
					successful++
				}
			}

			assert.Equal(t, tt.want.successful, successful)
			assert.Equal(t, tt.want.errcount, errcount)
		})
	}
}

func Test_Payouts(t *testing.T) {
	type want struct {
		err         bool
		errContains string
		payouts     *Payout
	}

	cases := []struct {
		name  string
		input gotezos.IFace
		want  want
	}{
		{
			"handles FrozenBalance failue",
			&gotezosMock{
				frozenBalanceErr: true,
			},
			want{
				true,
				"failed to get delegation earnings for cycle 100: failed to get frozen balance",
				nil,
			},
		},
		{
			"handles DelegatedContractsAtCycle failue",
			&gotezosMock{
				delegatedContractsErr: true,
			},
			want{
				true,
				"failed to get delegation earnings for cycle 100: failed to get delegated contracts at cycle",
				nil,
			},
		},
		{
			"handles Cycle failue",
			&gotezosMock{
				cycleErr: true,
			},
			want{
				true,
				"failed to get delegation earnings for cycle 100: failed to get cycle",
				nil,
			},
		},
		{
			"handles StakingBalance failue",
			&gotezosMock{
				stakingBalanceErr: true,
			},
			want{
				true,
				"failed to get delegation earnings for cycle 100: failed to get staking balance",
				nil,
			},
		},
		{
			"is successful",
			&gotezosMock{},
			want{
				false,
				"",
				&Payout{
					DelegationEarnings: []DelegationEarning{
						DelegationEarning{
							Delegation:   "tz1b",
							Fee:          big.NewInt(3500000),
							GrossRewards: big.NewInt(70000000),
							NetRewards:   big.NewInt(66500000),
							Share:        1,
						},
						DelegationEarning{
							Delegation:   "tz1c",
							Fee:          big.NewInt(3500000),
							GrossRewards: big.NewInt(70000000),
							NetRewards:   big.NewInt(66500000),
							Share:        1,
						},
						DelegationEarning{
							Delegation:   "tz1a",
							Fee:          big.NewInt(3500000),
							GrossRewards: big.NewInt(70000000),
							NetRewards:   big.NewInt(66500000),
							Share:        1,
						},
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			baker := &Baker{gt: tt.input}
			payout, err := baker.Payouts(goldenContext, 100)
			checkErr(t, tt.want.err, tt.want.errContains, err)

			if tt.want.payouts != nil {
				sort.Sort(tt.want.payouts.DelegationEarnings)
			}

			if payout != nil {
				sort.Sort(payout.DelegationEarnings)
			}

			assert.Equal(t, tt.want.payouts, payout)
		})
	}
}

func Test_PayoutsSort(t *testing.T) {
	delegationEarnings := DelegationEarnings{}
	delegationEarnings = append(delegationEarnings,
		[]DelegationEarning{
			DelegationEarning{
				Delegation: "tz1c",
			},
			DelegationEarning{
				Delegation: "tz1a",
			},
			DelegationEarning{
				Delegation: "tz1b",
			},
		}...,
	)
	sort.Sort(&delegationEarnings)

	want := DelegationEarnings{}
	want = append(want,
		[]DelegationEarning{
			DelegationEarning{
				Delegation: "tz1a",
			},
			DelegationEarning{
				Delegation: "tz1b",
			},
			DelegationEarning{
				Delegation: "tz1c",
			},
		}...,
	)

	assert.Equal(t, want, delegationEarnings)
}

func Test_ForgePayout(t *testing.T) {
	type input struct {
		payout Payout
		gt     gotezos.IFace
	}

	type want struct {
		err         bool
		errContains string
		forge       string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				Payout{
					DelegationEarnings: DelegationEarnings{
						DelegationEarning{
							Delegation:   "somedelegation",
							GrossRewards: big.NewInt(1000000),
							NetRewards:   big.NewInt(900000),
						},
						DelegationEarning{
							Delegation:   "someotherdelegation",
							GrossRewards: big.NewInt(1000000),
							NetRewards:   big.NewInt(950000),
						},
					},
				},
				&gotezosMock{},
			},
			want{
				false,
				"",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			baker := &Baker{gt: &gotezosMock{}}
			forge, err := baker.ForgePayout(goldenContext, tt.input.payout)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.forge, forge)
		})
	}
}

func Test_constructPayoutContents(t *testing.T) {
	type input struct {
		counter int
		payout  Payout
	}

	cases := []struct {
		name  string
		input input
		want  []gotezos.Contents
	}{
		{
			"is successful",
			input{
				100,
				Payout{
					DelegationEarnings: DelegationEarnings{
						DelegationEarning{
							Delegation:   "somedelegation",
							GrossRewards: big.NewInt(1000000),
							NetRewards:   big.NewInt(900000),
						},
						DelegationEarning{
							Delegation:   "someotherdelegation",
							GrossRewards: big.NewInt(1000000),
							NetRewards:   big.NewInt(950000),
						},
					},
				},
			},
			[]gotezos.Contents{
				gotezos.Contents{
					Kind:        "transaction",
					Source:      "",
					Fee:         gotezos.Int{Big: big.NewInt(100000)},
					Counter:     gotezos.Int{Big: big.NewInt(101)},
					GasLimit:    gotezos.Int{Big: big.NewInt(100000)},
					Amount:      gotezos.Int{Big: big.NewInt(900000)},
					Destination: "somedelegation",
				},
				gotezos.Contents{
					Kind:        "transaction",
					Source:      "",
					Fee:         gotezos.Int{Big: big.NewInt(100000)},
					Counter:     gotezos.Int{Big: big.NewInt(102)},
					GasLimit:    gotezos.Int{Big: big.NewInt(100000)},
					Amount:      gotezos.Int{Big: big.NewInt(950000)},
					Destination: "someotherdelegation",
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			baker := &Baker{gt: &gotezosMock{}}
			contents := baker.constructPayoutContents(goldenContext, tt.input.counter, tt.input.payout)
			assert.Equal(t, tt.want, contents)
		})
	}
}

func Test_isValidForge(t *testing.T) {
	type input struct {
		contents []gotezos.Contents
		forge    string
		gt       gotezos.IFace
	}

	type want struct {
		err         bool
		errContains string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles failed unforge",
			input{
				[]gotezos.Contents{},
				"",
				&gotezosMock{
					unforgeOperationErr: true,
				},
			},
			want{
				true,
				"failed to unforge forge",
			},
		},
		{
			"is successful",
			input{
				[]gotezos.Contents{
					gotezos.Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          gotezos.Int{Big: big.NewInt(10100)},
						Counter:      gotezos.Int{Big: big.NewInt(10)},
						GasLimit:     gotezos.Int{Big: big.NewInt(10100)},
						StorageLimit: gotezos.Int{Big: big.NewInt(0)},
						Amount:       gotezos.Int{Big: big.NewInt(30)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         gotezos.TRANSACTIONOP,
					},
				},
				"6c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e001e000008ba0cb2fad622697145cf1665124096d25bc31e00",
				&gotezosMock{
					unforgeOperationErr: false,
					unforgeOperationRtn: []gotezos.Contents{
						gotezos.Contents{
							Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
							Fee:          gotezos.Int{Big: big.NewInt(10100)},
							Counter:      gotezos.Int{Big: big.NewInt(10)},
							GasLimit:     gotezos.Int{Big: big.NewInt(10100)},
							StorageLimit: gotezos.Int{Big: big.NewInt(0)},
							Amount:       gotezos.Int{Big: big.NewInt(30)},
							Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
							Kind:         gotezos.TRANSACTIONOP,
						},
					},
				},
			},
			want{
				true,
				"failed to unforge forge",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			baker := &Baker{gt: tt.input.gt}
			err := baker.isValidForge(tt.input.contents, tt.input.forge)
			checkErr(t, tt.want.err, tt.want.errContains, err)
		})
	}
}

func Test_isEqualContents(t *testing.T) {
	type input struct {
		forge   []gotezos.Contents
		unforge []gotezos.Contents
	}

	cases := []struct {
		name  string
		input input
		want  bool
	}{
		{
			"is equal",
			input{
				forge: []gotezos.Contents{
					gotezos.Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          gotezos.Int{Big: big.NewInt(10100)},
						Counter:      gotezos.Int{Big: big.NewInt(10)},
						GasLimit:     gotezos.Int{Big: big.NewInt(10100)},
						StorageLimit: gotezos.Int{Big: big.NewInt(0)},
						Amount:       gotezos.Int{Big: big.NewInt(30)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         gotezos.TRANSACTIONOP,
					},
					gotezos.Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8iwYpA",
						Fee:          gotezos.Int{Big: big.NewInt(500)},
						Counter:      gotezos.Int{Big: big.NewInt(10)},
						GasLimit:     gotezos.Int{Big: big.NewInt(10100)},
						StorageLimit: gotezos.Int{Big: big.NewInt(0)},
						Amount:       gotezos.Int{Big: big.NewInt(30)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         gotezos.TRANSACTIONOP,
					},
				},
				unforge: []gotezos.Contents{
					gotezos.Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          gotezos.Int{Big: big.NewInt(10100)},
						Counter:      gotezos.Int{Big: big.NewInt(10)},
						GasLimit:     gotezos.Int{Big: big.NewInt(10100)},
						StorageLimit: gotezos.Int{Big: big.NewInt(0)},
						Amount:       gotezos.Int{Big: big.NewInt(30)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         gotezos.TRANSACTIONOP,
					},
					gotezos.Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8iwYpA",
						Fee:          gotezos.Int{Big: big.NewInt(500)},
						Counter:      gotezos.Int{Big: big.NewInt(10)},
						GasLimit:     gotezos.Int{Big: big.NewInt(10100)},
						StorageLimit: gotezos.Int{Big: big.NewInt(0)},
						Amount:       gotezos.Int{Big: big.NewInt(30)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         gotezos.TRANSACTIONOP,
					},
				},
			},
			true,
		},
		{
			"is unequal",
			input{
				forge: []gotezos.Contents{
					gotezos.Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          gotezos.Int{Big: big.NewInt(10100)},
						Counter:      gotezos.Int{Big: big.NewInt(10)},
						GasLimit:     gotezos.Int{Big: big.NewInt(10100)},
						StorageLimit: gotezos.Int{Big: big.NewInt(0)},
						Amount:       gotezos.Int{Big: big.NewInt(30)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         gotezos.TRANSACTIONOP,
					},
					gotezos.Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8iwYpA",
						Fee:          gotezos.Int{Big: big.NewInt(500)},
						Counter:      gotezos.Int{Big: big.NewInt(10)},
						GasLimit:     gotezos.Int{Big: big.NewInt(10100)},
						StorageLimit: gotezos.Int{Big: big.NewInt(0)},
						Amount:       gotezos.Int{Big: big.NewInt(30)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         gotezos.TRANSACTIONOP,
					},
				},
				unforge: []gotezos.Contents{
					gotezos.Contents{
						Source:       "tz118bwVksXci8gUC2YpA",
						Fee:          gotezos.Int{Big: big.NewInt(10100)},
						Counter:      gotezos.Int{Big: big.NewInt(10)},
						GasLimit:     gotezos.Int{Big: big.NewInt(10100)},
						StorageLimit: gotezos.Int{Big: big.NewInt(0)},
						Amount:       gotezos.Int{Big: big.NewInt(30)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         gotezos.TRANSACTIONOP,
					},
					gotezos.Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8iwYpA",
						Fee:          gotezos.Int{Big: big.NewInt(500)},
						Counter:      gotezos.Int{Big: big.NewInt(10)},
						GasLimit:     gotezos.Int{Big: big.NewInt(10100)},
						StorageLimit: gotezos.Int{Big: big.NewInt(0)},
						Amount:       gotezos.Int{Big: big.NewInt(30)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         gotezos.TRANSACTIONOP,
					},
				},
			},
			false,
		},
		{
			"is unequal length",
			input{
				forge: []gotezos.Contents{
					gotezos.Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          gotezos.Int{Big: big.NewInt(10100)},
						Counter:      gotezos.Int{Big: big.NewInt(10)},
						GasLimit:     gotezos.Int{Big: big.NewInt(10100)},
						StorageLimit: gotezos.Int{Big: big.NewInt(0)},
						Amount:       gotezos.Int{Big: big.NewInt(30)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         gotezos.TRANSACTIONOP,
					},
					gotezos.Contents{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8iwYpA",
						Fee:          gotezos.Int{Big: big.NewInt(500)},
						Counter:      gotezos.Int{Big: big.NewInt(10)},
						GasLimit:     gotezos.Int{Big: big.NewInt(10100)},
						StorageLimit: gotezos.Int{Big: big.NewInt(0)},
						Amount:       gotezos.Int{Big: big.NewInt(30)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         gotezos.TRANSACTIONOP,
					},
				},
				unforge: []gotezos.Contents{
					gotezos.Contents{
						Source:       "tz118bwVksXci8gUC2YpA",
						Fee:          gotezos.Int{Big: big.NewInt(10100)},
						Counter:      gotezos.Int{Big: big.NewInt(10)},
						GasLimit:     gotezos.Int{Big: big.NewInt(10100)},
						StorageLimit: gotezos.Int{Big: big.NewInt(0)},
						Amount:       gotezos.Int{Big: big.NewInt(30)},
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         gotezos.TRANSACTIONOP,
					},
				},
			},
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ok := isEqualContents(tt.input.forge, tt.input.unforge)
			assert.Equal(t, tt.want, ok)
		})
	}
}

func checkErr(t *testing.T, wantErr bool, errContains string, err error) {
	if wantErr {
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errContains)
	} else {
		assert.Nil(t, err)
	}
}

var goldenContext = context.WithValue(
	context.TODO(),
	enviroment.ENVIROMENTKEY,
	&enviroment.ContextEnviroment{
		BakersFee:      0.05,
		BlackList:      "somehash, somehash1",
		Delegate:       "somedelegate",
		GasLimit:       100000,
		HostNode:       "http://somenode.com:8732",
		MinimumPayment: 1000,
		NetworkFee:     100000,
	},
)

type gotezosMock struct {
	gotezos.IFace
	headErr               bool
	counterErr            bool
	forgeOperationErr     bool
	balanceErr            bool
	frozenBalanceErr      bool
	delegatedContractsErr bool
	cycleErr              bool
	stakingBalanceErr     bool
	unforgeOperationErr   bool
	unforgeOperationRtn   []gotezos.Contents
}

func (g *gotezosMock) Head() (*gotezos.Block, error) {
	if g.headErr {
		return &gotezos.Block{}, errors.New("failed to get block")
	}
	return &gotezos.Block{
		Hash: "BLfEWKVudXH15N8nwHZehyLNjRuNLoJavJDjSZ7nq8ggfzbZ18p",
	}, nil
}

func (g *gotezosMock) Counter(blockhash, pkh string) (*int, error) {
	counter := 0
	if g.counterErr {
		return &counter, errors.New("failed to get block")
	}
	counter = 100
	return &counter, nil
}

func (g *gotezosMock) Balance(blockhash, address string) (*big.Int, error) {
	if g.balanceErr {
		return big.NewInt(0), errors.New("failed to get balance")
	}
	return big.NewInt(10000000000), nil
}

func (g *gotezosMock) FrozenBalance(cycle int, delegate string) (*gotezos.FrozenBalance, error) {
	if g.frozenBalanceErr {
		return &gotezos.FrozenBalance{}, errors.New("failed to get frozen balance")
	}
	return &gotezos.FrozenBalance{
		Deposits: gotezos.Int{Big: big.NewInt(10000000000)},
		Fees:     gotezos.Int{Big: big.NewInt(3000)},
		Rewards:  gotezos.Int{Big: big.NewInt(70000000)},
	}, nil
}

func (g *gotezosMock) DelegatedContractsAtCycle(cycle int, delegate string) (*[]string, error) {
	if g.delegatedContractsErr {
		return &[]string{}, errors.New("failed to get delegated contracts at cycle")
	}
	return &[]string{
		"tz1a",
		"tz1b",
		"tz1c",
	}, nil
}

func (g *gotezosMock) Cycle(cycle int) (*gotezos.Cycle, error) {
	if g.cycleErr {
		return &gotezos.Cycle{}, errors.New("failed to get cycle")
	}
	return &gotezos.Cycle{
		RandomSeed:   "some_seed",
		RollSnapshot: 10,
		BlockHash:    "some_hash",
	}, nil
}

func (g *gotezosMock) StakingBalance(blockhash, delegate string) (*big.Int, error) {
	if g.stakingBalanceErr {
		return big.NewInt(0), errors.New("failed to get staking balance")
	}
	return big.NewInt(10000000000), nil
}
