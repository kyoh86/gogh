package gogh

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
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
	url url.URL
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
	n.url = *url
	n.raw = raw
	return nil
}

// URL will get a URL for a repository
func (n *RemoteName) URL(ctx Context, ssh bool) *url.URL {
	u := n.url // copy
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
