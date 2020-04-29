package cmd

import (
	"strconv"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v2/internal/enviroment"
	"github.com/goat-systems/tzpay/v2/internal/payouts"
	"github.com/goat-systems/tzpay/v2/internal/print"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DryRun configures and exposes functions to allow tzpay to simulate a payout without injecting it into the network.
type DryRun struct {
	payouts payouts.Payout
}

// DryRunInput is the input for NewDryRun
type DryRunInput struct {
	GoTezos gotezos.IFace
}

// NewDryRun returns a pointer to a DryRun
func NewDryRun(input DryRunInput) *DryRun {
	return &DryRun{
		payouts: payouts.NewBaker(input.GoTezos),
	}
}

// Command returns the cobra command for dryrun
func (d *DryRun) Command() *cobra.Command {
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
			d.execute(args[0], table)
		},
	}

	report.PersistentFlags().BoolVarP(&table, "table", "t", false, "formats result into a table (Default: json)")

	return report
}

func (d *DryRun) execute(arg string, table bool) {
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

	payouts, err := baker.Payouts(ctx, cycle)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed to get payouts.")
	}

	_, _, err = baker.ForgePayout(ctx, *payouts)
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
