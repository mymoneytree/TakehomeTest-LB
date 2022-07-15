PROJECT = decode-server-mac
BUILD_TIMESTAMP = $(shell date +%s)
COMMIT = $(shell git rev-parse HEAD)

build:
	GOPATH=$(GOPATH) go build -ldflags "-X main.SHA1=$(COMMIT) -X main.BUILD_TIMESTAMP=$(BUILD_TIMESTAMP)" -o $(PROJECT) .

fmt:
	go fmt ./...
