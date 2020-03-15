package gogh

import (
	"net/url"
	"path"
	"path/filepath"
)

// Repo specifies a repository in the GitHub
type Repo struct {
	scheme     string
	host       string        // host or host:port
	user       *url.Userinfo // username and password information
	owner      string
	name       string
	forceQuery bool   // append a query ('?') even if RawQuery is empty
	rawQuery   string // encoded query values, without '?'
	fragment   string // fragment for references, without '#'
}

// Owner returns a user name of an owner of the repository
func (r *Repo) Owner() string {
	return r.owner
}

// Name returns a name of the repository
func (r *Repo) Name() string {
	return r.name
}

// FullName returns a repository identifier that is formed like {Owner/Name}
func (r *Repo) FullName() string {
	return path.Join(r.owner, r.name)
}

// URL will get a URL for a repository
func (r *Repo) URL(ssh bool) *url.URL {
	if ssh {
		return &url.URL{
			Scheme: "ssh",
			User:   url.User("git"),
			Host:   r.host,
			Path:   path.Join("/", r.owner, r.name),
		}
	}
	return &url.URL{
		Scheme: r.scheme,
		User:   r.user,
		Host:   r.host,
		Path:   path.Join("/", r.owner, r.name),
	}
}

func (r Repo) String() string {
	return r.URL(false).String()
}

// Match with project.
func (r *Repo) Match(p *Project) bool {
	if r.host != p.PathParts[0] {
		return false
	}
	if r.owner != p.PathParts[1] {
		return false
	}
	return r.name == p.PathParts[2]
}

// RelPath get relative path from root directory
func (r *Repo) RelPath() string {
	return filepath.Join(r.host, r.owner, r.name)
}
