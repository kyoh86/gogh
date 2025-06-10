package hook_apply

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"maps"
	"os"
	"os/exec"

	"github.com/kyoh86/gogh/v4/app/hook_run"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
	"golang.org/x/sync/errgroup"
)

// UseCase for running hook scripts
type UseCase struct {
	hookService      hook.HookService
	referenceParser  repository.ReferenceParser
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
}

func NewUseCase(
	hookService hook.HookService,
	referenceParser repository.ReferenceParser,
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
) *UseCase {
	return &UseCase{
		hookService:      hookService,
		referenceParser:  referenceParser,
		workspaceService: workspaceService,
		finderService:    finderService,
	}
}

func (uc *UseCase) Execute(ctx context.Context, hookID string, refWithAlias string, globals map[string]any) error {
	ref, err := uc.referenceParser.ParseWithAlias(refWithAlias)
	if err != nil {
		return fmt.Errorf("parsing repository reference: %w", err)
	}
	loc, err := uc.finderService.FindByReference(ctx, uc.workspaceService, ref.Local())
	if err != nil {
		return fmt.Errorf("find repository location: %w", err)
	}

	hook, err := uc.hookService.GetHookByID(ctx, hookID)
	if err != nil {
		return fmt.Errorf("get hook by ID: %w", err)
	}
	src, err := uc.hookService.OpenHookScript(ctx, *hook)
	if err != nil {
		return fmt.Errorf("open hook script: %w", err)
	}
	defer src.Close()
	code, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("read script: %w", err)
	}

	g := make(map[string]any, len(globals)+1)
	maps.Copy(g, globals)
	// Add domain objects as maps
	g["repo"] = map[string]any{
		"full_path": loc.FullPath(),
		"path":      loc.Path(),
		"host":      loc.Host(),
		"owner":     loc.Owner(),
		"name":      loc.Name(),
	}

	cmd := exec.Command(os.Args[0], "hook", "run")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = loc.FullPath()
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	var eg errgroup.Group
	eg.SetLimit(2)

	eg.Go(func() error {
		gob.Register(map[string]any{})
		enc := gob.NewEncoder(stdin)
		defer stdin.Close()

		return enc.Encode(hook_run.Script{
			Code:    string(code),
			Globals: g,
		})
	})

	eg.Go(cmd.Run)

	return eg.Wait()
}
