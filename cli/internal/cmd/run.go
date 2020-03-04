package cmd

import (
	"fmt"
	"strconv"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v2/cli/internal/baker"
	"github.com/goat-systems/tzpay/v2/cli/internal/enviroment"
	"github.com/goat-systems/tzpay/v2/cli/internal/print"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewRunCommand returns a new run cobra command
func NewRunCommand() *cobra.Command {
	var table bool

	var report = &cobra.Command{
		Use:     "run",
		Short:   "run executes a batch payout",
		Long:    "run executes a batch payout and prints the result in json or a table",
		Example: `tzpay run <cycle>`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.WithFields(nil).Fatal("Missing cycle as argument.")
			}
			run(args[0], table)
		},
	}

	report.PersistentFlags().BoolVarP(&table, "table", "t", false, "formats result into a table (Default: json)")

	return report
}

func run(arg string, table bool) {
	cycle, err := strconv.Atoi(arg)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed to read cycle argument.")
	}

	ctx, err := enviroment.InitContext()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed to get enviroment and initialize context.")
	}

	base := enviroment.GetEnviromentFromContext(ctx)

	gt, err := gotezos.New(base.HostNode)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed to connect to node.")
	}
	baker := baker.NewBaker(gt)

	payouts, err := baker.Payouts(ctx, cycle)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Failed to get payouts.")
	}

	split := splitPayouts(payouts)

	var operations []string
	for _, p := range split {
		forge, err := baker.ForgePayout(ctx, *p)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to forge operation.")
		}

		fmt.Println(forge)

		signedop, err := base.Wallet.SignOperation(forge)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to sign operation.")
		}

		op, err := gt.InjectionOperation(&gotezos.InjectionOperationInput{
			Operation: &signedop,
		})
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to inject operation.")
		}

		operations = append(operations, string(*op))
	}

	if table {
		print.Table(ctx, payouts, operations...)
	} else {
		print.JSON(payouts, operations...)
	}
}

func splitPayouts(payout *baker.Payout) []*baker.Payout {
	var payouts []*baker.Payout
	size := 2
	for len(payout.DelegationEarnings) > size {
		p := &baker.Payout{
			Cycle:              payout.Cycle,
			FrozenBalance:      payout.FrozenBalance,
			StakingBalance:     payout.StakingBalance,
			Delegate:           payout.Delegate,
			DelegationEarnings: payout.DelegationEarnings[:size],
		}
		payouts = append(payouts, p)
		payout.DelegationEarnings = payout.DelegationEarnings[size:]
	}

	return payouts
}
