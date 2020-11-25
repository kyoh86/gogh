package main

//go:generate go run -tags generate ./env/generate

//go:generate go run "github.com/rjeczalik/interfaces/cmd/interfacer" -for github.com/kyoh86/gogh/env.Access -as gogh.Env -o gogh/env.go
//go:generate go run "github.com/golang/mock/mockgen" -source gogh/env.go -destination command/env_mock_test.go -package command_test
//go:generate go run "github.com/golang/mock/mockgen" -source gogh/env.go -destination gogh/env_mock_test.go -package gogh_test

//go:generate go run "github.com/rjeczalik/interfaces/cmd/interfacer" -for github.com/kyoh86/gogh/internal/hub.Client -as command.HubClient -o command/hub.go
//go:generate go run "github.com/golang/mock/mockgen" -source command/hub.go -destination command/hub_mock_test.go -package command_test

//go:generate go run "github.com/rjeczalik/interfaces/cmd/interfacer" -for github.com/kyoh86/gogh/internal/git.Client -as command.GitClient -o command/git.go
//go:generate go run "github.com/golang/mock/mockgen" -source command/git.go -destination command/git_mock_test.go -package command_test
