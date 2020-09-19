package config

import (
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
)

// Config encapsulates all configuration possibilities into a single structure
type Config struct {
	API           API
	Baker         Baker
	Key           Key
	Operations    Operations
	Notifications *Notifications
}

// Baker contains configurations related to the how a baker might run their baking operation
type Baker struct {
	Address                      string   `env:"TZPAY_BAKER" validate:"required"`
	Fee                          float64  `env:"TZPAY_BAKER_FEE" validate:"required"`
	MinimumPayment               int      `env:"TZPAY_BAKER_MINIMUM_PAYMENT"`
	EarningsOnly                 bool     `env:"TZPAY_BAKER_EARNINGS_ONLY"`
	DexterLiquidityContractsOnly bool     `env:"TZPAY_BAKER_LIQUIDITY_CONTRACTS_ONLY"`
	Blacklist                    []string `env:"TZPAY_BAKER_BLACK_LIST" envSeparator:","`
	DexterLiquidityContracts     []string `env:"TZPAY_BAKER_LIQUIDITY_CONTRACTS" envSeparator:","`
}

// API contains configurations for the tzkt API and a tezos node
type API struct {
	TZKT  string `env:"TZPAY_API_TZKT" envDefault:"https://api.tzkt.io" validate:"required"`
	Tezos string `env:"TZPAY_API_TEZOS" envDefault:"https://tezos.giganode.io/" validate:"required"`
}

// Operations contains configurations for modifying the actual operation to be injected into a node
type Operations struct {
	NetworkFee int `env:"TZPAY_OPERATIONS_NETWORK_FEE" envDefault:"2941"`
	GasLimit   int `env:"TZPAY_OPERATIONS_GAS_LIMIT" envDefault:"26283"`
	BatchSize  int `env:"TZPAY_OPERATIONS_BATCH_SIZE" envDefault:"125"`
}

// Key contains sensitive information regarding
type Key struct {
	Esk      string `env:"TZPAY_WALLET_ESK" validate:"required"`
	Password string `env:"TZPAY_WALLET_PASSWORD" validate:"required"`
}

// Notifications contains the configurations for notification features
type Notifications struct {
	SignMessage bool `env:"TZPAY_NOTIFICATIONS_SIGN" validate:"required"`
	Twitter     *Twitter
	Twilio      *Twilio
}

// Twitter contains twitter API information for automatic notifications
type Twitter struct {
	ConsumerKey    string `env:"TZPAY_TWITTER_CONSUMER_KEY" validate:"required"`
	ConsumerSecret string `env:"TZPAY_TWITTER_CONSUMER_SECRET" validate:"required"`
	AccessToken    string `env:"TZPAY_TWITTER_ACCESS_TOKEN" validate:"required"`
	AccessSecret   string `env:"TZPAY_TWITTER_ACCESS_SECRET" validate:"required"`
}

// Twilio contains twilio API information for automatic notifications
type Twilio struct {
	AccountSID string   `env:"TZPAY_TWILIO_ACCOUNT_SID" validate:"required"`
	AuthToken  string   `env:"TZPAY_TWILIO_AUTH_TOKEN" validate:"required"`
	From       string   `env:"TZPAY_TWILIO_FROM" validate:"required"`
	To         []string `env:"TZPAY_TWILIO_TO" envSeparator:"," validate:"required"`
}

// New loads enviroment variables into a Config struct
func New() (Config, error) {
	config := Config{}
	if err := env.Parse(&config); err != nil {
		return config, errors.Wrap(err, "failed to load enviroment variables")
	}

	config.Baker.Blacklist = cleanList(config.Baker.Blacklist)
	config.Baker.DexterLiquidityContracts = cleanList(config.Baker.DexterLiquidityContracts)

	if config.Notifications != nil {
		if config.Notifications.Twilio != nil {
			config.Notifications.Twilio.To = cleanList(config.Notifications.Twilio.To)
		}
	}

	err := validator.New().Struct(&config)
	if err != nil {
		return config, errors.Wrap(err, "invalid input")
	}

	return config, nil
}

func cleanList(list []string) []string {
	var out []string
	for _, element := range list {
		out = append(out, strings.Trim(element, " \n\t\r"))
	}

	return out
}
