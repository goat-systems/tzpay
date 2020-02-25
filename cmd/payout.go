package cmd

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
	gotezos "github.com/goat-systems/go-tezos/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// PayoutConf is the configuraion for a payman payout
type PayoutConf struct {
	Delegate       string `validate:"required" env:"PAYMAN_DELEGATE"`
	WalletSecret   string `validate:"required" env:"PAYMAN_WALLET_SECRET"`
	WalletPassword string `validate:"required" env:"PAYMAN_WALLET_PASSWORD"`
	HostNode       string `validate:"required" env:"PAYMAN_HOST_NODE"`
	NetworkFee     string `env:"PAYMAN_NETWORK_FEE" envDefault:"2941"`
	GasLimit       string `env:"PAYMAN_NETWORK_GAS_LIMIT" envDefault:"26283"`
	BakersFee      string `validate:"required" env:"PAYMAN_BAKERS_FEE"`
	MinimumPayment int    `env:"PAYMAN_MINIMUM_PAYMENT"`
	BlackList      string `env:"PAYMAN_BLACKLIST"`
}

func loadEnviroment() *PayoutConf {
	cfg := &PayoutConf{}
	if err := env.Parse(cfg); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Configuration error.")
	}

	return cfg
}

func validateEnviroment(cfg *PayoutConf) *PayoutConf {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Configuration error.")
	}

	return cfg
}

func payout(cycle int) {
	cfg := validateEnviroment(loadEnviroment())
	_, err := gotezos.ImportEncryptedWallet(cfg.WalletPassword, cfg.WalletSecret)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed to import wallet.")
	}

}

func newPayoutCommand() *cobra.Command {
	var cycle int

	var payout = &cobra.Command{
		Use:   "payout",
		Short: "Payout pays out rewards to delegations.",
		Long:  "Payout pays out rewards to delegations for the delegate passed.",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	payout.PersistentFlags().IntVarP(&cycle, "cycle", "c", 0, "cycle to payout for (e.g. 95)")
	return payout
}
