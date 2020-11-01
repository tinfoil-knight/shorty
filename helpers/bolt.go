package helpers

import (
	"log"

	"github.com/boltdb/bolt"
)

// InitDB connects to the BoltDB and returns a DB client.
func InitDB(path string) *bolt.DB {
	client, err := bolt.Open(path, 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
