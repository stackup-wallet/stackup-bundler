package testutils

import (
	"log"

	badger "github.com/dgraph-io/badger/v3"
)

func DBMock() *badger.DB {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLoggingLevel(badger.ERROR))
	if err != nil {
		log.Fatal(err)
	}

	return db
}
