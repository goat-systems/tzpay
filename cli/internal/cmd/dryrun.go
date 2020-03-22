package cmd

import (
	"strconv"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v2/cli/internal/baker"
	"github.com/goat-systems/tzpay/v2/cli/internal/enviroment"
	"github.com/goat-systems/tzpay/v2/cli/internal/print"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewDryRunCommand returns a new dryrun cobra command
func NewDryRunCommand() *cobra.Command {
	var table bool

	var report = &cobra.Command{
		Use:     "dryrun",
		Short:   "dryrun simulates a payout",
		Long:    "dryrun simulates a payout and prints the result in json or a table",
		Example: `tzpay dryrun <cycle>`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.WithFields(nil).Fatal("Missing cycle as argument.")
			}
			dryrun(args[0], table)
		},
	}

	report.PersistentFlags().BoolVarP(&table, "table", "t", false, "formats result into a table (Default: json)")

	return report
}

func dryrun(arg string, table bool) {
	cycle, err := strconv.Atoi(arg)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed to read cycle argument.")
	}

	ctx, err := enviroment.InitContext(nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed to get enviroment and initialize context.")
	}

	base := enviroment.GetEnviromentFromContext(ctx)

	gt, err := gotezos.New(base.HostNode)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed to connect to node.")
	}
	baker := baker.NewBaker(gt)

	payouts, err := baker.Payouts(ctx, cycle)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed to get payouts.")
	}

	_, err = baker.ForgePayout(ctx, *payouts)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed to forge operation.")
	}

	if table {
		print.Table(ctx, payouts)
	} else {
		print.JSON(payouts)
	}
}
