package overlay

import (
	"context"
	"errors"
	"io"
	"path/filepath"

	"github.com/spf13/afero"
)

// MockContentStore implements ContentStore using afero
type MockContentStore struct {
	fs      afero.Fs
	baseDir string
}

func NewMockContentStore() *MockContentStore {
	return &MockContentStore{
		fs:      afero.NewMemMapFs(),
		baseDir: "/content",
	}
}

func (a *MockContentStore) Save(ctx context.Context, overlayID string, content io.Reader) error {
	if content == nil {
		return errors.New("content is nil")
	}

	// Create base directory if it doesn't exist
	if err := a.fs.MkdirAll(a.baseDir, 0o755); err != nil {
		return err
	}

	// Generate a unique location
	location := filepath.Join(a.baseDir, overlayID)

	// Create the file
	file, err := a.fs.Create(location)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy content to file
	_, err = io.Copy(file, content)
	if err != nil {
		return err
	}

	return nil
}

func (a *MockContentStore) Open(ctx context.Context, overlayID string) (io.ReadCloser, error) {
	location := filepath.Join(a.baseDir, overlayID)
	return a.fs.Open(location)
}

func (a *MockContentStore) Remove(ctx context.Context, overlayID string) error {
	location := filepath.Join(a.baseDir, overlayID)
	return a.fs.Remove(location)
}
