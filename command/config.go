package command

import (
	"fmt"

	"github.com/kyoh86/gogh/env"
)

func ConfigGetAll(cfg *env.Config) error {
	for _, name := range env.ConfigNames() {
		opt, _ := cfg.Property(name) // ignore error: config.OptionNames covers all accessor
		value, err := opt.Get()
		if err != nil {
			return err
		}
		fmt.Printf("%s: %s\n", name, value)
	}
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
