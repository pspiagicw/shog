all: build

build:
	./scripts/fetch-project-readmes.sh
	env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod-cache go build .

format:
	gofmt -w app/*.go main.go

fetch-projects:
	./scripts/fetch-project-readmes.sh

.PHONY: build test format fetch-projects all
