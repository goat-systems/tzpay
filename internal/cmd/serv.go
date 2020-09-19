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
	"github.com/pkg/errors"
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
			server, err := newServer(verbose)
			if err != nil {
				log.WithField("error", err.Error()).Fatal("Failed to initialize server.")
			}
			quit := make(chan struct{})
			server.start()
			<-quit
		},
	}

	serv.PersistentFlags().BoolVarP(&verbose, "verbose", "v", true, "will print confirmations in between injections.")
	return serv
}

type server struct {
	rpcClient rpc.IFace
	cfg       config.Config
	verbose   bool
	notifier  notifier.PayoutNotifier
}

func newServer(verbose bool) (server, error) {
	config, err := config.New()
	if err != nil {
		return server{}, errors.Wrap(err, "failed to load configuration")
	}

	rpc, err := rpc.New(config.API.Tezos)
	if err != nil {
		return server{}, errors.Wrap(err, "failed to connect to tezos rpc")
	}

	var payoutMessengers []notifier.ClientIFace
	if config.Notifications != nil {
		if config.Notifications.Twilio != nil {
			payoutMessengers = append(payoutMessengers, twilio.New(twilio.Client{
				AccountSID: config.Notifications.Twilio.AccountSID,
				AuthToken:  config.Notifications.Twilio.AuthToken,
				From:       config.Notifications.Twilio.From,
				To:         config.Notifications.Twilio.To,
			}))
		}

		if config.Notifications.Twitter != nil {
			payoutMessengers = append(payoutMessengers, twitter.NewClient(
				config.Notifications.Twitter.ConsumerKey,
				config.Notifications.Twitter.ConsumerSecret,
				config.Notifications.Twitter.AccessToken,
				config.Notifications.Twitter.AccessSecret,
			))
		}
	}

	return server{
		rpcClient: rpc,
		cfg:       config,
		verbose:   verbose,
		notifier: notifier.NewPayoutNotifier(notifier.PayoutNotifierInput{
			Notifiers: payoutMessengers,
		}),
	}, nil
}

func (s *server) start() {
	log.Info("Starting tzpay payout server.")
	s.executePayouts(s.watchForCycle())
}

func (s *server) watchForCycle() chan int {
	block, err := s.rpcClient.Head()
	if err != nil {
		log.WithField("error", err.Error()).Error("Server failed to get current cycle.")
	}

	cycleChan := make(chan int, 2)

	go func() {
		currentCycle := block.Metadata.Level.Cycle
		log.Infof("Current cycle: %d.", currentCycle)
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			block, err := s.rpcClient.Head()
			if err != nil {
				log.WithField("error", err.Error()).Error("Server failed to get current cycle.")
			}

			if currentCycle < block.Metadata.Level.Cycle {
				log.Infof("Current cycle: %d.", block.Metadata.Level.Cycle)
				cycleChan <- block.Metadata.Level.Cycle
				currentCycle = block.Metadata.Level.Cycle
			}
		}
	}()

	return cycleChan
}

func (s *server) executePayouts(cycleChan chan int) {
	go func() {
		for cycle := range cycleChan {
			payout, err := payout.New(s.cfg, cycle, true, s.verbose)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
					"cycle": cycle,
				}).Error("Server failed to initalize payout.")
				continue
			}
			rewardsSplit, err := payout.Execute()
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
					"cycle": cycle,
				}).Error("Server failed to execute payout.")
				continue
			}

			err = s.notifier.Notify(fmt.Sprintf("#tezos payout for cycle %d: \n\n%s", cycle, rewardsSplit.OperationLink))
			if err != nil {
				log.WithField("error", err.Error()).Error("Failed to notify.")
			}
		}
	}()
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

	var payoutMessengers []notifier.ClientIFace

	if config.Notifications != nil {
		if config.Notifications.Twilio != nil {
			twilioClient := twilio.New(twilio.Client{
				AccountSID: config.Notifications.Twilio.AccountSID,
				AuthToken:  config.Notifications.Twilio.AuthToken,
				From:       config.Notifications.Twilio.From,
				To:         config.Notifications.Twilio.To,
			})
			payoutMessengers = append(payoutMessengers, twilioClient)

		}

		if config.Notifications.Twitter != nil {
			payoutMessengers = append(payoutMessengers, twitter.NewClient(
				config.Notifications.Twitter.ConsumerKey,
				config.Notifications.Twitter.ConsumerSecret,
				config.Notifications.Twitter.AccessToken,
				config.Notifications.Twitter.AccessSecret,
			))
		}
	}

	payoutNotifier := notifier.NewPayoutNotifier(notifier.PayoutNotifierInput{
		Notifiers: payoutMessengers,
	})

	err = payoutNotifier.Notify(fmt.Sprintf("[TZPAY] payout for cycle %d: \n%s\n #tezos #pos", cycle, rewardsSplit.OperationLink))
	if err != nil {
		log.WithField("error", err.Error()).Error("Failed to notify.")
	}

	print.JSON(rewardsSplit)
}
