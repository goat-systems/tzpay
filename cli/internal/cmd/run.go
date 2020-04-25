package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v2/cli/internal/baker"
	"github.com/goat-systems/tzpay/v2/cli/internal/db/model"
	"github.com/goat-systems/tzpay/v2/cli/internal/enviroment"
	"github.com/goat-systems/tzpay/v2/cli/internal/print"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	confirmationDurationInterval = time.Second * 1
	confirmationTimoutInterval   = time.Minute * 2
	confirm                      = true
)

// NewRunCommand returns a new run cobra command
func NewRunCommand() *cobra.Command {
	var table bool
	var batchSize int

	var report = &cobra.Command{
		Use:     "run",
		Short:   "run executes a batch payout",
		Long:    "run executes a batch payout and prints the result in json or a table",
		Example: `tzpay run <cycle>`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.WithFields(nil).Fatal("Missing cycle as argument.")
			}
			runner, err := newRunner(newRunnerInput{
				cycle:     args[0],
				table:     table,
				batchSize: batchSize,
			})
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Fatal("Failed to run payout.")
			}
			payout, err := runner.run()
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Fatal("Failed to run payout.")
			}

			runner.print(payout)
		},
	}

	report.PersistentFlags().BoolVarP(&table, "table", "t", false, "formats result into a table (Default: json)")
	report.PersistentFlags().IntVarP(&batchSize, "batch-size", "b", 125, "changes the size of the payout batches (too large may result in failure).")

	return report
}

type runner struct {
	ctx       context.Context
	base      *enviroment.ContextEnviroment
	cycle     int
	table     bool
	batchSize int
}

type newRunnerInput struct {
	cycle     string
	table     bool
	batchSize int
	gt        gotezos.IFace // only pass for testing
}

func newRunner(input newRunnerInput) (*runner, error) {
	cycle, err := strconv.Atoi(input.cycle)
	if err != nil {
		return nil, errors.New("failed to read cycle argument")
	}

	ctx, err := enviroment.InitContext(input.gt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize enviroment")
	}

	return &runner{
		ctx:       ctx,
		cycle:     cycle,
		table:     input.table,
		base:      enviroment.GetEnviromentFromContext(ctx),
		batchSize: input.batchSize,
	}, nil
}

func (r *runner) run() (*model.Payout, error) {
	baker := baker.NewBaker(r.base.GoTezos)

	payout, err := baker.Payouts(r.ctx, r.cycle)
	if err != nil {
		return nil, err
	}

	payouts := splitPayouts(*payout, r.batchSize)

	var operations []string
	for i, p := range payouts {
		forge, lastCounter, err := baker.ForgePayout(r.ctx, *p)
		if err != nil {
			return nil, err
		}

		signedop, err := r.base.Wallet.SignOperation(forge)
		if err != nil {
			return nil, err
		}

		op, err := r.base.GoTezos.InjectionOperation(&gotezos.InjectionOperationInput{
			Operation: &signedop,
		})
		if err != nil {
			return nil, err
		}

		if confirm {
			log.WithFields(log.Fields{
				"Injection": fmt.Sprintf("%d/%d", (i + 1), len(payouts)),
			}).Info("Confirming Injection.")

			ok := r.ConfirmInjection(lastCounter)
			if !ok {
				return p, errors.New("failed to confirm injection")
			}
		}

		operations = append(operations, string(*op))
	}

	payout.SetOperations(operations...)

	return payout, nil
}

func (r *runner) ConfirmInjection(lastCounter int) bool {
	timer := time.After(confirmationTimoutInterval)
	ticker := time.Tick(confirmationDurationInterval)
	for {
		select {
		case <-ticker:
			if head, err := r.base.GoTezos.Head(); err == nil {
				if counter, err := r.base.GoTezos.Counter(head.Hash, r.base.Wallet.Address); err == nil {
					if *counter == lastCounter {
						return true
					}
				}
			}
		case <-timer:
			return false
		}
	}
}

// func (r *runner) save(payout *model.Payout) error {
// 	err := r.base.BoltDB.SavePayout(*payout)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to save payout in tzpay.db")
// 	}
// 	return nil
// }

func (r *runner) print(payout *model.Payout) {
	if r.table {
		print.Table(r.ctx, payout)
	} else {
		print.JSON(payout)
	}
}

func splitPayouts(payout model.Payout, split int) []*model.Payout {
	var payouts []*model.Payout
	if len(payout.DelegationEarnings) <= split {
		payouts = append(payouts, &payout)
		return payouts
	}
	for len(payout.DelegationEarnings) >= split {
		p := &model.Payout{
			CycleHash:          payout.CycleHash,
			FrozenBalance:      payout.FrozenBalance,
			StakingBalance:     payout.StakingBalance,
			DelegateEarnings:   payout.DelegateEarnings,
			DelegationEarnings: payout.DelegationEarnings[:split],
		}
		payouts = append(payouts, p)
		payout.DelegationEarnings = payout.DelegationEarnings[split:]
	}

	return payouts
}
