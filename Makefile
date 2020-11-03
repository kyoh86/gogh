VERSION := `git vertag get`
COMMIT  := `git rev-parse HEAD`

generate-clear: gen-clear
.PHONY: generate-clear

gen-clear:
	rm ./**/*_gen.go
.PHONY: gen-clear

generate: gen
.PHONY: generate

gen: gen-clear
	go generate -x ./...
	go run "github.com/rjeczalik/interfaces/cmd/interfacer" -for github.com/kyoh86/gogh/env.Access -as gogh.Env -o gogh/env.go
	go run "github.com/golang/mock/mockgen" -source gogh/env.go -destination command/env_mock_test.go -package command_test
	go run "github.com/golang/mock/mockgen" -source gogh/env.go -destination gogh/env_mock_test.go -package gogh_test
	
	go run "github.com/rjeczalik/interfaces/cmd/interfacer" -for github.com/kyoh86/gogh/internal/hub.Client -as command.HubClient -o command/hub.go
	go run "github.com/golang/mock/mockgen" -source command/hub.go -destination command/hub_mock_test.go -package command_test
	
	go run "github.com/rjeczalik/interfaces/cmd/interfacer" -for github.com/kyoh86/gogh/internal/git.Client -as command.GitClient -o command/git.go
	go run "github.com/golang/mock/mockgen" -source command/git.go -destination command/git_mock_test.go -package command_test
.PHONY: gen

lint: gen
	golangci-lint run
.PHONY: lint

test: lint
	go test -tags mock -v --race ./...
.PHONY: test

install: test
	go install -a -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT)" ./...
.PHONY: install

man:
	go run . --help-man > gogh.1
.PHONY: man
