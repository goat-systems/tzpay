package payout

import (
	"encoding/json"
	"strconv"
	"strings"

	gotezos "github.com/goat-systems/go-tezos/v3"
	"github.com/goat-systems/tzpay/v2/internal/tzkt"
	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
)

/*
ExchangeContractV1 represents a liquidity pool contract

storage (pair (big_map %accounts (address :owner)
                                 (pair (nat :balance)
                                       (map (address :spender)
                                            (nat :allowance))))
              (pair (pair (bool :selfIsUpdatingTokenPool)
                          (pair (bool :freezeBaker)
                                (nat :lqtTotal)))
                    (pair (pair (address :manager)
                                (address :tokenAddress))
                          (pair (nat :tokenPool)
                                (mutez :xtzPool)))));
*/
type ExchangeContractV1 struct {
	Prim string `json:"prim"`
	Args []struct {
		Int  int    `json:"int,string,omitempty"`
		Prim string `json:"prim,omitempty"`
		Args []struct {
			Prim string `json:"prim"`
			Args []struct {
				Prim string `json:"prim"`
				Args []struct {
					Prim string `json:"prim,omitempty"`
					Int  int    `json:"int,string,omitempty"`
				} `json:"args,omitempty"`
			} `json:"args"`
		} `json:"args,omitempty"`
	} `json:"args"`
}

// BigMapV1 represents a big_map for ExchangeContractV1
type BigMapV1 struct {
	Prim string          `json:"prim"`
	Args json.RawMessage `json:"args"`
}

func (p *Payout) constructDexterContractPayout(delegator tzkt.Delegator) (tzkt.Delegator, error) {
	if isDex := p.isDexterContract(delegator.Address); isDex {
		var err error
		delegator, err = p.getLiquidityProvidersEarnings(delegator)
		if err != nil {
			return delegator, errors.Wrap(err, "failed to contruct dexter contract payout")
		}
	}

	return delegator, nil
}

func (p *Payout) getLiquidityProvidersEarnings(contract tzkt.Delegator) (tzkt.Delegator, error) {
	cycle, err := p.gt.Cycle(p.cycle)
	if err != nil {
		return contract, errors.Wrapf(err, "failed to get earnings for dexter liquidity providers")
	}

	exchangeContract, err := p.getContractStorage(cycle.BlockHash, contract.Address)
	if err != nil {
		return contract, errors.Wrapf(err, "failed to get earnings for dexter liquidity providers")
	}

	// safe because the contract will fail to marshal if changed
	totalLiquidity := exchangeContract.Args[1].Args[0].Args[1].Args[1].Int
	bigMap := exchangeContract.Args[0].Int

	liquidityProvidersAddresses, err := p.getLiquidityProvidersList(contract.Address)
	if err != nil {
		return contract, errors.Wrapf(err, "failed to get earnings for dexter liquidity providers")
	}

	var liquidityProviders []tzkt.LiquidityProvider
	for _, key := range liquidityProvidersAddresses {
		balance, err := p.getBalanceFromBigMap(key, bigMap, cycle.BlockHash)
		if err != nil {
			return contract, errors.Wrapf(err, "failed to get earnings for liquidity providers for contract '%s'", contract.Address)
		}

		lp := tzkt.LiquidityProvider{
			Address: key,
			Balance: balance,
			Share:   float64(balance) / float64(totalLiquidity),
		}

		lp.GrossRewards = int(lp.Share * float64(contract.GrossRewards))
		lp.Fee = int(float64(lp.GrossRewards) * p.bakerFee)
		lp.NetRewards = lp.GrossRewards - lp.Fee

		if p.isInBlacklist(lp.Address) {
			lp.BlackListed = true
		}

		liquidityProviders = append(liquidityProviders, lp)
	}
	contract.LiquidityProviders = liquidityProviders

	return contract, nil
}

func (p *Payout) getLiquidityProvidersList(target string) ([]string, error) {
	transactions, err := p.tzkt.GetTransactions([]tzkt.URLParameters{
		{
			Key:   "parameters.as",
			Value: "*addLiquidity*",
		},
		{
			Key:   "target",
			Value: target,
		},
	}...)
	if err != nil {
		return []string{}, errors.Wrapf(err, "failed to get list of liquidity providers for '%s'", target)
	}

	liquidityProvidersAddresses := map[string]struct{}{} // map to weed out duplicates
	for _, lp := range transactions {
		liquidityProvidersAddresses[lp.Sender.Address] = struct{}{}
	}

	out := []string{}
	for key := range liquidityProvidersAddresses {
		out = append(out, key)
	}

	return out, nil
}

func (p *Payout) getBalanceFromBigMap(key string, bigMapID int, blockhash string) (int, error) {
	scriptExp, err := gotezos.ForgeScriptExpressionForAddress(key)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get balance from big_map for '%s'", key)
	}

	bigMapResp, err := p.gt.BigMap(gotezos.BigMapInput{
		Blockhash:        blockhash,
		BigMapID:         bigMapID,
		ScriptExpression: scriptExp,
	})

	var bigmap BigMapV1
	if err := json.Unmarshal(bigMapResp, &bigmap); err != nil {
		return 0, errors.Wrapf(err, "failed to get balance from big_map for '%s'", key)
	}

	balance, err := parseBigMapForBalance(&bigmap.Args)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get balance from big_map for '%s'", key)
	}

	return balance, nil
}

func (p *Payout) getContractStorage(blockhash string, address string) (ExchangeContractV1, error) {
	storage, err := p.gt.ContractStorage(blockhash, address) //CHANGE TO cycle.Blockhash later
	if err != nil {
		return ExchangeContractV1{}, errors.Wrapf(err, "failed to get storage for contract '%s'", address)
	}

	var exchangeContract ExchangeContractV1
	if err := json.Unmarshal(storage, &exchangeContract); err != nil {
		return ExchangeContractV1{}, errors.Wrapf(err, "failed to get storage contract '%s'", address)
	}

	return exchangeContract, nil
}

func parseBigMapForBalance(msg *json.RawMessage) (int, error) {
	var p fastjson.Parser
	v, err := p.Parse(string(*msg))
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse as json blob")
	}

	args, err := v.Array()
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse args in json blob")
	}

	balanceObject, err := args[0].Object()
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse args in json blob")
	}

	return strconv.Atoi(strings.Trim(balanceObject.Get("int").String(), "\""))
}
