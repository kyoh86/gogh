package filesystem

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/v4/core/wfs"
)

// LocalWFS implements FS interface using the local file system
type LocalWFS struct {
	// Root directory for all operations
	root string
}

// NewLocalWFS creates a new LocalFS with the given root directory
func NewLocalWFS(root string) *LocalWFS {
	return &LocalWFS{root: root}
}

// Open implements fs.FS
func (l *LocalWFS) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(l.root, name))
}

// ReadDir implements fs.ReadDirFS
func (l *LocalWFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(filepath.Join(l.root, name))
}

// ReadFile implements fs.ReadFileFS
func (l *LocalWFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(l.root, name))
}

// Stat implements fs.StatFS
func (l *LocalWFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(filepath.Join(l.root, name))
}

// WriteFile writes data to the named file
func (l *LocalWFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(filepath.Join(l.root, name), data, perm)
}

// MkdirAll creates a directory and all necessary parents
func (l *LocalWFS) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(filepath.Join(l.root, path), perm)
}

// Remove removes the named file or directory
func (l *LocalWFS) Remove(name string) error {
	return os.Remove(filepath.Join(l.root, name))
}

// Create creates the named file
func (l *LocalWFS) Create(name string) (io.WriteCloser, error) {
	return os.Create(filepath.Join(l.root, name))
}

var _ wfs.WFS = (*LocalWFS)(nil)
