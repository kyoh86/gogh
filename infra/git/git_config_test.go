package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadUserExcludes(t *testing.T) {
	// Create a temporary directory for our test files
	tmpDir, err := os.MkdirTemp("", "git-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a mock repo path
	repoPath := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo dir: %v", err)
	}

	// Save original function and restore later
	originalUserExcludesFile := UserExcludesFile
	defer func() { UserExcludesFile = originalUserExcludesFile }()

	// Create a test gitignore file
	testExcludesPath := filepath.Join(tmpDir, "gitignore")
	testExcludesContent := `
# Comment line
*.log
/node_modules
build/
`
	if err := os.WriteFile(testExcludesPath, []byte(testExcludesContent), 0644); err != nil {
		t.Fatalf("Failed to write test excludes file: %v", err)
	}

	// Set the mock function
	UserExcludesFile = func() (string, error) {
		return testExcludesPath, nil
	}

	// Test the function
	patterns, err := LoadUserExcludes(repoPath)
	if err != nil {
		t.Fatalf("LoadUserExcludes failed: %v", err)
	}

	// Verify results
	if len(patterns) != 3 {
		t.Errorf("Expected 3 patterns, got %d", len(patterns))
	}

	// Test with non-existent file
	UserExcludesFile = func() (string, error) {
		return filepath.Join(tmpDir, "nonexistent"), nil
	}

	patterns, err = LoadUserExcludes(repoPath)
	if err != nil {
		t.Fatalf("LoadUserExcludes should not return error for non-existent file: %v", err)
	}
	if len(patterns) != 0 {
		t.Errorf("Expected 0 patterns for non-existent file, got %d", len(patterns))
	}
}

func TestLoadLocalExcludes(t *testing.T) {
	// Create a temporary directory for our test files
	tmpDir, err := os.MkdirTemp("", "git-local-excludes-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a mock repo with .git/info/exclude
	repoPath := filepath.Join(tmpDir, "repo")
	excludeDir := filepath.Join(repoPath, ".git", "info")
	if err := os.MkdirAll(excludeDir, 0755); err != nil {
		t.Fatalf("Failed to create exclude dir: %v", err)
	}

	// Create a test exclude file
	excludePath := filepath.Join(excludeDir, "exclude")
	excludeContent := `
# Local excludes
*.tmp
.DS_Store
`
	if err := os.WriteFile(excludePath, []byte(excludeContent), 0644); err != nil {
		t.Fatalf("Failed to write test exclude file: %v", err)
	}

	// Test the function
	patterns, err := LoadLocalExcludes(repoPath)
	if err != nil {
		t.Fatalf("LoadLocalExcludes failed: %v", err)
	}

	// Verify results
	if len(patterns) != 2 {
		t.Errorf("Expected 2 patterns, got %d", len(patterns))
	}

	// Test with non-existent file
	emptyRepoPath := filepath.Join(tmpDir, "empty-repo")
	if err := os.MkdirAll(emptyRepoPath, 0755); err != nil {
		t.Fatalf("Failed to create empty repo dir: %v", err)
	}

	patterns, err = LoadLocalExcludes(emptyRepoPath)
	if err != nil {
		t.Fatalf("LoadLocalExcludes should not return error for non-existent file: %v", err)
	}
	if len(patterns) != 0 {
		t.Errorf("Expected 0 patterns for non-existent file, got %d", len(patterns))
	}
}

func TestLoadLocalIgnore(t *testing.T) {
	// Create a temporary directory for our test files
	tmpDir, err := os.MkdirTemp("", "git-local-ignore-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a mock repo with .gitignore
	repoPath := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo dir: %v", err)
	}

	// Create a test .gitignore file
	ignoreContent := `
# Project specific ignores
/dist
/coverage
.env
`
	if err := os.WriteFile(filepath.Join(repoPath, ".gitignore"), []byte(ignoreContent), 0644); err != nil {
		t.Fatalf("Failed to write test .gitignore file: %v", err)
	}

	// Test the function
	patterns, err := LoadLocalIgnore(repoPath)
	if err != nil {
		t.Fatalf("LoadLocalIgnore failed: %v", err)
	}

	// Verify results
	if len(patterns) != 3 {
		t.Errorf("Expected 3 patterns, got %d", len(patterns))
	}

	// Test with non-existent file
	emptyRepoPath := filepath.Join(tmpDir, "empty-repo")
	if err := os.MkdirAll(emptyRepoPath, 0755); err != nil {
		t.Fatalf("Failed to create empty repo dir: %v", err)
	}

	patterns, err = LoadLocalIgnore(emptyRepoPath)
	if err != nil {
		t.Fatalf("LoadLocalIgnore should not return error for non-existent file: %v", err)
	}
	if len(patterns) != 0 {
		t.Errorf("Expected 0 patterns for non-existent file, got %d", len(patterns))
	}
}

func TestDefaultUserExcludesFile(t *testing.T) {
	// Save original functions and restore later
	originalUserConfigDir := osUserConfigDir
	originalUserHomeDir := osUserHomeDir
	defer func() {
		osUserConfigDir = originalUserConfigDir
		osUserHomeDir = originalUserHomeDir
	}()

	// Create a temporary directory for our test files
	tmpDir, err := os.MkdirTemp("", "git-user-excludes-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test directories
	configDir := filepath.Join(tmpDir, "config", "git")
	homeDir := filepath.Join(tmpDir, "home")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatalf("Failed to create home dir: %v", err)
	}

	// Set up mock functions
	osUserConfigDir = func() (string, error) {
		return filepath.Dir(configDir), nil
	}
	osUserHomeDir = func() (string, error) {
		return homeDir, nil
	}

	// Test 1: XDG config with excludesfile setting
	xdgConfigContent := `[core]
	excludesfile = /path/to/xdg/gitignore
`
	if err := os.WriteFile(filepath.Join(configDir, "config"), []byte(xdgConfigContent), 0644); err != nil {
		t.Fatalf("Failed to write XDG config file: %v", err)
	}

	excludesFile, err := defaultUserExcludesFile()
	if err != nil {
		t.Fatalf("defaultUserExcludesFile failed: %v", err)
	}
	if excludesFile != "/path/to/xdg/gitignore" {
		t.Errorf("Expected XDG gitignore path, got %q", excludesFile)
	}

	// Test 2: Home gitconfig overrides XDG config
	homeConfigContent := `[core]
	excludesfile = /path/to/home/gitignore
`
	if err := os.WriteFile(filepath.Join(homeDir, ".gitconfig"), []byte(homeConfigContent), 0644); err != nil {
		t.Fatalf("Failed to write home gitconfig file: %v", err)
	}

	excludesFile, err = defaultUserExcludesFile()
	if err != nil {
		t.Fatalf("defaultUserExcludesFile failed: %v", err)
	}
	if excludesFile != "/path/to/home/gitignore" {
		t.Errorf("Expected home gitignore path, got %q", excludesFile)
	}

	// Test 3: When no configs exist (should return default path)
	os.Remove(filepath.Join(configDir, "config"))
	os.Remove(filepath.Join(homeDir, ".gitconfig"))

	excludesFile, err = defaultUserExcludesFile()
	if err != nil {
		t.Fatalf("defaultUserExcludesFile failed: %v", err)
	}
	expectedDefault := filepath.Join(filepath.Dir(configDir), "git", "ignore")
	if excludesFile != expectedDefault {
		t.Errorf("Expected default gitignore path %q, got %q", expectedDefault, excludesFile)
	}

	// Test 4: Error handling for UserConfigDir
	osUserConfigDir = func() (string, error) { return "", os.ErrNotExist }

	_, err = defaultUserExcludesFile()
	if err == nil {
		t.Error("Expected error for UserConfigDir failure, got nil")
	}
}

func TestEnsureExcludeFile(t *testing.T) {
	// Create a temporary directory for our test files
	tmpDir, err := os.MkdirTemp("", "git-ensure-exclude-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test 1: Valid git config file with excludesfile
	configContent := `[core]
	excludesfile = /path/to/gitignore
`
	configPath := filepath.Join(tmpDir, "config")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	path, err := ensureExcludeFile(configPath)
	if err != nil {
		t.Fatalf("ensureExcludeFile failed: %v", err)
	}
	if path != "/path/to/gitignore" {
		t.Errorf("Expected /path/to/gitignore, got %q", path)
	}

	// Test 2: Git config without excludesfile
	configContent = `[core]
	editor = vim
`
	configPath = filepath.Join(tmpDir, "config2")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	path, err = ensureExcludeFile(configPath)
	if err != nil {
		t.Fatalf("ensureExcludeFile failed: %v", err)
	}
	if path != "" {
		t.Errorf("Expected empty path, got %q", path)
	}

	// Test 3: Non-existent file
	path, err = ensureExcludeFile(filepath.Join(tmpDir, "nonexistent"))
	if err != nil {
		t.Fatalf("ensureExcludeFile should not return error for non-existent file: %v", err)
	}
	if path != "" {
		t.Errorf("Expected empty path for non-existent file, got %q", path)
	}

	// Test 4: Invalid git config format
	configContent = `This is not a valid git config file`
	configPath = filepath.Join(tmpDir, "invalid")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write invalid config file: %v", err)
	}

	_, err = ensureExcludeFile(configPath)
	if err == nil {
		t.Error("Expected error for invalid config format, got nil")
	}
}

func TestReadIgnoreFile(t *testing.T) {
	// Create a temporary directory for our test files
	tmpDir, err := os.MkdirTemp("", "git-read-ignore-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test ignore file
	ignoreContent := `
# Comment line
*.log
/node_modules
build/

# Empty lines should be skipped

`
	ignoreFile := filepath.Join(tmpDir, "ignore")
	if err := os.WriteFile(ignoreFile, []byte(ignoreContent), 0644); err != nil {
		t.Fatalf("Failed to write test ignore file: %v", err)
	}

	// Test reading a valid file
	domain := []string{"example", "com", "repo"}
	patterns, err := readIgnoreFile(ignoreFile, domain)
	if err != nil {
		t.Fatalf("readIgnoreFile failed: %v", err)
	}

	// Verify patterns
	if len(patterns) != 3 {
		t.Errorf("Expected 3 patterns, got %d", len(patterns))
	}

	// Test with non-existent file
	patterns, err = readIgnoreFile(filepath.Join(tmpDir, "nonexistent"), domain)
	if err != nil {
		t.Fatalf("readIgnoreFile should not return error for non-existent file: %v", err)
	}
	if len(patterns) != 0 {
		t.Errorf("Expected 0 patterns for non-existent file, got %d", len(patterns))
	}
}
