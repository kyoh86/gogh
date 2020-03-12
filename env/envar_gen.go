// Code generated by main.go DO NOT EDIT.

package env

import (
	gostrcase "github.com/stoewer/go-strcase"
	"log"
	"os"
)

type Envar struct {
	Roots       *Roots
	GithubHost  *GithubHost
	GithubToken *GithubToken
}

func getEnvar(prefix string) (envar Envar, err error) {
	prefix = gostrcase.UpperSnakeCase(prefix)
	{
		v := os.Getenv(prefix + "ROOTS")
		if v == "" {
			log.Printf("info: there's no envar %sROOTS (%v)", prefix, err)
		} else {
			var value Roots
			if err = value.UnmarshalText([]byte(v)); err != nil {
				return envar, err
			}
			envar.Roots = &value
		}
	}
	{
		v := os.Getenv(prefix + "GITHUB_HOST")
		if v == "" {
			log.Printf("info: there's no envar %sGITHUB_HOST (%v)", prefix, err)
		} else {
			var value GithubHost
			if err = value.UnmarshalText([]byte(v)); err != nil {
				return envar, err
			}
			envar.GithubHost = &value
		}
	}
	{
		v := os.Getenv(prefix + "GITHUB_TOKEN")
		if v == "" {
			log.Printf("info: there's no envar %sGITHUB_TOKEN (%v)", prefix, err)
		} else {
			var value GithubToken
			if err = value.UnmarshalText([]byte(v)); err != nil {
				return envar, err
			}
			envar.GithubToken = &value
		}
	}
	return
}