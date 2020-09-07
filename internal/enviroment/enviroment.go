package enviroment

import (
	"encoding/json"
	"io/ioutil"

	"github.com/caarlos0/env/v6"
	validate "github.com/go-playground/validator/v10"
	gotezos "github.com/goat-systems/go-tezos/v3"
	"github.com/goat-systems/tzpay/v2/internal/tzkt"
	"github.com/pkg/errors"
)

var (
	newGoTezos   = gotezos.New
	readFile     = ioutil.ReadFile
	importWallet = gotezos.ImportEncryptedWallet
)

// DryRunEnviroment is a list of dry run enviroment variables used to configure the tzpay application
type DryRunEnviroment struct {
	BakersFee      float64 `env:"TZPAY_BAKERS_FEE" validate:"required" `
	Delegate       string  `env:"TZPAY_DELEGATE" validate:"required"`
	GasLimit       int     `env:"TZPAY_NETWORK_GAS_LIMIT" envDefault:"26283" validate:"required"`
	HostNode       string  `env:"TZPAY_HOST_NODE" validate:"required"`
	HostTzkt       string  `env:"TZPAY_HOST_TZKT" validate:"required"`
	MinimumPayment int     `env:"TZPAY_MINIMUM_PAYMENT" envDefault:"100" validate:"required"`
	NetworkFee     int     `env:"TZPAY_NETWORK_FEE" envDefault:"2941" validate:"required"`

	DexterLiquidiyContracts []string `env:"TZPAY_DEXTER_LIQUIDITY_CONTRACTS" envSeparator:","`
	BlackListFile           string   `env:"TZPAY_BLACKLIST"`
	EarningsOnly            bool     `env:"TZPAY_EARNINGS_ONLY"`

	Tzkt      tzkt.IFace    `env:"-"`
	GoTezos   gotezos.IFace `env:"-"`
	BlackList []string      `env:"-"`
}

// AccountSID string
// 	AuthToken  string
// 	From       string
// 	To         []string

// RunEnviroment is a list of dry run enviroment variables used to configure the tzpay application
type RunEnviroment struct {
	BakersFee      float64 `env:"TZPAY_BAKERS_FEE" validate:"required" `
	Delegate       string  `env:"TZPAY_DELEGATE" validate:"required"`
	GasLimit       int     `env:"TZPAY_NETWORK_GAS_LIMIT" envDefault:"26283" validate:"required"`
	HostNode       string  `env:"TZPAY_HOST_NODE" validate:"required"`
	HostTzkt       string  `env:"TZPAY_HOST_TZKT" validate:"required"`
	WalletPassword string  `env:"TZPAY_WALLET_PASSWORD" validate:"required"`
	MinimumPayment int     `env:"TZPAY_MINIMUM_PAYMENT" envDefault:"100" validate:"required"`
	NetworkFee     int     `env:"TZPAY_NETWORK_FEE" envDefault:"2941" validate:"required"`
	WalletSecret   string  `env:"TZPAY_WALLET_SECRET" validate:"required"`

	DexterLiquidiyContracts []string `env:"TZPAY_DEXTER_LIQUIDITY_CONTRACTS" envSeparator:","`
	BlackListFile           string   `env:"TZPAY_BLACKLIST"`
	EarningsOnly            bool     `env:"TZPAY_EARNINGS_ONLY"`

	// Notifications
	TwilioAccountSID      string   `env:"TZPAY_TWILIO_ACCOUNT_SID"`
	TwilioAuthToken       string   `env:"TZPAY_TWILIO_AUTH_TOKEN"`
	TwilioFrom            string   `env:"TZPAY_TWILIO_FROM"`
	TwilioTo              []string `env:"TZPAY_TWILIO_TO" envSeparator:","`
	TwitterConsumerKey    string   `env:"TZPAY_TWITTER_CONSUMER_KEY"`
	TwitterConsumerSecret string   `env:"TZPAY_TWITTER_CONSUMER_SECRET"`
	TwitterAccessToken    string   `env:"TZPAY_TWITTER_ACCESS_TOKEN"`
	TwitterAccessSecret   string   `env:"TZPAY_TWITTER_ACCESS_SECRET"`

	Tzkt      tzkt.IFace     `env:"-"`
	GoTezos   gotezos.IFace  `env:"-"`
	BlackList []string       `env:"-"`
	Wallet    gotezos.Wallet `env:"-"`
}

// NewDryRunEnviroment returns a new DryRunEnviroment
func NewDryRunEnviroment() (*DryRunEnviroment, error) {
	enviroment := &DryRunEnviroment{}
	if err := env.Parse(enviroment); err != nil {
		return nil, errors.Wrap(err, "failed to load paramters from enviroment")
	}

	err := validate.New().Struct(enviroment)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate required enviroment variables")
	}

	gt, err := newGoTezos(enviroment.HostNode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make connection to host node")
	}
	enviroment.GoTezos = gt
	enviroment.Tzkt = tzkt.NewTZKT(enviroment.HostTzkt)

	var blacklist []string
	if enviroment.BlackListFile != "" {
		byts, err := readFile(enviroment.BlackListFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open blacklist file")
		}

		err = json.Unmarshal(byts, &blacklist)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse blacklist file: expected json string array")
		}
	}
	enviroment.BlackList = blacklist

	return enviroment, nil
}

// NewRunEnviroment returns a new RunEnviroment
func NewRunEnviroment() (*RunEnviroment, error) {
	enviroment := &RunEnviroment{}
	if err := env.Parse(enviroment); err != nil {
		return nil, errors.Wrap(err, "failed to load paramters from enviroment")
	}

	err := validate.New().Struct(enviroment)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate required enviroment variables")
	}

	gt, err := newGoTezos(enviroment.HostNode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make connection to host node")
	}
	enviroment.GoTezos = gt
	enviroment.Tzkt = tzkt.NewTZKT(enviroment.HostTzkt)

	var blacklist []string
	if enviroment.BlackListFile != "" {
		byts, err := readFile(enviroment.BlackListFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open blacklist file")
		}

		err = json.Unmarshal(byts, &blacklist)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse blacklist file: expected json string array")
		}
	}
	enviroment.BlackList = blacklist

	wallet, err := importWallet(enviroment.WalletPassword, enviroment.WalletSecret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to import encrypted wallet")
	}
	enviroment.Wallet = *wallet

	return enviroment, nil
}
