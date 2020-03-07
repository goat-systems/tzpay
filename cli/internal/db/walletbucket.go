package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// InitWallet will initialize the storage of the edesk, and password.
// It will store the edesk, and password of the tezos wallet with a salted hash.
// This function call will unset TZPAY_WALLET_SECRET from the enviroment for safety
func (db *DB) InitWallet(password string, secret string) error {
	if err := db.bolt.Update(func(tx *bolt.Tx) error {
		return db.storeSecret(tx, password, secret)
	}); err != nil {
		return errors.Wrap(err, "failed to initialize wallet in tzpay bucket")
	}

	return nil
}

func (db *DB) storeSecret(tx *bolt.Tx, password, secret string) error {
	b := tx.Bucket([]byte(string(walletBucket)))
	if b == nil {
		return errors.New("failed to open wallet bucket")
	}

	encryptedSecret, err := encrypt([]byte(secret), password)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt edesk")
	}

	err = b.Put([]byte(string(walletBucketSecret)), encryptedSecret)
	if err != nil {
		return errors.Wrap(err, "failed to put edesk in wallet bucket")
	}

	return nil
}

// IsWalletInitialized will return an error if no wallet bucket exists, or
// the aes encrypted edesk is missing.
func (db *DB) IsWalletInitialized() bool {
	if err := db.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(string(walletBucket)))
		if b == nil {
			return errors.New("")
		}
		secretAES := b.Get([]byte(string(walletBucketSecret)))
		if secretAES == nil {
			return errors.New("")
		}

		return nil
	}); err != nil {
		return false
	}

	return true
}

// GetSecret returns the decrpted string of the aes encrypted edesk
func (db *DB) GetSecret(password string) (string, error) {
	var secret []byte
	if err := db.bolt.View(func(tx *bolt.Tx) error {
		var err error
		if secret, err = db.getSecret(tx, password); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", errors.Wrap(err, "failed to authorize")
	}

	return string(secret), nil
}

func (db *DB) getSecret(tx *bolt.Tx, password string) ([]byte, error) {
	b := tx.Bucket([]byte(string(walletBucket)))
	if b == nil {
		return nil, errors.New("failed to open wallet bucket")
	}

	secretAES := b.Get([]byte(string(walletBucketSecret)))
	if secretAES == nil {
		return nil, errors.New("failed to get aes encrypted secret")
	}

	secret, err := decrypt(secretAES, password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt secret")
	}

	return secret, nil
}

func encrypt(data []byte, passphrase string) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt data with aes")
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.Wrap(err, "failed to encrypt data with aes")
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt(data []byte, passphrase string) ([]byte, error) {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt aes data")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt aes data")
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt aes data")
	}
	return plaintext, nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
