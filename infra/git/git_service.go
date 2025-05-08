package git

import (
	"context"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	coregit "github.com/kyoh86/gogh/v3/core/git"
)

type GitService struct {
	// Dependencies
	auth transport.AuthMethod
}

// NewService creates a new Service instance with the given username and password
// for HTTP basic authentication.
func NewAuthenticatedService(username string, password string) *GitService {
	return &GitService{
		auth: &http.BasicAuth{
			Username: username,
			Password: password,
		},
	}
}

// NewService creates a new Service instance without authentication.
func NewService() *GitService {
	return &GitService{}
}

// Clone clones a remote repository to a local path.
func (s *GitService) Clone(ctx context.Context, remoteURL string, localPath string, options coregit.CloneOptions) error {
	_, err := git.PlainCloneContext(ctx, localPath, false, &git.CloneOptions{
		URL:  remoteURL,
		Auth: s.auth,
	})
	return err
}

func (s *GitService) Init(remoteURL, localPath string, isBare bool) error {
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

func (s *GitService) SetDefaultRemotes(
	ctx context.Context,
	localPath string,
	remotes []string,
) error {
	return s.SetRemotes(ctx, localPath, git.DefaultRemoteName, remotes)
}

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

func (s *GitService) GetDefaultRemotes(
	ctx context.Context,
	localPath string,
) ([]string, error) {
	return s.GetRemotes(ctx, localPath, git.DefaultRemoteName)
}

// Ensure GitService implements core.GitService
var _ coregit.GitService = (*GitService)(nil)
