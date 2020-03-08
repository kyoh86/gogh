// +build generate

package main

import (
	"log"

	"github.com/kyoh86/gogh/appenv/gen"
	"github.com/kyoh86/gogh/appenv/prop"
	"github.com/kyoh86/gogh/env"
)

//go:generate go run -tags generate ./main.go

func main() {
	gen := &gen.Generator{
		EnvarPrefix: "GOGH_",
	}

	if err := gen.Do(
		"github.com/kyoh86/gogh/env",
		"../",
		prop.Prop(new(env.Roots), prop.StoreFile(), prop.StoreEnvar()),
		prop.Prop(new(env.GithubHost), prop.StoreFile(), prop.StoreEnvar()),
		prop.Prop(new(env.GithubToken), prop.StoreKeyring(), prop.StoreEnvar()),
	); err != nil {
		log.Fatal(err)
	}
}
