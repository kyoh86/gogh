package git

import (
	"context"
	"errors"
	"io"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	coregit "github.com/kyoh86/gogh/v4/core/git"
)

type GitService struct {
	auth                transport.AuthMethod
	cloneProgressWriter io.Writer
}

type Option func(*GitService)

var (
	CloneProgressWriter = func(w io.Writer) Option {
		return func(s *GitService) {
			s.cloneProgressWriter = w
		}
	}
)

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
	_, err := git.PlainCloneContext(ctx, localPath, false, &git.CloneOptions{
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

// ListUntrackedFiles returns a list of untracked files in the repository.
func (s *GitService) ListUntrackedFiles(ctx context.Context, localPath string) ([]string, error) {
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return nil, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, err
	}

	var untrackedFiles []string
	for filePath, fileStatus := range status {
		if fileStatus.Worktree == '?' {
			untrackedFiles = append(untrackedFiles, filePath)
		}
	}

	return untrackedFiles, nil
}

var _ coregit.GitService = (*GitService)(nil)
