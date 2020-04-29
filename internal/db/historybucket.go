package db

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/goat-systems/tzpay/v2/internal/db/model"
	"github.com/pkg/errors"
)

// SavePayout stores an executed run into the tzpay db for further reference
func (db *DB) SavePayout(payout model.Payout) error {
	if err := db.bolt.Update(func(tx *bolt.Tx) error {
		return db.storePayout(tx, payout)
	}); err != nil {
		return errors.Wrap(err, "failed to save payout in tzpay's history bucket")
	}

	return nil
}

func (db *DB) storePayout(tx *bolt.Tx, payout model.Payout) error {
	b := tx.Bucket([]byte(string(historyBucket)))
	if b == nil {
		return errors.New("failed to open history bucket")
	}

	payoutBytes, err := json.Marshal(&payout)
	if err != nil {
		errors.New("failed to open history bucket")
	}

	fmt.Println(string(payoutBytes))
	err = b.Put([]byte(payout.CycleHash), payoutBytes)
	if err != nil {
		return errors.Wrap(err, "failed to put payout into history bucket")
	}

	return nil
}

// GetPayout retrieves a payout by cycle
func (db *DB) GetPayout(cycle int) (*model.Payout, error) {
	networkCycle, err := db.gt.Cycle(cycle)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get payout from history bucket")
	}

	var payoutBytes []byte
	if err := db.bolt.View(func(tx *bolt.Tx) error {
		var err error
		if payoutBytes, err = db.getPayout(tx, networkCycle.BlockHash); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed to get payout from history bucket")
	}

	var payout model.Payout
	err = json.Unmarshal(payoutBytes, &payout)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get payout from history bucket")
	}

	return &payout, nil
}

func (db *DB) getPayout(tx *bolt.Tx, blockhash string) ([]byte, error) {
	b := tx.Bucket([]byte(string(historyBucket)))
	if b == nil {
		return nil, errors.New("failed to open history bucket")
	}

	payout := b.Get([]byte(blockhash))
	if payout == nil {
		return nil, fmt.Errorf("failed to get %s from history bucket", blockhash)
	}

	return payout, nil
}
