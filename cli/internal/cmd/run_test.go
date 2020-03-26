package cmd

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v2/cli/internal/db/model"
	"github.com/goat-systems/tzpay/v2/cli/internal/enviroment"
	"github.com/goat-systems/tzpay/v2/cli/internal/test"
	"github.com/stretchr/testify/assert"
)

func Test_splitPayouts(t *testing.T) {
	type input struct {
		split  int
		payout *model.Payout
	}

	type want struct {
		payouts []model.Payout
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				split: 3,
				payout: &model.Payout{
					DelegationEarnings: []model.DelegationEarning{
						model.DelegationEarning{
							Address: "tz1a",
						},
						model.DelegationEarning{
							Address: "tz1b",
						},
						model.DelegationEarning{
							Address: "tz1c",
						},
						model.DelegationEarning{
							Address: "tz1d",
						},
						model.DelegationEarning{
							Address: "tz1e",
						},
						model.DelegationEarning{
							Address: "tz1f",
						},
						model.DelegationEarning{
							Address: "tz1g",
						},
						model.DelegationEarning{
							Address: "tz1h",
						},
						model.DelegationEarning{
							Address: "tz1i",
						},
						model.DelegationEarning{
							Address: "tz1j",
						},
						model.DelegationEarning{
							Address: "tz1k",
						},
						model.DelegationEarning{
							Address: "tz1l",
						},
						model.DelegationEarning{
							Address: "tz1m",
						},
						model.DelegationEarning{
							Address: "tz1n",
						},
						model.DelegationEarning{
							Address: "tz1o",
						},
						model.DelegationEarning{
							Address: "tz1p",
						},
						model.DelegationEarning{
							Address: "tz1q",
						},
						model.DelegationEarning{
							Address: "tz1r",
						},
						model.DelegationEarning{
							Address: "tz1s",
						},
					},
				},
			},
			want{
				payouts: []model.Payout{
					model.Payout{
						DelegationEarnings: []model.DelegationEarning{
							model.DelegationEarning{
								Address: "tz1a",
							},
							model.DelegationEarning{
								Address: "tz1b",
							},
							model.DelegationEarning{
								Address: "tz1c",
							},
						},
					},
					model.Payout{
						DelegationEarnings: []model.DelegationEarning{
							model.DelegationEarning{
								Address: "tz1d",
							},
							model.DelegationEarning{
								Address: "tz1e",
							},
							model.DelegationEarning{
								Address: "tz1f",
							},
						},
					},
					model.Payout{
						DelegationEarnings: []model.DelegationEarning{
							model.DelegationEarning{
								Address: "tz1g",
							},
							model.DelegationEarning{
								Address: "tz1h",
							},
							model.DelegationEarning{
								Address: "tz1i",
							},
						},
					},
					model.Payout{
						DelegationEarnings: []model.DelegationEarning{
							model.DelegationEarning{
								Address: "tz1j",
							},
							model.DelegationEarning{
								Address: "tz1k",
							},
							model.DelegationEarning{
								Address: "tz1l",
							},
						},
					},
					model.Payout{
						DelegationEarnings: []model.DelegationEarning{
							model.DelegationEarning{
								Address: "tz1m",
							},
							model.DelegationEarning{
								Address: "tz1n",
							},
							model.DelegationEarning{
								Address: "tz1o",
							},
						},
					},
					model.Payout{
						DelegationEarnings: []model.DelegationEarning{
							model.DelegationEarning{
								Address: "tz1p",
							},
							model.DelegationEarning{
								Address: "tz1q",
							},
							model.DelegationEarning{
								Address: "tz1r",
							},
						},
					},
					model.Payout{
						DelegationEarnings: []model.DelegationEarning{
							model.DelegationEarning{
								Address: "tz1s",
							},
						},
					},
				},
			},
		},
		{
			"is successful with 1 split",
			input{
				split: 1,
				payout: &model.Payout{
					DelegationEarnings: []model.DelegationEarning{
						model.DelegationEarning{
							Address: "tz1a",
						},
						model.DelegationEarning{
							Address: "tz1b",
						},
						model.DelegationEarning{
							Address: "tz1c",
						},
					},
				},
			},
			want{
				[]model.Payout{
					model.Payout{
						DelegationEarnings: []model.DelegationEarning{
							model.DelegationEarning{
								Address: "tz1a",
							},
						},
					},
					model.Payout{
						DelegationEarnings: []model.DelegationEarning{
							model.DelegationEarning{
								Address: "tz1b",
							},
						},
					},
					model.Payout{
						DelegationEarnings: []model.DelegationEarning{
							model.DelegationEarning{
								Address: "tz1c",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payouts := splitPayouts(*tt.input.payout, tt.input.split)
			for i := range payouts {
				assert.Equal(t, tt.want.payouts[i], *payouts[i])
			}
		})
	}
}

func Test_run(t *testing.T) {
	type input struct {
		runnerInput newRunnerInput
	}

	type want struct {
		err         bool
		errContains string
		payout      *model.Payout
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				newRunnerInput{
					"130",
					false,
					1,
					&test.GoTezosMock{},
				},
			},
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
						Address: "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						Fees:    big.NewInt(10500000),
						Share:   1,
						Rewards: big.NewInt(70000000),
						Net:     big.NewInt(80500000),
					},
					CycleHash:      "some_hash",
					Cycle:          130,
					FrozenBalance:  big.NewInt(70000000),
					StakingBalance: big.NewInt(10000000000),
					Operations:     []string{"ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2rfb1M", "ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2rfb1M", "ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2rfb1M"},
					OperationsLink: []string{"http://tzstats.com/ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2rfb1M", "http://tzstats.com/ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2rfb1M", "http://tzstats.com/ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2rfb1M"},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			initTestEnv()
			r, err := newRunner(tt.input.runnerInput)
			assert.NoError(t, err)

			confirm = false
			payout, err := r.run()
			if tt.want.err {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.want.errContains)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want.payout, payout)
			uninitTestEnv()
		})
	}
}

func Test_ConfirmInjection(t *testing.T) {
	type input struct {
		counter int
		gt      gotezos.IFace
	}
	cases := []struct {
		name  string
		input input
		want  bool
	}{
		{
			"is successful",
			input{
				100,
				&test.GoTezosMock{},
			},
			true,
		},
		{
			"handles timeout",
			input{
				100,
				&test.GoTezosMock{CounterErr: true},
			},
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			confirmationDurationInterval = time.Millisecond * 500
			confirmationTimoutInterval = time.Second * 1

			r := runner{
				base: &enviroment.ContextEnviroment{
					GoTezos: tt.input.gt,
				},
			}

			ok := r.ConfirmInjection(tt.input.counter)
			assert.Equal(t, tt.want, ok)
		})
	}
}

var goldenContext = context.WithValue(
	context.TODO(),
	enviroment.ENVIROMENTKEY,
	&enviroment.ContextEnviroment{
		BakersFee:      0.05,
		BlackList:      []string{"somehash", "somehash1"},
		Delegate:       "somedelegate",
		GasLimit:       100000,
		HostNode:       "http://somenode.com:8732",
		MinimumPayment: 1000,
		NetworkFee:     100000,
		Wallet: gotezos.Wallet{
			Address: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
		},
	},
)

var env = map[string]string{
	"TZPAY_BAKERS_FEE":        "0.05",
	"TZPAY_DELEGATE":          "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
	"TZPAY_NETWORK_GAS_LIMIT": "100000",
	"TZPAY_HOST_NODE":         "http://somenode.com:8732",
	"TZPAY_MINIMUM_PAYMENT":   "1000",
	"TZPAY_NETWORK_FEE":       "100000",
	"TZPAY_WALLET_SECRET":     "edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2",
	"TZPAY_WALLET_PASSWORD":   "password12345##",
}

func initTestEnv() {
	for key, elem := range env {
		os.Setenv(key, string(elem))
	}
}

func uninitTestEnv() {
	for key := range env {
		os.Unsetenv(key)
	}
}
