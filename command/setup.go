package command

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/sh"
)

// Setup shells in shell scipt
// Usage: eval "$(gogh setup)"
func Setup(ctx gogh.Context, _, shell string) error {
	_ = sh.Assets
	_, shellName := filepath.Split(shell)
	assetName := "/src/init." + shellName
	if !sh.Assets.Exists(assetName) {
		return fmt.Errorf("unsupported shell %q", shell)
	}
	file, err := sh.Assets.Open(assetName)
	if err != nil {
		return err
	}
	_, err = io.Copy(ctx.Stdout(), file)
	return err
}
