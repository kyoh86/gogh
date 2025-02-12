package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v2"
	"github.com/spf13/cobra"
)

type cwdFlagsStruct struct {
	Format ProjectFormat `yaml:"format,omitempty"`
}

var (
	cwdFlags   cwdFlagsStruct
	cwdCommand = &cobra.Command{
		Use:   "cwd",
		Short: "Print the project in current working directory",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := cwdFlags.Format.Formatter()
			if err != nil {
				return err
			}

			ctx := cmd.Context()
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			cwd = strings.ToLower(filepath.ToSlash(cwd))

			list := roots()
			for _, root := range list {
				local := gogh.NewLocalController(root)
				projects, err := local.List(ctx, &gogh.LocalListOption{Query: listFlags.Query})
				if err != nil {
					return err
				}
				log.FromContext(ctx).Debugf("found %d projects in %q", len(projects), root)
				for _, project := range projects {
					reg := strings.ToLower(filepath.ToSlash(project.FullFilePath()))
					if cwd == reg || strings.HasPrefix(cwd, reg+"/") {
						str, err := f.Format(project)
						if err != nil {
							log.FromContext(ctx).WithFields(log.Fields{
								"error":  err,
								"format": listFlags.Format.String(),
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
)

func init() {
	cwdFlags.Format = defaultFlag.Cwd.Format
	cwdCommand.Flags().VarP(&cwdFlags.Format, "format", "f", formatShortUsage)
	if err := cwdCommand.RegisterFlagCompletionFunc("format", completeFormat); err != nil {
		panic(err)
	}
	facadeCommand.AddCommand(cwdCommand)
}
