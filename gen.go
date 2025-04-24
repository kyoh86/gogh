package gogh

//go:generate go run go.uber.org/mock/mockgen -source ./infra/github/if.go -destination ./infra/github_mock/gen_mock.go -package github_mock
//go:generate go run -tags man ./cmd/gogh man
