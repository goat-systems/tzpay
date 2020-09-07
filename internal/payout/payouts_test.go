package payout

import (
	"testing"
	"time"

	gotezos "github.com/goat-systems/go-tezos/v3"
	"github.com/goat-systems/tzpay/v2/internal/test"
	"github.com/goat-systems/tzpay/v2/internal/tzkt"
	"github.com/stretchr/testify/assert"
)

func Test_Execute(t *testing.T) {
	type want struct {
		err          bool
		errContains  string
		rewardsSplit tzkt.RewardsSplit
	}

	type input struct {
		payout Payout
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				payout: Payout{
					gt:         &test.GoTezosMock{},
					tzkt:       &test.TzktMock{},
					bakerFee:   0.05,
					inject:     false,
					networkFee: 1345,
					gasLimit:   203999,
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
						{
							Address:        "tz1icdoLr8vof5oXiEKCFSyrVoouGiKDQ3Gd",
							Balance:        60545965782,
							CurrentBalance: 60739073316,
							Emptied:        false,
							NetRewards:     34665260,
							GrossRewards:   36489747,
							Share:          0.08175109509855863,
							Fee:            1824487,
							BlackListed:    false,
						},
						{
							Address:        "KT1FPyY6mAhnzyVGP8ApGvuRyF7SKcT9TDWy",
							Balance:        60075572992,
							CurrentBalance: 60267312348,
							Emptied:        false,
							NetRewards:     34395939,
							GrossRewards:   36206251,
							Share:          0.08111595574266121,
							Fee:            1810312,
							BlackListed:    false,
						},
						{
							Address:        "KT1LgkGigaMrnim3TonQWfwDHnM3fHkF1jMv",
							Balance:        57461165021,
							CurrentBalance: 57644560137,
							Emptied:        false,
							NetRewards:     32899074,
							GrossRewards:   34630604,
							Share:          0.07758589867109342,
							Fee:            1731530,
							BlackListed:    false,
						},
						{
							Address:        "KT1C8S2vLYbzgQHhdC8MBehunhcp1Q9hj6MC",
							Balance:        55305195039,
							CurrentBalance: 176566401,
							Emptied:        false,
							NetRewards:     31664685,
							GrossRewards:   33331247,
							Share:          0.07467483920161976,
							Fee:            1666562,
							BlackListed:    false,
						},
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rewardsSplit, err := tt.input.payout.Execute()
			test.CheckErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.rewardsSplit, rewardsSplit)
		})
	}
}

func Test_constructPayout(t *testing.T) {
	type input struct {
		payout Payout
	}

	type want struct {
		err         bool
		errContains string
		rewardSplit tzkt.RewardsSplit
	}
	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles failure to get rewards split",
			input{
				payout: Payout{
					gt: &test.GoTezosMock{},
					tzkt: &test.TzktMock{
						RewardsSplitErr: true,
					},
				},
			},
			want{
				true,
				"failed to contruct payout",
				tzkt.RewardsSplit{},
			},
		},
		{
			"is successful",
			input{
				payout: Payout{
					gt:       &test.GoTezosMock{},
					tzkt:     &test.TzktMock{},
					bakerFee: 0.05,
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
						{
							Address:        "tz1icdoLr8vof5oXiEKCFSyrVoouGiKDQ3Gd",
							Balance:        60545965782,
							CurrentBalance: 60739073316,
							Emptied:        false,
							NetRewards:     34665260,
							GrossRewards:   36489747,
							Share:          0.08175109509855863,
							Fee:            1824487,
							BlackListed:    false,
						},
						{
							Address:        "KT1FPyY6mAhnzyVGP8ApGvuRyF7SKcT9TDWy",
							Balance:        60075572992,
							CurrentBalance: 60267312348,
							Emptied:        false,
							NetRewards:     34395939,
							GrossRewards:   36206251,
							Share:          0.08111595574266121,
							Fee:            1810312,
							BlackListed:    false,
						},
						{
							Address:        "KT1LgkGigaMrnim3TonQWfwDHnM3fHkF1jMv",
							Balance:        57461165021,
							CurrentBalance: 57644560137,
							Emptied:        false,
							NetRewards:     32899074,
							GrossRewards:   34630604,
							Share:          0.07758589867109342,
							Fee:            1731530,
							BlackListed:    false,
						},
						{
							Address:        "KT1C8S2vLYbzgQHhdC8MBehunhcp1Q9hj6MC",
							Balance:        55305195039,
							CurrentBalance: 176566401,
							Emptied:        false,
							NetRewards:     31664685,
							GrossRewards:   33331247,
							Share:          0.07467483920161976,
							Fee:            1666562,
							BlackListed:    false,
						},
					},
				},
			},
		},
		{
			"is successful with Earning Only",
			input{
				payout: Payout{
					gt:           &test.GoTezosMock{},
					tzkt:         &test.TzktMock{},
					bakerFee:     0.05,
					earningsOnly: true,
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
						{
							Address:        "tz1icdoLr8vof5oXiEKCFSyrVoouGiKDQ3Gd",
							Balance:        60545965782,
							CurrentBalance: 60739073316,
							Emptied:        false,
							NetRewards:     27088824,
							GrossRewards:   28514551,
							Share:          0.08175109509855863,
							Fee:            1425727,
						},
						{
							Address:        "KT1FPyY6mAhnzyVGP8ApGvuRyF7SKcT9TDWy",
							Balance:        60075572992,
							CurrentBalance: 60267312348,
							Emptied:        false,
							NetRewards:     26878366,
							GrossRewards:   28293016,
							Share:          0.08111595574266121,
							Fee:            1414650,
							BlackListed:    false,
						},
						{
							Address:        "KT1LgkGigaMrnim3TonQWfwDHnM3fHkF1jMv",
							Balance:        57461165021,
							CurrentBalance: 57644560137,
							Emptied:        false,
							NetRewards:     25708655,
							GrossRewards:   27061742,
							Share:          0.07758589867109342,
							Fee:            1353087,
							BlackListed:    false,
						},
						{
							Address:            "KT1C8S2vLYbzgQHhdC8MBehunhcp1Q9hj6MC",
							Balance:            55305195039,
							CurrentBalance:     176566401,
							Emptied:            false,
							NetRewards:         24744055,
							GrossRewards:       26046373,
							Share:              0.07467483920161976,
							Fee:                1302318,
							LiquidityProviders: []tzkt.LiquidityProvider(nil),
							BlackListed:        false,
						},
					},
				},
			},
		},
		{
			"is successful on triggering blacklist",
			input{
				payout: Payout{
					gt:           &test.GoTezosMock{},
					tzkt:         &test.TzktMock{},
					bakerFee:     0.05,
					earningsOnly: true,
					blacklist: []string{
						"tz1icdoLr8vof5oXiEKCFSyrVoouGiKDQ3Gd",
					},
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
						{
							Address:        "tz1icdoLr8vof5oXiEKCFSyrVoouGiKDQ3Gd",
							Balance:        60545965782,
							CurrentBalance: 60739073316,
							Emptied:        false,
							NetRewards:     27088824,
							GrossRewards:   28514551,
							Share:          0.08175109509855863,
							Fee:            1425727,
							BlackListed:    true,
						},
						{
							Address:        "KT1FPyY6mAhnzyVGP8ApGvuRyF7SKcT9TDWy",
							Balance:        60075572992,
							CurrentBalance: 60267312348,
							Emptied:        false,
							NetRewards:     26878366,
							GrossRewards:   28293016,
							Share:          0.08111595574266121,
							Fee:            1414650,
							BlackListed:    false,
						},
						{
							Address:        "KT1LgkGigaMrnim3TonQWfwDHnM3fHkF1jMv",
							Balance:        57461165021,
							CurrentBalance: 57644560137,
							Emptied:        false,
							NetRewards:     25708655,
							GrossRewards:   27061742,
							Share:          0.07758589867109342,
							Fee:            1353087,
							BlackListed:    false,
						},
						{
							Address:            "KT1C8S2vLYbzgQHhdC8MBehunhcp1Q9hj6MC",
							Balance:            55305195039,
							CurrentBalance:     176566401,
							Emptied:            false,
							NetRewards:         24744055,
							GrossRewards:       26046373,
							Share:              0.07467483920161976,
							Fee:                1302318,
							LiquidityProviders: []tzkt.LiquidityProvider(nil),
							BlackListed:        false,
						},
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rewardSplit, err := tt.input.payout.constructPayout()
			test.CheckErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.rewardSplit, rewardSplit)
		})
	}
}

func Test_forge(t *testing.T) {
	type input struct {
		gt         gotezos.IFace
		delegators tzkt.Delegators
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
				tzkt.Delegators{
					{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: 1000000,
						NetRewards:   900000,
					},
					{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: 1000000,
						NetRewards:   950000,
					},
				},
			},
			want{
				false,
				"",
				[]string{
					"7cc601d2729c90b267e6a79d902f8b048d37fd990f2f7447efefb0cfb2f8e8a46c004b04ad1e57c2f13b61b3d2c95b3073d961a4132bc10a66dfb90c00f0fd390000056a59972593bdc74a5295671c8f5d43c21348da006c004b04ad1e57c2f13b61b3d2c95b3073d961a4132bc10a66dfb90c00f0fd390000056a59972593bdc74a5295671c8f5d43c21348da00",
				},
			},
		},
		{
			"handles failure to get head",
			input{
				&test.GoTezosMock{HeadErr: true},
				tzkt.Delegators{
					{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: 1000000,
						NetRewards:   900000,
					},
					{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: 1000000,
						NetRewards:   950000,
					},
				},
			},
			want{
				true,
				"failed to get operation hex string: failed to get block",
				[]string{},
			},
		},
		{
			"handles failure to get counter",
			input{
				&test.GoTezosMock{CounterErr: true},
				tzkt.Delegators{
					{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: 1000000,
						NetRewards:   900000,
					},
					{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: 1000000,
						NetRewards:   950000,
					},
				},
			},
			want{
				true,
				"failed to get operation hex string: failed to get counter",
				[]string{},
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
				networkFee: 1345,
				gasLimit:   203999,
			}
			ophexes, err := payout.forge(tt.input.delegators)
			test.CheckErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.ophexes, ophexes)
		})
	}
}

func Test_forgeOperation(t *testing.T) {
	type input struct {
		counter    int
		delegators tzkt.Delegators
		gt         gotezos.IFace
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
				tzkt.Delegators{
					{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: 1000000,
						NetRewards:   900000,
					},
					{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: 1000000,
						NetRewards:   950000,
					},
				},
				&test.GoTezosMock{},
			},
			want{
				err:         false,
				errContains: "",
				ophash:      "7cc601d2729c90b267e6a79d902f8b048d37fd990f2f7447efefb0cfb2f8e8a46c004b04ad1e57c2f13b61b3d2c95b3073d961a4132bc10a07dfb90c00f0fd390000056a59972593bdc74a5295671c8f5d43c21348da006c004b04ad1e57c2f13b61b3d2c95b3073d961a4132bc10a07dfb90c00f0fd390000056a59972593bdc74a5295671c8f5d43c21348da00",
				counter:     7,
			},
		},
		{
			"handles failure to get head",
			input{
				5,
				tzkt.Delegators{
					{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: 1000000,
						NetRewards:   900000,
					},
					{
						Address:      "tz1L8fUQLuwRuywTZUP5JUw9LL3kJa8LMfoo",
						GrossRewards: 1000000,
						NetRewards:   950000,
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
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := &Payout{
				gt: tt.input.gt,
				wallet: gotezos.Wallet{
					Address: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
				},
				networkFee: 1345,
				gasLimit:   203999,
			}
			ophash, counter, err := payout.forgeOperation(tt.input.counter, tt.input.delegators)
			test.CheckErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.counter, counter)
			assert.Equal(t, tt.want.ophash, ophash)
		})
	}

}

func Test_constructPayoutContents(t *testing.T) {
	type input struct {
		counter    int
		delegators tzkt.Delegators
	}

	cases := []struct {
		name  string
		input input
		want  []gotezos.Transaction
	}{
		{
			"is successful",
			input{
				100,
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
			[]gotezos.Transaction{
				{
					Source:       "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
					Fee:          0,
					Counter:      101,
					GasLimit:     0,
					StorageLimit: 0,
					Amount:       900000,
					Destination:  "somedelegation",
				},
				{
					Source:       "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
					Fee:          0,
					Counter:      102,
					GasLimit:     0,
					StorageLimit: 0,
					Amount:       950000,
					Destination:  "someotherdelegation",
				},
				{
					Source:       "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
					Fee:          0,
					Counter:      103,
					GasLimit:     0,
					StorageLimit: 0,
					Amount:       950000,
					Destination:  "liquidity_provider",
				},
				{
					Source:       "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
					Fee:          0,
					Counter:      104,
					GasLimit:     0,
					StorageLimit: 0,
					Amount:       950000,
					Destination:  "liquidity_provider1",
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
			}
			contents, _ := payout.constructPayoutContents(tt.input.counter, tt.input.delegators)
			assert.Equal(t, tt.want, contents)
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
			payout := Payout{dexterContracts: []string{
				"some_addr",
				"some_addr_1",
			}}

			actual := payout.isDexterContract(tt.input)
			assert.Equal(t, tt.want, actual)
		})
	}
}

func strToPointer(str string) *string {
	return &str
}
