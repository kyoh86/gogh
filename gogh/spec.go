package gogh

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var validName = regexp.MustCompile(`^([a-zA-Z0-9](?:(?:-[a-zA-Z0-9]+)*[a-zA-Z0-9]+)?)/([\w-]+)$`)

// var capital = regexp.MustCompile(`[A-Z]`) // UNDONE: warn if name contains capital cases

// LocalName is the name for the repository like <user>/<name>
type LocalName struct {
	user string
	name string
	text string
}

// RepoShared notices file shared flag like group, all, everybody, 0766.
type RepoShared string

var validShared = map[string]struct{}{
	"false":     {},
	"true":      {},
	"umask":     {},
	"group":     {},
	"all":       {},
	"world":     {},
	"everybody": {},
}

// Set text as RepoShared
func (s *RepoShared) Set(text string) error {
	if _, ok := validShared[text]; ok {
		*s = RepoShared(text)
		return nil
	}
	if _, err := strconv.ParseInt(text, 8, 16); err == nil {
		*s = RepoShared(text)
		return nil
	}
	return fmt.Errorf(`invalid shared value %q; shared can be specified with "false", "true", "umask", "group", "all", "world", "everybody" or "0xxx" (octed value)`, text)
}

func (s RepoShared) String() string {
	return string(s)
}

// Set text as LocalName
func (n *LocalName) Set(text string) error {
	matches := validName.FindStringSubmatch(text)
	if matches == nil {
		return fmt.Errorf("invalid repository name %q; repository should be specified as <user>/<name>", text)
	}
	n.user = matches[1]
	n.name = matches[2]
	n.text = text
	return nil
}

func (n LocalName) String() string {
	return n.text
}

// User is the part of the user in repo name like <user>/<name>.
func (n LocalName) User() string {
	return n.user
}

// Name is the part of the name in repo name like <user>/<name>.
func (n LocalName) Name() string {
	return n.name
}

// RemoteName specifies a repository in the GitHub
type RemoteName struct {
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

// ParseRemoteName parses a remote-name for a repository in the GitHub
func ParseRemoteName(rawName string) (*RemoteName, error) {
	name := new(RemoteName)
	if err := name.Set(rawName); err != nil {
		return nil, err
	}
	return name, nil
}

// Convert SCP-like URL to SSH URL(e.g. [user@]host.xz:path/to/repo.git/)
// ref. http://git-scm.com/docs/git-fetch#_git_urls
// (golang hasn't supported Perl-like negative look-behind match)
var hasSchemePattern = regexp.MustCompile("^[^:]+://")
var scpLikeURLPattern = regexp.MustCompile("^([^@]+@)?([^:]+):/?(.+)$")

// Set text as RemoteName
func (n *RemoteName) Set(rawName string) error {
	raw := rawName
	if !hasSchemePattern.MatchString(rawName) && scpLikeURLPattern.MatchString(rawName) {
		matched := scpLikeURLPattern.FindStringSubmatch(rawName)
		user := matched[1]
		host := matched[2]
		path := matched[3]

		rawName = fmt.Sprintf("ssh://%s%s/%s", user, host, path)
	}

	url, err := url.Parse(rawName)
	if err != nil {
		return err
	}

	if url.IsAbs() {
		n.scheme = url.Scheme
		n.host = url.Host
		n.user = url.User
	} else {
		n.scheme = "https"
		n.host = DefaultHost
		n.user = nil
	}
	n.forceQuery = url.ForceQuery
	n.rawQuery = url.RawQuery
	n.fragment = url.Fragment

	pp := strings.Split(strings.TrimPrefix(url.Path, "/"), "/")
	switch len(pp) {
	case 0:
		return errors.New("repository name has no local name")
	case 1:
		n.owner = "" // Use context.UserName() instead.
		n.name = pp[0]
	case 2:
		n.owner = pp[0]
		n.name = pp[1]
	default:
		return errors.New("repository name has too many slashes")
	}
	n.raw = raw
	return nil
}

// DefaultHost is the default host name of the GitHub
const DefaultHost = "github.com"

// Scheme returns scheme of the repository
func (n *RemoteName) Scheme(_ Context) string {
	return n.scheme
}

// Host returns host name of the repository
func (n *RemoteName) Host(_ Context) string {
	return n.host
}

// Owner returns a user name of an owner of the repository
func (n *RemoteName) Owner(ctx Context) string {
	if n.owner == "" {
		return ctx.UserName()
	}
	return n.owner
}

// Name returns a name of the repository
func (n *RemoteName) Name(_ Context) string {
	return n.name
}

// URL will get a URL for a repository
func (n *RemoteName) URL(ctx Context, ssh bool) *url.URL {
	if ssh {
		return &url.URL{
			Scheme: "ssh",
			User:   url.User("git"),
			Host:   n.Host(ctx),
			Path:   path.Join("/", n.Owner(ctx), n.Name(ctx)),
		}
	}
	return &url.URL{
		Scheme: n.Scheme(ctx),
		User:   n.user,
		Host:   n.Host(ctx),
		Path:   path.Join("/", n.Owner(ctx), n.Name(ctx)),
	}
}

func (n RemoteName) String() string {
	return n.raw
}

// RemoteNames is array of RemoteName
type RemoteNames []RemoteName

// Set will add a text to RemoteNames as a RemoteName
func (n *RemoteNames) Set(value string) error {
	name := new(RemoteName)
	if err := name.Set(value); err != nil {
		return err
	}
	*n = append(*n, *name)
	return nil
}

// String : Stringに変換する
func (n RemoteNames) String() string {
	if len(n) == 0 {
		return ""
	}
	strs := make([]string, 0, len(n))
	for _, name := range n {
		strs = append(strs, name.String())
	}
	return strings.Join(strs, ",")
}

// IsCumulative : 複数指定可能
func (s RemoteNames) IsCumulative() bool { return true }
