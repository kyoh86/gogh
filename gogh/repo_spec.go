package gogh

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// RepoSpec specifies a repository in the GitHub
type RepoSpec struct {
	raw  string
	repo Repo
}

// Convert SCP-like URL to SSH URL(e.g. [user@]host.xz:path/to/repo.git/)
// ref. http://git-scm.com/docs/git-fetch#_git_urls
// (golang hasn't supported Perl-like negative look-behind match)
var hasSchemePattern = regexp.MustCompile("^[^:]+://")
var scpLikeURLPattern = regexp.MustCompile("^([^@]+@)?([^:]+):/?(.+)$")

// Set text as RepoSpec
func (r *RepoSpec) Set(rawRepo string) error {
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
		r.repo.scheme = url.Scheme
		r.repo.host = url.Host
		r.repo.user = url.User
		path = strings.Trim(url.Path, "/")
	} else {
		r.repo.scheme = "https"
		r.repo.host = "" // use default value
		r.repo.user = nil
		path = url.Path
	}
	r.repo.forceQuery = url.ForceQuery
	r.repo.rawQuery = url.RawQuery
	r.repo.fragment = url.Fragment

	pp := strings.Split(path, "/")
	switch len(pp) {
	case 1:
		r.repo.owner = "" // To use context.UserName() instead.
		r.repo.name = strings.TrimSuffix(pp[0], ".git")
		if err := ValidateName(r.repo.name); err != nil {
			return err
		}
	case 2:
		r.repo.owner = pp[0]
		if err := ValidateOwner(r.repo.owner); err != nil {
			return err
		}
		r.repo.name = strings.TrimSuffix(pp[1], ".git")
		if err := ValidateName(r.repo.name); err != nil {
			return err
		}
	default:
		return errors.New("repository name has too many slashes")
	}
	r.raw = raw
	return nil
}

func (r RepoSpec) String() string {
	return r.raw
}

func (r *RepoSpec) Validate(env Env) (*Repo, error) {
	repo := r.repo // copy object
	if repo.host == "" {
		repo.host = env.GithubHost()
	} else if repo.host != env.GithubHost() {
		return nil, fmt.Errorf("not supported host: %q", repo.host)
	}
	if repo.owner == "" {
		repo.owner = "kyoh86" // TODO: cache.GithubUser() or the get 'me' with Github token
	}
	return &repo, nil
}

// RepoSpecs is array of RepoSpec
type RepoSpecs []RepoSpec

// Set will add a text to RepoSpecs as a RepoSpec
func (specs *RepoSpecs) Set(value string) error {
	repo := new(RepoSpec)
	if err := repo.Set(value); err != nil {
		return err
	}
	*specs = append(*specs, *repo)
	return nil
}

func (specs RepoSpecs) String() string {
	if len(specs) == 0 {
		return ""
	}
	strs := make([]string, 0, len(specs))
	for _, repo := range specs {
		strs = append(strs, repo.String())
	}
	return strings.Join(strs, ",")
}

func (specs RepoSpecs) IsCumulative() bool { return true }

// Repos will get repositories with GitHub host and user
func (specs RepoSpecs) Validate(env Env) ([]Repo, error) {
	repos := make([]Repo, 0, len(specs))
	for _, spec := range specs {
		repo, err := spec.Validate(env)
		if err != nil {
			return nil, err
		}
		repos = append(repos, *repo)
	}
	return repos, nil
}
