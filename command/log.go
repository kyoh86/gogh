package command

import (
	"os"

	"github.com/comail/colog"
)

func InitLog() {
	colog.SetMinLevel(colog.LInfo)
	colog.SetDefaultLevel(colog.LError)
	colog.SetOutput(os.Stderr)
	colog.Register()
}
