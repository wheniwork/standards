default: test

build:
	go build

test: build
	go get github.com/kr/pretty
	go test -v ./...