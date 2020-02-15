package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/gogh"
	_ "github.com/kyoh86/gogh/sh" //nolint
	"github.com/rakyll/statik/fs"
)

// Setup shells in shell scipt
// Usage: eval "$(gogh setup)"
func Setup(ctx gogh.Context, _, shell string) error {
	staticFs, err := fs.New()
	if err != nil {
		return err
	}
	_, shellName := filepath.Split(shell)
	assetName := "/init." + shellName
	file, err := staticFs.Open(assetName)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("unsupported shell %q", shell)
		}
		return err
	}
	defer file.Close()
	_, err = io.Copy(os.Stdout, file)
	return err
}
