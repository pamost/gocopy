# Basic go commands
GOCMD=go
GOGET=$(GOCMD) get
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

all: gofmt test linter build

get:
	go get -u github.com/vbauerster/mpb/v4

install:
	go build -o $(GOPATH)/bin/gocopy

gofmt:
	gofmt -w .

test:
	go test -count=1 -race -cover -v ./...

linter:
	golangci-lint run --enable-all

build:
	go build -o gocopy

clean:
	go clean -i ./...