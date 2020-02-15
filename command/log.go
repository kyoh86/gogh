package command

import (
	"io"
	"os"

	"github.com/comail/colog"
	"github.com/kyoh86/gogh/gogh"
)

func InitLog(ctx gogh.Context) {
	colog.SetMinLevel(colog.LInfo)
	colog.SetDefaultLevel(colog.LError)
	var stderr io.Writer = os.Stderr
	if ctx, ok := ctx.(gogh.IOContext); ok {
		stderr = ctx.Stderr()
	}
	colog.SetOutput(stderr)
	colog.Register()
}
