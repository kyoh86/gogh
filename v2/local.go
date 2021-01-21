package gogh

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

type Project struct {
	Root string
	Description
}

func (p Project) FullLevels() []string {
	return []string{p.Root, p.Host, p.User, p.Name}
}

func (p Project) RelLevels() []string {
	return []string{p.Host, p.User, p.Name}
}

func (p Project) FullPath() string {
	return filepath.Join(p.FullLevels()...)
}

func (p Project) RelPath() string {
	return filepath.Join(p.RelLevels()...)
}

func (p Project) URL() string {
	return "https://" + path.Join(p.RelLevels()...)
}

type Description struct {
	Host string
	User string
	Name string
}

const DefaultRootDirName = "Projects"

func NewLocal(ctx context.Context) (*Local, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get user home dir: %w", err)
	}
	return &Local{roots: []string{filepath.Join(home, DefaultRootDirName)}}, nil
}

type Local struct {
	roots []string
}

func (l *Local) SetRoot(roots ...string) {
	l.roots = roots
}

func (l *Local) Roots() []string {
	return l.roots
}

func (l *Local) Create(ctx context.Context, d Description) (*Project, error) {
	p := &Project{
		Root:        l.roots[0],
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
