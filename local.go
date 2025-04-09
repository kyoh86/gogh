package gogh

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
	"github.com/saracen/walker"
)

const DefaultRootDirName = "Projects"

func NewLocalController(root string) *LocalController {
	return &LocalController{root: root}
}

type LocalController struct {
	// UNDONE: support fs.FS
	// UNDONE: support fs.FS

	// NOTE: v1 -> v2 diferrence
	// if we wanna manage multiple root, create multiple controller instances.
	root string
}

type LocalExistOption struct {
}

func (l *LocalController) Exist(
	ctx context.Context,
	spec Spec,
	opt *LocalExistOption,
) (bool, error) {
	project := NewProject(l.root, spec)
	_, err := git.PlainOpen(project.FullFilePath())
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, git.ErrRepositoryNotExists):
		return false, nil
	default:
		return false, err
	}
}

type LocalCreateOption struct { // UNDONE: support isBare
}

func (l *LocalController) Create(
	ctx context.Context,
	spec Spec,
	opt *LocalCreateOption,
) (Project, error) {
	p := NewProject(l.root, spec)
	if err := CreateLocalProject(ctx, p, spec.URL(), opt); err != nil {
		return Project{}, err
	}
	return p, nil
}

func CreateLocalProject(
	_ context.Context,
	project Project,
	remoteURL string,
	_ *LocalCreateOption,
) error {
	repo, err := git.PlainInit(project.FullFilePath(), false)
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

func (l *LocalController) SetRemoteSpecs(
	ctx context.Context,
	spec Spec,
	remotes map[string][]Spec,
) error {
	urls := map[string][]string{}
	for name, specs := range remotes {
		for _, spec := range specs {
			urls[name] = append(urls[name], spec.URL())
		}
	}
	return l.SetRemoteURLs(ctx, spec, urls)
}

func (l *LocalController) SetRemoteURLs(
	ctx context.Context,
	spec Spec,
	remotes map[string][]string,
) error {
	return SetRemoteURLsOnLocalProject(ctx, NewProject(l.root, spec), remotes)
}

func SetRemoteURLsOnLocalProject(
	_ context.Context,
	project Project,
	remotes map[string][]string,
) error {
	repo, err := git.PlainOpen(project.FullFilePath())
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

func (l *LocalController) GetRemoteURLs(
	ctx context.Context,
	spec Spec,
	name string,
) ([]string, error) {
	return GetRemoteURLsFromLocalProject(ctx, NewProject(l.root, spec), name)
}

func GetRemoteURLsFromLocalProject(
	_ context.Context,
	project Project,
	name string,
) ([]string, error) {
	repo, err := git.PlainOpen(project.FullFilePath())
	if err != nil {
		return nil, fmt.Errorf("open local repository: %w", err)
	}
	remote, err := repo.Remote(name)
	if err != nil {
		return nil, fmt.Errorf("get remote %s: %w", name, err)
	}
	return remote.Config().URLs, nil
}

func GetDefaultRemoteURLFromLocalProject(_ context.Context, project Project) (string, error) {
	urls, err := GetRemoteURLsFromLocalProject(context.Background(), project, git.DefaultRemoteName)
	if err != nil {
		return "", err
	}
	return urls[0], nil
}

type LocalCloneOption struct {
	Alias *Spec
	// UNDONE: support isBare
	// UNDONE: support *git.CloneOptions
}

func (l *LocalController) Clone(
	ctx context.Context,
	spec Spec,
	token string,
	opt *LocalCloneOption,
) (Project, error) {
	var auth transport.AuthMethod
	if token != "" {
		auth = &http.BasicAuth{
			Username: spec.Owner(),
			Password: token,
		}
	}

	p := NewProject(l.root, spec)
	path := p.FullFilePath()
	url := spec.URL()
	if opt != nil && opt.Alias != nil {
		alias := NewProject(l.root, *opt.Alias)
		alias.spec.host = p.spec.host
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
		return Project{}, err
	}

	return p, nil
}

type LocalWalkFunc func(Project) error

type LocalWalkOption struct {
	Query string
}

var mu sync.Mutex

func (l *LocalController) Walk(
	ctx context.Context,
	opt *LocalWalkOption,
	walkFn LocalWalkFunc,
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
		spec, err := NewSpec(parts[0], parts[1], parts[2])
		if err != nil {
			log.FromContext(ctx).
				WithFields(log.Fields{"error": err, "rel": rel}).
				Debug("skip invalid entity")
			return nil
		}
		p := NewProject(l.root, spec)
		if opt != nil && !strings.Contains(p.RelPath(), opt.Query) {
			return nil
		}
		mu.Lock()
		defer mu.Unlock()
		return walkFn(p)
	})
}

type LocalListOption struct {
	Query string
}

func (l *LocalController) List(ctx context.Context, opt *LocalListOption) ([]Project, error) {
	var list []Project
	var woption *LocalWalkOption
	if opt != nil {
		woption = &LocalWalkOption{Query: opt.Query}
	}
	if err := l.Walk(ctx, woption, func(p Project) error {
		list = append(list, p)
		return nil
	}); err != nil {
		return nil, err
	}
	return list, nil
}

type LocalDeleteOption struct{}

func (l *LocalController) Delete(ctx context.Context, spec Spec, opt *LocalDeleteOption) error {
	return DeleteLocalProject(ctx, NewProject(l.root, spec), opt)
}

func DeleteLocalProject(_ context.Context, project Project, _ *LocalDeleteOption) error {
	if err := project.CheckEntity(); err != nil {
		return err
	}
	return os.RemoveAll(project.FullFilePath())
}
