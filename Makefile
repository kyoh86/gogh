VERSION ?= `git vertag get`
COMMIT  ?= `git rev-parse HEAD`
DATE    ?= `date --iso-8601`

.PHONY: clean
clean:
	$(MAKE) -C ./core clean
	$(MAKE) -C ./infra/githubv4 clean

# Alias for gen
.PHONY: generate
generate: gen

.PHONY: gen
gen:
	$(MAKE) -C ./infra/githubv4
	$(MAKE) -C ./core

lint: gen
	go tool golangci-lint run
	go tool arch-go
.PHONY: lint

test: gen
	go test -v --race ./...
.PHONY: test

man: gen
	rm -rf ./doc/usage/**.md
	rm -rf ./doc/man/*
	GOGH_FLAG_PATH=./dummy.yaml go run -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT) -X=main.date=$(DATE)" ./cmd/gogh man
.PHONY: man

install: test
	go install -a -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT) -X=main.date=$(DATE)" ./cmd/gogh/...
.PHONY: install

default: lint test
.DEFAULT_GOAL := default
