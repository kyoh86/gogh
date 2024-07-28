VERSION ?= `git vertag get`
COMMIT  ?= `git rev-parse HEAD`
DATE    ?= `date --iso-8601`

generate-clear: gen-clear
.PHONY: generate-clear

gen-clear: clear-sdl
	rm -rf ./**/*_gen.go
.PHONY: gen-clear

clear-sdl:
	rm -f ./internal/githubv4/schema.graphql
.PHONY: clear-sdl

get-sdl:
	curl -Lo ./internal/githubv4/schema.graphql https://docs.github.com/public/fpt/schema.docs.graphql
.PHONY: get-sdl

generate: gen
.PHONY: generate

gen: gen-clear get-sdl
	go generate -x ./...
.PHONY: gen

lint: gen
	golangci-lint run
.PHONY: lint

test: gen
	go test -tags man -v --race ./...
.PHONY: test

man: gen
	rm -rf ./usage/**.md
	go run -tags man -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT) -X=main.date=$(DATE)" ./cmd/gogh man
.PHONY: man

install: test
	go install -a -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT) -X=main.date=$(DATE)" ./cmd/gogh/...

.PHONY: install
