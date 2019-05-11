package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/DefinitelyNotAGoat/payman/reddit"
	"github.com/DefinitelyNotAGoat/payman/twitter"

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

		if conf.Secret == "" {
			errors = append(errors, "[payout][preflight] error: no secret key passed for payout wallet (e.g. --secret=<sk>)")
		}
		if conf.Password == "" {
			errors = append(errors, "[payout][preflight] error: no password passed for payout wallet (e.g. --password=<passwd>)")
		}

		if conf.PaymentsOverride.File == "" {
			if conf.Cycle == 0 && !conf.Service {
				errors = append(errors, "[payout][preflight] error: no cycle passed to payout for (e.g. --cycle=95)")
			}
			if conf.Fee == -1 {
				errors = append(errors, "[payout][preflight] error: no delegation fee passed for payout (e.g. --fee=0.05)")
			}
			if conf.Delegate == "" {
				errors = append(errors, "[payout][preflight] error: no delegate passed for payout (e.g. --delegate=<pkh>)")
			}
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

			gt, err := goTezos.NewGoTezos(conf.URL)
			if err != nil {
				reporter.Log(fmt.Sprintf("could not connect to network: %v\n", err))
			}

			if conf.PaymentsOverride.File != "" {
				conf.PaymentsOverride.Payments, err = conf.PaymentsOverride.ReadPaymentsOverride()
				if err != nil {
					reporter.Log(fmt.Sprintf("could not parse payments override into payments: %v", err))
					os.Exit(1)
				}
			}

			wallet, err := gt.Account.ImportEncryptedWallet(conf.Password, conf.Secret)
			if err != nil {
				reporter.Log(fmt.Sprintf("could not import wallet: %v", err))
				os.Exit(1)
			}

			var redditBot *reddit.Bot
			var redditBotStatus bool
			if conf.RedditAgent != "" {
				redditBot, err = reddit.NewRedditSession(conf.RedditAgent, "dng_delegation", conf.RedditTitle)
				if err != nil {
					reporter.Log(fmt.Sprintf("could not start reddit bot: %v", err))
				} else {
					redditBotStatus = true
				}
			}

			var twitterBot *twitter.Bot
			var twitterBotStatus bool
			if conf.Twitter {
				twitterBot, err = twitter.NewTwitterSession(conf.TwitterPath, conf.TwitterTitle)
				if err != nil {
					reporter.Log(fmt.Sprintf("could not start twitter bot: %v", err))
				} else {
					twitterBotStatus = true
				}
			}

			if conf.Service {

				serv := server.NewPayoutServer(gt, wallet, reporter, redditBot, twitterBot, &conf)
				serv.Serve()

			} else {
				payer := pay.NewPayer(gt, wallet, &conf)
				payouts, ops, err := payer.Payout()
				if err != nil {
					log.Fatal(err)
				}

				for _, op := range ops {
					reporter.Log("Successful operation: " + string(op))
					if conf.RedditAgent != "" && redditBotStatus {
						err := redditBot.Post(string(op), conf.Cycle)
						if err != nil {
							reporter.Log(fmt.Sprintf("could not post to reddit: %v", err))
						}
					}

					if conf.Twitter && twitterBotStatus {
						err := twitterBot.Post(string(op), conf.Cycle)
						if err != nil {
							reporter.Log(fmt.Sprintf("could not post to twitter: %v", err))
						}
					}
				}
				if len(conf.PaymentsOverride.Payments) == 0 {
					reporter.PrintPaymentsTable(payouts)
					reporter.WriteCSVReport(payouts)
				}
			}

			f.Close()
		},
	}

	payout.PersistentFlags().StringVarP(&conf.Delegate, "delegate", "d", "", "public key hash of the delegate that's paying out (e.g. --delegate=<phk>)")
	payout.PersistentFlags().StringVarP(&conf.Secret, "secret", "s", "", "encrypted secret key of the wallet paying (e.g. --secret=<sk>)")
	payout.PersistentFlags().StringVarP(&conf.Password, "password", "k", "", "password to the secret key of the wallet paying (e.g. --password=<passwd>)")
	payout.PersistentFlags().BoolVar(&conf.Service, "serve", false, "run service to payout for all new cycles going foward (default false)(e.g. --serve)")
	payout.PersistentFlags().IntVarP(&conf.Cycle, "cycle", "c", 0, "cycle to payout for (e.g. 95)")
	payout.PersistentFlags().StringVarP(&conf.URL, "node", "u", "http://127.0.0.1:8732", "address to the node to query (default http://127.0.0.1:8732)(e.g. https://mainnet-node.tzscan.io:443)")
	payout.PersistentFlags().Float32VarP(&conf.Fee, "fee", "f", -1, "fee for the delegate (e.g. 0.05 = 5%)")
	payout.PersistentFlags().IntVar(&conf.NetworkFee, "network-fee", 1270, "network fee for each transaction in mutez (default 1270)(e.g. 2000)")
	payout.PersistentFlags().IntVar(&conf.NetworkGasLimit, "gas-limit", 10200, "network gas limit for each transaction in mutez (default 10200)(e.g. 10300)")
	payout.PersistentFlags().StringVarP(&conf.File, "log-file", "l", "/dev/stdout", "file to log to (default stdout)(e.g. ./payman.log)")
	payout.PersistentFlags().StringVarP(&conf.RedditAgent, "reddit", "r", "", "path to reddit agent file (initiates reddit bot)(e.g. https://turnage.gitbooks.io/graw/content/chapter1.html)")
	payout.PersistentFlags().StringVar(&conf.RedditTitle, "reddit-title", "", "pre title for the reddit bot to post (e.g. DefinitelyNotABot: -- will read DefinitelyNotABot: Payout for Cycle <cycle>)")
	payout.PersistentFlags().StringVar(&conf.TwitterPath, "twitter-path", "", "path to twitter.yml file containing API keys if not in current dir (e.g. path/to/my/file/)")
	payout.PersistentFlags().StringVar(&conf.TwitterTitle, "twitter-title", "", "pre title for the twitter bot to post (e.g. DefinitelyNotABot: -- will read DefinitelyNotABot: Payout for Cycle <cycle>)")
	payout.PersistentFlags().BoolVarP(&conf.Twitter, "twitter", "t", false, "turn on twitter bot, will look for api keys in twitter.yml in current dir or --twitter-path (e.g. --twitter)")
	payout.PersistentFlags().StringVar(&conf.PaymentsOverride.File, "payments-override", "", "overrides the rewards calculation and allows you to pass in your own payments in a json file (e.g. path/to/my/file/payments.json)")
	return payout
}
