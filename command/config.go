package command

import (
	"fmt"
	"strings"

	"github.com/kyoh86/gogh/env"
	"github.com/kyoh86/gogh/gogh"
	keyring "github.com/zalando/go-keyring"
)

func ConfigGetAll(cfg *env.Config) error {
	for _, name := range env.PropertyNames() {
		opt, _ := cfg.Property(name) // ignore error: config.OptionNames covers all accessor
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
	opt, err := cfg.Property(optionName)
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
		host, user := ev.GithubHost(), ev.GithubUser()
		if err := keyring.Set(strings.Join([]string{host, env.KeyringService}, "."), user, optionValue); err != nil {
			return err
		}
		return nil
	}

	opt, err := cfg.Property(optionName)
	if err != nil {
		return err
	}
	return opt.Set(optionValue)
}

func ConfigUnset(ev gogh.Env, cfg *env.Config, optionName string) error {
	if optionName == "github.token" {
		host, user := ev.GithubHost(), ev.GithubUser()

		if err := keyring.Delete(strings.Join([]string{host, env.KeyringService}, "."), user); err != nil {
			return err
		}
		return nil
	}

	opt, err := cfg.Property(optionName)
	if err != nil {
		return err
	}
	opt.Unset()
	return nil
}
