FROM golang:1.15-alpine AS builder
WORKDIR /build
# CGO_ENABLED: Decides whether the resulting binary from being linked to any C libs.
ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64
# Copy go.mod, go.sum and download deps
COPY go.mod .
COPY go.sum .
RUN go mod download
# Copy the source code
COPY . .
# Build the application
RUN go build -o bin/main .

# Build a small image
FROM scratch

COPY --from=builder /build/main .
EXPOSE 8080
# Command to run
ENTRYPOINT ["/main"]
# Ref: https://levelup.gitconnected.com/complete-guide-to-create-docker-container-for-your-golang-application-80f3fb59a15e