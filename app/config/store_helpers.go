package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// loadTOMLFile loads TOML data from a file
func loadTOMLFile[T any](source string) (*T, error) {
	file, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data T
	if err := toml.NewDecoder(file).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode TOML: %w", err)
	}
	return &data, nil
}

// saveTOMLFile saves data to a TOML file, ensuring the directory exists
func saveTOMLFile[T any](source string, data T) error {
	if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	file, err := os.OpenFile(source, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(data); err != nil {
		return fmt.Errorf("encode TOML: %w", err)
	}
	return nil
}

// ensureDirectoryExists creates the directory for the given file path if it doesn't exist
func ensureDirectoryExists(filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create directory %s: %w", dir, err)
	}
	return nil
}

// openFileForWrite opens a file for writing, creating it if necessary
func openFileForWrite(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
}
