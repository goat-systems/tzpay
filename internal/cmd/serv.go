package cmd

import (
	"fmt"
	"time"

	"github.com/goat-systems/go-tezos/v3/rpc"
	"github.com/goat-systems/tzpay/v2/internal/config"
	"github.com/goat-systems/tzpay/v2/internal/notifier"
	"github.com/goat-systems/tzpay/v2/internal/notifier/twilio"
	"github.com/goat-systems/tzpay/v2/internal/notifier/twitter"
	"github.com/goat-systems/tzpay/v2/internal/payout"
	"github.com/goat-systems/tzpay/v2/internal/print"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ServCommand returns a new run cobra command
func ServCommand() *cobra.Command {
	var verbose bool

	var serv = &cobra.Command{
		Use:     "serv",
		Short:   "serv runs a service that will continously payout cycle by cycle",
		Example: `tzpay serv`,
		Run: func(cmd *cobra.Command, args []string) {
			start(verbose)
		},
	}

	serv.PersistentFlags().BoolVarP(&verbose, "verbose", "v", true, "will print confirmations in between injections.")
	return serv
}

func start(verbose bool) {
	log.Info("Starting tzpay payout server.")

	config, err := config.New()
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to load config.")
	}

	rpc, err := rpc.New(config.API.Tezos)
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to connect to tezos RPC.")
	}

	var cycle int
	ticker := time.NewTicker(time.Minute)
	block, err := rpc.Head()
	if err != nil {
		log.WithField("error", err.Error()).Error("Failed to parse get current cycle.")
	}
	cycle = block.Metadata.Level.Cycle

	for range ticker.C {
		block, err := rpc.Head()
		if err != nil {
			log.WithField("error", err.Error()).Error("Failed to parse get current cycle.")
		}

		if block.Metadata.Level.Cycle > cycle {
			executeServ(config, cycle, verbose)
			log.WithField("cycle", cycle).Infof("Executed payout for cycle '%s'", cycle)

			cycle = block.Metadata.Level.Cycle
			log.WithField("cycle", cycle).Info("Update to current cycle.")
		}
	}

}

func executeServ(config config.Config, cycle int, verbose bool) {
	payout, err := payout.New(config, cycle, true, verbose)
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to initialize payout.")
	}

	rewardsSplit, err := payout.Execute()
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to execute payout.")
	}

	var messengers []notifier.ClientIFace
	if config.Twilio != nil {
		messengers = append(messengers, twilio.NewTwilioClient(twilio.Client{
			AccountSID: config.Twilio.AccountSID,
			AuthToken:  config.Twilio.AuthToken,
			From:       config.Twilio.From,
			To:         config.Twilio.To,
		}))
	}

	if config.Twitter != nil {
		messengers = append(messengers, twitter.NewClient(
			config.Twitter.ConsumerKey,
			config.Twitter.ConsumerSecret,
			config.Twitter.AccessToken,
			config.Twitter.AccessSecret,
		))
	}

	payoutNotifier := notifier.NewPayoutNotifier(notifier.PayoutNotifierInput{
		Notifiers: messengers,
	})

	err = payoutNotifier.Notify(fmt.Sprintf("[TZPAY] payout for cycle %d: \n%s\n #tezos #pos", cycle, rewardsSplit.OperationLink))
	if err != nil {
		log.WithField("error", err.Error()).Error("Failed to notify.")
	}

	print.JSON(rewardsSplit)
}
