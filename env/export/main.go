// +build generate

package main

import (
	"log"

	"github.com/kyoh86/gogh/env"
	"github.com/kyoh86/gogh/env/generate"
	"github.com/kyoh86/gogh/env/props"
)

//go:generate go run -tags generate ./main.go

func main() {
	gen := &generate.Generator{
		EnvarPrefix: "GOGH_",
	}

	if err := gen.Do(
		"github.com/kyoh86/gogh/env",
		"../",
		props.Prop(new(env.Roots), props.StoreConfig(), props.StoreEnvar()),
		props.Prop(new(env.GithubHost), props.StoreConfig(), props.StoreEnvar()),
		props.Prop(new(env.GithubToken), props.StoreKeyring(), props.StoreEnvar()),
		props.Prop(new(env.GithubUser), props.StoreCache()),
	); err != nil {
		log.Fatal(err)
	}
}
