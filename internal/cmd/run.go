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

			executeRun(args[0], verbose, table)
		},
	}

	run.PersistentFlags().BoolVarP(&table, "table", "t", false, "formats result into a table (Default: json)")
	run.PersistentFlags().BoolVarP(&verbose, "verbose", "v", true, "will print confirmations in between injections.")

	return run
}

func executeRun(arg string, verbose, table bool) {
	config, err := config.New()
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to load config.")
	}

	cycle, err := strconv.Atoi(arg)
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to parse cycle argument into integer.")
	}

	payout, err := payout.New(config, cycle, true, verbose)
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to initialize payout.")
	}

	rewardsSplit, err := payout.Execute()
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to execute payout.")
	}

	var messengers []notifier.ClientIFace
	if config.Notifications != nil {
		if config.Notifications.Twilio != nil {
			messengers = append(messengers, twilio.New(twilio.Client{
				AccountSID: config.Notifications.Twilio.AccountSID,
				AuthToken:  config.Notifications.Twilio.AuthToken,
				From:       config.Notifications.Twilio.From,
				To:         config.Notifications.Twilio.To,
			}))
		}

		if config.Notifications.Twitter != nil {
			messengers = append(messengers, twitter.NewClient(
				config.Notifications.Twitter.ConsumerKey,
				config.Notifications.Twitter.ConsumerSecret,
				config.Notifications.Twitter.AccessToken,
				config.Notifications.Twitter.AccessSecret,
			))
		}
	}

	payoutNotifier := notifier.NewPayoutNotifier(notifier.PayoutNotifierInput{
		Notifiers: messengers,
	})

	err = payoutNotifier.Notify(fmt.Sprintf("[TZPAY] payout for cycle %d: \n%s\n #tezos #pos", cycle, rewardsSplit.OperationLink))
	if err != nil {
		log.WithField("error", err.Error()).Error("Failed to notify.")
	}

	if table {
		print.Table(cycle, config.Baker.Address, rewardsSplit)
	} else {
		print.JSON(rewardsSplit)
	}
}
