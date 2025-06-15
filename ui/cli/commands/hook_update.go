package commands

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/kyoh86/gogh/v4/app/hook_update"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewHookUpdateCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		name        string
		useCase     string
		event       string
		repoPattern string
		scriptPath  string
	}
	cmd := &cobra.Command{
		Use:   "update [flags] <hook-id>",
		Short: "Update an existing hook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			hookID := args[0]
			var content io.Reader
			if f.scriptPath != "" {
				fi, err := os.Stat(f.scriptPath)
				if err != nil || fi.IsDir() {
					return fmt.Errorf("invalid script path: %v", err)
				}
				c, err := os.Open(f.scriptPath)
				if err != nil {
					return err
				}
				content = c
				defer c.Close()
			}
			if err := hook_update.NewUseCase(svc.HookService).Execute(ctx, hookID, f.name, f.useCase, f.event, f.repoPattern, content); err != nil {
				return fmt.Errorf("updating hook metadata: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the hook")
	cmd.Flags().StringVar(&f.scriptPath, "script-path", "", "Path to the Lua script file")

	if err := enumFlag(cmd, &f.useCase, "use-case", "never", "Use case to hook automatically", "", "clone", "fork", "create", "never"); err != nil {
		return nil, fmt.Errorf("registering use-case flag: %w", err)
	}

	if err := enumFlag(cmd, &f.event, "event", "never", "event to hook automatically", "", "clone", "fork", "create", "never"); err != nil {
		return nil, fmt.Errorf("registering event flag: %w", err)
	}

	cmd.Flags().StringVar(&f.repoPattern, "repo-pattern", "", "Repository pattern")
	return cmd, nil
}
