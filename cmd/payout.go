package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	goTezos "github.com/DefinitelyNotAGoat/go-tezos"
	"github.com/DefinitelyNotAGoat/payman/logging"
	"github.com/DefinitelyNotAGoat/payman/payouts"
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
)

var payout = &cobra.Command{
	Use:   "payout",
	Short: "payout pays out rewards to delegations",
	Long: `payman is a simple golang example
				  that demonstartes go-tezos,
				  and also allows you to payout your delegations.`,
	Run: func(cmd *cobra.Command, args []string) {
		log, f := logging.GetLogging(file)

		if delegate == "" {
			log.Info("you need to specify the delegate that is paying out")
			os.Exit(0)
		} else if secret == "" {
			log.Info("you need to set the encrypted secret to the wallet paying")
			os.Exit(0)
		} else if password == "" {
			log.Info("you need to set the password to the wallet paying")
			os.Exit(0)
		}

		gt := goTezos.NewGoTezos()
		gt.AddNewClient(goTezos.NewTezosRPCClient(node, port))
		wallet, err := gt.ImportEncryptedWallet(password, secret)
		if err != nil {
			log.Fatalf("could not import wallet: %v", err)
			os.Exit(1)
		}

		payer := payouts.Payer{}
		if dry {
			payer = payouts.NewPayer(gt, wallet, delegate, fee, false)
		} else {
			payer = payouts.NewPayer(gt, wallet, delegate, fee, true)
		}

		if service == true {
			ticker := time.NewTicker(5 * time.Minute)
			quit := make(chan struct{})
			lastPaidCycle := -1
			for {
				select {
				case <-ticker.C:
					constants, err := gt.GetNetworkConstants()
					if err != nil {
						log.Error(err)
					}
					head, _, err := gt.GetBlockLevelHead()
					if err != nil {
						log.Error(err)
					}
					currentCycle := head / constants.BlocksPerCycle
					if currentCycle == cycle && lastPaidCycle == -1 {
						payouts, ops, err := payer.PayoutForCycle(currentCycle, networkFee, networkGasLimit)
						if err != nil {
							log.Fatal(err)
							close(quit)
						}
						for _, op := range ops {
							log.Info("Successful operation: ", op)
						}
						log.Info(payouts)
						lastPaidCycle = currentCycle
					}
					if (lastPaidCycle + 1) == currentCycle {
						payouts, ops, err := payer.PayoutForCycle(currentCycle, networkFee, networkGasLimit)
						if err != nil {
							log.Fatal(err)
							close(quit)
						}
						for _, op := range ops {
							log.Info("Successful operation: ", op)
						}
						log.Info(payouts)
						lastPaidCycle = currentCycle
					}
				case <-quit:
					ticker.Stop()
					return
				}
			}
		} else if cycles != "" {
			cycles, err := parseCyclesInput(cycles)
			if err != nil {
				log.Fatal(err)
			}
			payouts, ops, err := payer.PayoutForCycles(cycles[0], cycles[1], networkFee, networkGasLimit)
			if err != nil {
				log.Fatal(err)
			}

			for _, op := range ops {
				log.Info("Successful operation: ", op)
			}
			log.Info(payouts)

		} else if cycle > 0 {
			payouts, ops, err := payer.PayoutForCycle(cycle, networkFee, networkGasLimit)
			if err != nil {
				log.Fatal(err)
			}
			for _, op := range ops {
				log.Info("Successful operation: ", op)
			}
			log.Info(payouts)

		} else {
			log.Info("no cycles past to payout for.")
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

// func tablePayments(payments []goTezos.Payment) {
// 	var table [][]string
// 	for _, payment := range payments {
// 		table = append(table, []string{payment.Address, fmt.Sprintf("%f", payment.Amount)})
// 	}
// }

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
}

//Takes an interface v and returns a pretty json string.
func PrettyReport(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		return string(b)
	}
	return ""
}
