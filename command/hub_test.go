package command_test

//go:generate interfacer -for github.com/kyoh86/gogh/internal/hub.HubClient -as command.HubClient -o hub.go
//go:generate mockgen -source hub.go -destination hub_mock_test.go -package command_test
