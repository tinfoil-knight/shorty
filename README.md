# Shorty
[![Go Report Card](https://goreportcard.com/badge/github.com/tinfoil-knight/shorty)](https://goreportcard.com/report/github.com/tinfoil-knight/shorty)

Shorty is a simple URL shortening service. It uses a key-value store and has no user management itself but is built to integrate well into other services. Inspired from seeing microlinks (eg: kcd.im) used by [Kent C Dodds](https://kentcdodds.com/) all over his blog.

## Built With
- [Go](https://golang.org/)
- [BoltDB](https://github.com/boltdb/bolt)

## Getting Started

If you don't have `make` installed, you can just refer the `Makefile` and run the specified scripts.

### Development
- Run the server: `go run main.go`
- Running Tests: `make test`

### Usage
- Build the app from source: `make build`
- Place a `config.yaml` file in the same folder as the binary with the following variables:
  ```
  PORT: :<port you want to listen and serve from>
  BOLT-PATH: 'path of the bolt database file with extension .db'
  ```

## Author
Kunal Kundu [@tinfoil-knight](https://github.com/tinfoil-knight)

## Acknowledgements
[Educative](https://www.educative.io/) for their article: [Designing a URL Shortening service](https://www.educative.io/courses/grokking-the-system-design-interview/m2ygV4E81AR)

[OpenDNS](https://github.com/opendns) for their list of random domains. [public-domains-list](https://github.com/opendns/public-domain-lists/blob/master/opendns-random-domains.txt)