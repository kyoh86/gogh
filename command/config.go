package command

import (
	"github.com/kyoh86/gogh/gogh"
)

func ConfigGet(ctx gogh.Context, optionName string) error {
	return nil
}

func ConfigGetAll(ctx gogh.Context) error {
	return nil
}

func ConfigSet(config *gogh.Config, optionName, optionValue string) (*gogh.Config, error) {
	return nil, nil
}
