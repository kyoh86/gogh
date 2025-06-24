package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v4/app/script/add"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewScriptAddCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		name string
	}
	cmd := &cobra.Command{
		Use:   "add [flags] <lua-script-path>",
		Short: "Add an existing Lua script as script",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			path := args[0]
			content, err := os.Open(path)
			if err != nil {
				return err
			}
			defer content.Close()
			h, err := add.NewUseCase(svc.ScriptService).Execute(ctx, f.name, content)
			if err != nil {
				return fmt.Errorf("adding script: %w", err)
			}
			fmt.Printf("Script added %s\n", h.ID())
			return nil
		},
	}
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the script")
	return cmd, nil
}
