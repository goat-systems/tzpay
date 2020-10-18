package config

import (
	"os"
	"testing"

	"github.com/goat-systems/tzpay/v2/internal/test"
	"github.com/stretchr/testify/assert"
)

func Test_NewDryRunEnviroment(t *testing.T) {

	type want struct {
		err      bool
		contains string
		config   Config
	}

	cases := []struct {
		name  string
		input map[string]string
		want  want
	}{
		{
			"is successful",
			map[string]string{
				"TZPAY_BAKER":                     "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
				"TZPAY_BAKER_FEE":                 "0.05",
				"TZPAY_BAKER_MINIMUM_PAYMENT":     "1000",
				"TZPAY_BAKER_EARNINGS_ONLY":       "True",
				"TZPAY_BAKER_BLACK_LIST":          "some_address,        some_address_2",
				"TZPAY_BAKER_LIQUIDITY_CONTRACTS": "some_contract,        some_contract_2",
				"TZPAY_WALLET_ESK":                "some_esk",
				"TZPAY_WALLET_PASSWORD":           "some_pass",
			},

			want{
				false,
				"",
				Config{
					API{
						TZKT:  "https://api.tzkt.io",
						Tezos: "https://tezos.giganode.io/",
					},
					Baker{
						Address:        "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
						Fee:            0.05,
						MinimumPayment: 1000,
						EarningsOnly:   true,
						Blacklist: []string{
							"some_address",
							"some_address_2",
						},
						DexterLiquidityContracts: []string{
							"some_contract",
							"some_contract_2",
						},
					},
					Key{
						Esk:      "some_esk",
						Password: "some_pass",
					},
					Operations{
						NetworkFee: 2941,
						GasLimit:   26283,
						BatchSize:  125,
					},
					Notifications{},
				},
			},
		},

		{
			"handles required",
			map[string]string{
				"TZPAY_BAKER":                     "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
				"TZPAY_BAKER_MINIMUM_PAYMENT":     "1000",
				"TZPAY_BAKER_EARNINGS_ONLY":       "True",
				"TZPAY_BAKER_BLACK_LIST":          "some_address,        some_address_2",
				"TZPAY_BAKER_LIQUIDITY_CONTRACTS": "some_contract,        some_contract_2",
				"TZPAY_WALLET_ESK":                "some_esk",
				"TZPAY_WALLET_PASSWORD":           "some_pass",
			},
			want{
				true,
				"invalid input",
				Config{
					API{
						TZKT:  "https://api.tzkt.io",
						Tezos: "https://tezos.giganode.io/",
					},
					Baker{
						Address:        "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
						MinimumPayment: 1000,
						EarningsOnly:   true,
						Blacklist: []string{
							"some_address",
							"some_address_2",
						},
						DexterLiquidityContracts: []string{
							"some_contract",
							"some_contract_2",
						},
					},
					Key{
						Esk:      "some_esk",
						Password: "some_pass",
					},
					Operations{
						NetworkFee: 2941,
						GasLimit:   26283,
						BatchSize:  125,
					},
					Notifications{},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(tt.input)
			conf, err := New()
			test.CheckErr(t, tt.want.err, tt.want.contains, err)
			assert.Equal(t, tt.want.config, conf)
			unsetEnv(tt.input)
		})
	}
}

func setEnv(env map[string]string) {
	for key, element := range env {
		os.Setenv(key, element)
	}
}

func unsetEnv(env map[string]string) {
	for key := range env {
		os.Unsetenv(key)
	}
}
