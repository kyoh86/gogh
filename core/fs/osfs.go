package fs

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// OSFS implements FS interface using the local file system
type OSFS struct {
	// Root directory for all operations
	root string
}

// NewOSFS creates a new OSFS with the given root directory
func NewOSFS(root string) *OSFS {
	return &OSFS{root: root}
}

// Open implements fs.FS
func (l *OSFS) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(l.root, name))
}

// ReadDir implements fs.ReadDirFS
func (l *OSFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(filepath.Join(l.root, name))
}

// ReadFile implements fs.ReadFileFS
func (l *OSFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(l.root, name))
}

// Stat implements fs.StatFS
func (l *OSFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(filepath.Join(l.root, name))
}

// WriteFile writes data to the named file
func (l *OSFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(filepath.Join(l.root, name), data, perm)
}

// MkdirAll creates a directory and all necessary parents
func (l *OSFS) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(filepath.Join(l.root, path), perm)
}

// Remove removes the named file or directory
func (l *OSFS) Remove(name string) error {
	return os.Remove(filepath.Join(l.root, name))
}

// Create creates the named file
func (l *OSFS) Create(name string) (io.WriteCloser, error) {
	return os.Create(filepath.Join(l.root, name))
}

var _ FS = (*OSFS)(nil)
