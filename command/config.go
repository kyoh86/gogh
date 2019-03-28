package command

import (
	"fmt"

	"github.com/kyoh86/gogh/config"
)

func ConfigGetAll(cfg *config.Config) error {
	acc := config.DefaultAccessor
	for _, name := range acc.OptionNames() {
		opt, _ := acc.Option(name) // ignore error: acc.OptionNames covers all accessor
		value := opt.Get(cfg)
		fmt.Printf("%s = %s\n", name, value)
	}
	return nil
}

func ConfigGet(cfg *config.Config, optionName string) error {
	opt, err := config.DefaultAccessor.Option(optionName)
	if err != nil {
		return err
	}
	value := opt.Get(cfg)
	fmt.Printf("%s = %s\n", optionName, value)
	return nil
}

func ConfigPut(cfg *config.Config, optionName, optionValue string) error {
	opt, err := config.DefaultAccessor.Option(optionName)
	if err != nil {
		return err
	}
	return opt.Put(cfg, optionValue)
}

func ConfigUnset(cfg *config.Config, optionName string) error {
	opt, err := config.DefaultAccessor.Option(optionName)
	if err != nil {
		return err
	}
	return opt.Unset(cfg)
}
