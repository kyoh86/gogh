// +build generate

package main

import (
	"log"

	"github.com/kyoh86/gogh/appenv/gen"
	"github.com/kyoh86/gogh/env"
)

//go:generate go run -tags generate ./main.go

func main() {
	g := &gen.Generator{}

	if err := g.Do(
		"github.com/kyoh86/gogh/env",
		"../",
		gen.Prop(new(env.Roots), gen.File(), gen.Envar()),
		gen.Prop(new(env.GithubHost), gen.File(), gen.Envar()),
		gen.Prop(new(env.GithubToken), gen.Keyring(), gen.Envar()),
	); err != nil {
		log.Fatal(err)
	}
}
