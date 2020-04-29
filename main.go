package main

import (
	"github.com/caarlos0/env/v6"
	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v2/internal/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Enviroment is a list of enviroment variables used to configure the tzpay application
type Enviroment struct {
	BakersFee      float64 `env:"TZPAY_BAKERS_FEE" validate:"required" `
	Delegate       string  `env:"TZPAY_DELEGATE" validate:"required"`
	HostNode       string  `env:"TZPAY_HOST_NODE" validate:"required"`
	WalletPassword string  `env:"TZPAY_WALLET_PASSWORD" validate:"required"`

	GasLimit       int `env:"TZPAY_NETWORK_GAS_LIMIT" envDefault:"26283"`
	MinimumPayment int `env:"TZPAY_MINIMUM_PAYMENT" envDefault:"0"`
	NetworkFee     int `env:"TZPAY_NETWORK_FEE" envDefault:"2941"`

	BlackList    []string `env:"TZPAY_BLACKLIST"`
	BoltDB       string   `env:"TZPAY_BOLT_DB"`
	EarningsOnly bool     `env:"TZPAY_EARNINGS_ONLY"`
	WalletSecret string   `env:"TZPAY_WALLET_SECRET"` // Required on the first run
}

func main() {
	rootCommand := &cobra.Command{
		Use:   "tzpay",
		Short: "A bulk payout tool for bakers in the Tezos Ecosystem",
	}

	enviroment := &Enviroment{}
	if err := env.Parse(enviroment); err != nil {
		logrus.WithField("error", err.Error()).Error("Failed to load paramters from enviroment.")
	}

	gt, err := gotezos.New(enviroment.HostNode)
	if err != nil {
		logrus.WithField("error", err.Error()).Error("Failed to make connection to host node.")
	}

	rootCommand.AddCommand(
		cmd.NewDryRunCommand(),
		cmd.NewRunCommand(),
		cmd.NewVersionCommand(),
		cmd.NewSetupCommand(),
	)

	rootCommand.Execute()
}
