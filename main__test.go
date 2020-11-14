package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/tinfoil-knight/shorty/config"
	"github.com/tinfoil-knight/shorty/helpers"
)

func runServer(fn func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(fn))
}

func getApplication() *application {
	boltPath := config.Get("TEST-BOLT-PATH")
	db := helpers.InitDB(boltPath)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte("links"))
		if err != nil {
			if err.Error() == "bucket not found" {
				fmt.Println("No test bucket found at", boltPath)
			} else {
				return fmt.Errorf("delete bucket: %s", err)
			}

		}
		b, err := tx.CreateBucket([]byte("links"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		err = b.Put([]byte("1asdUv"), []byte("https://example.com"))
		if err != nil {
			return fmt.Errorf("store key in bucket: %s", err)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		db:     db,
		bucket: []byte("links"),
	}
	return app
}

func initDB() {

}

func Benchmark__GetWebsite(b *testing.B) {

}

func Benchmark__SetWebsite(b *testing.B) {

}

func Test__GetWebsite(t *testing.T) {
	app := getApplication()
	defer app.db.Close()

	ts := runServer(app.getHandler)

	ts.Close()
}

func Test__SetWebsite(t *testing.T) {

}

func Test__SetInvalidURL(t *testing.T) {

}

func Test__GetInvalidURL(t *testing.T) {

}
