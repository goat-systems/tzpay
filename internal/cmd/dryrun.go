package cmd

import (
	"strconv"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v2/internal/enviroment"
	"github.com/goat-systems/tzpay/v2/internal/payout"
	"github.com/goat-systems/tzpay/v2/internal/print"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DryRun configures and exposes functions to allow tzpay to simulate a payout without injecting it into the network.
type DryRun struct {
	gt             gotezos.IFace
	bakersFee      float64
	delegate       string
	gasLimit       int
	minimumPayment int
	networkFee     int
	blackList      []string
}

// DryRunInput is the input for NewDryRun
type DryRunInput struct {
	GoTezos        gotezos.IFace
	BakersFee      float64
	Delegate       string
	GasLimit       int
	MinimumPayment int
	NetworkFee     int
	BlackList      []string
}

// NewDryRun returns a pointer to a DryRun
func NewDryRun(input DryRunInput) *DryRun {
	return &DryRun{
		gt:             input.GoTezos,
		bakersFee:      input.BakersFee,
		delegate:       input.Delegate,
		gasLimit:       input.GasLimit,
		minimumPayment: input.MinimumPayment,
		networkFee:     input.NetworkFee,
		blackList:      input.BlackList,
	}
}

// DryRunCommand returns the cobra command for dryrun
func DryRunCommand() *cobra.Command {
	var table bool

	var dryrun = &cobra.Command{
		Use:     "dryrun",
		Short:   "dryrun simulates a payout",
		Long:    "dryrun simulates a payout and prints the result in json or a table",
		Example: `tzpay dryrun <cycle>`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Missing cycle as argument.")
			}

			env, err := enviroment.NewDryRunEnviroment()
			if err != nil {
				log.WithField("error", err.Error()).Fatal("Failed to load enviroment.")
			}

			NewDryRun(DryRunInput{
				GoTezos:        env.GoTezos,
				BakersFee:      env.BakersFee,
				Delegate:       env.Delegate,
				GasLimit:       env.GasLimit,
				MinimumPayment: env.MinimumPayment,
				NetworkFee:     env.NetworkFee,
				BlackList:      env.BlackList,
			}).execute(args[0], table)
		},
	}

	dryrun.PersistentFlags().BoolVarP(&table, "table", "t", false, "formats result into a table (Default: json)")

	return dryrun
}

func (d *DryRun) execute(arg string, table bool) {
	cycle, err := strconv.Atoi(arg)
	if err != nil {
		log.Fatal("Failed to parse cycle argument into integer.")
	}

	report, err := payout.NewPayout(payout.NewPayoutInput{
		GoTezos:    d.gt,
		Cycle:      cycle,
		Delegate:   d.delegate,
		BakerFee:   d.bakersFee,
		MinPayment: d.minimumPayment,
		BlackList:  d.blackList,
		BatchSize:  125,
		NetworkFee: d.networkFee,
		GasLimit:   d.gasLimit,
	}).Execute()

	if table {
		print.Table(d.delegate, "", report)
	} else {
		print.JSON(report)
	}
}
