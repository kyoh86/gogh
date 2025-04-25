package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/ui/cli/flags"
	"github.com/spf13/cobra"
)

func NewCwdCommand(conf *config.Config, defaults *config.Flags) *cobra.Command {
	var f config.CwdFlags
	cmd := &cobra.Command{
		Use:   "cwd",
		Short: "Print the project in current working directory",
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
				local := gogh.NewLocalController(root)
				projects, err := local.List(ctx, &gogh.LocalListOption{})
				if err != nil {
					return err
				}
				log.FromContext(ctx).Debugf("found %d projects in %q", len(projects), root)
				for _, project := range projects {
					reg := strings.ToLower(filepath.ToSlash(project.FullFilePath()))
					if cwd == reg || strings.HasPrefix(cwd, reg+"/") {
						str, err := formatter.Format(project)
						if err != nil {
							log.FromContext(ctx).WithFields(log.Fields{
								"error":  err,
								"format": f.Format.String(),
								"path":   project.FullFilePath(),
							}).Info("failed to format")
						}
						fmt.Println(str)
						return nil
					}
				}
			}
			log.FromContext(ctx).WithField("cwd", cwd).Info("it is not in any project")
			return nil
		},
	}

	f.Format = defaults.Cwd.Format
	cmd.Flags().VarP(&f.Format, "format", "f", flags.ProjectFormatShortUsage)
	if err := cmd.RegisterFlagCompletionFunc("format", flags.CompleteProjectFormat); err != nil {
		panic(err)
	}
	return cmd
}
