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
	root        string
	description Description
}

func (p Project) Root() string {
	return p.root
}

func (p Project) Host() string { return p.description.host }
func (p Project) User() string { return p.description.user }
func (p Project) Name() string { return p.description.name }

func (p Project) FullLevels() []string {
	return []string{p.root, p.description.host, p.description.user, p.description.name}
}

func (p Project) RelLevels() []string {
	return []string{p.description.host, p.description.user, p.description.name}
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

// CheckEntity checks the project is exist in the local file-system.
func (p Project) CheckEntity() error {
	path := p.FullPath()
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

func NewProject(root string, description Description) Project {
	return Project{root: root, description: description}
}
