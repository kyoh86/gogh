package gogh

//go:generate go run github.com/golang/mock/mockgen -source ./internal/github/if.go -destination ./internal/github_mock/gen_mock.go -package github_mock
//go:generate go run -tags man ./cmd/gogh man
