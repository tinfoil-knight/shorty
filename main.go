package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/tinfoil-knight/shorty/config"
	"github.com/tinfoil-knight/shorty/helpers"

	"github.com/boltdb/bolt"
	"github.com/go-playground/validator"
)

type application struct {
	db     *bolt.DB
	bucket []byte
}

func (app *application) getHandler(w http.ResponseWriter, r *http.Request) {
	var validate *validator.Validate
	validate = validator.New()

	if r.Method != "GET" {
		sendError(w, http.StatusMethodNotAllowed)
		return
	}
	shortCode := []byte(r.URL.Path[1:])

	err := validate.Var(string(shortCode), "required,alphanum,len=6")
	if err != nil {
		sendError(w, http.StatusBadRequest)
		return
	}
	err = app.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(app.bucket)
		url := b.Get([]byte(shortCode))
		if url != nil {
			fmt.Fprintf(w, "%s\n", url)
		} else {
			return errors.New("Not Found")
		}
		return nil
	})
	if err != nil {
		if err.Error() == "Not Found" {
			sendError(w, http.StatusNotFound)
		} else {
			sendError(w, http.StatusInternalServerError)
		}
	}
}

func (app *application) setHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendError(w, http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	url := body
	shortCode := helpers.GenerateShortString()
	err = app.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(app.bucket)
		if err != nil {
			return err
		}
		// TODO: Protect against crash/infinite loop by limiting num of iterations
		// TODO: Test this loop
		for {
			chk := b.Get(shortCode)
			if chk == nil {
				break
			}
			shortCode = helpers.GenerateShortString()
		}
		return b.Put(shortCode, url)
	})
	if err != nil {
		sendError(w, http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s\n", shortCode)
}

func sendError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		log.Printf("%s %s %s\n", r.Method, r.URL, time.Since(start).Round(time.Microsecond).String())
	})
}

func main() {
	httpPort := config.Get("PORT")
	boltPath := config.Get("BOLT-PATH")

	db := helpers.InitDB(boltPath)
	defer db.Close()

	app := &application{
		db:     db,
		bucket: []byte("links"),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/new", app.setHandler)
	mux.HandleFunc("/", app.getHandler)

	log.Printf("INFO: Starting server on %s", httpPort)
	log.Fatal(http.ListenAndServe(httpPort, logRequest(mux)))
}
