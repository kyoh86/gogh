package mock

import (
	"bytes"
	"io"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/kyoh86/gogh/v4/core/wfs"
)

// MockWFS implements wfs.MockWFS interface for testing
type MockWFS struct {
	files    map[string][]byte
	dirItems map[string][]fs.DirEntry
	errors   map[string]error // key is operation:path
}

// NewMockFS creates a new mock file system
func NewMockFS() *MockWFS {
	return &MockWFS{
		files:    make(map[string][]byte),
		dirItems: make(map[string][]fs.DirEntry),
		errors:   make(map[string]error),
	}
}

// SetError configures an error to be returned for a specific operation and path
func (m *MockWFS) SetError(operation, path string, err error) {
	m.errors[operation+":"+path] = err
}

// AddFile adds a file with content to the mock file system
func (m *MockWFS) AddFile(path string, content []byte) {
	m.files[path] = content
	dir := filepath.Dir(path)
	if _, exists := m.dirItems[dir]; !exists {
		m.dirItems[dir] = []fs.DirEntry{}
	}
	m.dirItems[dir] = append(m.dirItems[dir], &DirEntry{
		name:  filepath.Base(path),
		isDir: false,
	})
}

// Open implements fs.FS
func (m *MockWFS) Open(name string) (fs.File, error) {
	if err, exists := m.errors["Open:"+name]; exists {
		return nil, err
	}

	content, exists := m.files[name]
	if !exists {
		return nil, fs.ErrNotExist
	}

	return &File{
		name:    name,
		content: content,
		pos:     0,
		closed:  false,
	}, nil
}

// ReadDir implements fs.ReadDirFS
func (m *MockWFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if err, exists := m.errors["ReadDir:"+name]; exists {
		return nil, err
	}

	entries, exists := m.dirItems[name]
	if !exists {
		return nil, fs.ErrNotExist
	}

	return entries, nil
}

// ReadFile implements fs.ReadFileFS
func (m *MockWFS) ReadFile(name string) ([]byte, error) {
	if err, exists := m.errors["ReadFile:"+name]; exists {
		return nil, err
	}

	content, exists := m.files[name]
	if !exists {
		return nil, fs.ErrNotExist
	}

	return content, nil
}

// Stat implements fs.StatFS
func (m *MockWFS) Stat(name string) (fs.FileInfo, error) {
	if err, exists := m.errors["Stat:"+name]; exists {
		return nil, err
	}

	content, exists := m.files[name]
	if !exists {
		// Check if it's a directory
		if _, exists := m.dirItems[name]; exists {
			return &FileInfo{
				name:  filepath.Base(name),
				isDir: true,
			}, nil
		}
		return nil, fs.ErrNotExist
	}

	return &FileInfo{
		name:  filepath.Base(name),
		size:  int64(len(content)),
		isDir: false,
	}, nil
}

// WriteFile implements wfs.FS
func (m *MockWFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	if err, exists := m.errors["WriteFile:"+name]; exists {
		return err
	}

	m.AddFile(name, data)
	return nil
}

// MkdirAll implements wfs.FS
func (m *MockWFS) MkdirAll(path string, perm fs.FileMode) error {
	if err, exists := m.errors["MkdirAll:"+path]; exists {
		return err
	}

	// Add directory entry
	m.dirItems[path] = []fs.DirEntry{}

	// Add parent directories too
	parent := filepath.Dir(path)
	if parent != "." && parent != "/" && parent != path {
		if _, exists := m.dirItems[parent]; !exists {
			return m.MkdirAll(parent, perm)
		}
	}

	return nil
}

// Remove implements wfs.FS
func (m *MockWFS) Remove(name string) error {
	if err, exists := m.errors["Remove:"+name]; exists {
		return err
	}

	if _, exists := m.files[name]; !exists {
		if _, dirExists := m.dirItems[name]; !dirExists {
			return fs.ErrNotExist
		}
	}

	delete(m.files, name)

	// Remove from directory listing
	dir := filepath.Dir(name)
	if entries, exists := m.dirItems[dir]; exists {
		baseName := filepath.Base(name)
		newEntries := make([]fs.DirEntry, 0, len(entries))
		for _, entry := range entries {
			if entry.Name() != baseName {
				newEntries = append(newEntries, entry)
			}
		}
		m.dirItems[dir] = newEntries
	}

	// If it was a directory, remove it and its contents
	if entries, exists := m.dirItems[name]; exists {
		delete(m.dirItems, name)

		// Remove all files in this directory
		for _, entry := range entries {
			m.Remove(filepath.Join(name, entry.Name()))
		}
	}

	return nil
}

// Create implements wfs.FS
func (m *MockWFS) Create(name string) (io.WriteCloser, error) {
	if err, exists := m.errors["Create:"+name]; exists {
		return nil, err
	}

	// Create a buffer that will write to our files map on close
	buffer := &bytes.Buffer{}
	return &WriteCloser{
		WriteCloser: nopWriteCloser{buffer},
		onClose: func() error {
			m.AddFile(name, buffer.Bytes())
			return nil
		},
	}, nil
}

// File implements fs.File interface
type File struct {
	name    string
	content []byte
	pos     int
	closed  bool
}

// Stat implements fs.File
func (f *File) Stat() (fs.FileInfo, error) {
	if f.closed {
		return nil, fs.ErrClosed
	}
	return &FileInfo{
		name:  filepath.Base(f.name),
		size:  int64(len(f.content)),
		isDir: false,
	}, nil
}

// Read implements fs.File
func (f *File) Read(b []byte) (int, error) {
	if f.closed {
		return 0, fs.ErrClosed
	}
	if f.pos >= len(f.content) {
		return 0, io.EOF
	}
	n := copy(b, f.content[f.pos:])
	f.pos += n
	return n, nil
}

// Close implements fs.File
func (f *File) Close() error {
	if f.closed {
		return fs.ErrClosed
	}
	f.closed = true
	return nil
}

// FileInfo implements fs.FileInfo interface
type FileInfo struct {
	name  string
	size  int64
	isDir bool
}

// Name implements fs.FileInfo
func (fi *FileInfo) Name() string { return fi.name }

// Size implements fs.FileInfo
func (fi *FileInfo) Size() int64 { return fi.size }

// Mode implements fs.FileInfo
func (fi *FileInfo) Mode() fs.FileMode {
	if fi.isDir {
		return fs.ModeDir | 0755
	}
	return 0644
}

// ModTime implements fs.FileInfo
func (fi *FileInfo) ModTime() time.Time { return time.Now() }

// IsDir implements fs.FileInfo
func (fi *FileInfo) IsDir() bool { return fi.isDir }

// Sys implements fs.FileInfo
func (fi *FileInfo) Sys() any { return nil }

// DirEntry implements fs.DirEntry interface
type DirEntry struct {
	name  string
	isDir bool
	info  fs.FileInfo
}

// Name implements fs.DirEntry
func (e *DirEntry) Name() string { return e.name }

// IsDir implements fs.DirEntry
func (e *DirEntry) IsDir() bool { return e.isDir }

// Type implements fs.DirEntry
func (e *DirEntry) Type() fs.FileMode {
	if e.isDir {
		return fs.ModeDir
	}
	return 0
}

// Info implements fs.DirEntry
func (e *DirEntry) Info() (fs.FileInfo, error) {
	if e.info != nil {
		return e.info, nil
	}
	return &FileInfo{name: e.name, isDir: e.isDir}, nil
}

// Helper types for Create method
type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

// WriteCloser is a wrapper that calls onClose when Close is called
type WriteCloser struct {
	io.WriteCloser
	onClose func() error
}

// Close calls the underlying WriteCloser's Close and then onClose
func (w *WriteCloser) Close() error {
	err1 := w.WriteCloser.Close()
	err2 := w.onClose()
	if err1 != nil {
		return err1
	}
	return err2
}

// Ensure FS implements wfs.FS
var _ wfs.WFS = (*MockWFS)(nil)
