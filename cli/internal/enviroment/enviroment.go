package enviroment

import (
	"context"
	"strings"

	cenv "github.com/caarlos0/env/v6"
	"github.com/go-playground/validator"
	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/pkg/errors"
)

// ContextKey is a key to context for enviroment or wallet
type ContextKey string

const (
	// ENVIROMENTKEY is a key to get the tzpay enviroment off of context
	ENVIROMENTKEY ContextKey = "ffe0524c-a49c-401c-8b38-556b013fbdd4"
	// WALLETKEY is a key to get the tzpay wallet off of context
	WALLETKEY ContextKey = "43600f87-e1f7-4206-b1cc-7579a1532cc9"
)

// Enviroment is the enviroment for a tzpay baker
type Enviroment struct {
	BakersFee      float64 `validate:"required" env:"TZPAY_BAKERS_FEE"`
	BlackList      string  `env:"TZPAY_BLACKLIST"`
	Delegate       string  `validate:"required" env:"TZPAY_DELEGATE"`
	GasLimit       int     `env:"TZPAY_NETWORK_GAS_LIMIT" envDefault:"26283"`
	HostNode       string  `validate:"required" env:"TZPAY_HOST_NODE"`
	MinimumPayment int     `env:"TZPAY_MINIMUM_PAYMENT" envDefault:"0"`
	EarningsOnly   bool    `env:"TZPAY_EARNINGS_ONLY"` // If this is turned on, tzpay won't pay for missed blocks or endorsements
	NetworkFee     int     `env:"TZPAY_NETWORK_FEE" envDefault:"2941"`
	WalletSecret   string  `validate:"required" env:"TZPAY_WALLET_SECRET"`
	WalletPassword string  `validate:"required" env:"TZPAY_WALLET_PASSWORD"`
}

// ContextEnviroment is the enviroment to be attatched to context
type ContextEnviroment struct {
	BakersFee      float64
	BlackList      string
	Delegate       string
	GasLimit       int
	HostNode       string
	MinimumPayment int
	EarningsOnly   bool
	NetworkFee     int
	Wallet         gotezos.Wallet
}

// GetEnviromentFromContext gets the Enviroment off context
func GetEnviromentFromContext(ctx context.Context) *ContextEnviroment {
	val := ctx.Value(ENVIROMENTKEY)
	env, _ := val.(*ContextEnviroment)
	return env
}

// setEnviromentToContext sets tzpays enviroment to context
func setEnviromentToContext(ctx context.Context, env *Enviroment) (context.Context, error) {
	cenv, err := enviromentToContextEnviroment(*env)
	if err != nil {
		return ctx, errors.Wrap(err, "failed to set context enviroment")
	}

	return context.WithValue(ctx, ENVIROMENTKEY, cenv), nil
}

// InitContext returns and validates the enviroment for tzpay in context
func InitContext() (context.Context, error) {
	env, err := loadEnviroment()
	if err != nil {
		return context.Background(), err
	}
	err = validate(env)
	if err != nil {
		return context.Background(), err
	}

	ctx, err := setEnviromentToContext(context.Background(), env)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func loadEnviroment() (*Enviroment, error) {
	env := &Enviroment{}
	if err := cenv.Parse(env); err != nil {
		return env, errors.Wrap(err, "failed to load enviroment")
	}

	return env, nil
}

func validate(i interface{}) error {
	validate := validator.New()
	if err := validate.Struct(i); err != nil {
		return errors.Wrap(err, "failed to load required parameters")
	}

	return nil
}

// ParseBlackList -
func ParseBlackList(list string) []string {
	blacklist := strings.Split(list, ",")
	for i := range blacklist {
		blacklist[i] = strings.Trim(blacklist[i], " ")
	}

	return blacklist
}

func enviromentToContextEnviroment(env Enviroment) (ContextEnviroment, error) {
	wallet, err := gotezos.ImportEncryptedWallet(env.WalletPassword, env.WalletSecret) // TODO ImportEncryptedWallet second parameter name should be edesk
	if err != nil {
		return ContextEnviroment{}, errors.Wrap(err, "failed to set enviroment to context")
	}

	return ContextEnviroment{
		BakersFee:      env.BakersFee,
		BlackList:      env.BlackList,
		Delegate:       env.Delegate,
		GasLimit:       env.GasLimit,
		HostNode:       env.HostNode,
		MinimumPayment: env.MinimumPayment,
		EarningsOnly:   env.EarningsOnly,
		NetworkFee:     env.NetworkFee,
		Wallet:         *wallet,
	}, nil
}
