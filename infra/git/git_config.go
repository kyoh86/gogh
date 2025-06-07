package git

/**
* NOTE: I initially wanted to leverage the implementation from go-git/v5/plumbing/format/gitignore,
* but despite having relevant code, it doesn't actually provide the functionality I needed.
* Given go-git's lengthy development cycle and our limited bandwidth for contributing upstream,
* I've opted to implement our own solution instead.
 */

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/config"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/kyoh86/gogh/v4/core/fs"
)

// Function variables that can be mocked in tests
var (
	osUserConfigDir  = os.UserConfigDir
	osUserHomeDir    = os.UserHomeDir
	UserExcludesFile = defaultUserExcludesFile
)

// LoadUserExcludes loads the user's gitignore patterns from the core.excludesfile property.
// It reads the user's gitignore file, which is typically located at
// `$XDG_CONFIG_HOME/git/ignore` or `$HOME/.config/git/ignore`.
// If the core.excludesfile property is not declared in the user's git config,
// it defaults to the same path.
// See: UserExcludesFile for more details on how the file is determined.
func LoadUserExcludes(repoPath string) ([]gitignore.Pattern, error) {
	abs, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, fmt.Errorf("getting absolute path of repo: %w", err)
	}
	domain := strings.Split(filepath.ToSlash(abs), "/")
	uif, err := UserExcludesFile()
	if err != nil {
		return nil, fmt.Errorf("searching user gitignore file: %w", err)
	}
	ps, err := readIgnoreFile(uif, domain)
	if os.IsNotExist(err) {
		return nil, nil
	}
	return ps, err
}

// LoadLocalExcludes loads the local gitignore patterns from the .git/info/exclude file.
func LoadLocalExcludes(repoPath string) ([]gitignore.Pattern, error) {
	abs, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, fmt.Errorf("getting absolute path of repo: %w", err)
	}
	domain := strings.Split(filepath.ToSlash(abs), "/")
	excludeFile := filepath.Join(repoPath, ".git", "info", "exclude")
	ps, err := readIgnoreFile(excludeFile, domain)
	if os.IsNotExist(err) {
		return nil, nil
	}
	return ps, err
}

// LoadLocalIgnore loads the local gitignore patterns from the .gitignore file.
func LoadLocalIgnore(repoPath string) ([]gitignore.Pattern, error) {
	repoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, fmt.Errorf("getting absolute path of repo: %w", err)
	}
	domain := strings.Split(filepath.ToSlash(repoPath), "/")
	ignoreFile := filepath.Join(repoPath, ".gitignore")
	ps, err := readIgnoreFile(ignoreFile, domain)
	if os.IsNotExist(err) {
		return nil, nil
	}
	return ps, err
}

// defaultUserExcludesFile returns the path to the user's gitignore file.
// It reads the core.excludesfile property from the user-specific configuration files.
// (`$XDG_CONFIG_HOME/git/config` and `~/.gitconfig`)
// When the XDG_CONFIG_HOME environment variable is not set or empty, $HOME/.config/ is used as $XDG_CONFIG_HOME.
// These are also called "global" configuration files. If both files exist, both files are read in the order given above.
// If the core.excludesfile property is not declared in either of them, the function returns the default
// `$XDG_CONFIG_HOME/git/ignore` or `$HOME/.config/git/ignore` file.
// See:
// - Git documentation on configuration file locations: https://git-scm.com/docs/git-config#FILES
// - Git documentation on core.excludesFile: https://git-scm.com/docs/git-config#Documentation/git-config.txt-coreexcludesFile
func defaultUserExcludesFile() (filename string, _ error) {
	// Get the excludesfile property from the XDG git config file
	config, err := osUserConfigDir()
	if err != nil {
		return "", fmt.Errorf("searching user config dir: %w", err)
	}
	filename, err = ensureExcludeFile(filepath.Join(config, "git", "config"))
	if err != nil {
		return "", fmt.Errorf("searching user git config: %w", err)
	}

	// Override with the ~/.gitconfig if it exists
	{
		home, err := osUserHomeDir()
		if err != nil {
			return "", fmt.Errorf("searching user home dir: %w", err)
		}
		candidate, err := ensureExcludeFile(filepath.Join(home, ".gitconfig"))
		if err != nil {
			return "", fmt.Errorf("searching user git config: %w", err)
		}
		if candidate != "" {
			filename = candidate
		}
	}

	// If the core.excludesfile property is not declared in either of them,
	// return the default gitignore file path.
	if filename == "" {
		return filepath.Join(config, "git", "ignore"), nil
	}
	return filename, nil
}

func ensureExcludeFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	d := config.NewDecoder(bytes.NewBuffer(b))

	raw := config.New()
	if err = d.Decode(raw); err != nil {
		return "", err
	}

	s := raw.Section("core")
	efo := s.Options.Get("excludesfile")
	return efo, nil
}

// readIgnoreFile reads a specific git ignore file.
func readIgnoreFile(ignoreFile string, domain []string) (ps []gitignore.Pattern, err error) {
	ignoreFile, _ = fs.ReplaceTildeWithHome(ignoreFile)

	f, err := os.Open(ignoreFile)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err == nil {
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			s := scanner.Text()
			if !strings.HasPrefix(s, "#") && len(strings.TrimSpace(s)) > 0 {
				ps = append(ps, gitignore.ParsePattern(s, domain))
			}
		}
	}
	return
}
