package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/tinfoil-knight/shorty/config"
	"github.com/tinfoil-knight/shorty/helpers"
)

var bktName = []byte("links")
var tstLink = []byte("https://example.com")
var tstCode = []byte("1asdUv")

func runServer(fn func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(fn))
}

func getApplication() *application {
	boltPath := config.Get("TEST-BOLT-PATH")
	db := helpers.InitDB(boltPath)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bktName)
		if err != nil {
			if err.Error() == "bucket not found" {
				fmt.Println("No test bucket found at", boltPath)
			} else {
				return fmt.Errorf("delete bucket: %s", err)
			}
		}
		b, err := tx.CreateBucket(bktName)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		err = b.Put(tstCode, tstLink)
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
		bucket: bktName,
	}
	return app
}

func Benchmark__GetWebsite(b *testing.B) {

}

func Benchmark__SetWebsite(b *testing.B) {

}

func Test__GetWebsite(t *testing.T) {
	app := getApplication()
	defer app.db.Close()

	ts := runServer(app.getHandler)
	defer ts.Close()

	url := ts.URL + "/" + string(tstCode)

	res, err := http.Get(url)
	if err != nil {
		t.Errorf("%s", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("HTTPStatusCode | Expected: %v, Received: %v", http.StatusOK, res.StatusCode)
	}

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("%s", err)
	}

	if string(bodyBytes) != string(tstLink) {
		t.Errorf("ResponseBody | Expected: %v, Received: %v", string(bodyBytes), string(tstLink))
	}
}

func Test__SetWebsite(t *testing.T) {

}

func Test__SetInvalidURL(t *testing.T) {

}

func Test__GetInvalidURL(t *testing.T) {

}
