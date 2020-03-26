package enviroment

import (
	"context"
	"os"
	"strings"

	cenv "github.com/caarlos0/env/v6"
	"github.com/go-playground/validator"
	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v2/cli/internal/db"
	"github.com/pkg/errors"
)

// ContextKey is a key to context for enviroment or wallet
type ContextKey string

type BlackList []string

var newgt = gotezos.New

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
	WalletSecret   string  `env:"TZPAY_WALLET_SECRET"`
	WalletPassword string  `validate:"required" env:"TZPAY_WALLET_PASSWORD"`
	BoltDB         string  `env:"TZPAY_BOLT_DB"`
}

// ContextEnviroment is the enviroment to be attatched to context
type ContextEnviroment struct {
	BakersFee      float64
	BlackList      BlackList
	Delegate       string
	GasLimit       int
	HostNode       string
	MinimumPayment int
	EarningsOnly   bool
	NetworkFee     int
	BoltDB         *db.DB
	Password       string
	GoTezos        gotezos.IFace
	Wallet         gotezos.Wallet
}

// Contains returns true if the pkh is apart of BlackList
func (b *BlackList) Contains(pkh string) bool {
	for _, addr := range *b {
		if addr == pkh {
			return true
		}
	}

	return false
}

// GetEnviromentFromContext gets the Enviroment off context
func GetEnviromentFromContext(ctx context.Context) *ContextEnviroment {
	val := ctx.Value(ENVIROMENTKEY)
	env, _ := val.(*ContextEnviroment)
	return env
}

// setEnviromentToContext sets tzpays enviroment to context
func setEnviromentToContext(ctx context.Context, env *Enviroment, gt gotezos.IFace) (context.Context, error) {
	cenv, err := enviromentToContextEnviroment(*env, gt)
	if err != nil {
		return ctx, errors.Wrap(err, "failed to set context enviroment")
	}

	return context.WithValue(ctx, ENVIROMENTKEY, &cenv), nil
}

// InitContext returns and validates the enviroment for tzpay in context
func InitContext(gt gotezos.IFace) (context.Context, error) {
	env, err := loadEnviroment()
	if err != nil {
		return context.Background(), err
	}
	err = validate(env)
	if err != nil {
		return context.Background(), err
	}

	ctx, err := setEnviromentToContext(context.Background(), env, gt)
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

func enviromentToContextEnviroment(env Enviroment, gt gotezos.IFace) (ContextEnviroment, error) {
	if gt == nil {
		var err error
		gt, err = newgt(env.HostNode)
		if err != nil {
			return ContextEnviroment{}, errors.Wrap(err, "failed to connect to host node")
		}
	}

	// open tzpay db
	db, err := db.New(gt, env.BoltDB)
	if err != nil {
		return ContextEnviroment{}, errors.Wrap(err, "failed to open tzpay db")
	}

	// check if tzpay is initialized with a wallet by seeing if there is an edesk in the store
	if env.WalletSecret == "" {
		init := db.IsWalletInitialized()
		if !init {
			return ContextEnviroment{}, errors.New("failed to find existing wallet: initialize tzpay by passing TZPAY_WALLET_SECRET and TZPAY_WALLET_PASSWORD: please refer to the README")
		}
	}

	if env.WalletSecret != "" {
		err := db.InitWallet(env.WalletPassword, env.WalletSecret)
		if err != nil {
			return ContextEnviroment{}, errors.Wrap(err, "failed to initialize wallet")
		}

		err = os.Unsetenv("TZPAY_WALLET_SECRET")
		if err != nil {
			return ContextEnviroment{}, errors.Wrap(err, "failed to unset TZPAY_WALLET_SECRET from enviroment")
		}
	}

	secret, err := db.GetSecret(env.WalletPassword)
	if err != nil {
		return ContextEnviroment{}, errors.Wrap(err, "failed to get secret")
	}

	wallet, err := gotezos.ImportEncryptedWallet(env.WalletPassword, secret) // TODO ImportEncryptedWallet second parameter name should be edesk
	if err != nil {
		return ContextEnviroment{}, errors.Wrap(err, "failed to set enviroment to context")
	}

	blacklist := strings.Split(env.BlackList, " ,")

	return ContextEnviroment{
		BakersFee:      env.BakersFee,
		BlackList:      blacklist,
		Delegate:       env.Delegate,
		GasLimit:       env.GasLimit,
		HostNode:       env.HostNode,
		MinimumPayment: env.MinimumPayment,
		EarningsOnly:   env.EarningsOnly,
		NetworkFee:     env.NetworkFee,
		GoTezos:        gt,
		BoltDB:         db,
		Wallet:         *wallet,
	}, nil
}
