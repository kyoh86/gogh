package gogh

import (
	"net/url"
	"path"
	"regexp"
)

// RepoRef describes which repository is in a root.
type RepoRef struct {
	host  string
	owner string
	name  string
}

func (r RepoRef) Host() string  { return r.host }
func (r RepoRef) Owner() string { return r.owner }
func (r RepoRef) Name() string  { return r.name }

func (r RepoRef) RelLevels() []string {
	return []string{r.host, r.owner, r.name}
}

func (r RepoRef) URL() string {
	return "https://" + path.Join(r.RelLevels()...)
}

func (r RepoRef) String() string {
	return path.Join(r.RelLevels()...)
}

var (
	ErrEmptyHost  = ErrInvalidHost("empty host")
	ErrEmptyOwner = ErrInvalidOwner("empty owner")
	ErrEmptyName  = ErrInvalidName("empty name")
)

type ErrInvalidHost string

func (e ErrInvalidHost) Error() string {
	return string(e)
}

func ValidateHost(h string) error {
	if h == "" {
		return ErrEmptyHost
	}

	u, err := url.ParseRequestURI("https://" + h)
	if err != nil {
		return ErrInvalidHost("invalid host: " + h)
	}
	if u.Host != h {
		return ErrInvalidHost("invalid host: " + h)
	}
	return nil
}

type ErrInvalidName string

func (e ErrInvalidName) Error() string {
	return string(e)
}

var invalidNameRegexp = regexp.MustCompile(`[^\w\-\.]`)

func ValidateName(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	if name == "." {
		return ErrInvalidName("'.' is reserved name")
	}
	if name == ".." {
		return ErrInvalidName("'..' is reserved name")
	}
	if invalidNameRegexp.MatchString(name) {
		return ErrInvalidName("invalid name: " + name)
	}
	return nil
}

type ErrInvalidOwner string

func (e ErrInvalidOwner) Error() string {
	return string(e)
}

var validOwnerRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`)

func ValidateOwner(owner string) error {
	if owner == "" {
		return ErrEmptyOwner
	}
	if !validOwnerRegexp.MatchString(owner) {
		return ErrInvalidOwner("invalid owner: " + owner)
	}
	return nil
}

func NewRepoRef(host, owner, name string) (RepoRef, error) {
	if err := ValidateName(name); err != nil {
		return RepoRef{}, err
	}
	if err := ValidateOwner(owner); err != nil {
		return RepoRef{}, err
	}
	if err := ValidateHost(host); err != nil {
		return RepoRef{}, err
	}
	return RepoRef{
		host:  host,
		owner: owner,
		name:  name,
	}, nil
}
