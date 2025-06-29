package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.etcd.io/bbolt"
)

var (
	defaultDir    = ".wallet_data"
	jsonPath      = filepath.Join(defaultDir, "wallets.json")
	boltDBPath    = filepath.Join(defaultDir, "wallets.db")
	walletsBucket = []byte("wallets")
	fileLock      sync.Mutex
)

type StoredWallet struct {
	Name      string `json:"name"`
	Encrypted string `json:"encrypted"`
	CreatedAt int64  `json:"created_at"`
}

func ensureDir() error {
	if _, err := os.Stat(defaultDir); os.IsNotExist(err) {
		return os.MkdirAll(defaultDir, 0700)
	}
	return nil
}

func SaveWalletToJSON(wallet StoredWallet) error {
	fileLock.Lock()
	defer fileLock.Unlock()
	if err := ensureDir(); err != nil {
		return err
	}

	var wallets []StoredWallet
	_ = LoadWalletsFromJSON(&wallets)
	wallets = append(wallets, wallet)

	f, err := os.Create(jsonPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(wallets)
}

func LoadWalletsFromJSON(dest *[]StoredWallet) error {
	f, err := os.Open(jsonPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("close file error")
		}
	}(f)
	return json.NewDecoder(f).Decode(dest)
}

func SaveWalletToBolt(wallet StoredWallet) error {
	if err := ensureDir(); err != nil {
		return err
	}
	db, err := bbolt.Open(boltDBPath, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(walletsBucket)
		if err != nil {
			return err
		}
		data, err := json.Marshal(wallet)
		if err != nil {
			return err
		}
		return b.Put([]byte(wallet.Name), data)
	})
}

func LoadWalletFromBolt(name string) (*StoredWallet, error) {
	db, err := bbolt.Open(boltDBPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var w StoredWallet
	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(walletsBucket)
		if b == nil {
			return errors.New("bucket not found")
		}
		data := b.Get([]byte(name))
		if data == nil {
			return fmt.Errorf("wallet '%s' not found", name)
		}
		return json.Unmarshal(data, &w)
	})

	if err != nil {
		return nil, err
	}
	return &w, nil
}
