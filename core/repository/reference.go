package repository

import (
	"path"
)

// Reference is a interface for a repository Reference.
type Reference struct {
	host  string
	owner string
	name  string
}

// Host is a hostname (i.g.: "github.com")
func (r *Reference) Host() string { return r.host }

// Owner is a owner name (i.g.: "kyoh86")
func (r *Reference) Owner() string { return r.owner }

// Name of the repository (i.g.: "gogh")
func (r *Reference) Name() string { return r.name }

func (r *Reference) String() string {
	return path.Join(r.host, r.owner, r.name)
}

// NewReference will build a Reference with a host, owner and name.
func NewReference(host, owner, name string) Reference {
	return Reference{
		host:  host,
		owner: owner,
		name:  name,
	}
}

// ReferenceWithAlias is a struct that contains a Reference and an optional alias.
type ReferenceWithAlias struct {
	Reference Reference
	Alias     *Reference
}
