package gogh

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/saracen/walker"
	"github.com/wacul/ulog"
)

const DefaultRootDirName = "Projects"

func NewLocalController(root string) *LocalController {
	return &LocalController{root: root}
}

//UNDONE: provide DefaultLocalRoot for app? mainutil?
// func DefaultLocalRoot() (string, error) {
// 	home, err := os.UserHomeDir()
// 	if err != nil {
// 		return "", fmt.Errorf("get user home dir: %w", err)
// 	}
// 	return filepath.Join(home, DefaultRootDirName), nil
// }

type LocalController struct {
	// UNDONE: support fs.FS
	// UNDONE: support fs.FS

	// NOTE: v1 -> v2 diferrence
	// if we wanna manage mulstiple root, create multiple controller instances.
	root string
}

type CreateOption struct {
	//UNDONE: support isBare
}

func (l *LocalController) Create(ctx context.Context, d Description, _ *CreateOption) (Project, error) {
	p := NewProject(l.root, d)

	repo, err := git.PlainInit(p.FullPath(), false)
	if err != nil {
		return Project{}, err
	}

	if _, err := repo.CreateRemote(&config.RemoteConfig{
		Name: git.DefaultRemoteName,
		URLs: []string{p.URL()},
	}); err != nil {
		return Project{}, err
	}
	return p, nil
}

type CloneOption struct {
	//UNDONE: support authentication
	//UNDONE: support isBare
	//UNDONE: support *git.CloneOptions
}

func (l *LocalController) Clone(ctx context.Context, d Description, _ *CloneOption) (Project, error) {
	p := NewProject(l.root, d)

	if _, err := git.PlainCloneContext(ctx, p.FullPath(), false, &git.CloneOptions{
		URL: p.URL(),
	}); err != nil {
		return Project{}, err
	}

	return p, nil
}

type LocalListParam struct {
	Query string
}

type LocalWalkFunc func(Project) error

func (l *LocalController) Walk(ctx context.Context, query string, walkFn LocalWalkFunc) error {
	return walker.WalkWithContext(ctx, l.root, func(pathname string, info os.FileInfo) (retErr error) {
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

		p, err := l.newProjectFromEntity([3]string{parts[0], parts[1], parts[2]}, info)
		if err != nil {
			ulog.Logger(ctx).WithField("error", err).WithField("rel", rel).Debug("skip invalid entity")
			return nil
		}
		return walkFn(p)
	})
}

func (l *LocalController) newProjectFromEntity(parts [3]string, info os.FileInfo) (Project, error) {
	if !info.IsDir() {
		return Project{}, errors.New("not directory")
	}
	// NOTE: Case of len(parts) > 3 never happens because it returns filepath.SkipDir
	description, err := NewDescription(parts[0], parts[1], parts[2])
	if err != nil {
		return Project{}, err
	}
	return NewProject(l.root, description), nil
}

func (l *LocalController) List(ctx context.Context, query string) ([]Project, error) {
	var list []Project
	if err := l.Walk(ctx, query, func(p Project) error {
		list = append(list, p)
		return nil
	}); err != nil {
		return nil, err
	}
	return list, nil
}

func (l *LocalController) Remove(ctx context.Context, description Description) error {
	p := NewProject(l.root, description)
	if err := p.CheckEntity(); err != nil {
		return err
	}
	return os.RemoveAll(p.FullPath())
}
