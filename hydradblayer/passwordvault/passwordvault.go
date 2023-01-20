package passwordvault

import (
	"errors"

	"github.com/boltdb/bolt"
)

var ErrNilDB = errors.New("Database handler is nil...")

func ConnectPasswordVault() (*bolt.DB, error) {
	db, err := bolt.Open("hydra.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func AddToVault(db *bolt.DB, username, password string) error {
	if db == nil {
		return ErrNilDB
	}
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("PasswordVault"))
		if err != nil {
			return err
		}
		err = b.Put([]byte(username), []byte(password))
		return err
	})
}

func AddBytesToVault(db *bolt.DB, username string, password []byte) error {
	if db == nil {
		return ErrNilDB
	}
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("PasswordVault"))
		if err != nil {
			return err
		}
		err = b.Put([]byte(username), password)
		return err
	})
}

func GetPassword(db *bolt.DB, username string) (string, error) {
	if db == nil {
		return "", ErrNilDB
	}
	password := ""
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PasswordVault"))
		if b == nil {
			return errors.New("Could not find password vault bucket!!")
		}
		v := b.Get([]byte(username))
		password = string(v)
		return nil
	})
	return password, err
}

func GetPasswordBytes(db *bolt.DB, username string) ([]byte, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	password := []byte{}
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PasswordVault"))
		if b == nil {
			return errors.New("Could not find password vault bucket!!")
		}
		password = b.Get([]byte(username))
		return nil
	})
	return password, err
}
