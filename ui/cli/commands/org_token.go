package commands

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v3/domain/remote"
	"github.com/kyoh86/gogh/v3/domain/reporef"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/infra/github"
)

// RemoteControllerFor creates a GitHub adaptor and remote controller for the given repository reference
func RemoteControllerFor(ctx context.Context, tokens config.TokenStore, ref reporef.RepoRef) (github.Adaptor, *remote.Controller, error) {
	exactMatch, candidates, err := tokens.GetTokenForOwner(ref.Host(), ref.Owner())
	if err != nil {
		return nil, nil, err
	}

	// Try exact matching token if available
	if exactMatch != nil {
		adaptor, err := github.NewAdaptor(ctx, exactMatch.Host, &exactMatch.Token)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to build adaptor for %q: %w", exactMatch.Owner, err)
		}
		return adaptor, remote.NewController(adaptor), nil
	}

	// No exact match, try candidates that may have access to the org
	for _, candidate := range candidates {
		adaptor, err := github.NewAdaptor(ctx, candidate.Host, &candidate.Token)
		if err != nil {
			continue // Try next token if this one fails
		}

		// Check if this user is a member of the target organization
		ok, err := adaptor.MemberOf(ctx, ref.Owner())
		if err != nil {
			continue // Try next token if membership check fails
		}

		if ok {
			// Found a working token
			return adaptor, remote.NewController(adaptor), nil
		}
	}

	return nil, nil, fmt.Errorf("no valid token found for %s/%s", ref.Host(), ref.Owner())
}
