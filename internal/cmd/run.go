package cmd

import (
	"strconv"

	gotezos "github.com/goat-systems/go-tezos/v3"
	"github.com/goat-systems/tzpay/v2/internal/enviroment"
	"github.com/goat-systems/tzpay/v2/internal/payout"
	"github.com/goat-systems/tzpay/v2/internal/print"
	"github.com/goat-systems/tzpay/v2/internal/tzkt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Run configures and exposes functions to allow tzpay inject a payout into the tezos network.
type Run struct {
	GoTezos        gotezos.IFace
	Tzkt           tzkt.IFace
	BakersFee      float64
	Delegate       string
	GasLimit       int
	MinimumPayment int
	NetworkFee     int
	BlackList      []string
	Wallet         gotezos.Wallet
	EarningsOnly   bool
}

// RunCommand returns a new run cobra command
func RunCommand() *cobra.Command {
	var table bool
	var batchSize int
	var verbose bool

	var run = &cobra.Command{
		Use:     "run",
		Short:   "run executes a batch payout",
		Long:    "run executes a batch payout and prints the result in json or a table",
		Example: `tzpay run <cycle>`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Missing cycle as argument.")
			}

			env, err := enviroment.NewRunEnviroment()
			if err != nil {
				log.WithField("error", err.Error()).Fatal("Failed to load enviroment.")
			}

			runner := Run{
				GoTezos:        env.GoTezos,
				Tzkt:           env.Tzkt,
				BakersFee:      env.BakersFee,
				Delegate:       env.Delegate,
				GasLimit:       env.GasLimit,
				MinimumPayment: env.MinimumPayment,
				NetworkFee:     env.NetworkFee,
				BlackList:      env.BlackList,
				Wallet:         env.Wallet,
				EarningsOnly:   env.EarningsOnly,
			}

			runner.execute(args[0], batchSize, verbose, table)
		},
	}

	run.PersistentFlags().BoolVarP(&table, "table", "t", false, "formats result into a table (Default: json)")
	run.PersistentFlags().IntVarP(&batchSize, "batch-size", "b", 125, "changes the size of the payout batches (too large may result in failure).")
	run.PersistentFlags().BoolVarP(&verbose, "verbose", "v", true, "will print confirmations in between injections.")

	return run
}

func (r *Run) execute(arg string, batchSize int, verbose, table bool) {
	cycle, err := strconv.Atoi(arg)
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to parse cycle argument into integer.")
	}

	report, err := payout.NewPayout(payout.NewPayoutInput{
		GoTezos:      r.GoTezos,
		Tzkt:         r.Tzkt,
		Cycle:        cycle,
		Delegate:     r.Delegate,
		BakerFee:     r.BakersFee,
		MinPayment:   r.MinimumPayment,
		BlackList:    r.BlackList,
		BatchSize:    batchSize,
		NetworkFee:   r.NetworkFee,
		GasLimit:     r.GasLimit,
		Inject:       true,
		Verbose:      verbose,
		Wallet:       r.Wallet,
		EarningsOnly: r.EarningsOnly,
	}).Execute()
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to execute run.")
	}

	if table {
		print.Table(cycle, r.Delegate, report)
	} else {
		print.JSON(report)
	}
}
