package wfs_mock

import (
	"bytes"
	"io"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/kyoh86/gogh/v4/core/wfs"
)

// MockWFS is a mock implementation of wfs.WFS for testing
type MockWFS struct {
	files    map[string][]byte
	dirItems map[string][]fs.DirEntry
	errors   map[string]error
}

// NewMockWFS creates a new mock filesystem
func NewMockWFS() *MockWFS {
	return &MockWFS{
		files:    make(map[string][]byte),
		dirItems: make(map[string][]fs.DirEntry),
		errors:   make(map[string]error),
	}
}

// normalizePath ensures consistent path handling between "" and "."
func normalizePath(path string) string {
	if path == "." {
		return ""
	}
	return path
}

// SetError sets an error to be returned for a specific operation on a path
func (m *MockWFS) SetError(op, path string, err error) {
	path = normalizePath(path)
	m.errors[op+":"+path] = err
}

// addFile adds a file to the mock filesystem and updates directory entries
func (m *MockWFS) addFile(path string, content []byte) {
	path = normalizePath(path)
	m.files[path] = content

	// Update directory entries
	dir := normalizePath(filepath.Dir(path))

	if _, exists := m.dirItems[dir]; !exists {
		m.dirItems[dir] = []fs.DirEntry{}
	}

	// Check if entry already exists
	baseName := filepath.Base(path)
	for _, entry := range m.dirItems[dir] {
		if entry.Name() == baseName {
			return // Entry already exists
		}
	}

	// Add new entry
	m.dirItems[dir] = append(m.dirItems[dir], &MockDirEntry{
		name:  baseName,
		isDir: false,
	})
}

// Files returns the current files in the mock filesystem (for debugging)
func (m *MockWFS) Files() map[string][]byte {
	return m.files
}

// DirEntries returns the current directory entries (for debugging)
func (m *MockWFS) DirEntries() map[string][]fs.DirEntry {
	return m.dirItems
}

// Open implements wfs.WFS
func (m *MockWFS) Open(name string) (fs.File, error) {
	name = normalizePath(name)
	if err, exists := m.errors["Open:"+name]; exists && err != nil {
		return nil, err
	}

	content, exists := m.files[name]
	if !exists {
		return nil, fs.ErrNotExist
	}

	return &MockFile{
		name:    name,
		content: content,
		pos:     0,
		closed:  false,
	}, nil
}

// ReadDir implements wfs.WFS
func (m *MockWFS) ReadDir(name string) ([]fs.DirEntry, error) {
	name = normalizePath(name)
	if err, exists := m.errors["ReadDir:"+name]; exists && err != nil {
		return nil, err
	}

	entries, exists := m.dirItems[name]
	if !exists {
		return nil, fs.ErrNotExist
	}

	return entries, nil
}

// ReadFile implements wfs.WFS
func (m *MockWFS) ReadFile(name string) ([]byte, error) {
	name = normalizePath(name)
	if err, exists := m.errors["ReadFile:"+name]; exists && err != nil {
		return nil, err
	}

	content, exists := m.files[name]
	if !exists {
		return nil, fs.ErrNotExist
	}

	return content, nil
}

// Stat implements wfs.WFS
func (m *MockWFS) Stat(name string) (fs.FileInfo, error) {
	name = normalizePath(name)
	if err, exists := m.errors["Stat:"+name]; exists && err != nil {
		return nil, err
	}

	content, exists := m.files[name]
	if !exists {
		// Check if it's a directory
		if _, dirExists := m.dirItems[name]; dirExists {
			return &MockFileInfo{
				name:  filepath.Base(name),
				isDir: true,
			}, nil
		}
		return nil, fs.ErrNotExist
	}

	return &MockFileInfo{
		name:  filepath.Base(name),
		size:  int64(len(content)),
		isDir: false,
	}, nil
}

// WriteFile implements wfs.WFS
func (m *MockWFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	name = normalizePath(name)
	if err, exists := m.errors["WriteFile:"+name]; exists && err != nil {
		return err
	}

	m.addFile(name, data)
	return nil
}

// MkdirAll implements wfs.WFS
func (m *MockWFS) MkdirAll(path string, perm fs.FileMode) error {
	path = normalizePath(path)
	if err, exists := m.errors["MkdirAll:"+path]; exists && err != nil {
		return err
	}

	// Add directory entry
	if _, exists := m.dirItems[path]; !exists {
		m.dirItems[path] = []fs.DirEntry{}
	}

	// Add parent directories too
	parent := filepath.Dir(path)
	parent = normalizePath(parent)

	if parent != "" && parent != path {
		if _, exists := m.dirItems[parent]; !exists {
			return m.MkdirAll(parent, perm)
		}
	}

	return nil
}

// Remove implements wfs.WFS
func (m *MockWFS) Remove(name string) error {
	name = normalizePath(name)
	if err, exists := m.errors["Remove:"+name]; exists && err != nil {
		return err
	}

	if _, exists := m.files[name]; !exists {
		if _, dirExists := m.dirItems[name]; !dirExists {
			return fs.ErrNotExist
		}
	}

	delete(m.files, name)

	// Remove from directory listing
	dir := normalizePath(filepath.Dir(name))
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

// Create implements wfs.WFS
func (m *MockWFS) Create(name string) (io.WriteCloser, error) {
	name = normalizePath(name)
	if err, exists := m.errors["Create:"+name]; exists && err != nil {
		return nil, err
	}

	// Create a buffer that will write to our files map on close
	buffer := &bytes.Buffer{}
	return &MockWriteCloser{
		WriteCloser: NopWriteCloser{buffer},
		onClose: func() error {
			m.addFile(name, buffer.Bytes())
			return nil
		},
	}, nil
}

// MockFile implements fs.File interface
type MockFile struct {
	name    string
	content []byte
	pos     int
	closed  bool
}

// Stat implements fs.File
func (f *MockFile) Stat() (fs.FileInfo, error) {
	if f.closed {
		return nil, fs.ErrClosed
	}
	return &MockFileInfo{
		name:  filepath.Base(f.name),
		size:  int64(len(f.content)),
		isDir: false,
	}, nil
}

// Read implements fs.File
func (f *MockFile) Read(b []byte) (int, error) {
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
func (f *MockFile) Close() error {
	if f.closed {
		return fs.ErrClosed
	}
	f.closed = true
	return nil
}

// MockFileInfo implements fs.FileInfo interface
type MockFileInfo struct {
	name  string
	size  int64
	isDir bool
}

// Name implements fs.FileInfo
func (fi *MockFileInfo) Name() string { return fi.name }

// Size implements fs.FileInfo
func (fi *MockFileInfo) Size() int64 { return fi.size }

// Mode implements fs.FileInfo
func (fi *MockFileInfo) Mode() fs.FileMode {
	if fi.isDir {
		return fs.ModeDir | 0755
	}
	return 0644
}

// ModTime implements fs.FileInfo
func (fi *MockFileInfo) ModTime() time.Time { return time.Now() }

// IsDir implements fs.FileInfo
func (fi *MockFileInfo) IsDir() bool { return fi.isDir }

// Sys implements fs.FileInfo
func (fi *MockFileInfo) Sys() interface{} { return nil }

// MockDirEntry implements fs.DirEntry interface
type MockDirEntry struct {
	name  string
	isDir bool
	info  fs.FileInfo
}

// Name implements fs.DirEntry
func (e *MockDirEntry) Name() string { return e.name }

// IsDir implements fs.DirEntry
func (e *MockDirEntry) IsDir() bool { return e.isDir }

// Type implements fs.DirEntry
func (e *MockDirEntry) Type() fs.FileMode {
	if e.isDir {
		return fs.ModeDir
	}
	return 0
}

// Info implements fs.DirEntry
func (e *MockDirEntry) Info() (fs.FileInfo, error) {
	if e.info != nil {
		return e.info, nil
	}
	return &MockFileInfo{name: e.name, isDir: e.isDir}, nil
}

// NopWriteCloser is a WriteCloser with a no-op Close method
type NopWriteCloser struct {
	io.Writer
}

// Close implements io.Closer
func (NopWriteCloser) Close() error { return nil }

// MockWriteCloser is a wrapper that calls onClose when Close is called
type MockWriteCloser struct {
	io.WriteCloser
	onClose func() error
}

// Close implements io.Closer
func (w *MockWriteCloser) Close() error {
	err1 := w.WriteCloser.Close()
	err2 := w.onClose()
	if err1 != nil {
		return err1
	}
	return err2
}

// Ensure MockWFS implements wfs.WFS
var _ wfs.WFS = (*MockWFS)(nil)
