package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/cli/safeexec"
	shellquote "github.com/kballard/go-shellquote"
	"github.com/kyoh86/gogh/v4/app/hook_create"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func lookPath(name string) ([]string, error) {
	exe, err := safeexec.LookPath(name)
	if err != nil {
		return nil, err
	}
	return []string{exe}, nil
}

func edit(editor, fileName string) error {
	words, err := shellquote.Split(editor)
	if err != nil {
		return err
	}
	words = append(words, fileName)
	editorExe, err := lookPath(words[0])
	if err != nil {
		return err
	}
	words = append(editorExe, words[1:]...)

	cmdEdit := exec.Command(words[0], words[1:]...)
	cmdEdit.Env = os.Environ()
	cmdEdit.Stdin = os.Stdin
	cmdEdit.Stdout = os.Stdout
	cmdEdit.Stderr = os.Stderr
	return cmdEdit.Run()
}

func NewHookCreateCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		name        string
		useCase     string
		event       string
		repoPattern string
	}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new hook (edit with $EDITOR)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			tmpDir := os.TempDir()
			file, err := os.CreateTemp(tmpDir, "gogh_hook_*.lua")
			if err != nil {
				return err
			}
			if err := file.Close(); err != nil {
				return fmt.Errorf("close temporary file: %v", err)
			}
			defer os.Remove(file.Name())

			if err := edit(os.Getenv("EDITOR"), file.Name()); err != nil {
				return err
			}
			// Reopen the file to read the content after editing
			content, err := os.Open(file.Name())
			if err != nil {
				return err
			}
			defer content.Close()

			opts := hook_create.Options{
				Name:        f.name,
				UseCase:     f.useCase,
				Event:       f.event,
				RepoPattern: f.repoPattern,
			}
			return hook_create.NewUseCase(svc.HookService).Execute(ctx, opts, content)
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
