package logger

import (
	"context"
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
func (h *StdoutLogHandler) HandleLog(e *log.Entry) error {
	if e.Level >= log.ErrorLevel {
		return nil
	}

	return h.Handler.HandleLog(e)
}

// NewLogger creates a new logger instance.
func NewLogger() context.Context {
	errLog := level.New(cli.New(os.Stderr), log.ErrorLevel)
	stdLog := &StdoutLogHandler{Handler: cli.New(os.Stdout)}
	level := log.InfoLevel
	if os.Getenv("GOGH_DEBUG") == "1" {
		level = log.DebugLevel
	}
	return log.NewContext(context.Background(), &log.Logger{
		Handler: multi.New(stdLog, errLog),
		Level:   level,
	})
}
