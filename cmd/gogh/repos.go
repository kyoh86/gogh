package main

import (
	"fmt"

	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/app"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/spf13/cobra"
)

var reposFlags struct {
	query string
}

var reposCommand = &cobra.Command{
	Use:   "repos",
	Short: "List remote repositories",
	RunE: func(cmd *cobra.Command, _ []string) error {
		list, err := app.Servers().List()
		if err != nil {
			return err
		}
		ctx := cmd.Context()
		for _, server := range list {
			adaptor, err := github.NewAdaptor(ctx, server.Host(), server.Token())
			if err != nil {
				return err
			}
			remote := gogh.NewRemoteController(adaptor)
			specs, err := remote.List(ctx, &gogh.RemoteListOption{
				Query: reposFlags.query,
			})
			if err != nil {
				return err
			}
			for _, spec := range specs {
				fmt.Println(spec)
			}
		}
		return nil
	},
}

func init() {
	reposCommand.Flags().StringVarP(&listFlags.query, "query", "", "", "Query for selecting projects")
	facadeCommand.AddCommand(reposCommand)
}
