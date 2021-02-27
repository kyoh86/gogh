VERSION := `git vertag get`
COMMIT  := `git rev-parse HEAD`
DATE    := `date --iso-8601`

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

lint: gen
	golangci-lint run
.PHONY: lint

test: lint
	go test -tags mock -v --race ./...
.PHONY: test

install: test
	go install -a -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT) -X=main.date=$(DATE)" ./...
.PHONY: install
