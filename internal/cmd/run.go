package cmd

import (
	"fmt"
	"strconv"

	"github.com/goat-systems/tzpay/v2/internal/config"
	"github.com/goat-systems/tzpay/v2/internal/notifier"
	"github.com/goat-systems/tzpay/v2/internal/notifier/twilio"
	"github.com/goat-systems/tzpay/v2/internal/notifier/twitter"
	"github.com/goat-systems/tzpay/v2/internal/payout"
	"github.com/goat-systems/tzpay/v2/internal/print"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Run -
type Run struct {
	config   config.Config
	table    bool
	verbose  bool
	notifier notifier.PayoutNotifier
}

// NewRun returns a new Run
func NewRun(table bool, verbose bool) Run {
	config, err := config.New()
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to load config.")
	}

	var messengers []notifier.ClientIFace
	if config.Notifications.Twilio.AccountSID != "" && config.Notifications.Twilio.AuthToken != "" &&
		config.Notifications.Twilio.From != "" && config.Notifications.Twilio.To != nil {
		messengers = append(messengers, twilio.New(twilio.Client{
			AccountSID: config.Notifications.Twilio.AccountSID,
			AuthToken:  config.Notifications.Twilio.AuthToken,
			From:       config.Notifications.Twilio.From,
			To:         config.Notifications.Twilio.To,
		}))
	}

	if config.Notifications.Twitter.ConsumerKey != "" && config.Notifications.Twitter.ConsumerSecret != "" && config.Notifications.Twitter.AccessToken != "" && config.Notifications.Twitter.AccessSecret != "" {
		messengers = append(messengers, twitter.NewClient(
			config.Notifications.Twitter.ConsumerKey,
			config.Notifications.Twitter.ConsumerSecret,
			config.Notifications.Twitter.AccessToken,
			config.Notifications.Twitter.AccessSecret,
		))
	}

	return Run{
		config:  config,
		table:   table,
		verbose: verbose,
		notifier: notifier.NewPayoutNotifier(notifier.PayoutNotifierInput{
			Notifiers: messengers,
		}),
	}
}

// RunCommand returns a new run cobra command
func RunCommand() *cobra.Command {
	var table bool
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

			cycle, err := strconv.Atoi(args[0])
			if err != nil {
				log.WithField("error", err.Error()).Fatal("Failed to parse cycle argument into integer.")
			}

			run := NewRun(table, verbose)
			run.execute(cycle)
		},
	}

	run.PersistentFlags().BoolVarP(&table, "table", "t", false, "formats result into a table (Default: json)")
	run.PersistentFlags().BoolVarP(&verbose, "verbose", "v", true, "will print confirmations in between injections.")

	return run
}

func (r *Run) execute(cycle int) {
	payout, err := payout.New(r.config, cycle, true, r.verbose)
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to intialize payout.")
	}

	rewardsSplit, err := payout.Execute()
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to execute payout.")
	}

	err = r.notifier.Notify(fmt.Sprintf("[TZPAY] payout for cycle %d: \n%s\n #tezos #blockchain", cycle, rewardsSplit.OperationLink))
	if err != nil {
		log.WithField("error", err.Error()).Error("Failed to notify.")
	}

	if r.table {
		print.Table(cycle, r.config.Baker.Address, rewardsSplit)
	} else {
		err := print.JSON(rewardsSplit)
		if err != nil {
			log.WithField("error", err.Error()).Fatal("Failed to print JSON report.")
		}
	}
}
