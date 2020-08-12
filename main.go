package main

import (
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/tinfoil-knight/shorty/config"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	} else {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
		shortCode := []byte(r.URL.Path[1:])
		// RegEx: Alphanumeric string of length 6
		matched, err := regexp.Match("^[a-zA-Z0-9]{6}$", shortCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if matched != true {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
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
	http.ListenAndServe(httpPort, logRequest(mux))
}
