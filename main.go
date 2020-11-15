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

func (app *application) redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		sendError(w, http.StatusMethodNotAllowed)
		return
	}
	var validate *validator.Validate
	validate = validator.New()

	shortCode := []byte(r.URL.Path[1:])

	err := validate.Var(string(shortCode), "required,alphanum,len=6")
	if err != nil {
		sendError(w, http.StatusBadRequest)
		return
	}
	url, err := app.findOne(shortCode)
	if err != nil {
		if err.Error() == "Not Found" {
			sendError(w, http.StatusNotFound)
		} else {
			sendError(w, http.StatusInternalServerError)
		}
		return
	}
	http.Redirect(w, r, string(url), 302)
}

func (app *application) setHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendError(w, http.StatusMethodNotAllowed)
		return
	}
	var validate *validator.Validate
	validate = validator.New()

	url, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendError(w, http.StatusInternalServerError)
		return
	}
	err = validate.Var(string(url), "url")
	if err != nil {
		sendError(w, http.StatusBadRequest)
		return
	}
	shortCode, err := app.addAShortCode(url)
	if err != nil {
		sendError(w, http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", shortCode)
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

func (app *application) findOne(key []byte) ([]byte, error) {
	var val []byte
	err := app.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(app.bucket)
		val = b.Get([]byte(key))
		if val == nil {
			return errors.New("Not Found")
		}
		return nil
	})
	return val, err
}

func (app *application) addOne(key []byte, val []byte) error {
	err := app.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(app.bucket)
		if err != nil {
			return err
		}
		return b.Put(key, val)
	})
	return err
}

func (app *application) addAShortCode(url []byte) ([]byte, error) {
	shortCode := helpers.GenerateShortString()
	err := app.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(app.bucket)
		if err != nil {
			return err
		}

		for i := 0; i < 5; i++ {
			if i == 4 {
				return errors.New("not able to generate unique code")
			}
			chk := b.Get(shortCode)
			if chk == nil {
				break
			}
			shortCode = helpers.GenerateShortString()
		}
		return b.Put(shortCode, url)
	})
	return shortCode, err
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

	mux.HandleFunc("/app/new", app.setHandler)
	mux.HandleFunc("/", app.redirectHandler)

	log.Printf("INFO: Starting server on %s", httpPort)
	log.Fatal(http.ListenAndServe(httpPort, logRequest(mux)))
}
