package command

import (
	"fmt"
	"path/filepath"

	"github.com/kyoh86/gogh/gogh"
)

// Setup shells in shell scipt
// Usage: eval "$(gogh setup)"
func Setup(ctx gogh.Context, cdFuncName, shell string) error {
	_, shName := filepath.Split(shell)
	switch shName {
	case "zsh":
		fmt.Fprintf(ctx.Stdout(), `function %s { cd $(gogh find $@) }%s`, cdFuncName, "\n")
		fmt.Fprintf(ctx.Stdout(), `eval "$(gogh --completion-script-zsh)"%s`, "\n")
		return nil
	case "bash":
		fmt.Fprintf(ctx.Stdout(), `function %s { cd $(gogh find $@) }%s`, cdFuncName, "\n")
		fmt.Fprintf(ctx.Stdout(), `eval "$(gogh --completion-script-bash)"%s`, "\n")
		return nil
	default:
		return fmt.Errorf("unsupported shell %q", shell)
	}
}
