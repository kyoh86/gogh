package command_test

//go:generate interfacer -for github.com/kyoh86/gogh/internal/git.Client -as command.GitClient -o git.go
//go:generate mockgen -source git.go -destination git_mock_test.go -package command_test
