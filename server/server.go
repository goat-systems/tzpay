package server

import (
	"time"

	goTezos "github.com/DefinitelyNotAGoat/go-tezos"
	pay "github.com/DefinitelyNotAGoat/payman/payer"
	"github.com/DefinitelyNotAGoat/payman/reporting"
)

// PaymanServer is structure representing a payman payout server
type PaymanServer struct {
	delegate   string
	fee        float32
	networkFee int
	networkGas int
	gt         *goTezos.GoTezos
	wallet     goTezos.Wallet
	payer      pay.Payer
	reporter   reporting.Reporter
}

// NewPaymanServer contructs a new payman server
func NewPaymanServer(delegate string, fee float32, networkFee int, networkGas int, gt *goTezos.GoTezos, wallet goTezos.Wallet, payer pay.Payer, reporter reporting.Reporter) PaymanServer {
	return PaymanServer{
		delegate:   delegate,
		fee:        fee,
		networkFee: networkFee,
		networkGas: networkGas,
		gt:         gt,
		wallet:     wallet,
		payer:      payer,
		reporter:   reporter,
	}
}

// Serve starts the payman server
func (payman *PaymanServer) Serve(startCylce int) {
	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})
	lastPaidCycle := -1
	for {
		select {
		case <-ticker.C:
			constants, err := payman.gt.GetNetworkConstants()
			if err != nil {
				payman.reporter.Log(err)
			}
			head, _, err := payman.gt.GetBlockLevelHead()
			if err != nil {
				payman.reporter.Log(err)
			}
			currentCycle := head / constants.BlocksPerCycle
			if currentCycle == startCylce && lastPaidCycle == -1 {
				payouts, ops, err := payman.payer.PayoutForCycle(currentCycle, payman.networkFee, payman.networkGas)
				if err != nil {
					payman.reporter.Log(err)
					close(quit)
				}
				for _, op := range ops {
					payman.reporter.Log("Successful operation: " + string(op))
				}
				payman.reporter.PrintPaymentsTable(payouts)
				payman.reporter.WriteCSVReport(payouts)
				lastPaidCycle = currentCycle
			}
			if (lastPaidCycle + 1) == currentCycle {
				payouts, ops, err := payman.payer.PayoutForCycle(currentCycle, payman.networkFee, payman.networkGas)
				if err != nil {
					payman.reporter.Log(err)
					close(quit)
				}
				for _, op := range ops {
					payman.reporter.Log("Successful operation: " + string(op))
				}
				payman.reporter.PrintPaymentsTable(payouts)
				payman.reporter.WriteCSVReport(payouts)
				lastPaidCycle = currentCycle
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
