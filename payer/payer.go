package payer

import (
	"fmt"
	"strconv"

	goTezos "github.com/DefinitelyNotAGoat/go-tezos"
	"github.com/DefinitelyNotAGoat/payman/options"
)

// Payer is a structure to represent pay operations
type Payer struct {
	gt     *goTezos.GoTezos
	wallet goTezos.Wallet
	conf   *options.Options
}

// NewPayer returns is a contructor for Payer
func NewPayer(gt *goTezos.GoTezos, wallet goTezos.Wallet, conf *options.Options) Payer {
	return Payer{gt: gt, wallet: wallet, conf: conf}
}

// Payout pays out based of the configuration of the payer that calls it
func (payer *Payer) Payout() ([]goTezos.Payment, [][]byte, error) {
	var payment []goTezos.Payment
	var ops [][]byte

	if payer.conf.Cycle != 0 {
		payment, ops, err := payer.payoutForCycle()
		if err != nil {
			return payment, ops, err
		}
		return payment, ops, err
	} else if payer.conf.Cycles != "" {
		payment, ops, err := payer.payoutForCycles()
		if err != nil {
			return payment, ops, err
		}
		return payment, ops, err
	}
	return payment, ops, fmt.Errorf("no cycle configuration found to payout for")
}

// payoutForCycle uses the payer that calls to payout for the cycle passed with the network fee and gas limit specified
func (payer *Payer) payoutForCycle() ([]goTezos.Payment, [][]byte, error) {

	rewards, err := payer.gt.GetRewardsForDelegateCycle(payer.conf.Delegate, payer.conf.Cycle)
	if err != nil {
		return nil, nil, err
	}
	payments := payer.calcPayments(rewards, payer.conf.Fee)

	ops, err := payer.gt.CreateBatchPayment(payments, payer.wallet, payer.conf.NetworkFee, payer.conf.NetworkGasLimit)
	if err != nil {
		return payments, nil, err
	}

	responses := [][]byte{}
	if !payer.conf.Dry {
		for _, op := range ops {
			resp, err := payer.gt.InjectOperation(op)
			if err != nil {
				return payments, responses, err
			}
			responses = append(responses, resp)
		}
	}

	return payments, responses, nil
}

// payoutForCycles uses the payer that calls to payout for the cycles passed with the network fee and gas limit specified
func (payer *Payer) payoutForCycles() ([]goTezos.Payment, [][]byte, error) {
	cycles, err := payer.conf.ParseCyclesInput()
	rewards, err := payer.gt.GetRewardsForDelegateForCycles(payer.conf.Delegate, cycles[0], cycles[1])
	if err != nil {
		return nil, nil, err
	}
	payments := payer.calcPayments(rewards, payer.conf.Fee)
	ops, err := payer.gt.CreateBatchPayment(payments, payer.wallet, payer.conf.NetworkFee, payer.conf.NetworkGasLimit)
	if err != nil {
		return payments, nil, err
	}

	responses := [][]byte{}
	if !payer.conf.Dry {
		for _, op := range ops {
			resp, err := payer.gt.InjectOperation(op)
			if err != nil {
				return payments, responses, err
			}
			responses = append(responses, resp)
		}
	}

	return payments, responses, nil
}

// calcPayments iterates through the goTezos type DelegationServiceRewards, to form fill out the payment structure used
// for batch payments
func (payer *Payer) calcPayments(rewards goTezos.DelegationServiceRewards, fee float32) []goTezos.Payment {
	payments := []goTezos.Payment{}
	net := 1 - fee
	for _, cycle := range rewards.RewardsByCycle {
		for _, delegate := range cycle.Delegations {
			f, _ := strconv.ParseFloat(delegate.GrossRewards, 32)
			amount := f * float64(net)
			payment := goTezos.Payment{Address: delegate.DelegationPhk, Amount: amount}
			if amount > 1500 {
				payments = append(payments, payment)
			}
		}
	}
	return payments
}
