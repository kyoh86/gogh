.PHONY: gen lint test install man

VERSION := `git vertag get`
COMMIT  := `git rev-parse HEAD`

gen-clear:
	rm ./**/*_gen.go

gen:
	go generate -x ./...
	interfacer -for github.com/kyoh86/gogh/env.Access -as gogh.Env -o gogh/env.go
	mockgen -source gogh/env.go -destination command/env_mock_test.go -package command_test
	mockgen -source gogh/env.go -destination gogh/env_mock_test.go -package gogh_test
	
	interfacer -for github.com/kyoh86/gogh/internal/hub.Client -as command.HubClient -o command/hub.go
	mockgen -source command/hub.go -destination command/hub_mock_test.go -package command_test
	
	interfacer -for github.com/kyoh86/gogh/internal/git.Client -as command.GitClient -o command/git.go
	mockgen -source command/git.go -destination command/git_mock_test.go -package command_test

lint: gen
	golangci-lint run

test: lint
	go test -tags mock -v --race ./...

install: test
	go install -a -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT)" ./...

man:
	go run . --help-man > gogh.1
