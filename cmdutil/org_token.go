package cmdutil

import (
	"context"
	"errors"
	"fmt"

	"github.com/kyoh86/gogh/v3"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/infra/github"
)

func RemoteControllerFor(ctx context.Context, tokens config.TokenStore, ref gogh.RepoRef) (github.Adaptor, *gogh.RemoteController, error) {
	token, err := tokens.Get(ref.Host(), ref.Owner())
	switch {
	case err == nil:
		adaptor, err := github.NewAdaptor(ctx, ref.Host(), &token)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to build adaptor for %q: %w", ref.Owner(), err)
		}
		return adaptor, gogh.NewRemoteController(adaptor), nil
	case errors.Is(err, config.ErrNoHost):
		return nil, nil, err
	case errors.Is(err, config.ErrNoOwner):
		// Check each owners is member of the ref.Owner() organization
		owners := tokens.Hosts[ref.Host()].Owners
		for owner, token := range owners {
			adaptor, err := github.NewAdaptor(ctx, ref.Host(), &token)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to build adaptor for %q: %w", owner, err)
			}
			ctrl := gogh.NewRemoteController(adaptor)
			ok, err := ctrl.MemberOf(ctx, ref.Owner(), nil)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to check the member of %q: %w", owner, err)
			}
			if ok {
				return adaptor, ctrl, nil
			}
		}
	}

	tokenHost, _ := tokens.Hosts.TryGet(ref.Host())
	token, ok := tokenHost.Owners.TryGet(tokenHost.DefaultOwner)
	if !ok {
		return nil, nil, config.ErrNoOwner
	}
	adaptor, err := github.NewAdaptor(ctx, ref.Host(), &token)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build adaptor for %q: %w", ref.Owner(), err)
	}
	return adaptor, gogh.NewRemoteController(adaptor), nil
}
