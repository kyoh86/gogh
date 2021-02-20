package command

import (
	"context"

	"github.com/kyoh86/gogh/v2"
)

func Clone(ctx context.Context, root string, servers *gogh.Servers, rawSpec string, opt *gogh.LocalCloneOption) error {
	parser := gogh.NewSpecParser(servers)
	spec, server, err := parser.Parse(rawSpec)
	if err != nil {
		return err
	}

	local := gogh.NewLocalController(root)
	_, err = local.Clone(ctx, spec, server, opt)
	return err
}
