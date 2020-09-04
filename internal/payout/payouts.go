package payout

import (
	"fmt"
	"time"

	gotezos "github.com/goat-systems/go-tezos/v3"
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
	gt              gotezos.IFace
	tzkt            tzkt.IFace
	cycle           int
	delegate        string
	bakerFee        float64
	wallet          gotezos.Wallet
	minPayment      int
	blacklist       []string
	dexterContracts []string
	inject          bool
	networkFee      int
	gasLimit        int
	batchSize       int
	verbose         bool
	earningsOnly    bool
}

// NewPayoutInput is the input for NewPayout
type NewPayoutInput struct {
	GoTezos         gotezos.IFace
	Cycle           int
	Delegate        string
	BakerFee        float64
	Wallet          gotezos.Wallet
	MinPayment      int
	BlackList       []string
	DexterContracts []string
	Inject          bool // If false, nothing will be injected.
	NetworkFee      int
	GasLimit        int
	BatchSize       int
	Verbose         bool
}

// NewPayout returns a pointer to a new Baker
func NewPayout(input NewPayoutInput) *Payout {
	return &Payout{
		gt:              input.GoTezos,
		cycle:           input.Cycle,
		delegate:        input.Delegate,
		bakerFee:        input.BakerFee,
		wallet:          input.Wallet,
		minPayment:      input.MinPayment,
		inject:          input.Inject,
		networkFee:      input.NetworkFee,
		gasLimit:        input.GasLimit,
		verbose:         input.Verbose,
		batchSize:       input.BatchSize,
		blacklist:       input.BlackList,
		dexterContracts: input.DexterContracts,
	}
}

func (p *Payout) constructPayout() (tzkt.RewardsSplit, error) {
	rewardsSplit, err := p.tzkt.GetRewardsSplit(p.delegate, p.cycle)
	if err != nil {
		return rewardsSplit, errors.Wrap(err, "failed to contruct payout")
	}

	for i := range rewardsSplit.Delegators {
		if p.isInBlacklist(rewardsSplit.Delegators[i].Address) {
			rewardsSplit.Delegators[i].BlackListed = true
		}

		rewardsSplit.Delegators[i].Share = float64(rewardsSplit.Delegators[i].Balance) / float64(rewardsSplit.StakingBalance)
		if p.earningsOnly {
			totalRewards := float64(
				rewardsSplit.EndorsementRewards +
					rewardsSplit.RevelationRewards +
					rewardsSplit.OwnBlockFees +
					rewardsSplit.OwnBlockRewards +
					rewardsSplit.ExtraBlockFees +
					rewardsSplit.ExtraBlockRewards)
			rewardsSplit.Delegators[i].GrossRewards = int(rewardsSplit.Delegators[i].Share * totalRewards)
		} else {
			totalRewards := float64(rewardsSplit.EndorsementRewards +
				rewardsSplit.MissedEndorsementRewards +
				rewardsSplit.RevelationRewards +
				rewardsSplit.OwnBlockFees +
				rewardsSplit.MissedOwnBlockFees +
				rewardsSplit.OwnBlockRewards +
				rewardsSplit.MissedOwnBlockRewards +
				rewardsSplit.ExtraBlockFees +
				rewardsSplit.ExtraBlockRewards)
			rewardsSplit.Delegators[i].GrossRewards = int(rewardsSplit.Delegators[i].Share * totalRewards)
		}
		rewardsSplit.Delegators[i].Fee = int(float64(rewardsSplit.Delegators[i].GrossRewards) * p.bakerFee)
		rewardsSplit.Delegators[i].NetRewards = int(rewardsSplit.Delegators[i].GrossRewards - rewardsSplit.Delegators[i].Fee)

		if rewardsSplit.Delegators[i], err = p.constructDexterContractPayout(rewardsSplit.Delegators[i]); err != nil {
			return rewardsSplit, errors.Wrap(err, "failed to contruct payout")
		}
	}

	return rewardsSplit, nil
}

// Execute will execute a payout based off the Payout configuration
func (p *Payout) Execute() (tzkt.RewardsSplit, error) {
	payout, err := p.constructPayout()
	if err != nil {
		return payout, errors.Wrapf(err, "failed to execute payout for cycle %d", p.cycle)
	}

	forgedOperations, err := p.forge(payout.Delegators)
	if err != nil {
		return payout, errors.Wrapf(err, "failed to execute payout for cycle %d", p.cycle)
	}

	if p.inject {
		operations, err := p.injectOperations(forgedOperations)
		if err != nil {
			err = errors.Wrap(err, "failed to get inject payout for delegate")
		}

		for _, op := range operations {
			payout.OperationLink = append(payout.OperationLink, fmt.Sprintf("https://tzkt.io/%s", op))
		}
	}

	return payout, err
}

func (p *Payout) forge(delegators tzkt.Delegators) ([]string, error) {
	head, err := p.gt.Head()
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to get operation hex string")
	}

	counter, err := p.gt.Counter(head.Hash, p.wallet.Address)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to get operation hex string")
	}

	operations := []string{}
	for _, batch := range p.batch(delegators) {
		var op string
		var err error
		op, counter, err = p.forgeOperation(counter, batch)
		if err != nil {
			return []string{}, errors.Wrap(err, "failed to get operation hex string")
		}

		operations = append(operations, op)
	}

	return operations, nil
}

func (p *Payout) forgeOperation(counter int, delegators tzkt.Delegators) (string, int, error) {
	head, err := p.gt.Head()
	if err != nil {
		return "", counter, errors.Wrap(err, "failed to forge payout")
	}

	transactions, lastCounter := p.constructPayoutContents(counter, delegators)

	var contents []gotezos.OperationContents
	for _, transaction := range transactions {
		contents = append(contents, &transaction)
	}

	forge, err := gotezos.ForgeOperation(head.Hash, contents...)
	if err != nil {
		return "", lastCounter, errors.Wrap(err, "failed to forge payout")
	}

	return forge, lastCounter, nil
}

func (p *Payout) constructPayoutContents(counter int, delegators tzkt.Delegators) ([]gotezos.Transaction, int) {
	var transactions []gotezos.Transaction
	for _, delegation := range delegators {
		if delegation.LiquidityProviders != nil {
			for _, liquidityProvider := range delegation.LiquidityProviders {
				if delegation.NetRewards >= p.minPayment && !delegation.BlackListed { // don't payout to rewards smaller than minimal payment or that are blacklisted
					counter++
					transactions = append(transactions, gotezos.Transaction{
						Source:       p.wallet.Address,
						Destination:  liquidityProvider.Address,
						Amount:       int64(liquidityProvider.NetRewards),
						Fee:          int64(p.networkFee),
						GasLimit:     int64(p.gasLimit),
						Counter:      counter,
						StorageLimit: 0,
					})
				}
			}
		} else {
			if delegation.NetRewards >= p.minPayment && !delegation.BlackListed { // don't payout to rewards smaller than minimal payment or that are blacklisted
				counter++
				transactions = append(transactions, gotezos.Transaction{
					Source:       p.wallet.Address,
					Destination:  delegation.Address,
					Amount:       int64(delegation.NetRewards),
					Fee:          int64(p.networkFee),
					GasLimit:     int64(p.gasLimit),
					Counter:      counter,
					StorageLimit: 0,
				})
			}
		}
	}

	return transactions, counter
}

func (p *Payout) batch(delegators tzkt.Delegators) []tzkt.Delegators {
	var batch []tzkt.Delegators
	if len(delegators) <= p.batchSize {
		return append(batch, delegators)
	}

	for len(delegators) >= p.batchSize {
		batch = append(batch, delegators[:p.batchSize])
		delegators = delegators[p.batchSize:]
	}

	if len(delegators) != 0 {
		batch = append(batch, delegators)
	}

	return batch
}

func (p *Payout) injectOperations(operations []string) ([]string, error) {
	ophashes := []string{}
	for i, op := range operations {
		signedop, err := p.wallet.SignOperation(op)
		if err != nil {
			return ophashes, errors.Wrap(err, "failed to inject operation")
		}

		ophash, err := p.gt.InjectionOperation(gotezos.InjectionOperationInput{
			Operation: signedop.SignedOperation,
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
			if head, err := p.gt.Head(); err == nil {
				if ophashes, err := p.gt.OperationHashes(head.Hash); err == nil {
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
	for _, b := range p.blacklist {
		if b == delegation {
			return true
		}
	}

	return false
}

func (p *Payout) isDexterContract(address string) bool {
	for _, contract := range p.dexterContracts {
		if contract == address {
			return true
		}
	}

	return false
}
