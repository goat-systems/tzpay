package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/payman/v2/cmd/internal/delegates"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ReportConf is the configuraion for a payman report
type ReportConf struct {
	Delegate       string  `validate:"required" env:"PAYMAN_DELEGATE"`
	HostNode       string  `validate:"required" env:"PAYMAN_HOST_NODE"`
	BakersFee      float64 `validate:"required" env:"PAYMAN_BAKERS_FEE"`
	MinimumPayment int     `env:"PAYMAN_MINIMUM_PAYMENT"`
	BlackList      string  `env:"PAYMAN_BLACKLIST"`
}

func loadReportEnviroment() *ReportConf {
	cfg := &ReportConf{}
	if err := env.Parse(cfg); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Configuration error.")
	}

	return cfg
}

func validateReportEnviroment(cfg *ReportConf) *ReportConf {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Configuration error.")
	}

	return cfg
}

func report(cycle int) {
	cfg := validateReportEnviroment(loadReportEnviroment())

	gt, err := gotezos.New(cfg.HostNode)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed to initialize GoTezos library.")
	}

	baker := delegates.NewBaker(&delegates.NewBakerInput{
		GoTezos:   gt,
		Address:   cfg.Delegate,
		Fee:       cfg.BakersFee,
		BlackList: delegates.ParseBlackList(cfg.BlackList),
	})

	delegationEarnings, err := baker.GetDelegationEarnings(cycle)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed to get delegation earnings.")
	}

	prettyJSON, err := json.MarshalIndent(delegationEarnings, "", "    ")
	fmt.Println(prettyJSON)
}

func newReportCommand() *cobra.Command {
	var (
		cycle int
		table bool
	)

	var report = &cobra.Command{
		Use:   "report",
		Short: "report simulates a payout and generates a table and csv report",
		Run: func(cmd *cobra.Command, args []string) {
			report(cycle)
		},
	}

	report.PersistentFlags().IntVarP(&cycle, "cycle", "c", 0, "cycle to payout for (e.g. 95)")
	report.PersistentFlags().BoolVarP(&table, "table", "t", false, "formats result into a table (Default: json)")

	return report
}
