package gogh

import (
	"fmt"
	"path/filepath"
)

// Setup shells in shell scipt
// Usage: eval "$(gogh setup)"
func Setup(cdFuncName, shell string) error {
	_, shName := filepath.Split(shell)
	switch shName {
	case "zsh":
		if _, err := fmt.Printf(`function %s { cd $(gogh find $@) }%s`, cdFuncName, "\n"); err != nil {
			return err
		}
		if _, err := fmt.Printf(`eval "$(gogh --completion-script-zsh)"%s`, "\n"); err != nil {
			return err
		}
		return nil
	case "bash":
		if _, err := fmt.Printf(`function %s { cd $(gogh find $@) }%s`, cdFuncName, "\n"); err != nil {
			return err
		}
		if _, err := fmt.Printf(`eval "$(gogh --completion-script-bash)"%s`, "\n"); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported shell %q", shell)
	}
}
