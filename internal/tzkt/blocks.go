package tzkt

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type Head struct {
	Level      int       `json:"level"`
	Hash       string    `json:"hash"`
	Protocol   string    `json:"protocol"`
	Timestamp  time.Time `json:"timestamp"`
	KnownLevel int       `json:"knownLevel"`
	LastSync   time.Time `json:"lastSync"`
	Synced     bool      `json:"synced"`
	QuoteLevel int       `json:"quoteLevel"`
	QuoteBtc   int       `json:"quoteBtc"`
	QuoteEur   int       `json:"quoteEur"`
	QuoteUsd   int       `json:"quoteUsd"`
}

type Blocks []struct {
	Level         int       `json:"level"`
	Hash          string    `json:"hash"`
	Timestamp     time.Time `json:"timestamp"`
	Proto         int       `json:"proto"`
	Priority      int       `json:"priority"`
	Validations   int       `json:"validations"`
	Reward        int       `json:"reward"`
	Fees          int       `json:"fees"`
	NonceRevealed bool      `json:"nonceRevealed"`
	Baker         struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	} `json:"baker"`
	Endorsements []struct {
		Type      string    `json:"type"`
		ID        int       `json:"id"`
		Level     int       `json:"level"`
		Timestamp time.Time `json:"timestamp"`
		Block     string    `json:"block"`
		Hash      string    `json:"hash"`
		Delegate  struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"delegate"`
		Slots   int `json:"slots"`
		Rewards int `json:"rewards"`
		Quote   struct {
			Btc int `json:"btc"`
			Eur int `json:"eur"`
			Usd int `json:"usd"`
		} `json:"quote"`
	} `json:"endorsements"`
	Proposals []struct {
		Type      string    `json:"type"`
		ID        int       `json:"id"`
		Level     int       `json:"level"`
		Timestamp time.Time `json:"timestamp"`
		Block     string    `json:"block"`
		Hash      string    `json:"hash"`
		Period    struct {
			ID         int    `json:"id"`
			Kind       string `json:"kind"`
			StartLevel int    `json:"startLevel"`
			EndLevel   int    `json:"endLevel"`
		} `json:"period"`
		Proposal struct {
			Alias string `json:"alias"`
			Hash  string `json:"hash"`
		} `json:"proposal"`
		Delegate struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"delegate"`
		Rolls      int  `json:"rolls"`
		Duplicated bool `json:"duplicated"`
		Quote      struct {
			Btc int `json:"btc"`
			Eur int `json:"eur"`
			Usd int `json:"usd"`
		} `json:"quote"`
	} `json:"proposals"`
	Ballots []struct {
		Type      string    `json:"type"`
		ID        int       `json:"id"`
		Level     int       `json:"level"`
		Timestamp time.Time `json:"timestamp"`
		Block     string    `json:"block"`
		Hash      string    `json:"hash"`
		Period    struct {
			ID         int    `json:"id"`
			Kind       string `json:"kind"`
			StartLevel int    `json:"startLevel"`
			EndLevel   int    `json:"endLevel"`
		} `json:"period"`
		Proposal struct {
			Alias string `json:"alias"`
			Hash  string `json:"hash"`
		} `json:"proposal"`
		Delegate struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"delegate"`
		Rolls int    `json:"rolls"`
		Vote  string `json:"vote"`
		Quote struct {
			Btc int `json:"btc"`
			Eur int `json:"eur"`
			Usd int `json:"usd"`
		} `json:"quote"`
	} `json:"ballots"`
	Activations []struct {
		Type      string    `json:"type"`
		ID        int       `json:"id"`
		Level     int       `json:"level"`
		Timestamp time.Time `json:"timestamp"`
		Block     string    `json:"block"`
		Hash      string    `json:"hash"`
		Account   struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"account"`
		Balance int `json:"balance"`
		Quote   struct {
			Btc int `json:"btc"`
			Eur int `json:"eur"`
			Usd int `json:"usd"`
		} `json:"quote"`
	} `json:"activations"`
	DoubleBaking []struct {
		Type         string    `json:"type"`
		ID           int       `json:"id"`
		Level        int       `json:"level"`
		Timestamp    time.Time `json:"timestamp"`
		Block        string    `json:"block"`
		Hash         string    `json:"hash"`
		AccusedLevel int       `json:"accusedLevel"`
		Accuser      struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"accuser"`
		AccuserRewards int `json:"accuserRewards"`
		Offender       struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"offender"`
		OffenderLostDeposits int `json:"offenderLostDeposits"`
		OffenderLostRewards  int `json:"offenderLostRewards"`
		OffenderLostFees     int `json:"offenderLostFees"`
		Quote                struct {
			Btc int `json:"btc"`
			Eur int `json:"eur"`
			Usd int `json:"usd"`
		} `json:"quote"`
	} `json:"doubleBaking"`
	DoubleEndorsing []struct {
		Type         string    `json:"type"`
		ID           int       `json:"id"`
		Level        int       `json:"level"`
		Timestamp    time.Time `json:"timestamp"`
		Block        string    `json:"block"`
		Hash         string    `json:"hash"`
		AccusedLevel int       `json:"accusedLevel"`
		Accuser      struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"accuser"`
		AccuserRewards int `json:"accuserRewards"`
		Offender       struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"offender"`
		OffenderLostDeposits int `json:"offenderLostDeposits"`
		OffenderLostRewards  int `json:"offenderLostRewards"`
		OffenderLostFees     int `json:"offenderLostFees"`
		Quote                struct {
			Btc int `json:"btc"`
			Eur int `json:"eur"`
			Usd int `json:"usd"`
		} `json:"quote"`
	} `json:"doubleEndorsing"`
	NonceRevelations []struct {
		Type      string    `json:"type"`
		ID        int       `json:"id"`
		Level     int       `json:"level"`
		Timestamp time.Time `json:"timestamp"`
		Block     string    `json:"block"`
		Hash      string    `json:"hash"`
		Baker     struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"baker"`
		BakerRewards int `json:"bakerRewards"`
		Sender       struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"sender"`
		RevealedLevel int `json:"revealedLevel"`
		Quote         struct {
			Btc int `json:"btc"`
			Eur int `json:"eur"`
			Usd int `json:"usd"`
		} `json:"quote"`
	} `json:"nonceRevelations"`
	Delegations []struct {
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
		Nonce        int `json:"nonce"`
		GasLimit     int `json:"gasLimit"`
		GasUsed      int `json:"gasUsed"`
		BakerFee     int `json:"bakerFee"`
		PrevDelegate struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"prevDelegate"`
		NewDelegate struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"newDelegate"`
		Status string `json:"status"`
		Errors []struct {
			Type string `json:"type"`
		} `json:"errors"`
		Quote struct {
			Btc int `json:"btc"`
			Eur int `json:"eur"`
			Usd int `json:"usd"`
		} `json:"quote"`
	} `json:"delegations"`
	Originations []struct {
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
		Nonce           int `json:"nonce"`
		GasLimit        int `json:"gasLimit"`
		GasUsed         int `json:"gasUsed"`
		StorageLimit    int `json:"storageLimit"`
		StorageUsed     int `json:"storageUsed"`
		BakerFee        int `json:"bakerFee"`
		StorageFee      int `json:"storageFee"`
		AllocationFee   int `json:"allocationFee"`
		ContractBalance int `json:"contractBalance"`
		ContractManager struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"contractManager"`
		ContractDelegate struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"contractDelegate"`
		Status string `json:"status"`
		Errors []struct {
			Type string `json:"type"`
		} `json:"errors"`
		OriginatedContract struct {
			Kind    string `json:"kind"`
			Alias   string `json:"alias"`
			Address string `json:"address"`
		} `json:"originatedContract"`
		Quote struct {
			Btc int `json:"btc"`
			Eur int `json:"eur"`
			Usd int `json:"usd"`
		} `json:"quote"`
	} `json:"originations"`
	Transactions []struct {
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
	} `json:"transactions"`
	Reveals []struct {
		Type      string    `json:"type"`
		ID        int       `json:"id"`
		Level     int       `json:"level"`
		Timestamp time.Time `json:"timestamp"`
		Block     string    `json:"block"`
		Hash      string    `json:"hash"`
		Sender    struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"sender"`
		Counter  int    `json:"counter"`
		GasLimit int    `json:"gasLimit"`
		GasUsed  int    `json:"gasUsed"`
		BakerFee int    `json:"bakerFee"`
		Status   string `json:"status"`
		Errors   []struct {
			Type string `json:"type"`
		} `json:"errors"`
		Quote struct {
			Btc int `json:"btc"`
			Eur int `json:"eur"`
			Usd int `json:"usd"`
		} `json:"quote"`
	} `json:"reveals"`
	Quote struct {
		Btc int `json:"btc"`
		Eur int `json:"eur"`
		Usd int `json:"usd"`
	} `json:"quote"`
}

/*
GetHead -
See: https://api.tzkt.io/#operation/Head_Get
*/
func (t *Tzkt) GetHead() (Head, error) {
	resp, err := t.get("/v1/head")
	if err != nil {
		return Head{}, errors.Wrapf(err, "failed to get head")
	}

	var head Head
	if err := json.Unmarshal(resp, &head); err != nil {
		return Head{}, errors.Wrap(err, "failed to get head")
	}

	return head, nil
}

/*
GetBlocks -
See: https://api.tzkt.io/#operation/Blocks_Get
*/
func (t *Tzkt) GetBlocks(options ...URLParameters) (Blocks, error) {
	resp, err := t.get("/v1/blocks", options...)
	if err != nil {
		return Blocks{}, errors.Wrapf(err, "failed to get blocks")
	}

	var blocks Blocks
	if err := json.Unmarshal(resp, &blocks); err != nil {
		return Blocks{}, errors.Wrap(err, "failed to get blocks")
	}

	return blocks, nil
}
