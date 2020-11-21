.PHONY: run build format test bench cov clear

run:
	go run main.go

build:
	go build -o bin/main .
	cp config.yaml bin

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
	rm -rf *.db *.out tmp bin