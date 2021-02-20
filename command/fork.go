package command

import (
	"context"

	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/internal/github"
)

func Fork(ctx context.Context, root string, servers *gogh.Servers, rawSpec string, lopt *gogh.LocalCloneOption, ropt *gogh.RemoteForkOption) error {
	parser := gogh.NewSpecParser(servers)
	spec, server, err := parser.Parse(rawSpec)
	if err != nil {
		return err
	}
	adaptor, err := github.NewAdaptor(ctx, server.Host(), server.Token())
	if err != nil {
		return err
	}
	remote := gogh.NewRemoteController(adaptor)
	if _, err := remote.Fork(ctx, spec.Owner(), spec.Name(), ropt); err != nil {
		return err
	}
	local := gogh.NewLocalController(root)
	_, err = local.Clone(ctx, spec, server, lopt)
	return err
}
