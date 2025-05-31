package wfs

import (
	"io"
	"io/fs"
)

// WFS extends the standard fs.FS interface with write operations
// It combines several fs interfaces (fs.WFS, fs.ReadDirFS, fs.ReadFileFS, fs.StatFS)
// and adds write operations
type WFS interface {
	// Read operations from standard fs interfaces
	fs.FS
	fs.ReadDirFS
	fs.ReadFileFS
	fs.StatFS

	// Write operations
	WriteFile(name string, data []byte, perm fs.FileMode) error
	MkdirAll(path string, perm fs.FileMode) error
	Remove(name string) error
	Create(name string) (io.WriteCloser, error)
}
