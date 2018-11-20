package gogh

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var validName = regexp.MustCompile(`^([a-zA-Z0-9](?:(?:-[a-zA-Z0-9]+)*[a-zA-Z0-9])?)/(^[\w-]+)$`)

// var capital = regexp.MustCompile(`[A-Z]`) // UNDONE: warn if name contains capital cases

// RepoName is the name for the repository like <user>/<name>
type RepoName struct {
	user string
	name string
	text string
}

// Shared notices file shared flag like group, all, everybody, 0766.
type Shared string

var validShared = map[string]struct{}{
	"false":     {},
	"true":      {},
	"umask":     {},
	"group":     {},
	"all":       {},
	"world":     {},
	"everybody": {},
}

// Set text as Shared
func (s *Shared) Set(text string) error {
	if _, ok := validShared[text]; ok {
		*s = Shared(text)
		return nil
	}
	if _, err := strconv.ParseInt(text, 8, 8); err == nil {
		*s = Shared(text)
		return nil
	}
	return fmt.Errorf(`invalid shared value %q; shared can be specified with "false", "true", "umask", "group", "all", "world", "everybody" or "0xxx" (octed value)`, text)
}

func (s Shared) String() string {
	return string(s)
}

// Set text as RepoName
func (n *RepoName) Set(text string) error {
	matches := validName.FindStringSubmatch(text)
	if matches == nil {
		return fmt.Errorf("invalid repository name %q; repository should be specified as <user>/<name>", text)
	}
	n.user = matches[1]
	n.name = matches[2]
	n.text = text
	return nil
}

func (n RepoName) String() string {
	return n.text
}

// User is the part of the user in repo name like <user>/<name>.
func (n RepoName) User() string {
	return n.user
}

// Name is the part of the name in repo name like <user>/<name>.
func (n RepoName) Name() string {
	return n.name
}

// RepoSpec specifies a repository in the GitHub
type RepoSpec struct {
	ref    string
	refURL url.URL
}

// NewSpec parses ref string as a spacifier for a repository in the GitHub
func NewSpec(ref string) (*RepoSpec, error) {
	spec := new(RepoSpec)
	if err := spec.Set(ref); err != nil {
		return nil, err
	}
	return spec, nil
}

// Set text as RepoSpec
func (s *RepoSpec) Set(ref string) error {
	if !hasSchemePattern.MatchString(ref) && scpLikeURLPattern.MatchString(ref) {
		matched := scpLikeURLPattern.FindStringSubmatch(ref)
		user := matched[1]
		host := matched[2]
		path := matched[3]

		ref = fmt.Sprintf("ssh://%s%s/%s", user, host, path)
	}

	refURL, err := url.Parse(ref)
	if err != nil {
		return err
	}
	s.refURL = *refURL
	return nil
}

// URL will get a URL for a repository
func (s *RepoSpec) URL(ctx Context, ssh bool) *url.URL {
	u := s.refURL // copy
	httpsURL := &u
	if !httpsURL.IsAbs() {
		if !strings.Contains(httpsURL.Path, "/") {
			user := ctx.UserName()
			httpsURL.Path = user + "/" + httpsURL.Path
		}
		httpsURL.Scheme = "https"
		httpsURL.Host = "github.com"
		if httpsURL.Path[0] != '/' {
			httpsURL.Path = "/" + httpsURL.Path
		}
	}
	if ssh {
		sshURL, _ := url.Parse(fmt.Sprintf("ssh://git@%s%s", httpsURL.Host, httpsURL.Path))
		return sshURL
	}
	return httpsURL
}

func (s RepoSpec) String() string {
	return s.ref
}

// Convert SCP-like URL to SSH URL(e.g. [user@]host.xz:path/to/repo.git/)
// ref. http://git-scm.com/docs/git-fetch#_git_urls
// (golang hasn't supported Perl-like negative look-behind match)
var hasSchemePattern = regexp.MustCompile("^[^:]+://")
var scpLikeURLPattern = regexp.MustCompile("^([^@]+@)?([^:]+):/?(.+)$")

// Specs is array of RepoSpec
type Specs []RepoSpec

// Set will add a text to Specs as a RepoSpec
func (s *Specs) Set(value string) error {
	spec := new(RepoSpec)
	if err := spec.Set(value); err != nil {
		return err
	}
	*s = append(*s, *spec)
	return nil
}

// String : Stringに変換する
func (s Specs) String() string {
	if len(s) == 0 {
		return ""
	}
	strs := make([]string, 0, len(s))
	for _, spec := range s {
		strs = append(strs, spec.String())
	}
	return strings.Join(strs, ",")
}

// IsCumulative : 複数指定可能
func (s Specs) IsCumulative() bool { return true }
