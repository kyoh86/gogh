// +build generate

package main

import (
	"log"

	"github.com/kyoh86/appenv"
	"github.com/kyoh86/gogh/config"
)

func main() {
	if err := appenv.Generate(
		"github.com/kyoh86/gogh/config",
		"./config",
		appenv.Opt(new(config.GithubHost), appenv.StoreYAML(), appenv.StoreEnvar()),
		appenv.Opt(new(config.GithubUser), appenv.StoreYAML(), appenv.StoreEnvar()),
		appenv.Opt(new(config.Roots), appenv.StoreYAML(), appenv.StoreEnvar()),
		appenv.Opt(new(config.Hooks), appenv.StoreYAML(), appenv.StoreEnvar()),
	); err != nil {
		log.Fatalln(err)
	}
}
