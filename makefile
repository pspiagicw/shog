SHOG_PORT ?= 23234
export SHOG_PORT

all: build

build:
	./scripts/fetch-project-readmes.sh
	env GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod-cache go build .

run: build
	-@fuser -k $(SHOG_PORT)/tcp >/dev/null 2>&1 && sleep 1
	./shog

format:
	gofmt -w app/*.go main.go

fetch-projects:
	./scripts/fetch-project-readmes.sh

.PHONY: build run test format fetch-projects all
