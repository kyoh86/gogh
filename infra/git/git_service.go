package git

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	coregit "github.com/kyoh86/gogh/v4/core/git"
)

type GitService struct {
	auth                transport.AuthMethod
	cloneProgressWriter io.Writer
}

type Option func(*GitService)

var CloneProgressWriter = func(w io.Writer) Option {
	return func(s *GitService) {
		s.cloneProgressWriter = w
	}
}

// NewService creates a new Service instance without authentication.
func NewService(options ...Option) *GitService {
	s := &GitService{}
	for _, opt := range options {
		opt(s)
	}
	return s
}

// AuthenticateWithUsernamePassword implements git.GitService.
func (s *GitService) AuthenticateWithUsernamePassword(_ context.Context, username string, password string) (coregit.GitService, error) {
	return &GitService{
		auth: &http.BasicAuth{
			Username: username,
			Password: password,
		},
	}, nil
}

// Clone clones a remote repository to a local path.
func (s *GitService) Clone(ctx context.Context, remoteURL string, localPath string, opts coregit.CloneOptions) error {
	_, err := git.PlainCloneContext(ctx, localPath, opts.IsBare, &git.CloneOptions{
		URL:      remoteURL,
		Auth:     s.auth,
		Progress: s.cloneProgressWriter,
	})
	switch {
	case errors.Is(err, git.ErrRepositoryNotExists) || errors.Is(err, transport.ErrAuthenticationRequired) || errors.Is(err, transport.ErrAuthorizationFailed) || errors.Is(err, transport.ErrRepositoryNotFound):
		return coregit.ErrRepositoryNotExists
	case errors.Is(err, transport.ErrEmptyRemoteRepository):
		return coregit.ErrRepositoryEmpty
	}
	return err
}

// Init initializes a new git repository at the specified local path.
func (s *GitService) Init(_ context.Context, remoteURL, localPath string, isBare bool, _ coregit.InitOptions) error {
	repo, err := git.PlainInit(localPath, isBare)
	if err != nil {
		return err
	}
	if _, err := repo.CreateRemote(&config.RemoteConfig{
		Name: git.DefaultRemoteName,
		URLs: []string{remoteURL},
	}); err != nil {
		return err
	}
	return nil
}

// IsBare checks if a repository at the given path is a bare repository.
func (s *GitService) IsBare(_ context.Context, localPath string) (bool, error) {
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return false, err
	}
	// Try to get the worktree - if it's a bare repository, this will fail
	_, err = repo.Worktree()
	if err != nil {
		// If the error indicates it's a bare repository, return true
		return strings.Contains(err.Error(), "bare repository") || strings.Contains(err.Error(), "worktree not available"), nil
	}
	return false, nil
}

// AddWorktree creates a new worktree for an existing repository.
func (s *GitService) AddWorktree(_ context.Context, repoPath string, branch string, path string) error {
	// Use git worktree add command since go-git doesn't support AddWorktree yet
	worktreePath := filepath.Join(repoPath, path)

	// First, get the remote's default branch name
	symbolicRefCmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	symbolicRefCmd.Dir = repoPath
	refOutput, err := symbolicRefCmd.Output()
	if err != nil {
		return fmt.Errorf("getting remote HEAD: %w", err)
	}
	remoteHead := strings.TrimSpace(string(refOutput))
	// Extract branch name from refs/remotes/origin/<branch>
	remoteBranch := strings.TrimPrefix(remoteHead, "refs/remotes/origin/")

	// Create worktree with the branch
	// Try to create the branch first, ignore error if it already exists
	cmd := exec.Command("git", "worktree", "add", "-B", branch, worktreePath, "origin/"+remoteBranch)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		// If -B fails, try without -b (branch might already exist)
		cmd = exec.Command("git", "worktree", "add", worktreePath, branch)
		cmd.Dir = repoPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("adding worktree: %w, output: %s", err, string(output))
		}
	}
	return nil
}

// Fetch fetches updates from a remote repository.
func (s *GitService) Fetch(ctx context.Context, repoPath string, remote string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	err = repo.FetchContext(ctx, &git.FetchOptions{
		RemoteName: remote,
		Auth:       s.auth,
	})
	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		return coregit.ErrAlreadyUpToDate
	}
	return err
}

// CreateBranch creates a new branch from a starting point.
func (s *GitService) CreateBranch(_ context.Context, repoPath string, branchName string, startPoint string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	// Resolve the starting point to a hash
	hash, err := repo.ResolveRevision(plumbing.Revision(startPoint))
	if err != nil {
		return err
	}

	// Create the branch reference
	return repo.Storer.SetReference(plumbing.NewHashReference(plumbing.ReferenceName("refs/heads/"+branchName), *hash))
}

// SetRemoteHead sets the HEAD symbolic reference for a remote.
func (s *GitService) SetRemoteHead(_ context.Context, repoPath string, remote string) error {
	// Use git remote set-head to set the symbolic reference
	cmd := exec.Command("git", "remote", "set-head", remote, "-a")
	cmd.Dir = repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("setting remote head: %w, output: %s", err, string(output))
	}
	return nil
}

// SetRemotes sets the remote repositories for a local git repository.
func (s *GitService) SetRemotes(
	_ context.Context,
	localPath string,
	name string,
	remotes []string,
) error {
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return err
	}
	cfg, err := repo.Config()
	if err != nil {
		return err
	}
	if cfg.Remotes == nil {
		cfg.Remotes = map[string]*config.RemoteConfig{}
	}
	cfg.Remotes[name] = &config.RemoteConfig{
		Name: name,
		URLs: remotes,
	}
	return repo.SetConfig(cfg)
}

// SetDefaultRemotes sets the default remote repositories for a local git repository.
func (s *GitService) SetDefaultRemotes(
	ctx context.Context,
	localPath string,
	remotes []string,
) error {
	return s.SetRemotes(ctx, localPath, git.DefaultRemoteName, remotes)
}

// GetRemotes retrieves the remote repositories for a local git repository.
func (s *GitService) GetRemotes(
	ctx context.Context,
	localPath string,
	name string,
) ([]string, error) {
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return nil, err
	}
	cfg, err := repo.Config()
	if err != nil {
		return nil, err
	}
	remote, ok := cfg.Remotes[name]
	if !ok {
		return nil, nil
	}
	return remote.URLs, nil
}

// GetDefaultRemotes retrieves the default remote repositories for a local git repository.
func (s *GitService) GetDefaultRemotes(
	ctx context.Context,
	localPath string,
) ([]string, error) {
	return s.GetRemotes(ctx, localPath, git.DefaultRemoteName)
}

// ListExcludedFiles returns a list of excluded/ignored files in the repository.
func (s *GitService) ListExcludedFiles(
	ctx context.Context,
	localPath string,
	filePatterns []string,
) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		localPath, err := filepath.Abs(localPath)
		if err != nil {
			yield("", fmt.Errorf("getting absolute path of repo: %w", err))
			return
		}

		var inclusion gitignore.Matcher
		if len(filePatterns) > 0 {
			var ps []gitignore.Pattern
			domain := strings.Split(filepath.ToSlash(localPath), "/")
			for _, p := range filePatterns {
				ps = append(ps, gitignore.ParsePattern(p, domain))
			}
			inclusion = gitignore.NewMatcher(ps)
		}

		userExcludes, err := LoadUserExcludes(localPath)
		if err != nil {
			yield("", fmt.Errorf("loading user excludes: %w", err))
			return
		}
		localExcludes, err := LoadLocalExcludes(localPath)
		if err != nil {
			yield("", fmt.Errorf("loading local excludes: %w", err))
			return
		}

		ignores := map[string][]gitignore.Pattern{}
		var exDirs [][]string
		if err := filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if info.Name() == ".git" {
					return filepath.SkipDir
				}
				patterns, err := LoadLocalIgnore(path)
				if err != nil {
					return err
				}
				if len(patterns) > 0 {
					ignores[path] = patterns
				}
			}

			// build matcher to ignore excluded files in the current directory
			var excludes []gitignore.Pattern
			for traversePath := path; ; traversePath = filepath.Dir(traversePath) {
				patterns := ignores[traversePath]
				if len(patterns) > 0 {
					excludes = slices.Concat(patterns, excludes)
				}
				if traversePath == localPath {
					break
				}
				if traversePath == filepath.Dir(traversePath) {
					break
				}
			}
			exclusion := gitignore.NewMatcher(slices.Concat(userExcludes, localExcludes, excludes))

			pathWords := strings.Split(filepath.ToSlash(path), "/")
			// Check if the file is excluded by user or local excludes
			if exclusion.Match(pathWords, info.IsDir()) {
				if info.IsDir() {
					// If it's a directory, we need to remember it for later
					exDirs = append(exDirs, pathWords)
				} else if inclusion == nil || inclusion.Match(pathWords, info.IsDir()) {
					// If it's a file, we can add it directly to the result
					if !yield(path, nil) {
						return filepath.SkipAll
					}
				}
				return nil
			}
			if info.IsDir() {
				return nil
			}
			// If the file is not excluded, we check if it is in an ignored directory
			for _, dir := range exDirs {
				if len(pathWords) <= len(dir) {
					continue
				}
				if slices.Equal(pathWords[:len(dir)], dir) {
					if inclusion == nil || inclusion.Match(pathWords, info.IsDir()) {
						if !yield(path, nil) {
							return filepath.SkipAll
						}
					}
					return nil
				}
			}
			return nil
		}); err != nil {
			yield("", fmt.Errorf("walking repository path %q: %w", localPath, err))
			return
		}
	}
}

// ListAllFiles returns a list of all files in the repository.
// It includes `.git` directory.
func (s *GitService) ListAllFiles(
	ctx context.Context,
	localPath string,
	filePatterns []string,
) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		localPath, err := filepath.Abs(localPath)
		if err != nil {
			yield("", fmt.Errorf("getting absolute path of repo: %w", err))
			return
		}

		var inclusion gitignore.Matcher
		if len(filePatterns) > 0 {
			var ps []gitignore.Pattern
			domain := strings.Split(filepath.ToSlash(localPath), "/")
			for _, p := range filePatterns {
				ps = append(ps, gitignore.ParsePattern(p, domain))
			}
			inclusion = gitignore.NewMatcher(ps)
		}

		if err := filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if inclusion != nil && !inclusion.Match(strings.Split(filepath.ToSlash(path), "/"), info.IsDir()) {
				return nil
			}
			if !yield(path, nil) {
				return filepath.SkipAll
			}
			return nil
		}); err != nil {
			yield("", fmt.Errorf("walking repository path %q: %w", localPath, err))
			return
		}
	}
}

var _ coregit.GitService = (*GitService)(nil)
