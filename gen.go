package main

//go:generate go run -tags generate ./config/generate

//go:generate go run "github.com/rjeczalik/interfaces/cmd/interfacer" -for github.com/kyoh86/gogh/config.Access -as gogh.Env -o gogh/config.go
//go:generate go run "github.com/golang/mock/mockgen" -source gogh/config.go -destination command/config_mock_test.go -package command_test
//go:generate go run "github.com/golang/mock/mockgen" -source gogh/config.go -destination gogh/config_mock_test.go -package gogh_test

//go:generate go run "github.com/rjeczalik/interfaces/cmd/interfacer" -for github.com/kyoh86/gogh/internal/hub.Client -as command.HubClient -o command/hub.go
//go:generate go run "github.com/golang/mock/mockgen" -source command/hub.go -destination command/hub_mock_test.go -package command_test

//go:generate go run "github.com/rjeczalik/interfaces/cmd/interfacer" -for github.com/kyoh86/gogh/internal/git.Client -as command.GitClient -o command/git.go
//go:generate go run "github.com/golang/mock/mockgen" -source command/git.go -destination command/git_mock_test.go -package command_test
