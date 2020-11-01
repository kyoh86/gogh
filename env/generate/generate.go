// +build generate

package main

import (
	"log"

	"github.com/kyoh86/appenv"
	"github.com/kyoh86/gogh/env"
)

//go:generate go run -tags generate .

func main() {
	if err := appenv.Generate(
		"github.com/kyoh86/gogh/env",
		"../",
		appenv.Opt(new(env.GithubHost), appenv.StoreYAML(), appenv.StoreEnvar()),
		appenv.Opt(new(env.GithubUser), appenv.StoreYAML(), appenv.StoreEnvar()),
		appenv.Opt(new(env.Roots), appenv.StoreYAML(), appenv.StoreEnvar()),
		appenv.Opt(new(env.Hooks), appenv.StoreYAML(), appenv.StoreEnvar()),
	); err != nil {
		log.Fatalln(err)
	}
}
