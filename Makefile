run:
	go run main.go

build:
	go build -o main .

test:
	go test

cov:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

clear:
	go clean
	rm -rf *.db tmp
	rm *.out