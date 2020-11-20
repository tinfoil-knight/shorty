.PHONY: test

run:
	go run main.go

build:
	go build -o bin/main .

format:
	gofmt -d -e

test:
	go test -v

bench:
	go test -run=XXX -bench=. -benchtime 100x

cov:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

clear:
	go clean
	rm -rf *.db tmp
	rm *.out