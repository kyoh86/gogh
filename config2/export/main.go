// +build generate

package main

import (
	"log"

	config "github.com/kyoh86/gogh/config2"
	"github.com/kyoh86/gogh/config2/generate"
)

//go:generate go run -tags generate ./main.go

func main() {
	gen := &generate.Generator{
		PackageName: "config",
		EnvarPrefix: "GOGH_",
	}

	if err := gen.Do(
		"github.com/kyoh86/gogh/config2",
		"../",
		new(config.GithubToken),
		new(config.GithubHost),
		new(config.GithubUser),
		new(config.Roots),
	); err != nil {
		log.Fatal(err)
	}
}
