package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/DefinitelyNotAGoat/payman/reddit"

	pay "github.com/DefinitelyNotAGoat/payman/payer"
	"github.com/DefinitelyNotAGoat/payman/reporting"
	"github.com/DefinitelyNotAGoat/payman/server"

	goTezos "github.com/DefinitelyNotAGoat/go-tezos"
	"github.com/spf13/cobra"
)

var (
	delegate        string
	secret          string
	password        string
	service         bool
	cycles          string
	cycle           int
	reCycles        = regexp.MustCompile(`([0-9]+)`)
	node            string
	port            string
	fee             float32
	file            string
	networkFee      int
	networkGasLimit int
	dry             bool
	redditAgent     string
	title           string
)

var payout = &cobra.Command{
	Use:   "payout",
	Short: "payout pays out rewards to delegations",
	Long:  `payman is a simple golang example that demonstartes go-tezos, and also allows you to payout your delegations.`,
	Run: func(cmd *cobra.Command, args []string) {

		f, err := os.Create(file)
		if err != nil {
			fmt.Printf("could not open logging file: %v\n", err)
		}

		log := log.New(f, "", log.Ldate|log.Ltime|log.Lshortfile)

		reporter, err := reporting.NewReporter(log)
		if err != nil {
			reporter.Log(fmt.Sprintf("could not open file for reporting: %v\n", err))
		}

		if delegate == "" {
			reporter.Log("you need to specify the delegate that is paying out")
			os.Exit(0)
		} else if secret == "" {
			reporter.Log("you need to set the encrypted secret to the wallet paying")
			os.Exit(0)
		} else if password == "" {
			reporter.Log("you need to set the password to the wallet paying")
			os.Exit(0)
		}

		gt := goTezos.NewGoTezos()
		gt.AddNewClient(goTezos.NewTezosRPCClient(node, port))
		wallet, err := gt.ImportEncryptedWallet(password, secret)
		if err != nil {
			reporter.Log(fmt.Sprintf("could not import wallet: %v", err))
			os.Exit(1)
		}

		payer := pay.Payer{}
		if dry {
			payer = pay.NewPayer(gt, wallet, delegate, fee, false)
		} else {
			payer = pay.NewPayer(gt, wallet, delegate, fee, true)
		}

		var redditBot *reddit.Bot
		if redditAgent != "" {
			redditBot, err = reddit.NewRedditSession(redditAgent, "dng_delegation", title)
			if err != nil {
				reporter.Log(fmt.Sprintf("could not start reddit bot: %v", err))
			}
		}

		if service == true {
			var serv server.PaymanServer
			if redditBot != nil {
				serv = server.NewPaymanServer(delegate, fee, networkFee, networkGasLimit, gt, wallet, payer, reporter, redditBot)
			} else {
				serv = server.NewPaymanServer(delegate, fee, networkFee, networkGasLimit, gt, wallet, payer, reporter, nil)
			}
			serv.Serve(cycle)

		} else if cycles != "" {
			intCycles, err := parseCyclesInput(cycles)
			if err != nil {
				log.Fatal(err)
			}
			payouts, ops, err := payer.PayoutForCycles(intCycles[0], intCycles[1], networkFee, networkGasLimit)
			if err != nil {
				log.Fatal(err)
			}

			for _, op := range ops {
				reporter.Log("Successful operation: " + string(op))
				if redditAgent != "" {
					err := redditBot.Post(string(op), cycles)
					if err != nil {
						reporter.Log(fmt.Sprintf("could not post to reddit: %v", err))
					}
				}
			}
			reporter.PrintPaymentsTable(payouts)
			reporter.WriteCSVReport(payouts)

		} else if cycle > 0 {
			payouts, ops, err := payer.PayoutForCycle(cycle, networkFee, networkGasLimit)
			if err != nil {
				log.Fatal(err)
			}
			for _, op := range ops {
				reporter.Log("Successful operation: " + string(op))
				if redditAgent != "" {
					err := redditBot.Post(string(op), strconv.Itoa(cycle))
					if err != nil {
						reporter.Log(fmt.Sprintf("could not post to reddit: %v", err))
					}
				}
			}

			reporter.PrintPaymentsTable(payouts)
			reporter.WriteCSVReport(payouts)

		} else {
			reporter.Log("no cycles passed to payout for.")
			os.Exit(1)
		}
		f.Close()
	},
}

func parseCyclesInput(cycles string) ([2]int, error) {
	arrayCycles := reCycles.FindAllStringSubmatch(cycles, -1)
	if arrayCycles == nil || len(arrayCycles) > 2 {
		return [2]int{}, errors.New("unable to parse cycles flag. Example format 8-12")
	}
	var cycleRange [2]int

	if len(arrayCycles) == 1 {
		cycleRange[0], _ = strconv.Atoi(arrayCycles[0][1])
		cycleRange[1], _ = strconv.Atoi(arrayCycles[0][1])
	} else {
		cycleRange[0], _ = strconv.Atoi(arrayCycles[0][1])
		cycleRange[1], _ = strconv.Atoi(arrayCycles[1][1])
	}

	return cycleRange, nil
}

// Execute payout command
func Execute() {
	if err := payout.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	payout.PersistentFlags().StringVarP(&delegate, "delegate", "d", "", "public key hash of the delegate that's paying out")
	payout.PersistentFlags().StringVarP(&secret, "secret", "s", "", "encrypted secret key of the wallet paying")
	payout.PersistentFlags().StringVarP(&password, "password", "k", "", "password to the secret key of the wallet paying")
	payout.PersistentFlags().BoolVar(&service, "serve", false, "run service to payout for all new cycles going foward")
	payout.PersistentFlags().BoolVar(&dry, "dry", false, "run payout in simulation with report")
	payout.PersistentFlags().StringVar(&cycles, "cycles", "", "cycles to payout for, example 20-24")
	payout.PersistentFlags().IntVarP(&cycle, "cycle", "c", 0, "cycle to payout for, example 20")
	payout.PersistentFlags().StringVarP(&node, "node", "n", "http://127.0.0.1", "example mainnet-node.tzscan.io")
	payout.PersistentFlags().StringVarP(&port, "port", "p", "8732", "example 8732")
	payout.PersistentFlags().Float32VarP(&fee, "fee", "f", 0.05, "example 0.05")
	payout.PersistentFlags().IntVar(&networkFee, "network-fee", 1270, "network fee for each transaction in mutez")
	payout.PersistentFlags().IntVar(&networkGasLimit, "gas-limit", 10200, "network gas limit for each transaction in mutez")
	payout.PersistentFlags().StringVarP(&file, "log-file", "l", "/dev/stdout", "example ./payman.log")
	payout.PersistentFlags().StringVarP(&redditAgent, "reddit", "r", "", "example https://turnage.gitbooks.io/graw/content/chapter1.html")
	payout.PersistentFlags().StringVar(&title, "title", "", "example \"MyService:\"")
}
