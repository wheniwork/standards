default: test

build:
	go build

test: build
	go test -v ./...