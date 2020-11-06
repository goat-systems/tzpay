package cmd

import (
	"time"

	"github.com/goat-systems/go-tezos/v3/rpc"
	"github.com/goat-systems/tzpay/v3/internal/config"
	"github.com/goat-systems/tzpay/v3/internal/payout"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type server struct {
	queue     payout.Queue
	rpcClient rpc.IFace
	cfg       config.Config
	runner    Run
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

	runner := NewRun(false, verbose)
	queue := payout.NewQueue(&runner.notifier)
	queue.Start()

	return server{
		queue:     queue,
		rpcClient: rpc,
		cfg:       config,
		runner:    runner,
	}, nil
}

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
			server.start()
		},
	}

	serv.PersistentFlags().BoolVarP(&verbose, "verbose", "v", true, "will print confirmations in between injections.")
	return serv
}

func (s *server) start() {
	quit := make(chan struct{})
	log.Info("Starting tzpay payout server.")

	block, err := s.rpcClient.Head()
	if err != nil {
		log.WithField("error", err.Error()).Error("Server failed to get current cycle.")
	}

	go func() {
		currentCycle := block.Metadata.Level.Cycle
		log.Infof("Current cycle: %d.", currentCycle)
		ticker := time.NewTicker(time.Second * 30)
		for range ticker.C {
			b, err := s.rpcClient.Head()
			if err != nil {
				log.WithField("error", err.Error()).Error("Server failed to get current cycle.")
				continue
			}
			log.WithField("level", b.Header.Level).Debug("Found a new block.")

			if currentCycle < b.Metadata.Level.Cycle {
				log.WithField("cycle", currentCycle).Info("New current cycle found.", b.Metadata.Level.Cycle)
				payout, err := payout.New(s.runner.config, currentCycle, true, s.runner.verbose)
				if err != nil {
					log.WithField("error", err.Error()).Error("Failed to intialize payout.")
					continue
				}
				log.WithField("cycle", currentCycle).Info("Adding payout to queue.")
				s.queue.Enqueue(*payout)
				currentCycle = block.Metadata.Level.Cycle
			}
		}
	}()

	<-quit
}
