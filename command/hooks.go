package command

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/internal/delegate"
)

const (
	hookPostGetEach   = "post-get-each"
	hookPostCreate    = "post-create"
	hookPostFork      = "post-fork"
	hookPreRemoveEach = "pre-remove-each"
)

func execHookInDir(dir, hook string) error {
	if _, err := os.Stat(hook); err != nil {
		return nil // ignore file error
	}
	cmd := exec.Command(hook)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return delegate.ExecCommand(cmd)
}

func execHooks(ev gogh.Env, p *gogh.Project, name string) error {
	for _, hook := range ev.Hooks() {
		if err := execHookInDir(p.FullPath, filepath.Join(hook, name)); err != nil {
			return err
		}
	}
	return execHookInDir(p.FullPath, filepath.Join(p.FullPath, ".gogh", "hooks", name))
}
