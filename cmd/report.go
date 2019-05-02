package cmd

import (
	"fmt"
	"log"
	"os"

	goTezos "github.com/DefinitelyNotAGoat/go-tezos"
	"github.com/DefinitelyNotAGoat/payman/options"
	pay "github.com/DefinitelyNotAGoat/payman/payer"
	"github.com/DefinitelyNotAGoat/payman/reporting"
	"github.com/spf13/cobra"
)

func newReportCommand() *cobra.Command {
	var conf options.Options

	preflight := func(conf options.Options) {
		errors := []string{}
		if conf.Delegate == "" {
			errors = append(errors, "[payout][preflight] error: no delegate passed for payout (e.g. --delegate=<pkh>)")
		}
		if conf.Cycle == 0 {
			errors = append(errors, "[payout][preflight] error: no cycle passed to payout for (e.g. --cycle=95)")
		}
		if conf.Fee == -1 {
			errors = append(errors, "[payout][preflight] error: no delegation fee passed for payout (e.g. --fee=0.05)")
		}

		for _, err := range errors {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	var report = &cobra.Command{
		Use:   "report",
		Short: "report simulates a payout and generates a table and csv report",
		Run: func(cmd *cobra.Command, args []string) {

			preflight(conf)

			f, err := os.Create(conf.File)
			if err != nil {
				fmt.Printf("could not open logging file: %v\n", err)
			}

			log := log.New(f, "", log.Ldate|log.Ltime|log.Lshortfile)

			reporter, err := reporting.NewReporter(log)
			if err != nil {
				reporter.Log(fmt.Sprintf("could not open file for reporting: %v\n", err))
			}

			gt := goTezos.NewGoTezos()
			gt.AddNewClient(goTezos.NewTezosRPCClient(conf.Node, conf.Port))
			conf.Dry = true

			wallet := goTezos.Wallet{}
			payer := pay.NewPayer(gt, wallet, &conf)
			payouts, _, err := payer.Payout()
			if err != nil {
				log.Fatal(err)
			}

			reporter.PrintPaymentsTable(payouts)
			reporter.WriteCSVReport(payouts)

			f.Close()
		},
	}

	report.PersistentFlags().StringVarP(&conf.Delegate, "delegate", "d", "", "public key hash of the delegate that's paying out (e.g. --delegate=<phk>)")
	report.PersistentFlags().IntVarP(&conf.Cycle, "cycle", "c", 0, "cycle to payout for (e.g. 95)")
	report.PersistentFlags().StringVarP(&conf.Node, "node", "n", "http://127.0.0.1", "address to the node to query (default http://127.0.0.1)(e.g. mainnet-node.tzscan.io)")
	report.PersistentFlags().StringVarP(&conf.Port, "port", "p", "8732", "port to use for node (default 8732)(e.g. 443)")
	report.PersistentFlags().Float32VarP(&conf.Fee, "fee", "f", -1, "fee for the delegate (e.g. 0.05 = 5%)")
	report.PersistentFlags().StringVarP(&conf.File, "log-file", "l", "/dev/stdout", "file to log to (default stdout)(e.g. ./payman.log)")

	return report
}
