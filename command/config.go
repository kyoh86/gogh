package command

import (
	"fmt"

	"github.com/kyoh86/gogh/config"
)

func ConfigGetAll(cfg *config.Config) error {
	for _, name := range config.OptionNames() {
		opt, _ := config.Option(name) // ignore error: config.OptionNames covers all accessor
		value := opt.Get(cfg)
		fmt.Printf("%s: %s\n", name, value)
	}
	return nil
}

func ConfigGet(cfg *config.Config, optionName string) error {
	opt, err := config.Option(optionName)
	if err != nil {
		return err
	}
	value := opt.Get(cfg)
	fmt.Fprintln(cfg.Stdout(), value)
	return nil
}

func ConfigPut(cfg *config.Config, optionName, optionValue string) error {
	opt, err := config.Option(optionName)
	if err != nil {
		return err
	}
	return opt.Put(cfg, optionValue)
}

func ConfigUnset(cfg *config.Config, optionName string) error {
	opt, err := config.Option(optionName)
	if err != nil {
		return err
	}
	return opt.Unset(cfg)
}
