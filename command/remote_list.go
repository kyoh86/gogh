package command

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/internal/github"
)

func RemoteList(ctx context.Context, servers *gogh.Servers, query string) error {
	list, err := servers.List()
	if err != nil {
		return err
	}
	for _, server := range list {
		adaptor, err := github.NewAdaptor(ctx, server.Host(), server.Token())
		if err != nil {
			return err
		}
		remote := gogh.NewRemoteController(adaptor)
		specs, err := remote.List(ctx, &gogh.RemoteListOption{
			Query: query,
		})
		if err != nil {
			return err
		}
		for _, spec := range specs {
			fmt.Println(spec)
		}
	}
	return nil
}
