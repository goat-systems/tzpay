package server

import (
	"fmt"
	"time"

	goTezos "github.com/DefinitelyNotAGoat/go-tezos"
	"github.com/DefinitelyNotAGoat/payman/options"
	pay "github.com/DefinitelyNotAGoat/payman/payer"
	"github.com/DefinitelyNotAGoat/payman/reddit"
	"github.com/DefinitelyNotAGoat/payman/reporting"
)

// PayoutServer is structure representing a payout server
type PayoutServer struct {
	gt       *goTezos.GoTezos
	wallet   goTezos.Wallet
	reporter reporting.Reporter
	rbot     *reddit.Bot
	conf     *options.Options
}

// NewPayoutServer contructs a new payout server
func NewPayoutServer(gt *goTezos.GoTezos, wallet goTezos.Wallet, reporter reporting.Reporter, rbot *reddit.Bot, conf *options.Options) PayoutServer {
	return PayoutServer{
		gt:       gt,
		wallet:   wallet,
		reporter: reporter,
		rbot:     rbot,
		conf:     conf,
	}
}

// Serve starts the payout server
func (ps *PayoutServer) Serve() {
	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})
	lastCycle, _, err := ps.gt.GetBlockLevelHead()
	if err != nil {
		ps.reporter.Log(err)
	}

	payer := pay.NewPayer(ps.gt, ps.wallet, ps.conf)
	constants, err := ps.gt.GetNetworkConstants()
	if err != nil {
		ps.reporter.Log(err)
	}
	lastCycle = lastCycle / constants.BlocksPerCycle

	for {
		select {
		case <-ticker.C:
			head, _, err := ps.gt.GetBlockLevelHead()
			if err != nil {
				ps.reporter.Log(err)
			}
			currentCycle := head / constants.BlocksPerCycle
			if currentCycle > lastCycle {
				ps.conf.Cycle = currentCycle
				payouts, ops, err := payer.Payout()
				if err != nil {
					ps.reporter.Log(err)
					close(quit)
				}
				for _, op := range ops {
					ps.reporter.Log("Successful operation: " + string(op))
					if ps.rbot != nil {
						err := ps.rbot.Post(string(op), string(currentCycle))
						if err != nil {
							ps.reporter.Log(fmt.Sprintf("could not post to reddit: %v", err))
						}
					}
				}
				ps.reporter.PrintPaymentsTable(payouts)
				ps.reporter.WriteCSVReport(payouts)
				lastCycle = currentCycle
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
