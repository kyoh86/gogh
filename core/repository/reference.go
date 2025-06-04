package repository

import (
	"path"
)

// Reference is a struct that contains a host, owner and name of a repository.
// It is used to identify a repository in a hosting service or a location in a root.
type Reference struct {
	host  string
	owner string
	name  string
}

// Host is a hostname (e.g.: "github.com")
func (r *Reference) Host() string { return r.host }

// Owner is a owner name (e.g.: "kyoh86")
func (r *Reference) Owner() string { return r.owner }

// Name of the repository (e.g.: "gogh")
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
	// Reference is the main reference.
	Reference Reference
	// Alias is an optional alias for the reference if needed.
	Alias *Reference
}

// Local returns the local reference. If an alias is set, it returns the alias; otherwise, it returns the main reference.
func (r ReferenceWithAlias) Local() Reference {
	if r.Alias != nil {
		return *r.Alias
	}
	return r.Reference
}
