package checks

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/dbutils"
)

var (
	keyPrefix        = dbutils.JoinValues("checks")
	codeHashesPrefix = dbutils.JoinValues(keyPrefix, "codeHashes")
)

func getCodeHashesKey(userOpHash common.Hash) []byte {
	return []byte(dbutils.JoinValues(codeHashesPrefix, userOpHash.String()))
}

func saveCodeHashes(db *badger.DB, userOpHash common.Hash, codeHashes []codeHash) error {
	return db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(codeHashes)
		if err != nil {
			return err
		}

		return txn.Set(getCodeHashesKey(userOpHash), data)
	})
}

func getSavedCodeHashes(db *badger.DB, userOpHash common.Hash) ([]codeHash, error) {
	var ch []codeHash
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(getCodeHashesKey(userOpHash))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &ch)
		})
	})

	return ch, err
}

func removeSavedCodeHashes(db *badger.DB, userOpHashes ...common.Hash) error {
	return db.Update(func(txn *badger.Txn) error {
		for _, userOpHash := range userOpHashes {
			if err := txn.Delete(getCodeHashesKey(userOpHash)); err != nil {
				return err
			}
		}
		return nil
	})
}
