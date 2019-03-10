package gogh

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Repo specifies a repository in the GitHub
type Repo struct {
	raw string

	scheme string
	host   string        // host or host:port
	user   *url.Userinfo // username and password information
	owner  string
	name   string

	forceQuery bool   // append a query ('?') even if RawQuery is empty
	rawQuery   string // encoded query values, without '?'
	fragment   string // fragment for references, without '#'
}

// ParseRepo parses a repo-name for a repository in the GitHub
func ParseRepo(rawRepo string) (*Repo, error) {
	repo := new(Repo)
	if err := repo.Set(rawRepo); err != nil {
		return nil, err
	}
	return repo, nil
}

// CheckRepoHost that repo is in supported host
func CheckRepoHost(ctx Context, repo *Repo) error {
	host := repo.Host(ctx)
	if host == DefaultHost {
		return nil
	}
	for _, h := range ctx.GHEHosts() {
		if h == host {
			return nil
		}
	}
	return fmt.Errorf("not supported host: %q", host)
}

// Convert SCP-like URL to SSH URL(e.g. [user@]host.xz:path/to/repo.git/)
// ref. http://git-scm.com/docs/git-fetch#_git_urls
// (golang hasn't supported Perl-like negative look-behind match)
var hasSchemePattern = regexp.MustCompile("^[^:]+://")
var scpLikeURLPattern = regexp.MustCompile("^([^@]+@)?([^:]+):/?(.+)$")

//TODO: var validName = regexp.MustCompile(`^([a-zA-Z0-9](?:(?:-[a-zA-Z0-9]+)*[a-zA-Z0-9]+)?)/([\w-]+)$`)
//TODO: var capital = regexp.MustCompile(`[A-Z]`) // UNDONE: warn if name contains capital cases

// Set text as Repo
func (r *Repo) Set(rawRepo string) error {
	raw := rawRepo
	if !hasSchemePattern.MatchString(rawRepo) && scpLikeURLPattern.MatchString(rawRepo) {
		matched := scpLikeURLPattern.FindStringSubmatch(rawRepo)
		user := matched[1]
		host := matched[2]
		path := matched[3]

		rawRepo = fmt.Sprintf("ssh://%s%s/%s", user, host, path)
	}

	url, err := url.Parse(rawRepo)
	if err != nil {
		return err
	}

	if url.IsAbs() {
		r.scheme = url.Scheme
		r.host = url.Host
		r.user = url.User
	} else {
		r.scheme = "https"
		r.host = DefaultHost
		r.user = nil
	}
	r.forceQuery = url.ForceQuery
	r.rawQuery = url.RawQuery
	r.fragment = url.Fragment

	pp := strings.Split(strings.Trim(url.Path, "/"), "/")
	switch len(pp) {
	case 0:
		return errors.New("repository remote has no local name")
	case 1:
		r.owner = "" // To use context.UserName() instead.
		r.name = strings.TrimSuffix(pp[0], ".git")
	case 2:
		r.owner = pp[0]
		r.name = strings.TrimSuffix(pp[1], ".git")
	default:
		return errors.New("repository remote has too many slashes")
	}
	r.raw = raw
	return nil
}

// DefaultHost is the default host of the GitHub
const DefaultHost = "github.com"

// Scheme returns scheme of the repository
func (r *Repo) Scheme(_ Context) string {
	return r.scheme
}

// Host returns host of the repository
func (r *Repo) Host(_ Context) string {
	return r.host
}

// Owner returns a user name of an owner of the repository
func (r *Repo) Owner(ctx Context) string {
	if r.owner == "" {
		return ctx.UserName()
	}
	return r.owner
}

// Name returns a name of the repository
func (r *Repo) Name(_ Context) string {
	return r.name
}

// URL will get a URL for a repository
func (r *Repo) URL(ctx Context, ssh bool) *url.URL {
	if ssh {
		return &url.URL{
			Scheme: "ssh",
			User:   url.User("git"),
			Host:   r.Host(ctx),
			Path:   path.Join("/", r.Owner(ctx), r.Name(ctx)),
		}
	}
	return &url.URL{
		Scheme: r.Scheme(ctx),
		User:   r.user,
		Host:   r.Host(ctx),
		Path:   path.Join("/", r.Owner(ctx), r.Name(ctx)),
	}
}

// RelPath get relative path from root directory
func (r *Repo) RelPath(ctx Context) string {
	return filepath.Join(r.Host(ctx), r.Owner(ctx), r.Name(ctx))
}

func (r Repo) String() string {
	return r.raw
}

// Repos is array of Repo
type Repos []Repo

// Set will add a text to Repos as a Repo
func (r *Repos) Set(value string) error {
	repo := new(Repo)
	if err := repo.Set(value); err != nil {
		return err
	}
	*r = append(*r, *repo)
	return nil
}

// String : Stringに変換する
func (r Repos) String() string {
	if len(r) == 0 {
		return ""
	}
	strs := make([]string, 0, len(r))
	for _, repo := range r {
		strs = append(strs, repo.String())
	}
	return strings.Join(strs, ",")
}

// IsCumulative : 複数指定可能
func (r Repos) IsCumulative() bool { return true }
