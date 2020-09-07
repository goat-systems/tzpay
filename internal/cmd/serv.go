package cmd

import (
	"time"

	gotezos "github.com/goat-systems/go-tezos/v3"
	"github.com/goat-systems/tzpay/v2/internal/enviroment"
	"github.com/goat-systems/tzpay/v2/internal/payout"
	"github.com/goat-systems/tzpay/v2/internal/print"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Serv configures and exposes functions to allow tzpay inject a payout into the tezos network on a cycle by cycle basis.
type Serv struct {
	gt             gotezos.IFace
	bakersFee      float64
	delegate       string
	gasLimit       int
	minimumPayment int
	networkFee     int
	blackList      []string
	wallet         gotezos.Wallet
	batchSize      int
	verbose        bool
	table          bool
}

// ServInput is the input for NewServ
type ServInput struct {
	GoTezos        gotezos.IFace
	BakersFee      float64
	Delegate       string
	GasLimit       int
	MinimumPayment int
	NetworkFee     int
	BlackList      []string
	Wallet         gotezos.Wallet
	BatchSize      int
	Verbose        bool
	Table          bool
}

// NewServ returns a pointer to a Serv
func NewServ(input ServInput) *Serv {
	return &Serv{
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

// ServCommand returns a new run cobra command
func ServCommand() *cobra.Command {
	var table bool
	var batchSize int
	var verbose bool

	var serv = &cobra.Command{
		Use:     "serv",
		Short:   "serv runs a service that will continously payout cycle by cycle",
		Example: `tzpay serv`,
		Run: func(cmd *cobra.Command, args []string) {
			env, err := enviroment.NewRunEnviroment()
			if err != nil {
				log.WithField("error", err.Error()).Fatal("Failed to load enviroment.")
			}

			NewServ(ServInput{
				GoTezos:        env.GoTezos,
				BakersFee:      env.BakersFee,
				Delegate:       env.Delegate,
				GasLimit:       env.GasLimit,
				MinimumPayment: env.MinimumPayment,
				NetworkFee:     env.NetworkFee,
				BlackList:      env.BlackList,
				Wallet:         env.Wallet,
				BatchSize:      batchSize,
				Verbose:        verbose,
				Table:          table,
			}).Start()
		},
	}

	serv.PersistentFlags().BoolVarP(&table, "table", "t", false, "formats result into a table (Default: json)")
	serv.PersistentFlags().IntVarP(&batchSize, "batch-size", "b", 125, "changes the size of the payout batches (too large may result in failure).")
	serv.PersistentFlags().BoolVarP(&verbose, "verbose", "v", true, "will print confirmations in between injections.")
	return serv
}

// Start will start a payout server where all future cycles will be paid out automatically assuming the payout wallet is funded.
func (s *Serv) Start() {
	var cycle int

	ticker := time.NewTicker(time.Minute)
	block, err := s.gt.Head()
	if err != nil {
		log.WithField("error", err.Error()).Error("Failed to parse get current cycle.")
	}
	cycle = block.Metadata.Level.Cycle

	log.WithField("cycle", cycle).Info("Starting tzpay payout server.")

	for range ticker.C {
		block, err := s.gt.Head()
		if err != nil {
			log.WithField("error", err.Error()).Error("Failed to parse get current cycle.")
		}

		if block.Metadata.Level.Cycle > cycle {
			s.execute(cycle, s.batchSize, s.verbose, s.table)
			log.WithField("cycle", cycle).Info("tzpay executed a payout.")

			cycle = block.Metadata.Level.Cycle
			log.WithField("cycle", cycle).Info("Update to current cycle.")
		}
	}

}

func (s *Serv) execute(cycle int, batchSize int, verbose, table bool) {
	report, err := payout.NewPayout(payout.NewPayoutInput{
		GoTezos:    s.gt,
		Cycle:      cycle,
		Delegate:   s.delegate,
		BakerFee:   s.bakersFee,
		MinPayment: s.minimumPayment,
		BlackList:  s.blackList,
		BatchSize:  batchSize,
		NetworkFee: s.networkFee,
		GasLimit:   s.gasLimit,
		Inject:     true,
		Verbose:    verbose,
		Wallet:     s.wallet,
	}).Execute()

	if err != nil {
		log.WithField("error", err.Error()).Errorf("Failed to execute payout for cycle: %d", cycle)
	}

	if table {
		print.Table(cycle, s.delegate, report)
	} else {
		print.JSON(report)
	}
}
