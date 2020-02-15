package command

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kyoh86/gogh/gogh"
)

// Root prints a gogh.root
func Root(ctx gogh.Context, all bool) error {
	InitLog(ctx)

	var stdout io.Writer = os.Stdout
	if ctx, ok := ctx.(gogh.IOContext); ok {
		stdout = ctx.Stdout()
	}
	if !all {
		fmt.Fprintln(stdout, ctx.PrimaryRoot())
		return nil
	}
	log.Println("info: Finding all roots...")
	for _, root := range ctx.Root() {
		fmt.Fprintln(stdout, root)
	}
	return nil
}
