.PHONY: gen lint test install man e2e

VERSION := `git vertag get`
COMMIT  := `git rev-parse HEAD`

gen:
	go generate ./...

lint: gen
	golangci-lint run

test: lint
	go test v --race ./...

install: test
	go install -a -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT)" ./...

man: test
	go run main.go --help-man > gogh.1

e2e:
	circleci local execute --job e2e --volume $(PWD):/go/src/github.com/kyoh86/gogh

