package enviroment

import (
	"context"
	"strings"

	cenv "github.com/caarlos0/env/v6"
	"github.com/go-playground/validator"
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
	BakersFee      float64 `validate:"required" env:"PAYMAN_BAKERS_FEE"`
	BlackList      string  `env:"PAYMAN_BLACKLIST"`
	Delegate       string  `validate:"required" env:"PAYMAN_DELEGATE"`
	GasLimit       int     `env:"PAYMAN_NETWORK_GAS_LIMIT" envDefault:"26283"`
	HostNode       string  `validate:"required" env:"PAYMAN_HOST_NODE"`
	MinimumPayment int     `env:"PAYMAN_MINIMUM_PAYMENT"`
	NetworkFee     int     `env:"PAYMAN_NETWORK_FEE" envDefault:"2941"`
}

// Wallet is the enviroment for a tzpay baker's tezos wallet
type Wallet struct {
	Secret   string `validate:"required" env:"PAYMAN_WALLET_SECRET"`
	Password string `validate:"required" env:"PAYMAN_WALLET_PASSWORD"`
}

// GetEnviromentFromContext gets the Enviroment off context
func GetEnviromentFromContext(ctx context.Context) *Enviroment {
	val := ctx.Value(ENVIROMENTKEY)
	env, _ := val.(*Enviroment)
	return env
}

// GetWalletFromContext gets the Wallet off context
func GetWalletFromContext(ctx context.Context) *Wallet {
	val := ctx.Value(WALLETKEY)
	env, _ := val.(*Wallet)
	return env
}

// SetEnviromentToContext sets tzpays enviroment to context
func SetEnviromentToContext(ctx context.Context, env *Enviroment) context.Context {
	return context.WithValue(ctx, ENVIROMENTKEY, env)
}

// SetWalletToContext sets tzpays wallet to context
func SetWalletToContext(ctx context.Context, wallet *Wallet) context.Context {
	return context.WithValue(ctx, WALLETKEY, wallet)
}

// Parameters returns and validates the enviroment for tzpay
func Parameters() (*Enviroment, error) {
	env, err := loadEnviroment()
	if err != nil {
		return env, err
	}
	err = validate(env)
	if err != nil {
		return env, err
	}

	return env, nil
}

// ParametersWithWallet returns and validates the enviroment and wallet for tzpay
func ParametersWithWallet() (*Enviroment, *Wallet, error) {
	env, err := Parameters()
	if err != nil {
		return env, nil, err
	}
	wallet, err := loadWallet()
	if err != nil {
		return env, wallet, err
	}
	if err = validate(wallet); err != nil {
		return env, wallet, err
	}

	return env, wallet, nil
}

func loadEnviroment() (*Enviroment, error) {
	env := &Enviroment{}
	if err := cenv.Parse(env); err != nil {
		return env, errors.Wrap(err, "failed to load enviroment")
	}

	return env, nil
}

func loadWallet() (*Wallet, error) {
	wallet := &Wallet{}
	if err := cenv.Parse(wallet); err != nil {
		return wallet, errors.Wrap(err, "failed to load wallet")
	}

	return wallet, nil
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
