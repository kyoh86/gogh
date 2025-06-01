package overlay_apply

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// UseCase represents the create use case
type UseCase struct {
}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (uc *UseCase) Execute(ctx context.Context, repoPath string, relativePath string, content io.Reader) error {
	targetPath := filepath.Join(repoPath, relativePath)

	// Ensure the directory exists
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("creating directory '%s': %w", targetDir, err)
	}

	// Read content
	data, err := io.ReadAll(content)
	if err != nil {
		return fmt.Errorf("reading overlay content: %w", err)
	}

	// Write to target file
	if err := os.WriteFile(targetPath, data, 0644); err != nil {
		return fmt.Errorf("writing to file '%s': %w", targetPath, err)
	}
	return nil
}
