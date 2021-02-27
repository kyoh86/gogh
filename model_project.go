package gogh

import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

// Project is the location of a repository in the local.
// It is a valid location, that never means "exist".
type Project struct {
	root string
	spec Spec
}

func (p Project) Root() string {
	return p.root
}

func (p Project) Host() string  { return p.spec.host }
func (p Project) Owner() string { return p.spec.owner }
func (p Project) Name() string  { return p.spec.name }

func (p Project) FullLevels() []string {
	return []string{p.root, p.spec.host, p.spec.owner, p.spec.name}
}

func (p Project) RelLevels() []string {
	return []string{p.spec.host, p.spec.owner, p.spec.name}
}

func (p Project) FullFilePath() string {
	return filepath.Join(p.FullLevels()...)
}

func (p Project) RelFilePath() string {
	return filepath.Join(p.RelLevels()...)
}

func (p Project) RelPath() string {
	return path.Join(p.RelLevels()...)
}

func (p Project) URL() string {
	return "https://" + path.Join(p.RelLevels()...)
}

// CheckEntity checks the project is exist in the local file-system.
func (p Project) CheckEntity() error {
	path := p.FullFilePath()
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return errors.New("project is not dir")
	}
	return nil
}

// UNDONE: CheckEntityInFileSystem() support fs.FS

func NewProject(root string, spec Spec) Project {
	return Project{root: root, spec: spec}
}
