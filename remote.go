package gogh

import (
	"context"

	"github.com/kyoh86/gogh/v2/internal/github"
)

type RemoteController struct {
}

func GithubConnector(ctx context.Context, server Server) (github.Adaptor, error) {
}

type Connector func(context.Context, Server) (github.Adaptor, error)
