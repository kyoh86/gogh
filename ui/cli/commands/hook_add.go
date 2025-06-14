package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v4/app/hook_add"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewHookAddCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		name        string
		useCase     string
		event       string
		repoPattern string
	}
	cmd := &cobra.Command{
		Use:   "add [flags] <lua-script-path>",
		Short: "Add an existing Lua script as hook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			path := args[0]
			fi, err := os.Stat(path)
			if err != nil || fi.IsDir() {
				return fmt.Errorf("invalid script path: %v", err)
			}
			content, err := os.Open(path)
			if err != nil {
				return err
			}
			defer content.Close()
			opts := hook_add.Options{
				Name:        f.name,
				UseCase:     f.useCase,
				Event:       f.event,
				RepoPattern: f.repoPattern,
			}
			h, err := hook_add.NewUseCase(svc.HookService).Execute(ctx, opts, content)
			if err != nil {
				return fmt.Errorf("adding hook: %w", err)
			}
			fmt.Printf("Hook added %s\n", h.ID)
			return nil
		},
	}
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the hook")

	if err := enumFlag(cmd, &f.useCase, "use-case", "never", "Use case to hook automatically", "", "clone", "fork", "create", "never"); err != nil {
		return nil, fmt.Errorf("registering use-case flag: %w", err)
	}

	if err := enumFlag(cmd, &f.event, "event", "never", "event to hook automatically", "", "clone", "fork", "create", "never"); err != nil {
		return nil, fmt.Errorf("registering event flag: %w", err)
	}

	cmd.Flags().StringVar(&f.repoPattern, "repo-pattern", "", "Repository pattern")
	return cmd, nil
}
