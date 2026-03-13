all: build 

build:
	go build .

format:
	go format ./...
.PHONY: build test format all

