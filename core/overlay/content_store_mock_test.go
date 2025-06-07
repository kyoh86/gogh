package overlay

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/spf13/afero"
)

// MockContentStore implements ContentStore using afero
type MockContentStore struct {
	fs      afero.Fs
	baseDir string
	nextID  int
}

func NewMockContentStore() *MockContentStore {
	return &MockContentStore{
		fs:      afero.NewMemMapFs(),
		baseDir: "/content",
		nextID:  1,
	}
}

func (a *MockContentStore) SaveContent(ctx context.Context, ov Overlay, content io.Reader) (string, error) {
	if content == nil {
		return "", errors.New("content is nil")
	}

	// Create base directory if it doesn't exist
	if err := a.fs.MkdirAll(a.baseDir, 0755); err != nil {
		return "", err
	}

	// Generate a unique location
	location := filepath.Join(a.baseDir, fmt.Sprintf("content-%d", a.nextID))
	a.nextID++

	// Create the file
	file, err := a.fs.Create(location)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Copy content to file
	_, err = io.Copy(file, content)
	if err != nil {
		return "", err
	}

	return location, nil
}

func (a *MockContentStore) OpenContent(ctx context.Context, location string) (io.ReadCloser, error) {
	return a.fs.Open(location)
}

func (a *MockContentStore) RemoveContent(ctx context.Context, location string) error {
	return a.fs.Remove(location)
}
