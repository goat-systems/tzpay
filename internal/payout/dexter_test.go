package payout

import (
	"encoding/json"
	"testing"

	gotezos "github.com/goat-systems/go-tezos/v3"
	"github.com/goat-systems/tzpay/v2/internal/test"
	"github.com/goat-systems/tzpay/v2/internal/tzkt"
	"github.com/stretchr/testify/assert"
)

func Test_getLiquidityProvidersEarnings(t *testing.T) {
	type input struct {
		tzkt     tzkt.IFace
		gt       gotezos.IFace
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
				gt: &test.GoTezosMock{
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
				gt: &test.GoTezosMock{
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
				gt: &test.GoTezosMock{},
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
				gt: &test.GoTezosMock{
					BigMapErr: true,
				},
				tzkt: &test.TzktMock{},
			},
			want{
				true,
				"failed to get balance from big_map",
				tzkt.Delegator{},
			},
		},
		{
			"is successful",
			input{
				gt:   &test.GoTezosMock{},
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
				gt:       tt.input.gt,
				tzkt:     tt.input.tzkt,
				bakerFee: 0.05,
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
		gt  gotezos.IFace
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
				gt: &test.GoTezosMock{
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
				gt: &test.GoTezosMock{
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
				gt: &test.GoTezosMock{
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
				gt: tt.input.gt,
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
		gt gotezos.IFace
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
				gt: &test.GoTezosMock{
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
				gt: &test.GoTezosMock{
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
				gt: tt.input.gt,
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
