package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func runServer(fn func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(fn))
}

func Benchmark__GetWebsite(b *testing.B) {

}

func Benchmark__SetWebsite(b *testing.B) {

}

func Test__GetWebsite(t *testing.T) {

}

func Test__SetWebsite(t *testing.T) {

}

func Test__SetInvalidURL(t *testing.T) {

}

func Test__GetInvalidURL(t *testing.T) {

}
