.PHONY: test

run:
	go run main.go

build:
	go build -o bin/main .

test:
	go test -v

cov:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

clear:
	go clean
	rm -rf *.db tmp
	rm *.out