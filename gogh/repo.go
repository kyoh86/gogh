package gogh

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"
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

// ParseProject parses a repo-name for a repository in the GitHub
func ParseProject(p *Project) (*Repo, error) {
	repo := new(Repo)
	if err := repo.Set(p.Subpaths()[1]); err != nil {
		return nil, err
	}
	return repo, nil
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
func CheckRepoHost(ctx GitHubContext, repo *Repo) error {
	return SupportedHost(ctx, repo.Host(ctx))
}

// SupportedHost checks that a host is supported
func SupportedHost(ctx GitHubContext, host string) error {
	if host == ctx.GitHubHost() {
		return nil
	}
	return fmt.Errorf("not supported host: %q", host)
}

// Convert SCP-like URL to SSH URL(e.g. [user@]host.xz:path/to/repo.git/)
// ref. http://git-scm.com/docs/git-fetch#_git_urls
// (golang hasn't supported Perl-like negative look-behind match)
var hasSchemePattern = regexp.MustCompile("^[^:]+://")
var scpLikeURLPattern = regexp.MustCompile("^([^@]+@)?([^:]+):/?(.+)$")

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

	var path string
	if url.IsAbs() {
		r.scheme = url.Scheme
		r.host = url.Host
		r.user = url.User
		path = strings.Trim(url.Path, "/")
	} else {
		r.scheme = "https"
		r.host = "" // use default value
		r.user = nil
		path = url.Path
	}
	r.forceQuery = url.ForceQuery
	r.rawQuery = url.RawQuery
	r.fragment = url.Fragment

	pp := strings.Split(path, "/")
	switch len(pp) {
	case 1:
		r.owner = "" // To use context.UserName() instead.
		r.name = strings.TrimSuffix(pp[0], ".git")
		if err := ValidateName(r.name); err != nil {
			return err
		}
	case 2:
		r.owner = pp[0]
		if err := ValidateOwner(r.owner); err != nil {
			return err
		}
		r.name = strings.TrimSuffix(pp[1], ".git")
		if err := ValidateName(r.name); err != nil {
			return err
		}
	default:
		return errors.New("repository name has too many slashes")
	}
	r.raw = raw
	return nil
}

// Scheme returns scheme of the repository
func (r *Repo) Scheme(_ GitHubContext) string {
	return r.scheme
}

// Host returns host of the repository
func (r *Repo) Host(ctx GitHubContext) string {
	if r.host == "" {
		return ctx.GitHubHost()
	}
	return r.host
}

// ExplicitOwner returns a user name of an explicit owner of the repository
func (r *Repo) ExplicitOwner(_ GitHubContext) string {
	return r.owner
}

// Owner returns a user name of an owner of the repository
func (r *Repo) Owner(ctx GitHubContext) string {
	if r.owner == "" {
		return ctx.GitHubUser()
	}
	return r.owner
}

// Name returns a name of the repository
func (r *Repo) Name(_ GitHubContext) string {
	return r.name
}

// FullName returns a repository identifier that is formed like {Owner/Name}
func (r *Repo) FullName(ctx GitHubContext) string {
	return path.Join(r.Owner(ctx), r.Name(ctx))
}

// URL will get a URL for a repository
func (r *Repo) URL(ctx GitHubContext, ssh bool) *url.URL {
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

// Check if a GitHub repo is public (we can access the repo without token or auth)
func (r *Repo) IsPublic(ctx GitHubContext) (bool, error) {
	url := r.URL(ctx, false)
	res, err := http.Head(url.String())
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("invalid status code: %d", res.StatusCode)
	}
}

// Match with project.
func (r *Repo) Match(p *Project) bool {
	if r.host != "" && r.host != p.PathParts[0] {
		return false
	}
	if r.owner != "" && r.owner != p.PathParts[1] {
		return false
	}
	return r.name == p.PathParts[2]
}

// RelPath get relative path from root directory
func (r *Repo) RelPath(ctx GitHubContext) string {
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
