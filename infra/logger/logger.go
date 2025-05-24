package logger

import (
	"context"
	"io"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/level"
	"github.com/apex/log/handlers/multi"
)

// StdoutLogHandler implementation.
type StdoutLogHandler struct {
	Handler log.Handler
}

// HandleLog implements log.Handler.
// It filters out error logs.
func (h *StdoutLogHandler) HandleLog(e *log.Entry) error {
	if e.Level >= log.ErrorLevel {
		return nil
	}

	return h.Handler.HandleLog(e)
}

// NewLogger creates a new logger instance.
func NewLogger(ctx context.Context, outWriter io.Writer, errWriter io.Writer) context.Context {
	errLog := level.New(cli.New(errWriter), log.WarnLevel)
	outLog := &StdoutLogHandler{Handler: cli.New(outWriter)}
	level := log.InfoLevel
	if os.Getenv("GOGH_DEBUG") != "" {
		level = log.DebugLevel
	}
	return log.NewContext(ctx, &log.Logger{
		Handler: multi.New(outLog, errLog),
		Level:   level,
	})
}
