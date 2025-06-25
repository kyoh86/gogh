package invoke

import (
	"context"
	"encoding/gob"
	"os"
	"os/exec"

	"github.com/kyoh86/gogh/v4/app/script/run"
	"github.com/kyoh86/gogh/v4/core/repository"
	"golang.org/x/sync/errgroup"
)

// InvokeInstant executes a script directly without storing it
func InvokeInstant(ctx context.Context, location *repository.Location, code string, globals map[string]any) error {
	g := make(map[string]any, len(globals)+1)
	for k, v := range globals {
		g[k] = v
	}
	// Add domain objects as maps
	g["repo"] = map[string]any{
		"full_path": location.FullPath(),
		"path":      location.Path(),
		"host":      location.Host(),
		"owner":     location.Owner(),
		"name":      location.Name(),
	}

	// Get the executable path in a cross-platform way
	exePath := os.Args[0]
	if exe, err := os.Executable(); err == nil {
		exePath = exe
	}

	cmd := exec.Command(exePath, "script", "run")
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

		return enc.Encode(run.Script{
			Code:    code,
			Globals: g,
		})
	})

	eg.Go(cmd.Run)

	return eg.Wait()
}
