package gogh_test

//go:generate interfacer -for github.com/kyoh86/gogh/env.Access -as gogh.Env -o env.go
//go:generate mockgen -source env.go -destination env_mock_test.go -package gogh_test
