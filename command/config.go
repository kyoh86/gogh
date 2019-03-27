package command

import (
	"fmt"

	"github.com/kyoh86/gogh/config"
	"github.com/kyoh86/gogh/gogh"
)

func ConfigGetAll(ctx gogh.Context) error {
	cfg := config.DefaultConfigAccessor
	for _, name := range cfg.OptionNames() {
		acc, _ := cfg.Accessor(name) // ignore error: cfg.OptionNames covers all accessor
		value := acc.Get(ctx)
		fmt.Printf("%s = %s\n", name, value)
	}
	return nil
}

func ConfigGet(ctx gogh.Context, optionName string) error {
	acc, err := config.DefaultConfigAccessor.Accessor(optionName)
	if err != nil {
		return err
	}
	value := acc.Get(ctx)
	fmt.Printf("%s = %s\n", optionName, value)
	return nil
}

func ConfigSet(cfg *config.Config, optionName, optionValue string) error {
	acc, err := config.DefaultConfigAccessor.Accessor(optionName)
	if err != nil {
		return err
	}
	return acc.Set(cfg, optionValue)
}

func ConfigUnset(cfg *config.Config, optionName string) error {
	acc, err := config.DefaultConfigAccessor.Accessor(optionName)
	if err != nil {
		return err
	}
	return acc.Unset(cfg)
}
