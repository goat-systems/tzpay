package db

import (
	"fmt"
	"os"
	"os/user"

	"github.com/boltdb/bolt"
	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/pkg/errors"
)

type bucketKey string

const (
	walletBucket       bucketKey = "WALLETBUCKET"
	historyBucket      bucketKey = "HISTORYBUCKET"
	walletBucketSecret bucketKey = "edesk"
)

// DB wraps bolt db functions
type DB struct {
	bolt *bolt.DB
	gt   gotezos.IFace
}

// New will open or create the tzpay boltdb profile
func New(gt gotezos.IFace, path string) (*DB, error) {

	if path == "" {
		usr, err := user.Current()
		if err != nil {
			return &DB{}, errors.Wrap(err, "failed to open boltdb")
		}

		if _, err := os.Stat(fmt.Sprintf("%s/.tzpay", usr.HomeDir)); os.IsNotExist(err) {
			err = os.Mkdir(fmt.Sprintf("%s/.tzpay", usr.HomeDir), 0755)
			if err != nil {
				return &DB{}, errors.Wrap(err, "failed to open boltdb")
			}
		}

		path = fmt.Sprintf("%s/.tzpay/tzpay.db", usr.HomeDir)
	}
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return &DB{}, errors.Wrap(err, "failed to open boltdb")
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(walletBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(historyBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return &DB{
		bolt: db,
		gt:   gt,
	}, nil
}
