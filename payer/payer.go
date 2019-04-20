package payer

import (
	"strconv"

	goTezos "github.com/DefinitelyNotAGoat/go-tezos"
)

type Payer struct {
	gt       *goTezos.GoTezos
	wallet   goTezos.Wallet
	delegate string
	fee      float32
	enabled  bool
}

func NewPayer(gt *goTezos.GoTezos, wallet goTezos.Wallet, delegate string, fee float32, enabled bool) Payer {
	return Payer{gt: gt, wallet: wallet, delegate: delegate, fee: fee, enabled: enabled}
}

func (payer *Payer) PayoutForCycle(cycle int, networkFee int, networkGasLimit int) ([]goTezos.Payment, [][]byte, error) {
	rewards, err := payer.gt.GetRewardsForDelegateCycle(payer.delegate, cycle)
	if err != nil {
		return nil, nil, err
	}
	payments := payer.calcPayments(rewards, payer.fee)

	ops, err := payer.gt.CreateBatchPayment(payments, payer.wallet, networkFee, networkGasLimit)
	if err != nil {
		return payments, nil, err
	}

	responses := [][]byte{}
	if payer.enabled {
		responses, err = payer.gt.InjectOperation(ops)
		if err != nil {
			return payments, responses, err
		}
	}

	return payments, responses, nil
}

func (payer *Payer) PayoutForCycles(cycleStart, cycleEnd int, networkFee int, networkGasLimit int) ([]goTezos.Payment, [][]byte, error) {
	rewards, err := payer.gt.GetRewardsForDelegateForCycles(payer.delegate, cycleStart, cycleEnd)
	if err != nil {
		return nil, nil, err
	}
	payments := payer.calcPayments(rewards, payer.fee)
	ops, err := payer.gt.CreateBatchPayment(payments, payer.wallet, networkFee, networkGasLimit)
	if err != nil {
		return payments, nil, err
	}

	responses := [][]byte{}
	if payer.enabled {
		responses, err = payer.gt.InjectOperation(ops)
		if err != nil {
			return payments, responses, err
		}
	}

	return payments, responses, nil
}

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
