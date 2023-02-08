package start

import (
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

func runDBGarbageCollection(db *badger.DB) {
	go func(db *badger.DB) {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
		again:
			err := db.RunValueLogGC(0.7)
			if err == nil {
				goto again
			}
		}
	}(db)
}
