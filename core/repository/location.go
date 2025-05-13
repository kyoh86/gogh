package repository

import (
	"path"
)

// Location is a struct that contains information about a repository location.
type Location struct {
	fullPath string
	path     string
	host     string
	owner    string
	name     string
}

// FullPath is a full path of the repository (i.g.: "/path/to/root/github.com/kyoh86/gogh")
func (r *Location) FullPath() string { return r.fullPath }

// Host is a hostname (i.g.: "github.com")
func (r *Location) Host() string { return r.host }

// Owner is a owner name (i.g.: "kyoh86")
func (r *Location) Owner() string { return r.owner }

// Name of the repository (i.g.: "gogh")
func (r *Location) Name() string { return r.name }

// Path returns the path from a root of the repository (i.g.: "github.com/kyoh86/gogh")
func (r *Location) Path() string {
	return r.path
}

// NewLocation will build a repository location with a full path, host, owner and name.
func NewLocation(fullPath string, host, owner, name string) Location {
	return Location{
		fullPath: fullPath,
		path:     path.Join(host, owner, name),
		host:     host,
		owner:    owner,
		name:     name,
	}
}
