package payout

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/goat-systems/go-tezos/v3/forge"
	"github.com/goat-systems/go-tezos/v3/keys"
	"github.com/goat-systems/go-tezos/v3/rpc"
	"github.com/goat-systems/tzpay/v2/internal/config"
	"github.com/goat-systems/tzpay/v2/internal/tzkt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	confirmationDurationInterval = time.Second * 1
	confirmationTimoutInterval   = time.Minute * 2
)

// Payout represents a payout and payout operations.
type Payout struct {
	config  config.Config
	rpc     rpc.IFace
	tzkt    tzkt.IFace
	key     keys.Key
	cycle   int
	inject  bool
	verbose bool
}

// New returns a pointer to a new Baker
func New(config config.Config, cycle int, inject, verbose bool) (*Payout, error) {
	payout := &Payout{
		config:  config,
		tzkt:    tzkt.NewTZKT(config.API.TZKT),
		cycle:   cycle,
		inject:  inject,
		verbose: verbose,
	}

	var err error
	payout.rpc, err = rpc.New(config.API.Tezos)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize tezos rpc client")
	}

	if inject {
		payout.key, err = keys.NewKey(keys.NewKeyInput{
			Kind:     keys.Ed25519,
			Esk:      config.Key.Esk,
			Password: config.Key.Password,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to initialize import key")
		}

		config.Key.Esk = ""
		config.Key.Password = ""
	}

	return payout, nil
}

func (p *Payout) constructPayout() (tzkt.RewardsSplit, error) {
	rewardsSplit, err := p.tzkt.GetRewardsSplit(p.config.Baker.Address, p.cycle)
	if err != nil {
		return rewardsSplit, errors.Wrap(err, "failed to contruct payout")
	}

	totalRewards := p.calculateTotals(rewardsSplit)

	bakerBalance, err := p.rpc.Balance(rpc.BalanceInput{
		Cycle:   p.cycle,
		Address: p.config.Baker.Address,
	})
	if err != nil {
		return rewardsSplit, errors.Wrap(err, "failed to contruct payout")
	}

	rewardsSplit.BakerShare = float64(bakerBalance) / float64(rewardsSplit.StakingBalance)
	rewardsSplit.BakerRewards = int(rewardsSplit.BakerShare * float64(totalRewards))

	for i := range rewardsSplit.Delegators {
		if p.isInBlacklist(rewardsSplit.Delegators[i].Address) {
			rewardsSplit.Delegators[i].BlackListed = true
		}

		rewardsSplit.Delegators[i].Share = float64(rewardsSplit.Delegators[i].Balance) / float64(rewardsSplit.StakingBalance)
		if p.config.Baker.EarningsOnly {
			rewardsSplit.Delegators[i].GrossRewards = int(rewardsSplit.Delegators[i].Share * float64(totalRewards))
		} else {
			rewardsSplit.Delegators[i].GrossRewards = int(rewardsSplit.Delegators[i].Share * float64(totalRewards))
		}
		rewardsSplit.Delegators[i].Fee = int(float64(rewardsSplit.Delegators[i].GrossRewards) * p.config.Baker.Fee)
		rewardsSplit.BakerCollectedFees += rewardsSplit.Delegators[i].Fee
		rewardsSplit.Delegators[i].NetRewards = int(rewardsSplit.Delegators[i].GrossRewards - rewardsSplit.Delegators[i].Fee)

		if rewardsSplit.Delegators[i], err = p.constructDexterContractPayout(rewardsSplit.Delegators[i]); err != nil {
			return rewardsSplit, errors.Wrap(err, "failed to contruct payout")
		}
	}

	return rewardsSplit, nil
}

func (p *Payout) calculateTotals(rewards tzkt.RewardsSplit) int {
	if p.config.Baker.EarningsOnly {
		return rewards.EndorsementRewards +
			rewards.RevelationRewards +
			rewards.OwnBlockFees +
			rewards.OwnBlockRewards +
			rewards.ExtraBlockFees +
			rewards.ExtraBlockRewards
	}

	return rewards.EndorsementRewards +
		rewards.MissedEndorsementRewards +
		rewards.RevelationRewards +
		rewards.OwnBlockFees +
		rewards.MissedOwnBlockFees +
		rewards.OwnBlockRewards +
		rewards.MissedOwnBlockRewards +
		rewards.ExtraBlockFees +
		rewards.ExtraBlockRewards
}

// Execute will execute a payout based off the Payout configuration
func (p *Payout) Execute() (tzkt.RewardsSplit, error) {
	payout, err := p.constructPayout()
	if err != nil {
		return payout, errors.Wrapf(err, "failed to execute payout for cycle %d", p.cycle)
	}

	if p.inject {
		operations, err := p.apply(payout.Delegators)
		if err != nil {
			return payout, errors.Wrapf(err, "failed to execute payout for cycle %d", p.cycle)
		}

		for _, op := range operations {
			payout.OperationLink = append(payout.OperationLink, fmt.Sprintf("https://tzkt.io/%s", op))
		}
	}

	return payout, err
}

func (p *Payout) apply(delegators tzkt.Delegators) ([]string, error) {
	head, err := p.rpc.Head()
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to apply payout")
	}

	var operationStrings []string
	for _, transactions := range p.constructTransactionBatches(head.Hash, delegators) {
		if operation, err := forge.Encode(head.Hash, transactions...); err == nil {
			operationStrings = append(operationStrings, operation)
		} else {
			return []string{}, errors.Wrap(err, "failed to forge operation")
		}
	}

	operationHashes, err := p.injectOperations(operationStrings)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to forge operation")
	}

	return operationHashes, nil
}

func (p *Payout) constructTransactionBatches(blockhash string, delegators tzkt.Delegators) []rpc.Contents {
	var transactionBatches []rpc.Contents

	counter, err := p.rpc.Counter(blockhash, p.key.PubKey.GetPublicKeyHash())
	if err != nil {
		return transactionBatches
	}

	for _, batch := range p.batch(delegators) {
		var transactions rpc.Contents
		for _, delegation := range batch {
			if delegation.LiquidityProviders != nil {
				for _, liquidityProvider := range delegation.LiquidityProviders {
					if delegation.NetRewards >= p.config.Baker.MinimumPayment && !delegation.BlackListed { // don't payout to rewards smaller than minimal payment or that are blacklisted
						counter++
						transactions = append(transactions, rpc.Content{
							Kind:         rpc.TRANSACTION,
							Source:       p.key.PubKey.GetPublicKeyHash(),
							Destination:  liquidityProvider.Address,
							Amount:       int64(liquidityProvider.NetRewards),
							Fee:          int64(p.config.Operations.NetworkFee),
							GasLimit:     int64(p.config.Operations.GasLimit),
							Counter:      counter,
							StorageLimit: int64(0),
						})
					}
				}
			} else {
				if delegation.NetRewards >= p.config.Baker.MinimumPayment && !delegation.BlackListed { // don't payout to rewards smaller than minimal payment or that are blacklisted
					counter++
					transactions = append(transactions, rpc.Content{
						Kind:         rpc.TRANSACTION,
						Source:       p.key.PubKey.GetPublicKeyHash(),
						Destination:  delegation.Address,
						Amount:       int64(delegation.NetRewards),
						Fee:          int64(p.config.Operations.NetworkFee),
						GasLimit:     int64(p.config.Operations.GasLimit),
						Counter:      counter,
						StorageLimit: int64(0),
					})
				}
			}
		}

		transactionBatches = append(transactionBatches, transactions)
	}

	return transactionBatches
}

func (p *Payout) batch(delegators tzkt.Delegators) []tzkt.Delegators {
	var batch []tzkt.Delegators
	if len(delegators) <= p.config.Operations.BatchSize {
		return append(batch, delegators)
	}

	for len(delegators) >= p.config.Operations.BatchSize {
		batch = append(batch, delegators[:p.config.Operations.BatchSize])
		delegators = delegators[p.config.Operations.BatchSize:]
	}

	if len(delegators) != 0 {
		batch = append(batch, delegators)
	}

	return batch
}

func (p *Payout) injectOperations(operations []string) ([]string, error) {
	ophashes := []string{}
	for i, op := range operations {
		signedop, err := p.key.Sign(keys.SignInput{
			Message: op,
		})
		if err != nil {
			return ophashes, errors.Wrap(err, "failed to inject operation")
		}

		ophash, err := p.rpc.InjectionOperation(rpc.InjectionOperationInput{
			Operation: fmt.Sprintf("%s%s", op, hex.EncodeToString(signedop.Bytes)),
		})
		if err != nil {
			return ophashes, errors.Wrap(err, "failed to inject operation")
		}

		ophashes = append(ophashes, ophash)

		if p.verbose {
			logrus.WithFields(logrus.Fields{
				"hash":      ophash,
				"operation": fmt.Sprintf("%d/%d", (i + 1), len(operations)),
			}).Info("Confirming injection.")
		}

		if !p.confirmOperation(ophash) {
			return ophashes, errors.Wrap(err, "failed to inject operation: failed to confirm operation")
		}

		if p.verbose {
			logrus.WithFields(logrus.Fields{
				"hash":      ophash,
				"operation": fmt.Sprintf("%d/%d", (i + 1), len(operations)),
			}).Info("Injection confirmed.")
		}
	}

	return ophashes, nil
}

func (p *Payout) confirmOperation(operation string) bool {
	timer := time.After(confirmationTimoutInterval)
	ticker := time.Tick(confirmationDurationInterval)
	for {
		select {
		case <-ticker:
			if head, err := p.rpc.Head(); err == nil {
				if ophashes, err := p.rpc.OperationHashes(head.Hash); err == nil {
					for _, out := range ophashes {
						for _, in := range out {
							if in == operation {
								return true
							}
						}
					}
				}
			}
		case <-timer:
			return false
		}
	}
}

func (p *Payout) isInBlacklist(delegation string) bool {
	for _, b := range p.config.Baker.Blacklist {
		if b == delegation {
			return true
		}
	}

	return false
}

func (p *Payout) isDexterContract(address string) bool {
	for _, contract := range p.config.Baker.DexterLiquidityContracts {
		if contract == address {
			return true
		}
	}

	return false
}
