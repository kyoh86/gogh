package command

import (
	"fmt"

	"github.com/kyoh86/gogh/env"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/internal/hub"
)

var TokenManager = hub.NewKeyring

func ConfigGetAll(cfg *env.Config) error {
	for _, name := range env.OptionNames() {
		opt, _ := cfg.Option(name) // ignore error: config.OptionNames covers all accessor
		value, err := opt.Get()
		if err != nil {
			return err
		}
		if value == "" {
			// NOTE: to avoid a bug in the example test...
			// https://github.com/golang/go/issues/26460
			fmt.Printf("%s:\n", name)
		} else {
			fmt.Printf("%s: %s\n", name, value)
		}
	}
	fmt.Println("github.token: *****")
	return nil
}

func ConfigGet(cfg *env.Config, optionName string) error {
	opt, err := cfg.Option(optionName)
	if err != nil {
		return err
	}
	value, err := opt.Get()
	if err != nil {
		return err
	}
	fmt.Println(value)
	return nil
}

func ConfigSet(ev gogh.Env, cfg *env.Config, optionName, optionValue string) error {
	if optionName == "github.token" {
		tm, err := TokenManager(ev.GithubHost())
		if err != nil {
			return err
		}
		return tm.SetGithubToken(ev.GithubUser(), optionValue)
	}

	opt, err := cfg.Option(optionName)
	if err != nil {
		return err
	}
	return opt.Set(optionValue)
}

func ConfigUnset(ev gogh.Env, cfg *env.Config, optionName string) error {
	if optionName == "github.token" {
		host, user := ev.GithubHost(), ev.GithubUser()

		tm, err := TokenManager(host)
		if err != nil {
			return err
		}
		if err := tm.DeleteGithubToken(user); err != nil {
			return err
		}
		return nil
	}

	opt, err := cfg.Option(optionName)
	if err != nil {
		return err
	}
	opt.Unset()
	return nil
}
