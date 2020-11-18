package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/go-playground/validator"
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

// loadSitesInDB : Loads websites from a text file in the given bolt bucket with short codes as keys
// Returns a list of websites it loaded
func (app *application) loadSitesInDB() []string {
	file, err := os.Open("./testdata/sitelst.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lst []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := "https://" + scanner.Text()
		lst = append(lst, url)
		_, err := app.addAShortCode([]byte(url))
		if err != nil {
			log.Fatal(err)
		}
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lst
}

func Benchmark__GetWebsite(b *testing.B) {

}

func Benchmark__SetWebsite(b *testing.B) {

}

func Test__GetWebsite(t *testing.T) {
	app := getApplication()
	defer app.db.Close()

	ts := runServer(app.redirectHandler)
	defer ts.Close()

	url := ts.URL + "/" + string(tstCode)

	res, err := http.Get(url)
	if err != nil {
		t.Errorf("%s", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("HTTPStatusCode | Expected: %v ; Received: %v", http.StatusOK, res.StatusCode)
	}
	finalURL := res.Request.URL.String()
	if finalURL != string(tstLink) {
		t.Errorf("RequestURL | Expected: %s ; Received: %s", tstLink, finalURL)
	}
}

func Test__SetWebsite(t *testing.T) {
	app := getApplication()
	defer app.db.Close()

	ts := runServer(app.setHandler)
	defer ts.Close()

	url := ts.URL

	res, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewReader(tstLink))
	if err != nil {
		t.Errorf("%s", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Errorf("HTTPStatusCode | Expected: %v ; Received: %v", http.StatusCreated, res.StatusCode)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("%s", err)
	}
	validate := validator.New()
	err = validate.Var(string(string(body)), "required,alphanum,len=6")
	if err != nil {
		t.Errorf("ResponseBody | Validation: Alphanumeric 6 Character ; Received: %s", body)
	}
}

func Test__SetInvalidURL(t *testing.T) {

}

func Test__GetInvalidURL(t *testing.T) {

}
