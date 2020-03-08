// Code generated by main.go DO NOT EDIT.

package env

import (
	"log"
	"os"
)

type Envar struct {
	Roots       *Roots
	GithubHost  *GithubHost
	GithubToken *GithubToken
}

func GetEnvar() (envar Envar, err error) {
	{
		v := os.Getenv("GOGH_ROOTS")
		if v == "" {
			log.Printf("info: there's no envar GOGH_ROOTS (%v)", err)
		} else {
			var value Roots
			if err = value.UnmarshalText([]byte(v)); err != nil {
				return envar, err
			}
			envar.Roots = &value
		}
	}
	{
		v := os.Getenv("GOGH_GITHUB_HOST")
		if v == "" {
			log.Printf("info: there's no envar GOGH_GITHUB_HOST (%v)", err)
		} else {
			var value GithubHost
			if err = value.UnmarshalText([]byte(v)); err != nil {
				return envar, err
			}
			envar.GithubHost = &value
		}
	}
	{
		v := os.Getenv("GOGH_GITHUB_TOKEN")
		if v == "" {
			log.Printf("info: there's no envar GOGH_GITHUB_TOKEN (%v)", err)
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
