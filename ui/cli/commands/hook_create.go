package commands

import (
	"context"
	"os"
	"os/exec"

	"github.com/kyoh86/gogh/v4/app/hook_create"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewHookCreateCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		name        string
		description string
		event       string
		repoPattern string
	}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new hook (edit with $EDITOR)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			tmpDir := os.TempDir()
			file, err := os.CreateTemp(tmpDir, "gogh_hook_*.lua")
			if err != nil {
				return err
			}
			defer os.Remove(file.Name())
			// 空テンプレートやサンプルLuaを書き込んでもよい
			file.WriteString(`-- Gogh Hook Lua script template
function main()
  -- your code here
end
`)
			file.Close()

			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vi"
			}
			cmdEdit := exec.Command(editor, file.Name())
			cmdEdit.Stdin = os.Stdin
			cmdEdit.Stdout = os.Stdout
			cmdEdit.Stderr = os.Stderr
			if err := cmdEdit.Run(); err != nil {
				return err
			}
			// 編集後内容を読み込む
			content, err := os.Open(file.Name())
			if err != nil {
				return err
			}
			defer content.Close()

			opts := hook_create.Options{
				Name:        f.name,
				Description: f.description,
				Event:       f.event,
				RepoPattern: f.repoPattern,
			}
			return hook_create.NewUseCase(svc.HookService).Execute(ctx, opts, content)
		},
	}
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the hook")
	cmd.Flags().StringVar(&f.description, "description", "", "Description")
	cmd.Flags().StringVar(&f.event, "event", "", "Event (before-clone, after-clone, etc.)")
	cmd.Flags().StringVar(&f.repoPattern, "repo-pattern", "", "Repository pattern")
	return cmd, nil
}

