VERSION ?= `git vertag get`
COMMIT  ?= `git rev-parse HEAD`
DATE    ?= `date --iso-8601`

generate-clear: gen-clear
.PHONY: generate-clear

gen-clear:
	rm -rf ./**/*_gen.go
.PHONY: gen-clear

generate: gen
.PHONY: generate

gen: gen-clear
	go generate -x ./...
.PHONY: gen

clear-sdl:
	rm -f ./infra/githubv4/schema.graphql
.PHONY: clear-sdl

get-sdl:
	curl -Lo ./infra/githubv4/schema.graphql https://docs.github.com/public/fpt/schema.docs.graphql
.PHONY: get-sdl

gen-gql: clear-sdl get-sdl
	go generate -tags gengraphql -x ./infra/githubv4
.PHONY: gen-gql

lint: gen
	golangci-lint run
.PHONY: lint

test: gen
	go test -tags man -v --race ./...
.PHONY: test

man: gen
	rm -rf ./doc/usage/**.md
	go run -tags man -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT) -X=main.date=$(DATE)" ./cmd/gogh man
.PHONY: man

install: test
	go install -a -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT) -X=main.date=$(DATE)" ./cmd/gogh/...

.PHONY: install
