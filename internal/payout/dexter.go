package payout

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/goat-systems/go-tezos/v3/rpc"
	"github.com/goat-systems/tzpay/v3/internal/tzkt"
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

type ExchangeContractV15 struct {
	Prim string `json:"prim"`
	Args []Args `json:"args"`
}

type Args struct {
	Int  int    `json:"int,string,omitempty"`
	Prim string `json:"prim,omitempty"`
	Args []Args `json:"args,omitempty"`
}

// BigMapV1 represents a big_map for ExchangeContractV1
type BigMapV1 struct {
	Prim string          `json:"prim"`
	Args json.RawMessage `json:"args"`
}

func (p *Payout) constructDexterContractPayout(delegator tzkt.Delegator) (tzkt.Delegator, error) {
	delegator, err := p.getLiquidityProvidersEarnings(delegator)
	if err != nil {
		return delegator, errors.Wrap(err, "failed to contruct dexter contract payout")
	}

	return delegator, nil
}

func (p *Payout) getLiquidityProvidersEarnings(contract tzkt.Delegator) (tzkt.Delegator, error) {
	cycle, err := p.rpc.Cycle(p.cycle)
	if err != nil {
		return contract, errors.Wrapf(err, "failed to get earnings for dexter liquidity providers")
	}

	totalLiquidity, bigMap, err := p.processContract(cycle.BlockHash, contract.Address)
	if err != nil {
		return contract, err
	}

	liquidityProvidersAddresses, err := p.getLiquidityProvidersList(contract.Address)
	if err != nil {
		return contract, errors.Wrapf(err, "failed to get earnings for dexter liquidity providers")
	}

	var liquidityProviders []tzkt.LiquidityProvider
	for _, key := range liquidityProvidersAddresses {
		found := true
		balance, err := p.getBalanceFromBigMap(key, bigMap, cycle.BlockHash)
		if err != nil {
			if !strings.Contains(err.Error(), "not found in big map") {
				return contract, errors.Wrapf(err, "failed to get earnings for liquidity providers for contract '%s'", contract.Address)
			}
			found = false
		}

		if found {
			lp := tzkt.LiquidityProvider{
				Address: key,
				Balance: balance,
				Share:   float64(balance) / float64(totalLiquidity),
			}

			lp.GrossRewards = int(lp.Share * float64(contract.GrossRewards))
			lp.Fee = int(float64(lp.GrossRewards) * p.config.Baker.Fee)
			lp.NetRewards = lp.GrossRewards - lp.Fee

			if lp.NetRewards < p.config.Baker.MinimumPayment {
				lp.BlackListed = true
			}

			if p.isInBlacklist(lp.Address) {
				lp.BlackListed = true
			}

			if !p.config.Baker.BakerPaysBurnFees {
				requiresBurnFee, err := p.requiresBurnFee(lp.Address)
				if err != nil {
					return contract, errors.Wrapf(err, "failed to get earnings for liquidity providers for contract '%s'", contract.Address)
				}
				if requiresBurnFee {
					lp.BlackListed = true
				}
			}

			liquidityProviders = append(liquidityProviders, lp)
		}
	}
	contract.LiquidityProviders = liquidityProviders

	return contract, nil
}

func (p *Payout) processContract(blockHash, address string) (int, int, error) {
	totalLiquidity, bigMap, err := p.processV1Contract(blockHash, address)
	if err != nil {
		var v15err error
		totalLiquidity, bigMap, v15err = p.processV15Contract(blockHash, address)
		if v15err != nil {
			return 0, 0, errors.Wrap(errors.Wrap(err, v15err.Error()), "failed to get earnings for dexter liquidity providers")
		}
	}

	return totalLiquidity, bigMap, nil
}

func (p *Payout) processV15Contract(blockHash, address string) (int, int, error) {
	contract, err := p.getContractStorageV15(blockHash, address)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to parse v1 contract")
	}

	liquidity, err := getLiquidityV15(contract)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to parse v1 contract")
	}

	bigMap, err := getBigMapV15(contract)
	if err != nil {
		return liquidity, 0, errors.Wrapf(err, "failed to parse v1 contract")
	}

	return liquidity, bigMap, nil
}

func (p *Payout) processV1Contract(blockHash, address string) (int, int, error) {
	contract, err := p.getContractStorageV1(blockHash, address)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to parse v1 contract")
	}

	liquidity, err := getLiquidityV1(contract)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to parse v1 contract")
	}

	bigMap, err := getBigMapV1(contract)
	if err != nil {
		return liquidity, 0, errors.Wrapf(err, "failed to parse v1 contract")
	}

	return liquidity, bigMap, nil
}

func getLiquidityV1(contract ExchangeContractV1) (int, error) {
	if len(contract.Args) >= 2 {
		if len(contract.Args[1].Args) >= 1 {
			if len(contract.Args[1].Args[0].Args) >= 2 {
				if len(contract.Args[1].Args[0].Args[1].Args) >= 2 {
					if contract.Args[1].Args[0].Args[1].Args[1].Int == 0 {
						return 0, errors.New("no liquidity")
					}
					return contract.Args[1].Args[0].Args[1].Args[1].Int, nil
				}
			}
		}
	}

	return 0, errors.New("failed to get liquidity from v1 contract: contract may not be v1")
}

func getBigMapV1(contract ExchangeContractV1) (int, error) {
	if len(contract.Args) >= 1 {
		if contract.Args[0].Int == 0 {
			return 0, errors.New("invalid big map id")
		}

		return contract.Args[0].Int, nil
	}
	return 0, errors.New("failed to get big_map from v1 contract: contract may not be v1")
}

func getLiquidityV15(contract ExchangeContractV15) (int, error) {
	// {"prim":"Pair","args":[{"int":"541"},{"prim":"Pair","args":[{"prim":"False"},{"prim":"False"},{"int":"49707523463"}]},{"prim":"Pair","args":[{"string":"KT1B5VTw8ZSMnrjhy337CEvAm4tnT8Gu8Geu"},{"string":"KT1PWx2mnDueood7fEmfbBDKx1D9BAnnXitn"}]},{"int":"382997319"},{"int":"47813915032"}]}`)
	if len(contract.Args) >= 5 {
		if len(contract.Args[1].Args) == 3 {
			if contract.Args[1].Args[2].Int == 0 {
				return 0, errors.New("no liquidity")
			}
		}
		return contract.Args[1].Args[2].Int, nil
	}

	return 0, errors.New("failed to get liquidity from v1 contract: contract may not be v1")
}

func getBigMapV15(contract ExchangeContractV15) (int, error) {
	if len(contract.Args) >= 1 {
		if contract.Args[0].Int == 0 {
			return 0, errors.New("invalid big map id")
		}

		return contract.Args[0].Int, nil
	}
	return 0, errors.New("failed to get big_map from v1 contract: contract may not be v1")
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
		{
			Key:   "limit",
			Value: "10000",
		},
	}...)
	if err != nil || len(transactions) >= 10000 {
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
	scriptExp, err := rpc.ForgeScriptExpressionForAddress(key)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get balance from big_map for '%s'", key)
	}

	bigMapResp, err := p.rpc.BigMap(rpc.BigMapInput{
		Blockhash:        blockhash,
		BigMapID:         bigMapID,
		ScriptExpression: scriptExp,
	})

	if err != nil {
		if len(bigMapResp) == 0 {
			return 0, errors.Wrapf(err, "key '%s' not found in big map", key)
		}
		return 0, errors.Wrapf(err, "failed to get balance from big_map for '%s'", key)
	}

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

func (p *Payout) getContractStorageV1(blockhash string, address string) (ExchangeContractV1, error) {
	storage, err := p.rpc.ContractStorage(blockhash, address) //CHANGE TO cycle.Blockhash later
	if err != nil {
		return ExchangeContractV1{}, errors.Wrapf(err, "failed to get storage for contract '%s'", address)
	}

	var exchangeContract ExchangeContractV1
	if err := json.Unmarshal(storage, &exchangeContract); err != nil {
		return ExchangeContractV1{}, errors.Wrapf(err, "failed to get storage contract '%s'", address)
	}

	return exchangeContract, nil
}

func (p *Payout) getContractStorageV15(blockhash string, address string) (ExchangeContractV15, error) {
	storage, err := p.rpc.ContractStorage(blockhash, address) //CHANGE TO cycle.Blockhash later
	if err != nil {
		return ExchangeContractV15{}, errors.Wrapf(err, "failed to get storage for contract '%s'", address)
	}

	var exchangeContract ExchangeContractV15
	if err := json.Unmarshal(storage, &exchangeContract); err != nil {
		return ExchangeContractV15{}, errors.Wrapf(err, "failed to get storage contract '%s'", address)
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
