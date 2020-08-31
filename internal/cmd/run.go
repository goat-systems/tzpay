package cmd

import (
	"strconv"

	gotezos "github.com/goat-systems/go-tezos/v3"
	"github.com/goat-systems/tzpay/v2/internal/enviroment"
	"github.com/goat-systems/tzpay/v2/internal/payout"
	"github.com/goat-systems/tzpay/v2/internal/print"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Run configures and exposes functions to allow tzpay inject a payout into the tezos network.
type Run struct {
	gt             gotezos.IFace
	bakersFee      float64
	delegate       string
	gasLimit       int
	minimumPayment int
	networkFee     int
	blackList      []string
	wallet         gotezos.Wallet
}

// RunInput is the input for NewDryRun
type RunInput struct {
	GoTezos        gotezos.IFace
	BakersFee      float64
	Delegate       string
	GasLimit       int
	MinimumPayment int
	NetworkFee     int
	BlackList      []string
	Wallet         gotezos.Wallet
}

// NewRun returns a pointer to a Run
func NewRun(input RunInput) *Run {
	return &Run{
		gt:             input.GoTezos,
		bakersFee:      input.BakersFee,
		delegate:       input.Delegate,
		gasLimit:       input.GasLimit,
		minimumPayment: input.MinimumPayment,
		networkFee:     input.NetworkFee,
		blackList:      input.BlackList,
		wallet:         input.Wallet,
	}
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

			NewRun(RunInput{
				GoTezos:        env.GoTezos,
				BakersFee:      env.BakersFee,
				Delegate:       env.Delegate,
				GasLimit:       env.GasLimit,
				MinimumPayment: env.MinimumPayment,
				NetworkFee:     env.NetworkFee,
				BlackList:      env.BlackList,
				Wallet:         env.Wallet,
			}).execute(args[0], batchSize, verbose, table)
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
		GoTezos:    r.gt,
		Cycle:      cycle,
		Delegate:   r.delegate,
		BakerFee:   r.bakersFee,
		MinPayment: r.minimumPayment,
		BlackList:  r.blackList,
		BatchSize:  batchSize,
		NetworkFee: r.networkFee,
		GasLimit:   r.gasLimit,
		Inject:     true,
		Verbose:    verbose,
		Wallet:     r.wallet,
	}).Execute()

	if table {
		print.Table(r.delegate, r.wallet.Address, report)
	} else {
		print.JSON(report)
	}
}
