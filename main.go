package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/tinfoil-knight/shorty/config"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	var validate *validator.Validate
	validate = validator.New()
	if r.URL.Path == "/" {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}

	} else {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
		shortCode := []byte(r.URL.Path[1:])

		err := validate.Var(string(shortCode), "required,alphanum,len=6")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "%s\n", shortCode)
	}
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", apiHandler)
	log.Printf("INFO: Starting server on %s", httpPort)
	log.Fatal(http.ListenAndServe(httpPort, logRequest(mux)))
}
