package command

import (
	"context"
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/internal/github"
)

func Create(ctx context.Context, root string, servers gogh.Servers, rawSpec string, lopt *gogh.LocalCreateOption, ropt *gogh.RemoteCreateOption) error {
	parser := gogh.NewSpecParser(servers)
	spec, server, err := parser.Parse(rawSpec)
	if err != nil {
		return err
	}
	_ = server

	local := gogh.NewLocalController(root)
	if _, err = local.Create(ctx, spec, nil); err != nil {
		if !errors.Is(err, git.ErrRepositoryAlreadyExists) {
			return err
		}
	}

	adaptor, err := github.NewAdaptor(ctx, server.Host(), server.Token())
	if err != nil {
		return err
	}
	remote := gogh.NewRemoteController(adaptor)

	// check repo has already existed
	if _, err := remote.Get(ctx, spec.User(), spec.Name(), nil); err == nil {
		return nil
	}

	var org string
	if server.User() != spec.User() {
		org = spec.User()
	}
	if _, err := remote.Create(ctx, spec.Name(), &gogh.RemoteCreateOption{
		Organization: org,
	}); err != nil {
		return err
	}
	return nil
}
