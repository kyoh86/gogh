package script_invoke

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"os/exec"

	"github.com/kyoh86/gogh/v4/app/script_run"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/script"
	"github.com/kyoh86/gogh/v4/core/workspace"
	"golang.org/x/sync/errgroup"
)

// UseCase for running script scripts
type UseCase struct {
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
	scriptService    script.ScriptService
	referenceParser  repository.ReferenceParser
}

func NewUseCase(
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
	scriptService script.ScriptService,
	referenceParser repository.ReferenceParser,
) *UseCase {
	return &UseCase{
		workspaceService: workspaceService,
		finderService:    finderService,
		scriptService:    scriptService,
		referenceParser:  referenceParser,
	}
}

func (uc *UseCase) Execute(ctx context.Context, refStr string, scriptID string, globals map[string]any) error {
	refWithAlias, err := uc.referenceParser.ParseWithAlias(refStr)
	if err != nil {
		return fmt.Errorf("parsing repository reference: %w", err)
	}
	match, err := uc.finderService.FindByReference(ctx, uc.workspaceService, refWithAlias.Local())
	if err != nil {
		return fmt.Errorf("find repository location: %w", err)
	}
	return uc.Invoke(ctx, match, scriptID, globals)
}

func (uc *UseCase) Invoke(ctx context.Context, location *repository.Location, scriptID string, globals map[string]any) error {
	if location == nil {
		return errors.New("repository not found")
	}
	src, err := uc.scriptService.Open(ctx, scriptID)
	if err != nil {
		return fmt.Errorf("open script script: %w", err)
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
		"full_path": location.FullPath(),
		"path":      location.Path(),
		"host":      location.Host(),
		"owner":     location.Owner(),
		"name":      location.Name(),
	}

	cmd := exec.Command(os.Args[0], "script", "run")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = location.FullPath()
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

		return enc.Encode(script_run.Script{
			Code:    string(code),
			Globals: g,
		})
	})

	eg.Go(cmd.Run)

	return eg.Wait()
}
