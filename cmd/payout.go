package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/DefinitelyNotAGoat/payman/reddit"

	"github.com/DefinitelyNotAGoat/payman/options"
	pay "github.com/DefinitelyNotAGoat/payman/payer"
	"github.com/DefinitelyNotAGoat/payman/reporting"
	"github.com/DefinitelyNotAGoat/payman/server"

	goTezos "github.com/DefinitelyNotAGoat/go-tezos"
	"github.com/spf13/cobra"
)

func newPayoutCommand() *cobra.Command {
	var conf options.Options

	preflight := func(conf options.Options) {
		errors := []string{}
		warnings := []string{}
		if conf.Delegate == "" {
			errors = append(errors, "[payout][preflight] error: no delegate passed for payout (e.g. --delegate=<pkh>)")
		}
		if conf.Secret == "" {
			errors = append(errors, "[payout][preflight] error: no secret key passed for payout wallet (e.g. --secret=<sk>)")
		}
		if conf.Password == "" {
			errors = append(errors, "[payout][preflight] error: no password passed for payout wallet (e.g. --password=<passwd>)")
		}
		if conf.Cycles == "" && conf.Cycle == 0 {
			errors = append(errors, "[payout][preflight] error: no cycle(s) passed to payout for (e.g. --cycle=95 || --cycles=95-100)")
		}
		if conf.Cycles != "" && conf.Cycle > 0 {
			errors = append(errors, "[payout][preflight] error: cannot pass both --cycles and --cycle, it's either or (e.g. --cycle=95 || --cycles=95-100)")
		}
		if conf.Fee == -1 {
			errors = append(errors, "[payout][preflight] error: no delegation fee passed for payout (e.g. --fee=0.05)")
		}
		if conf.NetworkFee == 1270 {
			warnings = append(warnings, "[payout][preflight] warning: no network fee passed for payout, using default 1270 mutez")
		}
		if conf.NetworkFee == 1270 {
			warnings = append(warnings, "[payout][preflight] warning: no gas limit passed for payout, using default 10200 mutez")
		}

		for _, err := range errors {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, warning := range warnings {
			fmt.Println(warning)
		}
	}

	var payout = &cobra.Command{
		Use:   "payout",
		Short: "Payout pays out rewards to delegations.",
		Long:  "Payout pays out rewards to delegations for the delegate passed.",
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
			wallet, err := gt.ImportEncryptedWallet(conf.Password, conf.Secret)
			if err != nil {
				reporter.Log(fmt.Sprintf("could not import wallet: %v", err))
				os.Exit(1)
			}

			var redditBot *reddit.Bot
			if conf.RedditAgent != "" {
				redditBot, err = reddit.NewRedditSession(conf.RedditAgent, "dng_delegation", conf.RedditTitle)
				if err != nil {
					reporter.Log(fmt.Sprintf("could not start reddit bot: %v", err))
				}
			}

			if conf.Service == true {
				var serv server.PayoutServer
				if redditBot != nil {
					serv = server.NewPayoutServer(gt, wallet, reporter, redditBot, &conf)
				} else {
					serv = server.NewPayoutServer(gt, wallet, reporter, nil, &conf)
				}
				serv.Serve()

			} else {
				payer := pay.NewPayer(gt, wallet, &conf)
				payouts, ops, err := payer.Payout()
				if err != nil {
					log.Fatal(err)
				}

				for _, op := range ops {
					reporter.Log("Successful operation: " + string(op))
					if conf.RedditAgent != "" {
						err := redditBot.Post(string(op), conf.Cycles)
						if err != nil {
							reporter.Log(fmt.Sprintf("could not post to reddit: %v", err))
						}
					}
				}
				reporter.PrintPaymentsTable(payouts)
				reporter.WriteCSVReport(payouts)
			}

			f.Close()
		},
	}

	payout.PersistentFlags().StringVarP(&conf.Delegate, "delegate", "d", "", "public key hash of the delegate that's paying out (e.g. --delegate=<phk>)")
	payout.PersistentFlags().StringVarP(&conf.Secret, "secret", "s", "", "encrypted secret key of the wallet paying (e.g. --secret=<sk>)")
	payout.PersistentFlags().StringVarP(&conf.Password, "password", "k", "", "password to the secret key of the wallet paying (e.g. --password=<passwd>)")
	payout.PersistentFlags().BoolVar(&conf.Service, "serve", false, "run service to payout for all new cycles going foward (default false)(e.g. --serve)")
	payout.PersistentFlags().BoolVar(&conf.Dry, "dry", false, "run payout in simulation with report (default false)(e.g. --dry)")
	payout.PersistentFlags().StringVar(&conf.Cycles, "cycles", "", "cycles to payout for (e.g. 95-100)")
	payout.PersistentFlags().IntVarP(&conf.Cycle, "cycle", "c", 0, "cycle to payout for (e.g. 95)")
	payout.PersistentFlags().StringVarP(&conf.Node, "node", "n", "http://127.0.0.1", "address to the node to query (default http://127.0.0.1)(e.g. mainnet-node.tzscan.io)")
	payout.PersistentFlags().StringVarP(&conf.Port, "port", "p", "8732", "port to use for node (default 8732)(e.g. 443)")
	payout.PersistentFlags().Float32VarP(&conf.Fee, "fee", "f", -1, "fee for the delegate (e.g. 0.05 = 5%)")
	payout.PersistentFlags().IntVar(&conf.NetworkFee, "network-fee", 1270, "network fee for each transaction in mutez (default 1270)(e.g. 2000)")
	payout.PersistentFlags().IntVar(&conf.NetworkGasLimit, "gas-limit", 10200, "network gas limit for each transaction in mutez (default 10200)(e.g. 10300)")
	payout.PersistentFlags().StringVarP(&conf.File, "log-file", "l", "/dev/stdout", "file to log to (default stdout)(e.g. ./payman.log)")
	payout.PersistentFlags().StringVarP(&conf.RedditAgent, "reddit", "r", "", "path to reddit agent file (initiates reddit bot)(e.g. https://turnage.gitbooks.io/graw/content/chapter1.html)")
	payout.PersistentFlags().StringVar(&conf.RedditTitle, "reddit-title", "", "pre title for the reddit bot to post (e.g. DefinitelyNotABot: -- will read DefinitelyNotABot: Payout for Cycle(s) <cycles>)")

	return payout
}
