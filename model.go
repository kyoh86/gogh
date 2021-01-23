package gogh

import (
	"path"
	"path/filepath"
)

type Project struct {
	root string
	Description
}

func (p Project) Root() string {
	return p.root
}

func (p Project) FullLevels() []string {
	return []string{p.root, p.host, p.user, p.name}
}

func (p Project) RelLevels() []string {
	return []string{p.host, p.user, p.name}
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
	host string
	user string
	name string
}

func (d Description) Host() string { return d.host }
func (d Description) User() string { return d.user }
func (d Description) Name() string { return d.name }
