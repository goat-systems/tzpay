package payout

import (
	"fmt"
	"math/big"
	"sort"
	"time"
	"unicode"

	gotezos "github.com/goat-systems/go-tezos/v3"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	confirmationDurationInterval = time.Second * 1
	confirmationTimoutInterval   = time.Minute * 2
)

// Payout represents a payout and payout operations.
type Payout struct {
	gt         gotezos.IFace
	cycle      int
	delegate   string
	bakerFee   float64
	wallet     gotezos.Wallet
	minPayment int
	blacklist  []string
	inject     bool
	networkFee int
	gasLimit   int
	batchSize  int
	verbose    bool
}

// NewPayoutInput is the input for NewPayout
type NewPayoutInput struct {
	GoTezos    gotezos.IFace
	Cycle      int
	Delegate   string
	BakerFee   float64
	Wallet     gotezos.Wallet
	MinPayment int
	BlackList  []string
	Inject     bool // If false, nothing will be injected.
	NetworkFee int
	GasLimit   int
	BatchSize  int
	Verbose    bool
}

// Report contains all needed information for a payout
type Report struct {
	DelegationEarnings DelegationEarnings `json:"delegaions"`
	DelegateEarnings   DelegateEarnings   `json:"delegate"`
	CycleHash          string             `json:"cycle_hash"`
	Cycle              int                `json:"cycle"`
	FrozenBalance      *big.Int           `json:"rewards"`
	StakingBalance     *big.Int           `json:"staking_balance"`
	Operations         []string           `json:"operation"`
	OperationsLink     []string           `json:"operation_link"`
}

// DelegationEarning -
type DelegationEarning struct {
	Address      string
	Fee          *big.Int
	GrossRewards *big.Int
	NetRewards   *big.Int
	Share        float64
}

// DelegateEarnings -
type DelegateEarnings struct {
	Address string
	Fees    *big.Int
	Share   float64
	Rewards *big.Int
	Net     *big.Int
}

// DelegationEarnings contains list of DelegationEarning and implements sort.
type DelegationEarnings []DelegationEarning

func (d DelegationEarnings) Len() int { return len(d) }
func (d DelegationEarnings) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d DelegationEarnings) Less(i, j int) bool {
	iRunes := []rune(d[i].Address)
	jRunes := []rune(d[j].Address)

	max := len(iRunes)
	if max > len(jRunes) {
		max = len(jRunes)
	}

	for idx := 0; idx < max; idx++ {
		ir := iRunes[idx]
		jr := jRunes[idx]

		lir := unicode.ToLower(ir)
		ljr := unicode.ToLower(jr)

		if lir != ljr {
			return lir < ljr
		}

		if ir != jr {
			return ir < jr
		}
	}

	return false
}

type processDelegationsInput struct {
	delegations          []*string
	stakingBalance       *big.Int
	frozenBalanceRewards gotezos.FrozenBalance
	blockHash            string
}

type processDelegateInput struct {
	delegate             string
	delegations          []DelegationEarning
	stakingBalance       *big.Int
	frozenBalanceRewards gotezos.FrozenBalance
	blockHash            string
}

type processDelegationsOutput struct {
	delegationEarning DelegationEarning
	err               error
}

type processDelegationInput struct {
	delegation           string
	stakingBalance       *big.Int
	frozenBalanceRewards gotezos.FrozenBalance
	blockHash            string
}

// NewPayout returns a pointer to a new Baker
func NewPayout(input NewPayoutInput) *Payout {
	return &Payout{
		gt:         input.GoTezos,
		cycle:      input.Cycle,
		delegate:   input.Delegate,
		bakerFee:   input.BakerFee,
		wallet:     input.Wallet,
		minPayment: input.MinPayment,
		inject:     input.Inject,
		networkFee: input.NetworkFee,
		gasLimit:   input.GasLimit,
		verbose:    input.Verbose,
		batchSize:  input.BatchSize,
		blacklist:  input.BlackList,
	}
}

// Execute will execute a payout based off the Payout configuration
func (p *Payout) Execute() (Report, error) {
	frozenBalanceRewards, err := p.gt.FrozenBalance(p.cycle, p.delegate)
	if err != nil {
		return Report{}, errors.Wrapf(err, "failed to get delegation earnings for cycle %d", p.cycle)
	}

	rpcDelegations, err := p.gt.DelegatedContractsAtCycle(p.cycle, p.delegate)
	if err != nil {
		return Report{}, errors.Wrapf(err, "failed to get delegation earnings for cycle %d", p.cycle)
	}

	var delegations []*string
	for _, delegation := range rpcDelegations {
		if !p.isInBlacklist(*delegation) {
			delegations = append(delegations, delegation)
		}
	}

	networkCycle, err := p.gt.Cycle(p.cycle)
	if err != nil {
		return Report{}, errors.Wrapf(err, "failed to get delegation earnings for cycle %d", p.cycle)
	}

	stakingBalance, err := p.gt.StakingBalance(networkCycle.BlockHash, p.delegate)
	if err != nil {
		return Report{}, errors.Wrapf(err, "failed to get delegation earnings for cycle %d", p.cycle)
	}

	delegationsOutput := p.proccessDelegations(processDelegationsInput{
		delegations:          delegations,
		stakingBalance:       stakingBalance,
		frozenBalanceRewards: frozenBalanceRewards,
		blockHash:            networkCycle.BlockHash,
	})

	report := Report{
		StakingBalance: stakingBalance,
		CycleHash:      networkCycle.BlockHash,
		Cycle:          p.cycle,
		FrozenBalance:  frozenBalanceRewards.Rewards.Big,
	}

	for _, delegation := range delegationsOutput {
		if delegation.err != nil {
			err = errors.Wrapf(delegation.err, "failed to get payout for delegation %s", delegation.delegationEarning.Address)
		} else {
			report.DelegationEarnings = append(report.DelegationEarnings, delegation.delegationEarning)
		}
	}
	sort.Sort(report.DelegationEarnings)

	if report.DelegateEarnings, err = p.processDelegate(processDelegateInput{
		delegate:             p.delegate,
		delegations:          report.DelegationEarnings,
		stakingBalance:       stakingBalance,
		frozenBalanceRewards: frozenBalanceRewards,
		blockHash:            networkCycle.BlockHash,
	}); err != nil {
		err = errors.Wrap(err, "failed to get contruct payout info for delegate")
	}

	operations, err := p.getOperationHexStrings(report.DelegationEarnings)
	if err != nil {
		err = errors.Wrap(err, "failed to get contruct payout for delegate")
	}

	if p.inject {
		operationHashes, err := p.injectOperations(operations)
		if err != nil {
			err = errors.Wrap(err, "failed to get inject payout for delegate")
		}
		report.Operations = operationHashes
		for _, op := range operationHashes {
			report.OperationsLink = append(report.OperationsLink, fmt.Sprintf("https://tzstats.com/%s", op))
		}
	}

	return report, err
}

func (p *Payout) processDelegate(input processDelegateInput) (DelegateEarnings, error) {
	delegateEarning := DelegateEarnings{
		Address: input.delegate,
		Net:     big.NewInt(0),
	}

	balance, err := p.gt.Balance(input.blockHash, input.delegate)
	if err != nil {
		return delegateEarning, errors.Wrapf(err, "failed to process delegate earnings for %s", input.delegate)
	}

	delegateEarning.Share = float64(balance.Int64()) / float64(input.stakingBalance.Int64())
	rewardsFloat := delegateEarning.Share * float64(input.frozenBalanceRewards.Rewards.Big.Int64())
	delegateEarning.Rewards = big.NewInt(int64(rewardsFloat))

	fees := big.NewInt(0)
	for _, delegation := range input.delegations {
		fees.Add(fees, delegation.Fee)
	}

	delegateEarning.Fees = fees
	delegateEarning.Net.Add(delegateEarning.Fees, delegateEarning.Rewards)

	return delegateEarning, nil
}

func (p *Payout) proccessDelegations(input processDelegationsInput) []processDelegationsOutput {
	numJobs := len(input.delegations)
	jobs := make(chan processDelegationInput, numJobs)
	results := make(chan processDelegationsOutput, numJobs)

	for i := 0; i < 50; i++ {
		go p.proccessDelegationWorker(jobs, results)
	}

	for _, delegation := range input.delegations {
		jobs <- processDelegationInput{
			delegation:           *delegation,
			stakingBalance:       input.stakingBalance,
			frozenBalanceRewards: input.frozenBalanceRewards,
			blockHash:            input.blockHash,
		}
	}
	close(jobs)

	var out []processDelegationsOutput
	for i := 1; i <= numJobs; i++ {
		out = append(out, <-results)
	}

	close(results)
	return out
}

func (p *Payout) proccessDelegationWorker(jobs <-chan processDelegationInput, results chan<- processDelegationsOutput) {
	for j := range jobs {
		d, err := p.processDelegation(j)
		if err != nil {
			results <- processDelegationsOutput{
				err: err,
			}
		} else {
			results <- processDelegationsOutput{
				delegationEarning: *d,
			}
		}
	}
}

func (p *Payout) processDelegation(input processDelegationInput) (*DelegationEarning, error) {
	delegationEarning := &DelegationEarning{Address: input.delegation}
	balance, err := p.gt.Balance(input.blockHash, input.delegation)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to process delegation earnings for delegation %s", input.delegation)
	}

	delegationEarning.Share = float64(balance.Int64()) / float64(input.stakingBalance.Int64())
	grossRewardsFloat := delegationEarning.Share * float64(input.frozenBalanceRewards.Rewards.Big.Int64())
	feeFloat := grossRewardsFloat * p.bakerFee

	delegationEarning.GrossRewards = big.NewInt(int64(grossRewardsFloat))
	delegationEarning.Fee = big.NewInt(int64(feeFloat))
	delegationEarning.NetRewards = big.NewInt(0).Sub(delegationEarning.GrossRewards, delegationEarning.Fee)

	return delegationEarning, nil
}

func (p *Payout) getOperationHexStrings(delegationEarnings DelegationEarnings) ([]string, error) {
	head, err := p.gt.Head()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get operation hex string")
	}

	counter, err := p.gt.Counter(head.Hash, p.wallet.Address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get operation hex string")
	}

	operations := []string{}
	for _, batch := range p.batch(delegationEarnings) {
		var op string
		var err error
		op, counter, err = p.forgeOperation(counter, batch)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get operation hex string")
		}

		operations = append(operations, op)
	}

	return operations, nil
}

func (p *Payout) forgeOperation(counter int, delegationEarnings DelegationEarnings) (string, int, error) {
	head, err := p.gt.Head()
	if err != nil {
		return "", counter, errors.Wrap(err, "failed to forge payout")
	}

	transactions, lastCounter := p.constructPayoutContents(counter, delegationEarnings)

	forge, err := gotezos.ForgeOperation(head.Hash, transactions)
	if err != nil {
		return "", lastCounter, errors.Wrap(err, "failed to forge payout")
	}

	return forge, lastCounter, nil
}

func (p *Payout) constructPayoutContents(counter int, delegationEarnings DelegationEarnings) (gotezos.Contents, int) {
	var transactions []gotezos.Transaction
	for _, delegation := range delegationEarnings {
		if delegation.NetRewards.Int64() >= int64(p.minPayment) {
			counter++
			transactions = append(transactions, gotezos.Transaction{
				Kind:         gotezos.TRANSACTION,
				Source:       p.wallet.Address,
				Destination:  delegation.Address,
				Amount:       &gotezos.Int{Big: delegation.NetRewards},
				Fee:          gotezos.NewInt(p.networkFee),
				GasLimit:     gotezos.NewInt(p.gasLimit),
				Counter:      counter,
				StorageLimit: gotezos.NewInt(0),
			})
		}
	}

	return gotezos.Contents{
		Transactions: transactions,
	}, counter
}

func (p *Payout) batch(delegationEarnings DelegationEarnings) []DelegationEarnings {
	var delegationEarningsBatch []DelegationEarnings
	if len(delegationEarnings) <= p.batchSize {
		delegationEarningsBatch = append(delegationEarningsBatch, delegationEarnings)
		return delegationEarningsBatch
	}

	for len(delegationEarnings) >= p.batchSize {
		delegationEarningsBatch = append(delegationEarningsBatch, delegationEarnings[:p.batchSize])
		delegationEarnings = delegationEarnings[p.batchSize:]
	}

	if len(delegationEarnings) != 0 {
		delegationEarningsBatch = append(delegationEarningsBatch, delegationEarnings)
	}

	return delegationEarningsBatch
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
