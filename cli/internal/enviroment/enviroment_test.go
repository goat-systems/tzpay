package enviroment

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetEnviromentFromContext(t *testing.T) {
	env := &ContextEnviroment{}
	ctx := context.WithValue(context.Background(), ENVIROMENTKEY, env)
	outenv := GetEnviromentFromContext(ctx)
	assert.Equal(t, env, outenv)
}
func Test_validate(t *testing.T) {
	type Test struct {
		Someval int `validate:"required"`
	}

	cases := []struct {
		name    string
		input   Test
		wantErr bool
	}{
		{
			"it is successful",
			Test{10},
			false,
		},
		{
			"it handles error",
			Test{},
			true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(tt.input)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_loadEnviroment(t *testing.T) {

	type want struct {
		err bool
		env *Enviroment
	}

	cases := []struct {
		name string
		env  map[string]string
		want want
	}{
		{
			"it is successful",
			map[string]string{
				"TZPAY_BAKERS_FEE":        "0.05",
				"TZPAY_BLACKLIST":         "somehash, somehash1",
				"TZPAY_DELEGATE":          "somedelegate",
				"TZPAY_NETWORK_GAS_LIMIT": "100000",
				"TZPAY_HOST_NODE":         "http://somenode.com:8732",
				"TZPAY_MINIMUM_PAYMENT":   "1000",
				"TZPAY_NETWORK_FEE":       "100000",
				"TZPAY_WALLET_SECRET":     "edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2",
				"TZPAY_WALLET_PASSWORD":   "password12345",
			},
			want{
				false,
				&Enviroment{
					BakersFee:      0.05,
					BlackList:      "somehash, somehash1",
					Delegate:       "somedelegate",
					GasLimit:       100000,
					HostNode:       "http://somenode.com:8732",
					MinimumPayment: 1000,
					EarningsOnly:   false,
					NetworkFee:     100000,
					WalletSecret:   "edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2",
					WalletPassword: "password12345",
					BoltDB:         "",
				},
			},
		},
		{
			"is sucessful with missing fields",
			map[string]string{
				"TZPAY_BAKERS_FEE": "0.05",
			},
			want{
				false,
				&Enviroment{
					BakersFee:  0.05,
					GasLimit:   26283,
					NetworkFee: 2941,
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			for key, elem := range tt.env {
				os.Setenv(key, string(elem))
			}

			env, err := loadEnviroment()
			if tt.want.err {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.env, env)

			for key := range tt.env {
				os.Unsetenv(key)
			}
		})
	}
}

func Test_ParseBlackList(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			"is successful",
			"some_address, some_other_address, yet_another_address",
			[]string{
				"some_address",
				"some_other_address",
				"yet_another_address",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			out := ParseBlackList(tt.input)
			assert.Equal(t, tt.want, out)
		})
	}
}
