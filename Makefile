.PHONY: build build-linux build-macos build-windows build-all test test-unit test-e2e fmt

build:
	go build -o bin/pic2video ./cmd/pic2video

build-linux:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/pic2video-linux-amd64 ./cmd/pic2video

build-macos:
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -o bin/pic2video-darwin-amd64 ./cmd/pic2video

build-windows:
	mkdir -p bin
	GOOS=windows GOARCH=amd64 go build -o bin/pic2video-windows-amd64.exe ./cmd/pic2video

build-all: build-linux build-macos build-windows

test: test-unit test-e2e

test-unit:
	go test ./... -run Test -count=1

test-e2e:
	go test ./tests/e2e -count=1

fmt:
	gofmt -w $(shell find . -name '*.go' -not -path './.git/*')
