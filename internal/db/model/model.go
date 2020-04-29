package model

import (
	"fmt"
	"math/big"
	"unicode"
)

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

// Payout contains all needed information for a payout
type Payout struct {
	DelegationEarnings DelegationEarnings `json:"delegaions"`
	DelegateEarnings   DelegateEarnings   `json:"delegate"`
	CycleHash          string             `json:"cycle_hash"`
	Cycle              int                `json:"cycle"`
	FrozenBalance      *big.Int           `json:"rewards"`
	StakingBalance     *big.Int           `json:"staking_balance"`
	Operations         []string           `json:"operation"`
	OperationsLink     []string           `json:"operation_link"`
}

// SetOperations will set Payout's operation string and link
func (p *Payout) SetOperations(operations ...string) {
	p.Operations = append(p.Operations, operations...)
	for _, operation := range operations {
		p.OperationsLink = append(p.OperationsLink, fmt.Sprintf("http://tzstats.com/%s", operation))
	}
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
