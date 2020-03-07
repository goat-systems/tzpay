package db

import (
	"fmt"
	"os/user"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

type bucketKey string

const (
	walletBucket       bucketKey = "WALLETBUCKET"
	walletBucketSecret bucketKey = "edesk"
)

// DB wraps bolt db functions
type DB struct {
	bolt   *bolt.DB
	bucket *bolt.Bucket
}

// Open will open or create the tzpay boltdb profile
func Open(path string) (*DB, error) {
	if path == "" {
		usr, err := user.Current()
		if err != nil {
			return &DB{}, errors.Wrap(err, "failed to open boltdb")
		}
		path = fmt.Sprintf("%s.tzpay/tzpay.db", usr)
	}
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return &DB{}, errors.Wrap(err, "failed to open boltdb")
	}
	var b *bolt.Bucket
	err = db.Update(func(tx *bolt.Tx) error {
		b, err = tx.CreateBucketIfNotExists([]byte(walletBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return &DB{
		bolt:   db,
		bucket: b,
	}, nil
}
