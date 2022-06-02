## 
GO111MODULE=on
## image registry host
REGISTRY?=docker.pkg.github.com
## git commit
COMMIT?=$(shell git rev-parse HEAD)
## build date
DATE?=$(shell date "+%Y/%m/%dT%H:%M:%S")
## app name
APP_NAME?=$(shell basename ${PWD})

## build: go build the application
.PHONY: build
build: clean
	go build -o api .

## run: run the application server arg: PORT 
.PHONY: run
run:
	go run main.go server -p ${PORT}

## clean: golang clena
.PHONY: clean
clean:
	go clean

## test: go test the application
.PHONY: test
test:
  go test -v -count=1 ./...

## help: help range
.PHONY: help
help:
	@echo "Useage:\n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

