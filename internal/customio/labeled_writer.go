package customio

import (
	"io"
	"sync"
)

type LabeledWriter struct {
	Label string
	Once  *sync.Once
	Base  io.Writer
}

func (w *LabeledWriter) Write(p []byte) (retN int, retErr error) {
	w.Once.Do(func() {
		_, err := w.Base.Write([]byte(w.Label + "\n"))
		retErr = err
	})
	if retErr != nil {
		return 0, retErr
	}
	n, err := w.Base.Write(p)
	retN = n
	retErr = err
	return
}
