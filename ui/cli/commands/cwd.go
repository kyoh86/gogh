package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"github.com/kyoh86/gogh/v3/domain/local"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/ui/cli/flags"
	"github.com/spf13/cobra"
)

func NewCwdCommand(conf *config.ConfigStore, defaults *config.FlagStore, finderService workspace.FinderService) *cobra.Command {
	var f config.CwdFlags
	cmd := &cobra.Command{
		Use:   "cwd",
		Short: "Print the local reposiotry in current working directory",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			formatter, err := f.Format.Formatter()
			if err != nil {
				return err
			}

			ctx := cmd.Context()
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			cwd = strings.ToLower(filepath.ToSlash(cwd))

			list := conf.GetRoots()
			for _, root := range list {
				ctrl := local.NewController(root)
				repos, err := ctrl.List(ctx, &local.ListOption{})
				if err != nil {
					return err
				}
				log.FromContext(ctx).Debugf("found %d local repositories in %q", len(repos), root)
				for _, repo := range repos {
					reg := strings.ToLower(filepath.ToSlash(repo.FullFilePath()))
					if cwd == reg || strings.HasPrefix(cwd, reg+"/") {
						str, err := formatter.Format(repo)
						if err != nil {
							log.FromContext(ctx).WithFields(log.Fields{
								"error":  err,
								"format": f.Format.String(),
								"path":   repo.FullFilePath(),
							}).Info("failed to format")
						}
						fmt.Println(str)
						return nil
					}
				}
			}
			log.FromContext(ctx).WithField("cwd", cwd).Info("it is not in any local repository")
			return nil
		},
	}

	f.Format = defaults.Cwd.Format
	cmd.Flags().VarP(&f.Format, "format", "f", flags.LocalRepoFormatShortUsage)
	if err := cmd.RegisterFlagCompletionFunc("format", flags.CompleteLocalRepoFormat); err != nil {
		panic(err)
	}
	return cmd
}
