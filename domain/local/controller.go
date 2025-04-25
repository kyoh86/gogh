package local

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/apex/log"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/kyoh86/gogh/v3/domain/reporef"
	"github.com/saracen/walker"
)

const DefaultRootDirName = "Projects"

func NewController(root string) *Controller {
	return &Controller{root: root}
}

type Controller struct {
	// UNDONE: support fs.FS
	// UNDONE: support fs.FS

	// NOTE: v1 -> v2 diferrence
	// if we wanna manage multiple root, create multiple controller instances.
	root string
}

type ExistOption struct {
}

func (l *Controller) Exist(
	ctx context.Context,
	ref reporef.RepoRef,
	opt *ExistOption,
) (bool, error) {
	repo := NewRepo(l.root, ref)
	_, err := git.PlainOpen(repo.FullFilePath())
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, git.ErrRepositoryNotExists):
		return false, nil
	default:
		return false, err
	}
}

type CreateOption struct { // UNDONE: support isBare
}

func (l *Controller) Create(
	ctx context.Context,
	ref reporef.RepoRef,
	opt *CreateOption,
) (Repo, error) {
	p := NewRepo(l.root, ref)
	if err := CreateRepo(ctx, p, ref.URL(), opt); err != nil {
		return Repo{}, err
	}
	return p, nil
}

func CreateRepo(
	_ context.Context,
	localRepo Repo,
	remoteURL string,
	_ *CreateOption,
) error {
	repo, err := git.PlainInit(localRepo.FullFilePath(), false)
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

func (l *Controller) SetRemoteRefs(
	ctx context.Context,
	newRef reporef.RepoRef,
	remotes map[string][]reporef.RepoRef,
) error {
	urls := map[string][]string{}
	for name, refs := range remotes {
		for _, ref := range refs {
			urls[name] = append(urls[name], ref.URL())
		}
	}
	return l.SetRemoteURLs(ctx, newRef, urls)
}

func (l *Controller) SetRemoteURLs(
	ctx context.Context,
	newRef reporef.RepoRef,
	remotes map[string][]string,
) error {
	return SetRemoteURLs(ctx, NewRepo(l.root, newRef), remotes)
}

func SetRemoteURLs(
	_ context.Context,
	locRepo Repo,
	remotes map[string][]string,
) error {
	repo, err := git.PlainOpen(locRepo.FullFilePath())
	if err != nil {
		return err
	}
	cfg, err := repo.Config()
	if err != nil {
		return err
	}
	cfg.Remotes = map[string]*config.RemoteConfig{}
	for name, urls := range remotes {
		cfg.Remotes[name] = &config.RemoteConfig{
			Name: name,
			URLs: urls,
		}
	}
	return repo.SetConfig(cfg)
}

func (l *Controller) GetRemoteURLs(
	ctx context.Context,
	ref reporef.RepoRef,
	name string,
) ([]string, error) {
	return GetRemoteURLs(ctx, NewRepo(l.root, ref), name)
}

func GetRemoteURLs(
	_ context.Context,
	locRepo Repo,
	name string,
) ([]string, error) {
	repo, err := git.PlainOpen(locRepo.FullFilePath())
	if err != nil {
		return nil, fmt.Errorf("open local repository: %w", err)
	}
	remote, err := repo.Remote(name)
	if err != nil {
		return nil, fmt.Errorf("get remote %s: %w", name, err)
	}
	return remote.Config().URLs, nil
}

func GetDefaultRemoteURL(_ context.Context, locRepo Repo) (string, error) {
	urls, err := GetRemoteURLs(context.Background(), locRepo, git.DefaultRemoteName)
	if err != nil {
		return "", err
	}
	return urls[0], nil
}

type CloneOption struct {
	Alias *reporef.RepoRef
	// UNDONE: support isBare
	// UNDONE: support *git.CloneOptions
}

func (l *Controller) Clone(
	ctx context.Context,
	ref reporef.RepoRef,
	token string,
	opt *CloneOption,
) (Repo, error) {
	var auth transport.AuthMethod
	if token != "" {
		auth = &http.BasicAuth{
			Username: ref.Owner(),
			Password: token,
		}
	}

	p := NewRepo(l.root, ref)
	path := p.FullFilePath()
	url := ref.URL()
	if opt != nil && opt.Alias != nil {
		aliasRef, err := reporef.NewRepoRef(p.ref.Host(), opt.Alias.Owner(), opt.Alias.Name())
		if err != nil {
			return Repo{}, err
		}
		alias := NewRepo(l.root, aliasRef)
		path = alias.FullFilePath()
		p = alias
	}

	log.FromContext(ctx).
		WithFields(log.Fields{"path": p, "url": url}).
		Debug("clone a repository")
	if _, err := git.PlainCloneContext(ctx, path, false, &git.CloneOptions{
		URL:  url,
		Auth: auth,
	}); err != nil {
		return Repo{}, err
	}

	return p, nil
}

type WalkFunc func(Repo) error

type WalkOption struct {
	Query string
}

var mu sync.Mutex

func (l *Controller) Walk(
	ctx context.Context,
	opt *WalkOption,
	walkFn WalkFunc,
) error {
	if _, err := os.Lstat(l.root); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return walker.WalkWithContext(ctx, l.root, func(pathname string, _ os.FileInfo) (retErr error) {
		rel, _ := filepath.Rel(l.root, pathname)
		parts := strings.Split(rel, string(filepath.Separator))

		if len(parts) < 3 {
			return nil
		}
		defer func() {
			if retErr == nil {
				retErr = filepath.SkipDir
			}
		}()

		if _, err := git.PlainOpen(pathname); err != nil {
			log.FromContext(ctx).
				WithFields(log.Fields{"error": err, "rel": rel}).
				Debug("skip a dir that is not a git directory")
			return nil
		}

		// NOTE: Case of len(parts) > 3 never happens because it returns filepath.SkipDir
		ref, err := reporef.NewRepoRef(parts[0], parts[1], parts[2])
		if err != nil {
			log.FromContext(ctx).
				WithFields(log.Fields{"error": err, "rel": rel}).
				Debug("skip invalid entity")
			return nil
		}
		p := NewRepo(l.root, ref)
		if opt != nil && !strings.Contains(p.RelPath(), opt.Query) {
			return nil
		}
		mu.Lock()
		defer mu.Unlock()
		return walkFn(p)
	})
}

type ListOption struct {
	Query string
}

func (l *Controller) List(ctx context.Context, opt *ListOption) ([]Repo, error) {
	var list []Repo
	var woption *WalkOption
	if opt != nil {
		woption = &WalkOption{Query: opt.Query}
	}
	if err := l.Walk(ctx, woption, func(p Repo) error {
		list = append(list, p)
		return nil
	}); err != nil {
		return nil, err
	}
	return list, nil
}

type DeleteOption struct{}

func (l *Controller) Delete(ctx context.Context, ref reporef.RepoRef, opt *DeleteOption) error {
	return DeleteRepository(ctx, NewRepo(l.root, ref), opt)
}

func DeleteRepository(_ context.Context, repo Repo, _ *DeleteOption) error {
	if err := repo.CheckEntity(); err != nil {
		return err
	}
	return os.RemoveAll(repo.FullFilePath())
}
