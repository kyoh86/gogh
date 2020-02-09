package command

import (
	"log"

	"github.com/comail/colog"
	"github.com/kyoh86/gogh/gogh"
)

func InitLog(ctx gogh.Context) {
	rawLevel := ctx.LogLevel()
	if rawLevel != "" {
		lvl, err := colog.ParseLevel(rawLevel)
		if err != nil {
			defer log.Println("warn: could not parse level %s", rawLevel)
		}
		colog.SetMinLevel(lvl)
	}
	colog.SetDefaultLevel(colog.LError)
	colog.SetFormatter(&colog.StdFormatter{
		Flag:        ctx.LogFlags(),
		HeaderPlain: plainLabels,
		HeaderColor: colorLabels,
	})
	colog.SetOutput(ctx.Stderr())
	colog.Register()
}

var plainLabels = colog.LevelMap{
	colog.LTrace:   []byte("[ trace ] "),
	colog.LDebug:   []byte("\u2699 "),
	colog.LInfo:    []byte("\u24d8 "),
	colog.LWarning: []byte("\u26a0 "),
	colog.LError:   []byte("\u2622 "),
	colog.LAlert:   []byte("\u2620 "),
}

var colorLabels = colog.LevelMap{
	colog.LTrace:   []byte("[ trace ] "),
	colog.LDebug:   []byte("\x1b[0;36m\u2699 \x1b[0m"),
	colog.LInfo:    []byte("\x1b[0;32m\u24d8 \x1b[0m"),
	colog.LWarning: []byte("\x1b[0;33m\u26a0 \x1b[0m"),
	colog.LError:   []byte("\x1b[0;31m\u2622 \x1b[0m"),
	colog.LAlert:   []byte("\x1b[0;37;41m\u2620 \x1b[0m"),
}
