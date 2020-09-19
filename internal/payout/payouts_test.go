package payout

import (
	"errors"
	"testing"
	"time"

	"github.com/goat-systems/go-tezos/v3/keys"
	"github.com/goat-systems/go-tezos/v3/rpc"
	"github.com/goat-systems/tzpay/v2/internal/config"
	"github.com/goat-systems/tzpay/v2/internal/test"
	"github.com/goat-systems/tzpay/v2/internal/tzkt"
	"github.com/stretchr/testify/assert"
)

func Test_Execute(t *testing.T) {
	type input struct {
		constructPayoutFunc func() (tzkt.RewardsSplit, error)
		applyFunc           func(delegators tzkt.Delegators) ([]string, error)
		inject              bool
	}

	type want struct {
		err          bool
		contains     string
		rewardsSplit tzkt.RewardsSplit
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles failure to construct payout",
			input{
				constructPayoutFunc: func() (tzkt.RewardsSplit, error) {
					return tzkt.RewardsSplit{}, errors.New("failed to construct")
				},
			},
			want{
				true,
				"failed to construct",
				tzkt.RewardsSplit{},
			},
		},
		{
			"handles failure to apply payout",
			input{
				constructPayoutFunc: func() (tzkt.RewardsSplit, error) {
					return tzkt.RewardsSplit{}, nil
				},
				applyFunc: func(delegators tzkt.Delegators) ([]string, error) {
					return []string{}, errors.New("failed to apply")
				},
				inject: true,
			},
			want{
				true,
				"failed to apply",
				tzkt.RewardsSplit{},
			},
		},
		{
			"is successful",
			input{
				constructPayoutFunc: func() (tzkt.RewardsSplit, error) {
					return tzkt.RewardsSplit{}, nil
				},
				applyFunc: func(delegators tzkt.Delegators) ([]string, error) {
					return []string{}, nil
				},
				inject: true,
			},
			want{
				false,
				"",
				tzkt.RewardsSplit{},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := Payout{
				constructPayoutFunc: tt.input.constructPayoutFunc,
				applyFunc:           tt.input.applyFunc,
				inject:              tt.input.inject,
			}
			rewardsSplit, err := payout.Execute()
			test.CheckErr(t, tt.want.err, tt.want.contains, err)
			assert.Equal(t, tt.want.rewardsSplit, rewardsSplit)
		})
	}

}

func Test_constructPayout(t *testing.T) {
	type input struct {
		rpcClient                         rpc.IFace
		tzktClient                        tzkt.IFace
		constructDexterContractPayoutFunc func(delegator tzkt.Delegator) (tzkt.Delegator, error)
		dexterOnly                        bool
	}

	type want struct {
		err          bool
		contains     string
		rewardsSplit tzkt.RewardsSplit
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles failure to get rewards split",
			input{
				rpcClient: &test.RPCMock{},
				tzktClient: &test.TzktMock{
					RewardsSplitErr: true,
				},
				constructDexterContractPayoutFunc: func(delegator tzkt.Delegator) (tzkt.Delegator, error) {
					return delegator, nil
				},
			},
			want{
				true,
				"failed to get rewards split",
				tzkt.RewardsSplit{},
			},
		},
		{
			"handles failure to get balance",
			input{
				rpcClient: &test.RPCMock{
					BalanceErr: true,
				},
				tzktClient: &test.TzktMock{},
				constructDexterContractPayoutFunc: func(delegator tzkt.Delegator) (tzkt.Delegator, error) {
					return delegator, nil
				},
			},
			want{
				true,
				"failed to get balance",
				tzkt.RewardsSplit{
					Cycle:                    270,
					StakingBalance:           740613513605,
					DelegatedBalance:         555430526884,
					NumDelegators:            107,
					ExpectedBlocks:           4.43,
					ExpectedEndorsements:     141.71,
					OwnBlocks:                5,
					OwnBlockRewards:          191250000,
					MissedOwnBlocks:          2,
					MissedOwnBlockRewards:    77500000,
					BlockDeposits:            2560000000,
					Endorsements:             126,
					EndorsementRewards:       157500000,
					MissedEndorsements:       16,
					MissedEndorsementRewards: 20000000,
					EndorsementDeposits:      8064000000,
					OwnBlockFees:             47180,
					MissedOwnBlockFees:       54607,
					Delegators: tzkt.Delegators{
						tzkt.Delegator{
							Address:        "tz1icdoLr8vof5oXiEKCFSyrVoouGiKDQ3Gd",
							Balance:        60545965782,
							CurrentBalance: 60739073316,
						},
						tzkt.Delegator{
							Address:        "KT1FPyY6mAhnzyVGP8ApGvuRyF7SKcT9TDWy",
							Balance:        60075572992,
							CurrentBalance: 60267312348,
						},
						tzkt.Delegator{
							Address:        "KT1LgkGigaMrnim3TonQWfwDHnM3fHkF1jMv",
							Balance:        57461165021,
							CurrentBalance: 57644560137,
						}, tzkt.Delegator{
							Address:        "KT1C8S2vLYbzgQHhdC8MBehunhcp1Q9hj6MC",
							Balance:        55305195039,
							CurrentBalance: 176566401,
						},
					},
				},
			},
		},
		{
			"handles failure to contruct payout for dexter contract",
			input{
				rpcClient:  &test.RPCMock{},
				tzktClient: &test.TzktMock{},
				constructDexterContractPayoutFunc: func(delegator tzkt.Delegator) (tzkt.Delegator, error) {
					return delegator, errors.New("failed to contruct dexter")
				},
			},
			want{
				true,
				"failed to contruct dexter",
				tzkt.RewardsSplit{},
			},
		},
		{
			"is successful",
			input{
				rpcClient:  &test.RPCMock{},
				tzktClient: &test.TzktMock{},
				constructDexterContractPayoutFunc: func(delegator tzkt.Delegator) (tzkt.Delegator, error) {
					return delegator, nil
				},
			},
			want{
				false,
				"",
				tzkt.RewardsSplit{
					Cycle:                    270,
					StakingBalance:           740613513605,
					DelegatedBalance:         555430526884,
					NumDelegators:            107,
					ExpectedBlocks:           4.43,
					ExpectedEndorsements:     141.71,
					OwnBlocks:                5,
					OwnBlockRewards:          191250000,
					MissedOwnBlocks:          2,
					MissedOwnBlockRewards:    77500000,
					BlockDeposits:            2560000000,
					Endorsements:             126,
					EndorsementRewards:       157500000,
					MissedEndorsements:       16,
					MissedEndorsementRewards: 20000000,
					EndorsementDeposits:      8064000000,
					OwnBlockFees:             47180,
					MissedOwnBlockFees:       54607,
					Delegators: tzkt.Delegators{
						tzkt.Delegator{
							Address:        "tz1icdoLr8vof5oXiEKCFSyrVoouGiKDQ3Gd",
							Balance:        60545965782,
							CurrentBalance: 60739073316,
							NetRewards:     34665260,
							GrossRewards:   36489747,
							Share:          0.08175109509855863,
							Fee:            1824487,
						}, tzkt.Delegator{
							Address:        "KT1FPyY6mAhnzyVGP8ApGvuRyF7SKcT9TDWy",
							Balance:        60075572992,
							CurrentBalance: 60267312348,
							NetRewards:     34395939,
							GrossRewards:   36206251,
							Share:          0.08111595574266121,
							Fee:            1810312,
						}, tzkt.Delegator{
							Address:        "KT1LgkGigaMrnim3TonQWfwDHnM3fHkF1jMv",
							Balance:        57461165021,
							CurrentBalance: 57644560137, NetRewards: 32899074,
							GrossRewards: 34630604,
							Share:        0.07758589867109342,
							Fee:          1731530,
						}, tzkt.Delegator{
							Address:        "KT1C8S2vLYbzgQHhdC8MBehunhcp1Q9hj6MC",
							Balance:        55305195039,
							CurrentBalance: 176566401,
							NetRewards:     31664685,
							GrossRewards:   33331247,
							Share:          0.07467483920161976,
							Fee:            1666562,
						},
					},
					BakerRewards:       3013,
					BakerShare:         6.751159556435947e-06,
					BakerCollectedFees: 7032891,
				},
			},
		},
		{
			"is successful dexter only",
			input{
				rpcClient:  &test.RPCMock{},
				tzktClient: &test.TzktMock{},
				constructDexterContractPayoutFunc: func(delegator tzkt.Delegator) (tzkt.Delegator, error) {
					return delegator, nil
				},
				dexterOnly: true,
			},
			want{
				false,
				"",
				tzkt.RewardsSplit{
					Cycle:                    270,
					StakingBalance:           740613513605,
					DelegatedBalance:         555430526884,
					NumDelegators:            107,
					ExpectedBlocks:           4.43,
					ExpectedEndorsements:     141.71,
					OwnBlocks:                5,
					OwnBlockRewards:          191250000,
					MissedOwnBlocks:          2,
					MissedOwnBlockRewards:    77500000,
					BlockDeposits:            2560000000,
					Endorsements:             126,
					EndorsementRewards:       157500000,
					MissedEndorsements:       16,
					MissedEndorsementRewards: 20000000,
					EndorsementDeposits:      8064000000,
					OwnBlockFees:             47180,
					MissedOwnBlockFees:       54607,
					Delegators: tzkt.Delegators{
						tzkt.Delegator{
							Address:        "KT1LgkGigaMrnim3TonQWfwDHnM3fHkF1jMv",
							Balance:        57461165021,
							CurrentBalance: 57644560137, NetRewards: 32899074,
							GrossRewards: 34630604,
							Share:        0.07758589867109342,
							Fee:          1731530,
						}, tzkt.Delegator{
							Address:        "KT1C8S2vLYbzgQHhdC8MBehunhcp1Q9hj6MC",
							Balance:        55305195039,
							CurrentBalance: 176566401,
							NetRewards:     31664685,
							GrossRewards:   33331247,
							Share:          0.07467483920161976,
							Fee:            1666562,
						},
					},
					BakerRewards:       3013,
					BakerShare:         6.751159556435947e-06,
					BakerCollectedFees: 3398092,
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := &Payout{
				rpc:  tt.input.rpcClient,
				tzkt: tt.input.tzktClient,
				config: config.Config{
					Baker: config.Baker{
						DexterLiquidityContracts: []string{
							"KT1LgkGigaMrnim3TonQWfwDHnM3fHkF1jMv",
							"KT1C8S2vLYbzgQHhdC8MBehunhcp1Q9hj6MC",
						},
						Fee:                          0.05,
						DexterLiquidityContractsOnly: tt.input.dexterOnly,
					},
				},
				constructDexterContractPayoutFunc: tt.input.constructDexterContractPayoutFunc,
			}
			rewardsSplit, err := payout.constructPayout()
			test.CheckErr(t, tt.want.err, tt.want.contains, err)
			assert.Equal(t, tt.want.rewardsSplit, rewardsSplit)
		})
	}
}

func Test_splitDelegationsAndDexterContracts(t *testing.T) {
	type input struct {
		cfg     config.Config
		rewards tzkt.RewardsSplit
	}

	type want struct {
		delegations tzkt.Delegators
		contracts   tzkt.Delegators
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				config.Config{
					Baker: config.Baker{
						DexterLiquidityContracts: []string{
							"KT1a",
							"KT1b",
						},
					},
				},
				tzkt.RewardsSplit{
					Delegators: tzkt.Delegators{
						{
							Address: "KT1a",
						},
						{
							Address: "KT1b",
						},
						{
							Address: "tz1",
						},
						{
							Address: "tz2",
						},
					},
				},
			},
			want{
				delegations: tzkt.Delegators{
					{
						Address: "tz1",
					},
					{
						Address: "tz2",
					},
				},
				contracts: tzkt.Delegators{
					{
						Address: "KT1a",
					},
					{
						Address: "KT1b",
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := &Payout{
				config: tt.input.cfg,
			}
			delegations, contracts := payout.splitDelegationsAndDexterContracts(tt.input.rewards)
			assert.Equal(t, tt.want.delegations, delegations)
			assert.Equal(t, tt.want.contracts, contracts)
		})
	}
}

func Test_constructDelegation(t *testing.T) {
	type input struct {
		delegator      tzkt.Delegator
		totalRewards   int
		stakingBalance int
	}

	cases := []struct {
		name  string
		input input
		want  tzkt.Delegator
	}{
		{
			"handles blacklist marking",
			input{
				tzkt.Delegator{
					Address: "some_blacklisted_address",
					Balance: 5000000,
				},
				10000000,
				1000000000,
			},
			tzkt.Delegator{
				Address:      "some_blacklisted_address",
				Balance:      5000000,
				NetRewards:   47500,
				GrossRewards: 50000,
				Share:        0.005,
				Fee:          2500,
				BlackListed:  true,
			},
		},
		{
			"is successful",
			input{
				tzkt.Delegator{
					Address: "some_address",
					Balance: 5000000,
				},
				10000000,
				1000000000,
			},
			tzkt.Delegator{
				Address:      "some_address",
				Balance:      5000000,
				NetRewards:   47500,
				GrossRewards: 50000,
				Share:        0.005,
				Fee:          2500,
				BlackListed:  false,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := &Payout{
				config: config.Config{
					Baker: config.Baker{
						Blacklist: []string{
							"some_blacklisted_address",
						},
						Fee: 0.05,
					},
				},
			}
			delegation := payout.constructDelegation(tt.input.delegator, tt.input.totalRewards, tt.input.stakingBalance)
			assert.Equal(t, tt.want, delegation)
		})
	}
}

func Test_calculateTotals(t *testing.T) {
	type input struct {
		earningsOnly bool
		rewardsSplit tzkt.RewardsSplit
	}

	cases := []struct {
		name  string
		input input
		want  int
	}{
		{
			"handles earnings only",
			input{
				true,
				tzkt.RewardsSplit{
					EndorsementRewards:       1000,
					RevelationRewards:        1012,
					OwnBlockFees:             213441,
					OwnBlockRewards:          24124321,
					ExtraBlockFees:           32321,
					ExtraBlockRewards:        234123,
					MissedEndorsementRewards: 2134423,
					MissedOwnBlockFees:       21234,
					MissedOwnBlockRewards:    3214312,
				},
			},
			24606218,
		},
		{
			"handles earnings only false",
			input{
				false,
				tzkt.RewardsSplit{
					EndorsementRewards:       1000,
					RevelationRewards:        1012,
					OwnBlockFees:             213441,
					OwnBlockRewards:          24124321,
					ExtraBlockFees:           32321,
					ExtraBlockRewards:        234123,
					MissedEndorsementRewards: 2134423,
					MissedOwnBlockFees:       21234,
					MissedOwnBlockRewards:    3214312,
				},
			},
			29976187,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := Payout{
				config: config.Config{
					Baker: config.Baker{
						EarningsOnly: tt.input.earningsOnly,
					},
				},
			}

			total := payout.calculateTotals(tt.input.rewardsSplit)
			assert.Equal(t, tt.want, total)
		})
	}

}

func Test_apply(t *testing.T) {
	type input struct {
		rpcClient  rpc.IFace
		delegators tzkt.Delegators
	}

	type want struct {
		err        bool
		contains   string
		operations []string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles failure to get head",
			input{
				rpcClient: &test.RPCMock{
					HeadErr: true,
				},
				delegators: tzkt.Delegators{},
			},
			want{
				true,
				"failed to get block",
				[]string{},
			},
		},
		{
			"handles failure to contruct transaction batches",
			input{
				rpcClient: &test.RPCMock{
					CounterErr: true,
				},
				delegators: tzkt.Delegators{},
			},
			want{
				true,
				"failed to get counter",
				[]string{},
			},
		},
		{
			"handles failure forge",
			input{
				rpcClient: &test.RPCMock{},
				delegators: tzkt.Delegators{
					{
						Address:      "somedelegation",
						GrossRewards: 1000000,
						NetRewards:   900000,
					},
					{
						Address:      "someotherdelegation",
						GrossRewards: 1000000,
						NetRewards:   950000,
					},
					{
						Address:      "delegation_dexter",
						GrossRewards: 1000000,
						NetRewards:   950000,
						LiquidityProviders: []tzkt.LiquidityProvider{
							{
								Address:      "liquidity_provider",
								GrossRewards: 1000000,
								NetRewards:   950000,
							},
							{
								Address:      "liquidity_provider1",
								GrossRewards: 1000000,
								NetRewards:   950000,
							},
						},
					},
				},
			},
			want{
				true,
				"failed to forge",
				[]string{},
			},
		},
		{
			"handles failure to inject",
			input{
				rpcClient: &test.RPCMock{
					InjectionOperationErr: true,
				},
				delegators: tzkt.Delegators{
					{
						Address:      "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
						GrossRewards: 1000000,
						NetRewards:   900000,
					},
					{
						Address:      "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
						GrossRewards: 1000000,
						NetRewards:   950000,
					},
					{
						Address:      "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
						GrossRewards: 1000000,
						NetRewards:   950000,
						LiquidityProviders: []tzkt.LiquidityProvider{
							{
								Address:      "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
								GrossRewards: 1000000,
								NetRewards:   950000,
							},
							{
								Address:      "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
								GrossRewards: 1000000,
								NetRewards:   950000,
							},
						},
					},
				},
			},
			want{
				true,
				"failed to inject operation",
				[]string{},
			},
		},
		{
			"is successful",
			input{
				rpcClient: &test.RPCMock{},
				delegators: tzkt.Delegators{
					{
						Address:      "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
						GrossRewards: 1000000,
						NetRewards:   900000,
					},
					{
						Address:      "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
						GrossRewards: 1000000,
						NetRewards:   950000,
					},
					{
						Address:      "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
						GrossRewards: 1000000,
						NetRewards:   950000,
						LiquidityProviders: []tzkt.LiquidityProvider{
							{
								Address:      "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
								GrossRewards: 1000000,
								NetRewards:   950000,
							},
							{
								Address:      "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
								GrossRewards: 1000000,
								NetRewards:   950000,
							},
						},
					},
				},
			},
			want{
				false,
				"",
				[]string{"ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2rfb1M"},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			key, err := keys.NewKey(keys.NewKeyInput{
				Esk:      "edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2",
				Password: "password12345##",
				Kind:     keys.Ed25519,
			})
			assert.Nil(t, err)

			payout := Payout{
				rpc: tt.input.rpcClient,
				config: config.Config{
					Operations: config.Operations{
						GasLimit:   10000,
						NetworkFee: 3000,
						BatchSize:  100,
					},
				},
				key: key,
			}

			ops, err := payout.apply(tt.input.delegators)
			test.CheckErr(t, tt.want.err, tt.want.contains, err)
			assert.Equal(t, tt.want.operations, ops)

		})
	}
}

func Test_constructTransactionBatches(t *testing.T) {
	type input struct {
		counter    int
		rpcClient  rpc.IFace
		delegators tzkt.Delegators
	}

	type want struct {
		err      bool
		contains string
		contents []rpc.Contents
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles cycle error",
			input{
				100,
				&test.RPCMock{
					CounterErr: true,
				},
				tzkt.Delegators{},
			},
			want{
				true,
				"failed to",
				nil,
			},
		},
		{
			"is successful",
			input{
				100,
				&test.RPCMock{
					CounterErr: false,
				},
				tzkt.Delegators{
					{
						Address:      "somedelegation",
						GrossRewards: 1000000,
						NetRewards:   900000,
					},
					{
						Address:      "someotherdelegation",
						GrossRewards: 1000000,
						NetRewards:   950000,
					},
					{
						Address:      "delegation_dexter",
						GrossRewards: 1000000,
						NetRewards:   950000,
						LiquidityProviders: []tzkt.LiquidityProvider{
							{
								Address:      "liquidity_provider",
								GrossRewards: 1000000,
								NetRewards:   950000,
							},
							{
								Address:      "liquidity_provider1",
								GrossRewards: 1000000,
								NetRewards:   950000,
							},
						},
					},
				},
			},
			want{
				false,
				"",
				[]rpc.Contents{
					{
						{
							Kind:         rpc.TRANSACTION,
							Source:       "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
							Fee:          0,
							Counter:      101,
							GasLimit:     0,
							StorageLimit: 0,
							Amount:       900000,
							Destination:  "somedelegation",
						},
						{
							Kind:         rpc.TRANSACTION,
							Source:       "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
							Fee:          0,
							Counter:      102,
							GasLimit:     0,
							StorageLimit: 0,
							Amount:       950000,
							Destination:  "someotherdelegation",
						},
						{
							Kind:         rpc.TRANSACTION,
							Source:       "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
							Fee:          0,
							Counter:      103,
							GasLimit:     0,
							StorageLimit: 0,
							Amount:       950000,
							Destination:  "liquidity_provider",
						},
						{
							Kind:         rpc.TRANSACTION,
							Source:       "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
							Fee:          0,
							Counter:      104,
							GasLimit:     0,
							StorageLimit: 0,
							Amount:       950000,
							Destination:  "liquidity_provider1",
						},
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			key, err := keys.NewKey(keys.NewKeyInput{
				Esk:      "edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2",
				Password: "password12345##",
				Kind:     keys.Ed25519,
			})
			assert.Nil(t, err)
			payout := &Payout{
				config: config.Config{
					Operations: config.Operations{
						BatchSize: 100,
					},
				},
				rpc: tt.input.rpcClient,
				key: key,
			}
			contents, err := payout.constructTransactionBatches("some_hash", tt.input.delegators)
			test.CheckErr(t, tt.want.err, tt.want.contains, err)
			assert.Equal(t, tt.want.contents, contents)
		})
	}
}

func Test_batch(t *testing.T) {
	cases := []struct {
		name  string
		input tzkt.Delegators
		want  []tzkt.Delegators
	}{
		{
			"is successful with multiple batches",
			tzkt.Delegators{
				{Address: "some_addr"},
				{Address: "some_addr1"},
				{Address: "some_addr2"},
				{Address: "some_addr3"},
				{Address: "some_addr4"},
			},
			[]tzkt.Delegators{
				{
					{Address: "some_addr"},
					{Address: "some_addr1"},
				},
				{
					{Address: "some_addr2"},
					{Address: "some_addr3"}},
				{
					{Address: "some_addr4"},
				},
			},
		},
		{
			"is successful with one batch",
			tzkt.Delegators{
				{Address: "some_addr"},
				{Address: "some_addr1"},
			},
			[]tzkt.Delegators{
				{
					{Address: "some_addr"},
					{Address: "some_addr1"},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := Payout{
				config: config.Config{
					Operations: config.Operations{
						BatchSize: 2,
					},
				},
			}
			batch := payout.batch(tt.input)
			assert.Equal(t, tt.want, batch)
		})
	}
}

func Test_injectOperations(t *testing.T) {
	type input struct {
		rpcClient  rpc.IFace
		operations []string
	}

	type want struct {
		err      bool
		contains string
		ophashes []string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles failure to sign",
			input{
				rpcClient: &rpc.Client{},
				operations: []string{
					"non_hex_string",
				},
			},
			want{
				true,
				"failed to hex decode message",
				[]string{},
			},
		},
		{
			"handles failure to inject",
			input{
				rpcClient: &test.RPCMock{
					InjectionOperationErr: true,
				},
				operations: []string{
					"5aff622d53d32a8bae591627718c60a35b16737e301c57a13b6f1765483d88ff6c007fd82c06cf5a203f18faaf562447ed1efcc6c010830a07c350008090dfc04a0000a31e81ac3425310e3274a4698a793b2839dc0afa00",
				},
			},
			want{
				true,
				"failed to inject operation",
				[]string{},
			},
		},
		{
			"is successful",
			input{
				rpcClient: &test.RPCMock{},
				operations: []string{
					"5aff622d53d32a8bae591627718c60a35b16737e301c57a13b6f1765483d88ff6c007fd82c06cf5a203f18faaf562447ed1efcc6c010830a07c350008090dfc04a0000a31e81ac3425310e3274a4698a793b2839dc0afa00",
				},
			},
			want{
				false,
				"",
				[]string{"ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2rfb1M"},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			key, err := keys.NewKey(keys.NewKeyInput{
				Esk:      "edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2",
				Password: "password12345##",
				Kind:     keys.Ed25519,
			})
			assert.Nil(t, err)

			payout := Payout{
				rpc: tt.input.rpcClient,
				key: key,
			}

			ophashes, err := payout.injectOperations(tt.input.operations)
			test.CheckErr(t, tt.want.err, tt.want.contains, err)
			assert.Equal(t, tt.want.ophashes, ophashes)

		})
	}
}

func Test_confirmOperation(t *testing.T) {
	type input struct {
		operation string
		rpcClient rpc.IFace
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
				&test.RPCMock{},
			},
			true,
		},
		{
			"handles timeout",
			input{
				"ooYympR9wfV98X4MUHtE78NjXYRDeMTAD4ei7zEZDqoHv2safdj",
				&test.RPCMock{OperationHashesErr: true},
			},
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			confirmationDurationInterval = time.Millisecond * 500
			confirmationTimoutInterval = time.Second * 1

			payout := Payout{
				rpc: tt.input.rpcClient,
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
			payout := Payout{
				config: config.Config{
					Baker: config.Baker{
						Blacklist: []string{
							"some_addr",
							"some_addr_1",
						},
					},
				},
			}

			actual := payout.isInBlacklist(tt.input)
			assert.Equal(t, tt.want, actual)
		})
	}
}

func Test_isDexterContract(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{
			"finds dexter contract",
			"some_addr",
			true,
		},
		{
			"does not find dexter contract",
			"some_other_addr",
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := Payout{
				config: config.Config{
					Baker: config.Baker{
						DexterLiquidityContracts: []string{
							"some_addr",
							"some_addr_1",
						},
					},
				},
			}

			actual := payout.isDexterContract(tt.input)
			assert.Equal(t, tt.want, actual)
		})
	}
}

func strToPointer(str string) *string {
	return &str
}
