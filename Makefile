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

lint: gen
	golangci-lint run
.PHONY: lint

test: gen
	go test -tags man -v --race ./...
.PHONY: test

man: gen
	go run -tags man -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT) -X=main.date=$(DATE)" ./cmd/gogh man
.PHONY: man

ver:
	echo "VERSION=\"$(VERSION)\"" > .version
	echo "COMMIT=\"$(COMMIT)\"" >> .version

pkg:
	include .version
	export
	:
	envsubst '$$VERSION $$COMMIT' < $< > $@
	updpkgsums
	makepkg -f
	namcap
.PHONY: pkg

install: test
	go install -a -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT) -X=main.date=$(DATE)" ./cmd/gogh/...

.PHONY: install
