package command

import (
	"fmt"
	"strings"

	"github.com/kyoh86/gogh/env"
	keyring "github.com/zalando/go-keyring"
)

func ConfigGetAll(cfg *env.Config) error {
	for _, name := range env.PropertyNames() {
		opt, _ := cfg.Property(name) // ignore error: config.OptionNames covers all accessor
		value, err := opt.Get()
		if err != nil {
			return err
		}
		fmt.Printf("%s: %s\n", name, value)
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

func ConfigSet(cfg *env.Config, optionName, optionValue string) error {
	if optionName == "github.token" {
		hostCfg, err := cfg.Property("github.host")
		if err != nil {
			return err
		}
		host, err := hostCfg.Get()
		if err != nil {
			return err
		}
		userCfg, err := cfg.Property("github.user")
		if err != nil {
			return err
		}
		user, err := userCfg.Get()
		if err != nil {
			return err
		}
		keyring.Set(strings.Join([]string{host, env.KeyringService}, "."), user, optionValue)
	}

	opt, err := cfg.Property(optionName)
	if err != nil {
		return err
	}
	return opt.Set(optionValue)
}

func ConfigUnset(cfg *env.Config, optionName string) error {
	opt, err := cfg.Property(optionName)
	if err != nil {
		return err
	}
	opt.Unset()
	return nil
}
