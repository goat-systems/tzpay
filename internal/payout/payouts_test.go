package payout

import (
	"math/big"
	"sort"
	"testing"
	"time"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v2/internal/db/model"
	"github.com/goat-systems/tzpay/v2/internal/test"
	"github.com/stretchr/testify/assert"
)

// func Test_PayoutsSort(t *testing.T) {
// 	delegationEarnings := model.DelegationEarnings{}
// 	delegationEarnings = append(delegationEarnings,
// 		[]model.DelegationEarning{
// 			model.DelegationEarning{
// 				Address: "tz1c",
// 			},
// 			model.DelegationEarning{
// 				Address: "tz1a",
// 			},
// 			model.DelegationEarning{
// 				Address: "tz1b",
// 			},
// 		}...,
// 	)
// 	sort.Sort(&delegationEarnings)

// 	want := model.DelegationEarnings{}
// 	want = append(want,
// 		[]model.DelegationEarning{
// 			model.DelegationEarning{
// 				Address: "tz1a",
// 			},
// 			model.DelegationEarning{
// 				Address: "tz1b",
// 			},
// 			model.DelegationEarning{
// 				Address: "tz1c",
// 			},
// 		}...,
// 	)

// 	assert.Equal(t, want, delegationEarnings)
// }

func Test_Execute(t *testing.T) {
	type want struct {
		err         bool
		errContains string
		payouts     *model.Payout
	}

	cases := []struct {
		name  string
		input gotezos.IFace
		want  want
	}{
		{
			"handles FrozenBalance failue",
			&test.GoTezosMock{
				FrozenBalanceErr: true,
			},
			want{
				true,
				"failed to get delegation earnings for cycle 100: failed to get frozen balance",
				nil,
			},
		},
		{
			"handles DelegatedContractsAtCycle failue",
			&test.GoTezosMock{
				DelegatedContractsErr: true,
			},
			want{
				true,
				"failed to get delegation earnings for cycle 100: failed to get delegated contracts at cycle",
				nil,
			},
		},
		{
			"handles Cycle failue",
			&test.GoTezosMock{
				CycleErr: true,
			},
			want{
				true,
				"failed to get delegation earnings for cycle 100: failed to get cycle",
				nil,
			},
		},
		{
			"handles StakingBalance failue",
			&test.GoTezosMock{
				StakingBalanceErr: true,
			},
			want{
				true,
				"failed to get delegation earnings for cycle 100: failed to get staking balance",
				nil,
			},
		},
		{
			"is successful",
			&test.GoTezosMock{},
			want{
				false,
				"",
				&model.Payout{
					DelegationEarnings: model.DelegationEarnings{
						model.DelegationEarning{
							Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
							Fee:          big.NewInt(3500000),
							GrossRewards: big.NewInt(70000000),
							NetRewards:   big.NewInt(66500000),
							Share:        1,
						},
						model.DelegationEarning{
							Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
							Fee:          big.NewInt(3500000),
							GrossRewards: big.NewInt(70000000),
							NetRewards:   big.NewInt(66500000),
							Share:        1,
						},
						model.DelegationEarning{
							Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
							Fee:          big.NewInt(3500000),
							GrossRewards: big.NewInt(70000000),
							NetRewards:   big.NewInt(66500000),
							Share:        1,
						},
					},
					DelegateEarnings: model.DelegateEarnings{
						Address: "somedelegate",
						Fees:    big.NewInt(10500000),
						Share:   1,
						Rewards: big.NewInt(70000000),
						Net:     big.NewInt(80500000),
					},
					CycleHash:      "some_hash",
					Cycle:          100,
					FrozenBalance:  big.NewInt(70000000),
					StakingBalance: big.NewInt(10000000000),
					Operations:     nil,
					OperationsLink: nil,
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := &Payout{gt: tt.input}
			report, err := payout.Execute()
			checkErr(t, tt.want.err, tt.want.errContains, err)

			if tt.want.payouts != nil {
				sort.Sort(tt.want.payouts.DelegationEarnings)
			}

			if payout != nil {
				sort.Sort(report.DelegationEarnings)
			}

			assert.Equal(t, tt.want.payouts, payout)
		})
	}
}
func Test_processDelegate(t *testing.T) {
	type input struct {
		gt            gotezos.IFace
		delegateInput processDelegateInput
	}

	type want struct {
		err              bool
		errContains      string
		delegateEarnings DelegateEarnings
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				gt: &test.GoTezosMock{},
				delegateInput: processDelegateInput{
					delegate: "some_delegate",
					delegations: []DelegationEarning{
						{
							Fee: big.NewInt(1000),
						},
						{
							Fee: big.NewInt(1000),
						},
						{
							Fee: big.NewInt(1000),
						},
					},
					stakingBalance: big.NewInt(10000000),
					frozenBalanceRewards: gotezos.FrozenBalance{
						Rewards: gotezos.NewInt(700),
					},
					blockHash: "block_hash",
				},
			},
			want{
				false,
				"",
				DelegateEarnings{Address: "some_delegate", Fees: big.NewInt(3000), Share: 1000, Rewards: big.NewInt(700000), Net: big.NewInt(703000)},
			},
		},
		{
			"handles failure to get balance",
			input{
				gt: &test.GoTezosMock{BalanceErr: true},
				delegateInput: processDelegateInput{
					delegate: "some_delegate",
					delegations: []DelegationEarning{
						{
							Fee: big.NewInt(1000),
						},
						{
							Fee: big.NewInt(1000),
						},
						{
							Fee: big.NewInt(1000),
						},
					},
					stakingBalance: big.NewInt(10000000),
					frozenBalanceRewards: gotezos.FrozenBalance{
						Rewards: gotezos.NewInt(700),
					},
					blockHash: "block_hash",
				},
			},
			want{
				true,
				"failed to process delegate earnings",
				DelegateEarnings{Address: "some_delegate", Fees: nil, Share: 0, Rewards: nil, Net: big.NewInt(0)},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := Payout{
				gt: tt.input.gt,
			}
			delegateEarnings, err := payout.processDelegate(tt.input.delegateInput)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.delegateEarnings, delegateEarnings)
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
		input processDelegationsInput
		want  want
	}{
		{
			"is successful",
			processDelegationsInput{
				delegations: []*string{
					strToPointer("some_delegation"),
					strToPointer("some_delegation1"),
					strToPointer("some_delegation2"),
				},
				stakingBalance: big.NewInt(100000000000),
				frozenBalanceRewards: gotezos.FrozenBalance{
					Rewards: gotezos.NewInt(800000000),
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
			processDelegationsInput{
				delegations: []*string{
					strToPointer("some_delegation"),
					strToPointer("some_delegation1"),
					strToPointer("some_delegation2"),
				},
				stakingBalance: big.NewInt(100000000000),
				frozenBalanceRewards: gotezos.FrozenBalance{
					Rewards: gotezos.NewInt(800000000),
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
			payout := &Payout{
				gt: &test.GoTezosMock{
					BalanceErr: tt.want.err,
				},
			}

			out := payout.proccessDelegations(tt.input)
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

func Test_processDelegation(t *testing.T) {
	type want struct {
		err                bool
		errContains        string
		delegationEarnings *DelegationEarning
	}

	cases := []struct {
		name  string
		input processDelegationInput
		want  want
	}{
		{
			"is successful",
			processDelegationInput{
				stakingBalance: big.NewInt(100000000000),
				frozenBalanceRewards: gotezos.FrozenBalance{
					Rewards: gotezos.NewInt(800000000),
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
			processDelegationInput{
				stakingBalance: big.NewInt(100000000000),
				frozenBalanceRewards: gotezos.FrozenBalance{
					Rewards: gotezos.NewInt(800000000),
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
			payout := &Payout{
				gt: &test.GoTezosMock{
					BalanceErr: tt.want.err,
				},
				bakerFee: 0.05,
			}
			out, err := payout.processDelegation(tt.input)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.delegationEarnings, out)
		})
	}
}

func Test_getOperationHexStrings(t *testing.T) {
	type input struct {
		gt                 gotezos.IFace
		delegationEarnings DelegationEarnings
	}

	type want struct {
		err         bool
		errContains string
		ophexes     []string
	}
	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				&test.GoTezosMock{},
				DelegationEarnings{
					DelegationEarning{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(900000),
					},
					DelegationEarning{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(950000),
					},
				},
			},
			want{
				false,
				"",
				[]string{
					"7cc601d2729c90b267e6a79d902f8b048d37fd990f2f7447efefb0cfb2f8e8a46c004b04ad1e57c2f13b61b3d2c95b3073d961a4132b00650000a0f7360000056a59972593bdc74a5295671c8f5d43c21348da006c004b04ad1e57c2f13b61b3d2c95b3073d961a4132b00660000f0fd390000056a59972593bdc74a5295671c8f5d43c21348da00",
				},
			},
		},
		{
			"handles failure to get head",
			input{
				&test.GoTezosMock{HeadErr: true},
				DelegationEarnings{
					DelegationEarning{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(900000),
					},
					DelegationEarning{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(950000),
					},
				},
			},
			want{
				true,
				"failed to get operation hex string: failed to get block",
				nil,
			},
		},
		{
			"handles failure to get counter",
			input{
				&test.GoTezosMock{CounterErr: true},
				DelegationEarnings{
					DelegationEarning{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(900000),
					},
					DelegationEarning{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(950000),
					},
				},
			},
			want{
				true,
				"failed to get operation hex string: failed to get counter",
				nil,
			},
		},
		{
			"handles failure to forge",
			input{
				&test.GoTezosMock{},
				DelegationEarnings{
					DelegationEarning{
						Address:      "invalid_addr",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(900000),
					},
					DelegationEarning{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(950000),
					},
				},
			},
			want{
				true,
				"failed to get operation hex string: failed to forge payout",
				nil,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := Payout{
				gt:        tt.input.gt,
				batchSize: 10,
				wallet: gotezos.Wallet{
					Address: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
				},
			}
			ophexes, err := payout.getOperationHexStrings(tt.input.delegationEarnings)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.ophexes, ophexes)
		})
	}
}

func Test_forgeOperation(t *testing.T) {
	type input struct {
		counter            int
		delegationEarnings DelegationEarnings
		gt                 gotezos.IFace
	}

	type want struct {
		err         bool
		errContains string
		ophash      string
		counter     int
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				5,
				DelegationEarnings{
					DelegationEarning{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(900000),
					},
					DelegationEarning{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(950000),
					},
				},
				&test.GoTezosMock{},
			},
			want{
				err:         false,
				errContains: "",
				ophash:      "7cc601d2729c90b267e6a79d902f8b048d37fd990f2f7447efefb0cfb2f8e8a46c004b04ad1e57c2f13b61b3d2c95b3073d961a4132b00060000a0f7360000056a59972593bdc74a5295671c8f5d43c21348da006c004b04ad1e57c2f13b61b3d2c95b3073d961a4132b00070000f0fd390000056a59972593bdc74a5295671c8f5d43c21348da00",
				counter:     7,
			},
		},
		{
			"handles failure to get head",
			input{
				5,
				DelegationEarnings{
					DelegationEarning{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(900000),
					},
					DelegationEarning{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(950000),
					},
				},
				&test.GoTezosMock{HeadErr: true},
			},
			want{
				err:         true,
				errContains: "failed to forge payout: failed to get block",
				ophash:      "",
				counter:     5,
			},
		},
		{
			"handles failure to forge",
			input{
				5,
				DelegationEarnings{
					DelegationEarning{
						Address:      "invalid_addr",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(900000),
					},
					DelegationEarning{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(950000),
					},
				},
				&test.GoTezosMock{},
			},
			want{
				err:         true,
				errContains: "failed to forge payout: failed to forge operation",
				ophash:      "",
				counter:     7,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := &Payout{
				gt: tt.input.gt,
				wallet: gotezos.Wallet{
					Address: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
				},
			}
			ophash, counter, err := payout.forgeOperation(tt.input.counter, tt.input.delegationEarnings)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.counter, counter)
			assert.Equal(t, tt.want.ophash, ophash)
		})
	}

}

func Test_constructPayoutContents(t *testing.T) {
	type input struct {
		counter            int
		blacklist          []string
		delegationEarnings DelegationEarnings
	}

	cases := []struct {
		name  string
		input input
		want  []gotezos.ForgeTransactionOperationInput
	}{
		{
			"is successful",
			input{
				100,
				[]string{},
				DelegationEarnings{
					DelegationEarning{
						Address:      "somedelegation",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(900000),
					},
					DelegationEarning{
						Address:      "someotherdelegation",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(950000),
					},
				},
			},
			[]gotezos.ForgeTransactionOperationInput{
				{
					Source:       "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
					Fee:          gotezos.NewInt(0),
					Counter:      101,
					GasLimit:     gotezos.NewInt(0),
					StorageLimit: gotezos.NewInt(0),
					Amount:       gotezos.NewInt(900000),
					Destination:  "somedelegation",
				},
				{
					Source:       "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
					Fee:          gotezos.NewInt(0),
					Counter:      102,
					GasLimit:     gotezos.NewInt(0),
					StorageLimit: gotezos.NewInt(0),
					Amount:       gotezos.NewInt(950000),
					Destination:  "someotherdelegation",
				},
			},
		},
		{
			"is successful in respecting blacklist",
			input{
				100,
				[]string{
					"someotherdelegation",
				},
				DelegationEarnings{
					DelegationEarning{
						Address:      "somedelegation",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(900000),
					},
					DelegationEarning{
						Address:      "someotherdelegation",
						GrossRewards: big.NewInt(1000000),
						NetRewards:   big.NewInt(950000),
					},
				},
			},
			[]gotezos.ForgeTransactionOperationInput{
				{
					Source:       "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
					Fee:          gotezos.NewInt(0),
					Counter:      101,
					GasLimit:     gotezos.NewInt(0),
					StorageLimit: gotezos.NewInt(0),
					Amount:       gotezos.NewInt(900000),
					Destination:  "somedelegation",
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := &Payout{
				wallet: gotezos.Wallet{
					Address: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
				},
				blacklist: tt.input.blacklist,
			}
			contents, _ := payout.constructPayoutContents(tt.input.counter, tt.input.delegationEarnings)
			assert.Equal(t, tt.want, contents)
		})
	}
}

func Test_batch(t *testing.T) {
	cases := []struct {
		name  string
		input DelegationEarnings
		want  []DelegationEarnings
	}{
		{
			"is successful with multiple batches",
			DelegationEarnings{
				DelegationEarning{Address: "some_addr"},
				DelegationEarning{Address: "some_addr1"},
				DelegationEarning{Address: "some_addr2"},
				DelegationEarning{Address: "some_addr3"},
				DelegationEarning{Address: "some_addr4"},
			},
			[]DelegationEarnings{
				{
					DelegationEarning{Address: "some_addr"},
					DelegationEarning{Address: "some_addr1"},
				},
				{
					DelegationEarning{Address: "some_addr2"},
					DelegationEarning{Address: "some_addr3"}},
				{
					DelegationEarning{Address: "some_addr4"},
				},
			},
		},
		{
			"is successful with one batch",
			DelegationEarnings{
				DelegationEarning{Address: "some_addr"},
				DelegationEarning{Address: "some_addr1"},
			},
			[]DelegationEarnings{
				{
					DelegationEarning{Address: "some_addr"},
					DelegationEarning{Address: "some_addr1"},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := Payout{batchSize: 2}
			batch := payout.batch(tt.input)
			assert.Equal(t, tt.want, batch)
		})
	}
}

func Test_confirmOperation(t *testing.T) {
	type input struct {
		operation string
		gt        gotezos.IFace
	}
	cases := []struct {
		name  string
		input input
		want  bool
	}{
		{
			"is successful",
			input{
				"ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2rfb1M",
				&test.GoTezosMock{},
			},
			true,
		},
		{
			"handles timeout",
			input{
				"ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2safdj",
				&test.GoTezosMock{OperationHashesErr: true},
			},
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			confirmationDurationInterval = time.Millisecond * 500
			confirmationTimoutInterval = time.Second * 1

			payout := Payout{
				gt: tt.input.gt,
			}

			ok := payout.confirmOperation(tt.input.operation)
			assert.Equal(t, tt.want, ok)
		})
	}
}

func Test_isInBlacklist(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{
			"finds address in blacklist",
			"some_addr",
			true,
		},
		{
			"does not find address in blacklist",
			"some_other_addr",
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := Payout{blacklist: []string{
				"some_addr",
				"some_addr_1",
			}}

			actual := payout.isInBlacklist(tt.input)
			assert.Equal(t, tt.want, actual)
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

func strToPointer(str string) *string {
	return &str
}
