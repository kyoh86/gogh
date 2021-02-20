package command

import (
	"context"

	"github.com/kyoh86/gogh/v2"
)

func Fork(ctx context.Context, root string, servers *gogh.Servers, from string, to string, opt *gogh.LocalCloneOption) error {
	parser := gogh.NewSpecParser(servers)
	spec, server, err := parser.Parse(from)
	if err != nil {
		return err
	}
	local := gogh.NewLocalController(root)
	_, err = local.Clone(ctx, spec, server, opt)
	return err
}
