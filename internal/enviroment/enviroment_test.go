package enviroment

import (
	"errors"
	"os"
	"testing"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v2/internal/test"
	"github.com/stretchr/testify/assert"
)

func Test_NewDryRunEnviroment(t *testing.T) {
	type input struct {
		newGoTezos func(host string) (*gotezos.GoTezos, error)
		readFile   func(filename string) ([]byte, error)
		enviroment map[string]string
	}

	type want struct {
		err         bool
		errContains string
		dryrun      *DryRunEnviroment
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				newGoTezos: func(host string) (*gotezos.GoTezos, error) {
					return &gotezos.GoTezos{}, nil
				},
				readFile: func(filename string) ([]byte, error) {
					return []byte(`["a","b"]`), nil
				},
				enviroment: map[string]string{
					"TZPAY_BAKERS_FEE": "0.05",
					"TZPAY_DELEGATE":   "some_delegate",
					"TZPAY_HOST_NODE":  "http://127.0.0.1:8732",
					"TZPAY_BLACKLIST":  "some_blacklist_file.json",
				},
			},
			want{
				false,
				"",
				&DryRunEnviroment{
					BakersFee:      0.05,
					Delegate:       "some_delegate",
					GasLimit:       26283,
					HostNode:       "http://127.0.0.1:8732",
					MinimumPayment: 100,
					NetworkFee:     2941,
					BlackListFile:  "some_blacklist_file.json",
					BlackList: []string{
						"a",
						"b",
					},
					GoTezos: &gotezos.GoTezos{},
				},
			},
		},
		{
			"handles missing field",
			input{
				newGoTezos: func(host string) (*gotezos.GoTezos, error) {
					return &gotezos.GoTezos{}, nil
				},
				readFile: func(filename string) ([]byte, error) {
					return []byte(`["a","b"]`), nil
				},
				enviroment: map[string]string{
					"TZPAY_DELEGATE":  "some_delegate",
					"TZPAY_HOST_NODE": "http://127.0.0.1:8732",
					"TZPAY_BLACKLIST": "some_blacklist_file.json",
				},
			},
			want{
				true,
				"failed to validate required enviroment variables",
				nil,
			},
		},
		{
			"handles failure to initialize go tezos",
			input{
				newGoTezos: func(host string) (*gotezos.GoTezos, error) {
					return &gotezos.GoTezos{}, errors.New("some err")
				},
				readFile: func(filename string) ([]byte, error) {
					return []byte(`["a","b"]`), nil
				},
				enviroment: map[string]string{
					"TZPAY_BAKERS_FEE": "0.05",
					"TZPAY_DELEGATE":   "some_delegate",
					"TZPAY_HOST_NODE":  "http://127.0.0.1:8732",
					"TZPAY_BLACKLIST":  "some_blacklist_file.json",
				},
			},
			want{
				true,
				"failed to make connection to host node",
				nil,
			},
		},
		{
			"handles failure to read black list file",
			input{
				newGoTezos: func(host string) (*gotezos.GoTezos, error) {
					return &gotezos.GoTezos{}, nil
				},
				readFile: func(filename string) ([]byte, error) {
					return []byte(`["a","b"]`), errors.New("some_err")
				},
				enviroment: map[string]string{
					"TZPAY_BAKERS_FEE": "0.05",
					"TZPAY_DELEGATE":   "some_delegate",
					"TZPAY_HOST_NODE":  "http://127.0.0.1:8732",
					"TZPAY_BLACKLIST":  "some_blacklist_file.json",
				},
			},
			want{
				true,
				"failed to open blacklist file",
				nil,
			},
		},
		{
			"handles failure unmarshal blacklist file",
			input{
				newGoTezos: func(host string) (*gotezos.GoTezos, error) {
					return &gotezos.GoTezos{}, nil
				},
				readFile: func(filename string) ([]byte, error) {
					return []byte(`junk`), nil
				},
				enviroment: map[string]string{
					"TZPAY_BAKERS_FEE": "0.05",
					"TZPAY_DELEGATE":   "some_delegate",
					"TZPAY_HOST_NODE":  "http://127.0.0.1:8732",
					"TZPAY_BLACKLIST":  "some_blacklist_file.json",
				},
			},
			want{
				true,
				"failed to parse blacklist file",
				nil,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(tt.input.enviroment)

			newGoTezos = tt.input.newGoTezos
			readFile = tt.input.readFile

			env, err := NewDryRunEnviroment()
			test.CheckErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.dryrun, env)

			unsetEnv(tt.input.enviroment)
		})
	}
}

func Test_NewRunEnviroment(t *testing.T) {
	type input struct {
		newGoTezos   func(host string) (*gotezos.GoTezos, error)
		readFile     func(filename string) ([]byte, error)
		importWallet func(password string, esk string) (*gotezos.Wallet, error)
		enviroment   map[string]string
	}

	type want struct {
		err         bool
		errContains string
		dryrun      *RunEnviroment
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				newGoTezos: func(host string) (*gotezos.GoTezos, error) {
					return &gotezos.GoTezos{}, nil
				},
				readFile: func(filename string) ([]byte, error) {
					return []byte(`["a","b"]`), nil
				},
				importWallet: func(password string, esk string) (*gotezos.Wallet, error) {
					return &gotezos.Wallet{}, nil
				},
				enviroment: map[string]string{
					"TZPAY_BAKERS_FEE":      "0.05",
					"TZPAY_DELEGATE":        "some_delegate",
					"TZPAY_HOST_NODE":       "http://127.0.0.1:8732",
					"TZPAY_BLACKLIST":       "some_blacklist_file.json",
					"TZPAY_WALLET_SECRET":   "secret",
					"TZPAY_WALLET_PASSWORD": "password",
				},
			},
			want{
				false,
				"",
				&RunEnviroment{
					BakersFee:      0.05,
					Delegate:       "some_delegate",
					GasLimit:       26283,
					HostNode:       "http://127.0.0.1:8732",
					MinimumPayment: 100,
					NetworkFee:     2941,
					BlackListFile:  "some_blacklist_file.json",
					BlackList: []string{
						"a",
						"b",
					},
					GoTezos:        &gotezos.GoTezos{},
					Wallet:         gotezos.Wallet{},
					WalletPassword: "password",
					WalletSecret:   "secret",
				},
			},
		},
		{
			"handles missing field",
			input{
				newGoTezos: func(host string) (*gotezos.GoTezos, error) {
					return &gotezos.GoTezos{}, nil
				},
				readFile: func(filename string) ([]byte, error) {
					return []byte(`["a","b"]`), nil
				},
				importWallet: func(password string, esk string) (*gotezos.Wallet, error) {
					return &gotezos.Wallet{}, nil
				},
				enviroment: map[string]string{
					"TZPAY_DELEGATE":        "some_delegate",
					"TZPAY_HOST_NODE":       "http://127.0.0.1:8732",
					"TZPAY_BLACKLIST":       "some_blacklist_file.json",
					"TZPAY_WALLET_SECRET":   "secret",
					"TZPAY_WALLET_PASSWORD": "password",
				},
			},
			want{
				true,
				"failed to validate required enviroment variables",
				nil,
			},
		},
		{
			"handles failure to initialize go tezos",
			input{
				newGoTezos: func(host string) (*gotezos.GoTezos, error) {
					return &gotezos.GoTezos{}, errors.New("some err")
				},
				readFile: func(filename string) ([]byte, error) {
					return []byte(`["a","b"]`), nil
				},
				importWallet: func(password string, esk string) (*gotezos.Wallet, error) {
					return &gotezos.Wallet{}, nil
				},
				enviroment: map[string]string{
					"TZPAY_BAKERS_FEE":      "0.05",
					"TZPAY_DELEGATE":        "some_delegate",
					"TZPAY_HOST_NODE":       "http://127.0.0.1:8732",
					"TZPAY_BLACKLIST":       "some_blacklist_file.json",
					"TZPAY_WALLET_SECRET":   "secret",
					"TZPAY_WALLET_PASSWORD": "password",
				},
			},
			want{
				true,
				"failed to make connection to host node",
				nil,
			},
		},
		{
			"handles failure to read black list file",
			input{
				newGoTezos: func(host string) (*gotezos.GoTezos, error) {
					return &gotezos.GoTezos{}, nil
				},
				readFile: func(filename string) ([]byte, error) {
					return []byte(`["a","b"]`), errors.New("some_err")
				},
				importWallet: func(password string, esk string) (*gotezos.Wallet, error) {
					return &gotezos.Wallet{}, nil
				},
				enviroment: map[string]string{
					"TZPAY_BAKERS_FEE":      "0.05",
					"TZPAY_DELEGATE":        "some_delegate",
					"TZPAY_HOST_NODE":       "http://127.0.0.1:8732",
					"TZPAY_BLACKLIST":       "some_blacklist_file.json",
					"TZPAY_WALLET_SECRET":   "secret",
					"TZPAY_WALLET_PASSWORD": "password",
				},
			},
			want{
				true,
				"failed to open blacklist file",
				nil,
			},
		},
		{
			"handles failure unmarshal blacklist file",
			input{
				newGoTezos: func(host string) (*gotezos.GoTezos, error) {
					return &gotezos.GoTezos{}, nil
				},
				readFile: func(filename string) ([]byte, error) {
					return []byte(`junk`), nil
				},
				importWallet: func(password string, esk string) (*gotezos.Wallet, error) {
					return &gotezos.Wallet{}, nil
				},
				enviroment: map[string]string{
					"TZPAY_BAKERS_FEE":      "0.05",
					"TZPAY_DELEGATE":        "some_delegate",
					"TZPAY_HOST_NODE":       "http://127.0.0.1:8732",
					"TZPAY_BLACKLIST":       "some_blacklist_file.json",
					"TZPAY_WALLET_SECRET":   "secret",
					"TZPAY_WALLET_PASSWORD": "password",
				},
			},
			want{
				true,
				"failed to parse blacklist file",
				nil,
			},
		},
		{
			"handles failure import wallet",
			input{
				newGoTezos: func(host string) (*gotezos.GoTezos, error) {
					return &gotezos.GoTezos{}, nil
				},
				readFile: func(filename string) ([]byte, error) {
					return []byte(`["a","b"]`), nil
				},
				importWallet: func(password string, esk string) (*gotezos.Wallet, error) {
					return &gotezos.Wallet{}, errors.New("some_err")
				},
				enviroment: map[string]string{
					"TZPAY_BAKERS_FEE":      "0.05",
					"TZPAY_DELEGATE":        "some_delegate",
					"TZPAY_HOST_NODE":       "http://127.0.0.1:8732",
					"TZPAY_BLACKLIST":       "some_blacklist_file.json",
					"TZPAY_WALLET_SECRET":   "secret",
					"TZPAY_WALLET_PASSWORD": "password",
				},
			},
			want{
				true,
				"failed to import encrypted wallet",
				nil,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(tt.input.enviroment)

			newGoTezos = tt.input.newGoTezos
			readFile = tt.input.readFile
			importWallet = tt.input.importWallet

			env, err := NewRunEnviroment()
			test.CheckErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.dryrun, env)

			unsetEnv(tt.input.enviroment)
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
