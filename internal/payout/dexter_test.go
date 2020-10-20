package payout

import (
	"encoding/json"
	"testing"

	"github.com/goat-systems/go-tezos/v3/rpc"
	"github.com/goat-systems/tzpay/v3/internal/config"
	"github.com/goat-systems/tzpay/v3/internal/test"
	"github.com/goat-systems/tzpay/v3/internal/tzkt"
	"github.com/stretchr/testify/assert"
)

// func Test_constructPayoutX(t *testing.T) {
// 	tz := tzkt.NewTZKT("https://api.tzkt.io/")
// 	r, err := rpc.New("https://mainnet-tezos.giganode.io")
// 	assert.Nil(t, err)

// 	payout := Payout{
// 		rpc:   r,
// 		tzkt:  tz,
// 		cycle: 289,
// 		config: config.Config{
// 			Baker: config.Baker{
// 				Fee: 0.05,
// 				DexterLiquidityContracts: []string{
// 					"KT1DrJV8vhkdLEj76h1H9Q4irZDqAkMPo1Qf",
// 					"KT1Puc9St8wdNoGtLiD2WXaHbWU7styaxYhD",
// 				},
// 				Address: "tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q",
// 			},
// 		},
// 	}
// 	payout.constructDexterContractPayoutFunc = payout.constructDexterContractPayout

// 	rewardsSplit, err := tz.GetRewardsSplit(payout.config.Baker.Address, payout.cycle)
// 	assert.Nil(t, err)

// 	totalRewards := payout.calculateTotals(rewardsSplit)

// 	bakerBalance, err := payout.rpc.Balance(rpc.BalanceInput{
// 		Cycle:   payout.cycle,
// 		Address: payout.config.Baker.Address,
// 	})
// 	assert.Nil(t, err)

// 	rewardsSplit.BakerShare = float64(bakerBalance) / float64(rewardsSplit.StakingBalance)
// 	rewardsSplit.BakerRewards = int(rewardsSplit.BakerShare * float64(totalRewards))

// 	_, dexterContracts := payout.splitDelegationsAndDexterContracts(rewardsSplit)
// 	rewardsSplit.Delegators = tzkt.Delegators{}

// 	for _, contract := range dexterContracts {
// 		contract = payout.constructDelegation(contract, totalRewards, rewardsSplit.StakingBalance)
// 		rewardsSplit.BakerCollectedFees += contract.Fee

// 		var err error
// 		contract, err = payout.constructDexterContractPayoutFunc(contract)
// 		assert.Nil(t, err)
// 		fmt.Println(contract)

// 		rewardsSplit.Delegators = append(rewardsSplit.Delegators, contract)
// 	}

// 	t.Fail()
// }

func Test_constructDexterContractPayout(t *testing.T) {
	type input struct {
		payout   Payout
		contract tzkt.Delegator
	}

	type want struct {
		err         bool
		errContains string
		delegator   tzkt.Delegator
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles failure to getLiquidityProvidersEarnings",
			input{
				payout: Payout{
					rpc: &test.RPCMock{
						CycleErr: true,
					},
					config: config.Config{
						Baker: config.Baker{
							Fee: 0.05,
							DexterLiquidityContracts: []string{
								"tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
							},
						},
					},
				},
				contract: tzkt.Delegator{
					Address:      "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
					GrossRewards: 149992399,
				},
			},
			want{
				true,
				"failed to get earnings for dexter liquidity providers",
				tzkt.Delegator{
					Address:      "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
					GrossRewards: 149992399,
				},
			},
		},
		{
			"is successful",
			input{
				payout: Payout{
					rpc:  &test.RPCMock{},
					tzkt: &test.TzktMock{},
					config: config.Config{
						Baker: config.Baker{
							Fee: 0.05,
							DexterLiquidityContracts: []string{
								"tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
							},
						},
					},
				},
				contract: tzkt.Delegator{
					Address:      "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
					GrossRewards: 149992399,
				},
			},
			want{
				false,
				"",
				tzkt.Delegator{
					Address:      "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
					GrossRewards: 149992399,
					LiquidityProviders: []tzkt.LiquidityProvider{
						{
							Address:      "tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV",
							Balance:      23567891,
							NetRewards:   142492780,
							GrossRewards: 149992399,
							Share:        1,
							Fee:          7499619,
							BlackListed:  false,
						},
					},
					BlackListed: false},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			delegator, err := tt.input.payout.constructDexterContractPayout(tt.input.contract)
			test.CheckErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.delegator, delegator)
		})
	}
}

func Test_getLiquidityProvidersEarnings(t *testing.T) {
	type input struct {
		tzkt     tzkt.IFace
		rpc      rpc.IFace
		contract tzkt.Delegator
	}

	type want struct {
		err         bool
		errContains string
		delegator   tzkt.Delegator
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles failure to get cycle",
			input{
				rpc: &test.RPCMock{
					CycleErr: true,
				},
			},
			want{
				true,
				"failed to get earnings for dexter liquidity providers",
				tzkt.Delegator{},
			},
		},
		{
			"handles failure to get contract storage",
			input{
				rpc: &test.RPCMock{
					ContractStorageErr: true,
				},
			},
			want{
				true,
				"failed to get storage for contract",
				tzkt.Delegator{},
			},
		},
		{
			"handles failure to get liquidity provider list",
			input{
				rpc: &test.RPCMock{},
				tzkt: &test.TzktMock{
					TransactionsErr: true,
				},
			},
			want{
				true,
				"failed to get list of liquidity providers",
				tzkt.Delegator{},
			},
		},
		{
			"handles failure to get get balance from big map",
			input{
				rpc: &test.RPCMock{
					BigMapErr: true,
				},
				tzkt: &test.TzktMock{},
			},
			want{
				true,
				"failed to get earnings for liquidity providers for contract",
				tzkt.Delegator{},
			},
		},
		{
			"is successful",
			input{
				rpc:  &test.RPCMock{},
				tzkt: &test.TzktMock{},
				contract: tzkt.Delegator{
					GrossRewards: 149992399,
				},
			},
			want{
				false,
				"",
				tzkt.Delegator{
					GrossRewards: 149992399,
					LiquidityProviders: []tzkt.LiquidityProvider{
						{
							Address:      "tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV",
							Balance:      23567891,
							NetRewards:   142492780,
							GrossRewards: 149992399,
							Share:        1,
							Fee:          7499619,
							BlackListed:  false,
						},
					},
					BlackListed: false},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := Payout{
				rpc:  tt.input.rpc,
				tzkt: tt.input.tzkt,
				config: config.Config{
					Baker: config.Baker{
						Fee: 0.05,
					},
				},
			}
			delegator, err := payout.getLiquidityProvidersEarnings(tt.input.contract)
			test.CheckErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.delegator, delegator)
		})
	}
}

func Test_getLiquidityProvidersList(t *testing.T) {
	type input struct {
		tzkt tzkt.IFace
	}

	type want struct {
		err         bool
		errContains string
		list        []string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				tzkt: &test.TzktMock{
					TransactionsErr: false,
				},
			},
			want{
				false,
				"",
				[]string{"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV"},
			},
		},
		{
			"handles tzkt failure",
			input{
				tzkt: &test.TzktMock{
					TransactionsErr: true,
				},
			},
			want{
				true,
				"failed to get list of liquidity providers",
				[]string{},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := Payout{
				tzkt: tt.input.tzkt,
			}
			list, err := payout.getLiquidityProvidersList("some_target")
			test.CheckErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.list, list)
		})
	}
}

func Test_getBalanceFromBigMap(t *testing.T) {
	type input struct {
		rpc rpc.IFace
		key string
	}

	type want struct {
		err         bool
		errContains string
		balance     int
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				rpc: &test.RPCMock{
					BigMapErr: false,
				},
				key: "tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV",
			},
			want{
				false,
				"",
				23567891,
			},
		},
		{
			"handles failure to forge script expression",
			input{
				rpc: &test.RPCMock{
					BigMapErr: false,
				},
				key: "dfdsafjj",
			},
			want{
				true,
				"failed to get balance from big_map for",
				0,
			},
		},
		{
			"handles gotezos failure",
			input{
				rpc: &test.RPCMock{
					BigMapErr: true,
				},
				key: "tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV",
			},
			want{
				true,
				"failed to get balance from big_map for 'tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV'",
				0,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := Payout{
				rpc: tt.input.rpc,
			}
			balance, err := payout.getBalanceFromBigMap(tt.input.key, 20999, "address")
			test.CheckErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.balance, balance)
		})
	}
}

func Test_getContractStorage(t *testing.T) {
	storageJSON := []byte(`{"prim":"Pair","args":[{"int":"16033"},{"prim":"Pair","args":[{"prim":"Pair","args":[{"prim":"False"},{"prim":"Pair","args":[{"prim":"False"},{"int":"23567891"}]}]},{"prim":"Pair","args":[{"prim":"Pair","args":[{"string":"tz1S82rGFZK8cVbNDpP1Hf9VhTUa4W8oc2WV"},{"string":"KT1GQcLae1ve1ZEPNfD9z1dyv5ev9ki39SNW"}]},{"prim":"Pair","args":[{"int":"123456"},{"int":"23567891"}]}]}]}]}`)
	var exchangeContractV1 ExchangeContractV1
	err := json.Unmarshal(storageJSON, &exchangeContractV1)
	test.CheckErr(t, false, "", err)

	type input struct {
		rpc rpc.IFace
	}

	type want struct {
		err                bool
		errContains        string
		exchangeContractV1 ExchangeContractV1
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				rpc: &test.RPCMock{
					ContractStorageErr: false,
				},
			},
			want{
				false,
				"",
				exchangeContractV1,
			},
		},
		{
			"handles gotezos failure",
			input{
				rpc: &test.RPCMock{
					ContractStorageErr: true,
				},
			},
			want{
				true,
				"failed to get storage for contract 'address'",
				ExchangeContractV1{},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			payout := Payout{
				rpc: tt.input.rpc,
			}
			exchangeContractV1, err := payout.getContractStorage("block_hash", "address")
			test.CheckErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.exchangeContractV1, exchangeContractV1)
		})
	}
}

func Test_parseBigMapForBalance(t *testing.T) {
	type input struct {
		msg json.RawMessage
	}

	type want struct {
		err         bool
		errContains string
		balance     int
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				[]byte(`[{"int":"23567891"},[]]`),
			},
			want{
				false,
				"",
				23567891,
			},
		},
		{
			"handles malformed object",
			input{
				[]byte(`{{"int":"23567891"},[]]`),
			},
			want{
				true,
				"failed to parse as json blob",
				0,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			balance, err := parseBigMapForBalance(&tt.input.msg)
			test.CheckErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.balance, balance)
		})
	}
}
