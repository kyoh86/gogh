.PHONY: test gen lint release

VERSION := `git vertag get`
COMMIT  := `git rev-parse HEAD`

ifeq ($(XDG_CONFIG_HOME),)
	XDG_CONFIG_HOME := $(HOME)/.config
endif

test:
	go test --race ./...

gen:
	go generate ./...

lint:
	gometalinter --config $(XDG_CONFIG_HOME)/gometalinter/config.json ./...

install:
	go install -a -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT)" ./...

man:
	man.sh

seg = patch
release:
	git-vertag $(seg)
	goreleaser --rm-dist
