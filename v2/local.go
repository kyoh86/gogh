package gogh

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/saracen/walker"
	"github.com/wacul/ulog"
)

const DefaultRootDirName = "Projects"

func NewLocalController(ctx context.Context, root string) *LocalController {
	return &LocalController{root: root}
}

func DefaultLocalRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get user home dir: %w", err)
	}
	return filepath.Join(home, DefaultRootDirName), nil
}

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

func (l *LocalController) Create(ctx context.Context, d Description, _ *CreateOption) (*Project, error) {
	p := &Project{
		root:        l.root,
		Description: d,
	}

	repo, err := git.PlainInit(p.FullPath(), false)
	if err != nil {
		return nil, err
	}

	if _, err := repo.CreateRemote(&config.RemoteConfig{
		Name: git.DefaultRemoteName,
		URLs: []string{p.URL()},
	}); err != nil {
		return nil, err
	}
	return p, nil
}

type CloneOption struct {
	//UNDONE: support authentication
	//UNDONE: support isBare
	//UNDONE: support *git.CloneOptions
}

func (l *LocalController) Clone(ctx context.Context, d Description, _ *CloneOption) (*Project, error) {
	p := &Project{
		root:        l.root,
		Description: d,
	}

	if _, err := git.PlainCloneContext(ctx, p.FullPath(), false, &git.CloneOptions{
		URL: p.URL(),
	}); err != nil {
		return nil, err
	}

	return p, nil
}

type LocalListParam struct {
	Query string
}

type LocalWalkFunc func(*Project) error

func (l *LocalController) Walk(ctx context.Context, query string, callback LocalWalkFunc) error {
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
		if !info.IsDir() {
			return nil
		}
		desc, err := ValidateDescription(parts[0], parts[1], parts[2])
		if err != nil {
			ulog.Logger(ctx).WithField("error", err).WithField("rel", rel).Debug("invalid path is skipped")
			return nil
		}
		return callback(&Project{root: l.root, Description: *desc})
	})
}

func (l *LocalController) List(ctx context.Context, query string) ([]Project, error) {
	var list []Project
	if err := l.Walk(ctx, query, func(p *Project) error {
		list = append(list, *p)
		return nil
	}); err != nil {
		return nil, err
	}
	return list, nil
}

func (l *LocalController) Remove(ctx context.Context, description Description) error {
	p, err := l.Get(ctx, description)
	if err != nil {
		return err
	}
	return os.RemoveAll(p.FullPath())
}

func (l *LocalController) Get(ctx context.Context, description Description) (*Project, error) {
	p := &Project{root: l.root, Description: description}
	path := p.FullPath()
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, errors.New("project is not dir")
	}
	return p, nil
}
