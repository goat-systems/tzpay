package payer

import (
	goTezos "github.com/DefinitelyNotAGoat/go-tezos"
	"github.com/DefinitelyNotAGoat/payman/options"
)

// Payer is a structure to represent pay operations
type Payer struct {
	gt     *goTezos.GoTezos
	wallet goTezos.Wallet
	conf   *options.Options
}

// PayoutResults is a helper structure to describe results of a payout
type PayoutResults struct {
	OpHashes []string
	Payouts  []Payout
}

// Payout describes a single payout to a single address
type Payout struct {
	Address  string
	Share    float64
	Gross    float64
	Fee      float64
	Total    float64
	USDValue float64
}

// Node describes the node's total in PayoutResults
type Node struct {
	Address           string
	TotalFees         float64
	SelfBaked         float64
	TotalFeesUSD      float64
	TotalSelfBakedUSD float64
}

// NewPayer returns is a contructor for Payer
func NewPayer(gt *goTezos.GoTezos, wallet goTezos.Wallet, conf *options.Options) Payer {
	return Payer{gt: gt, wallet: wallet, conf: conf}
}

// Payout uses the payers configuration that calls it, to pay out for the cycle in the conf
func (payer *Payer) Payout() (goTezos.DelegateReport, [][]byte, error) {
	var payments []goTezos.Payment
	rewards := &goTezos.DelegateReport{}

	if len(payer.conf.PaymentsOverride.Payments) > 0 {
		payments = payer.conf.PaymentsOverride.Payments
	} else {
		var err error
		rewards, err = payer.gt.Delegate.GetReport(payer.conf.Delegate, payer.conf.Cycle, float64(payer.conf.Fee))
		if err != nil {
			return *rewards, nil, err
		}

		payments = rewards.GetPayments()
	}

	responses := [][]byte{}
	if !payer.conf.Dry {
		ops, err := payer.gt.Operation.CreateBatchPayment(payments, payer.wallet, payer.conf.NetworkFee, payer.conf.NetworkGasLimit)
		if err != nil {
			return *rewards, nil, err
		}

		for _, op := range ops {
			resp, err := payer.gt.Operation.InjectOperation(op)
			if err != nil {
				return *rewards, responses, err
			}
			responses = append(responses, resp)
		}
	}

	return *rewards, responses, nil
}
