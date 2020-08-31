package cmd

import (
	"github.com/goat-systems/tzpay/v2/internal/enviroment"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ServCommand returns a new run cobra command
func ServCommand() *cobra.Command {
	var table bool
	var batchSize int
	var verbose bool

	var run = &cobra.Command{
		Use:     "serv",
		Short:   "serv runs a service that will continously payout cycle by cycle",
		Example: `tzpay serv`,
		Run: func(cmd *cobra.Command, args []string) {
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

	run.PersistentFlags().IntVarP(&batchSize, "batch-size", "b", 125, "changes the size of the payout batches (too large may result in failure).")
	run.PersistentFlags().BoolVarP(&verbose, "verbose", "v", true, "will print confirmations in between injections.")
	return run
}
