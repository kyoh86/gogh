package gogh

//go:generate go tool mockgen -source ./infra/github/if.go -destination ./infra/github_mock/gen_mock.go -package github_mock
//go:generate go run ./cmd/gogh man
