package testutil

import (
	"errors"
	"io"
)

type ErrorWriter struct {
	Error error
}

var (
	DefaultErrorWriter io.Writer = &ErrorWriter{
		Error: errors.New("error writer"),
	}
)

func (w *ErrorWriter) Write(b []byte) (int, error) {
	return 0, w.Error
}
