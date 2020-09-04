package tzkt

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

/*
Transaction -
see: https://api.tzkt.io/#operation/Operations_GetTransactions
*/
type Transaction struct {
	Type      string    `json:"type"`
	ID        int       `json:"id"`
	Level     int       `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Block     string    `json:"block"`
	Hash      string    `json:"hash"`
	Counter   int       `json:"counter"`
	Initiator struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	} `json:"initiator"`
	Sender struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	} `json:"sender"`
	Nonce         int `json:"nonce"`
	GasLimit      int `json:"gasLimit"`
	GasUsed       int `json:"gasUsed"`
	StorageLimit  int `json:"storageLimit"`
	StorageUsed   int `json:"storageUsed"`
	BakerFee      int `json:"bakerFee"`
	StorageFee    int `json:"storageFee"`
	AllocationFee int `json:"allocationFee"`
	Target        struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	} `json:"target"`
	Amount     int    `json:"amount"`
	Parameters string `json:"parameters"`
	Status     string `json:"status"`
	Errors     []struct {
		Type string `json:"type"`
	} `json:"errors"`
	HasInternals bool `json:"hasInternals"`
	Quote        struct {
		Btc int `json:"btc"`
		Eur int `json:"eur"`
		Usd int `json:"usd"`
	} `json:"quote"`
}

/*
GetTransactions -
see: https://api.tzkt.io/#operation/Operations_GetTransactions
*/
func (t *Tzkt) GetTransactions(options ...URLParameters) ([]Transaction, error) {
	resp, err := t.get("/v1/operations/transactions", options...)
	if err != nil {
		return []Transaction{}, errors.Wrapf(err, "failed to get transactions")
	}

	var transactions []Transaction
	if err := json.Unmarshal(resp, &transactions); err != nil {
		return []Transaction{}, errors.Wrap(err, "failed to get transactions")
	}

	return transactions, nil
}
