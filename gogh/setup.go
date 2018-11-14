package gogh

import (
	"fmt"
	"path/filepath"
)

func Setup(cdFuncName, shell string) error {
	_, shName := filepath.Split(shell)
	switch shName {
	case "zsh", "bash":
		_, err := fmt.Printf(`function %s { cd $(gogh find $@) }`, cdFuncName)
		return err
	default:
		return fmt.Errorf("unsupported shell %q", shell)
	}
}
