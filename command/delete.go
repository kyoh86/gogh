package command

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/internal/github"
)

func Delete(ctx context.Context, root string, servers *gogh.Servers, rawSpec string, lopt *gogh.LocalDeleteOption, ropt *gogh.RemoteDeleteOption) error {
	parser := gogh.NewSpecParser(servers)
	spec, server, err := parser.Parse(rawSpec)
	if err != nil {
		return err
	}

	local := gogh.NewLocalController(root)
	if err := local.Delete(ctx, spec, lopt); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("delete local: %w", err)
		}
	}

	adaptor, err := github.NewAdaptor(ctx, server.Host(), server.Token())
	if err != nil {
		return err
	}
	return gogh.NewRemoteController(adaptor).Delete(ctx, spec.Owner(), spec.Name(), ropt)
}
