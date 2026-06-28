package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v4/app/script/add"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewScriptCreateCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		name string
	}
	cmd := &cobra.Command{
		Use:   "create [flags]",
		Short: "Create a new script (with $GOGH_EDITOR or $EDITOR)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Create a temporary file for editing
			tmpFile, err := os.CreateTemp("", "gogh_script_create_*.lua")
			if err != nil {
				return err
			}
			defer os.Remove(tmpFile.Name())

			// Write initial template
			initialContent := `-- Gogh script
---@class gogh.Repo
---@field full_path string Local repository full-path
---@field path string Local repository path
---@field host string Repository host
---@field owner string Repository owner
---@field name string Repository name

---@class gogh.Hook
---@field id string Hook UUID
---@field name string Hook name
---@field repoPattern string Pattern that matched
---@field triggerEvent string Event that triggered the hook
---@field operationType string Type of operation that triggered the hook: "script" always.
---@field operationId string Operation UUID (script UUID)

---@class gogh
---@field repo gogh.Repo
---@field hook gogh.Hook
---@field parent? gogh.Repo|nil

---@type gogh
_G.gogh = _G.gogh or {}

`
			if _, err := tmpFile.WriteString(initialContent); err != nil {
				return err
			}
			tmpFile.Close()

			// Open editor
			if err := edit(tmpFile.Name()); err != nil {
				return err
			}

			// Read the edited content
			content, err := os.Open(tmpFile.Name())
			if err != nil {
				return err
			}
			defer content.Close()

			// Add the script
			h, err := add.NewUsecase(svc.ScriptService).Execute(ctx, f.name, content)
			if err != nil {
				return fmt.Errorf("adding script: %w", err)
			}
			fmt.Printf("Script created %s\n", h.ID())
			return nil
		},
	}
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the script")
	return cmd, nil
}
