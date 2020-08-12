package helpers

import (
	"math/rand"
	"time"
)

// GenerateShortString : Creates a short string of alphanumeric characters.
func GenerateShortString() string {
	const base62Corpus = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	var short []byte
	short = make([]byte, 6, 6)
	for i := 0; i < 6; i++ {
		short[i] = base62Corpus[rand.Intn(len(base62Corpus))]
	}
	return string(short)
}
