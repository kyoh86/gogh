package repo

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/kyoh86/gogh/internal/git"
)

var validName = regexp.MustCompile(`^([a-zA-Z0-9](?:(?:-[a-zA-Z0-9]+)*[a-zA-Z0-9])?)/(^[\w-]+)$`)

// var capital = regexp.MustCompile(`[A-Z]`) // UNDONE: warn if name contains capital cases

type Name struct {
	user string
	name string
	text string
}

type Shared string

var validShared = map[string]struct{}{
	"false":     struct{}{},
	"true":      struct{}{},
	"umask":     struct{}{},
	"group":     struct{}{},
	"all":       struct{}{},
	"world":     struct{}{},
	"everybody": struct{}{},
}

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

func (n *Name) Set(text string) error {
	matches := validName.FindStringSubmatch(text)
	if matches == nil {
		return fmt.Errorf("invalid repository name %q; repository should be specified as <user>/<name>", text)
	}
	n.user = matches[1]
	n.name = matches[2]
	n.text = text
	return nil
}

func (n Name) String() string {
	return n.text
}

func (n Name) User() string {
	return n.user
}

func (n Name) Name() string {
	return n.name
}

type Spec struct {
	ref   string
	https *url.URL
	ssh   *url.URL
}

func NewSpec(ref string) (*Spec, error) {
	spec := new(Spec)
	if err := spec.Set(ref); err != nil {
		return nil, err
	}
	return spec, nil
}

func (s *Spec) Set(ref string) error {
	if !hasSchemePattern.MatchString(ref) && scpLikeUrlPattern.MatchString(ref) {
		matched := scpLikeUrlPattern.FindStringSubmatch(ref)
		user := matched[1]
		host := matched[2]
		path := matched[3]

		ref = fmt.Sprintf("ssh://%s%s/%s", user, host, path)
	}

	httpsURL, err := url.Parse(ref)
	if err != nil {
		return err
	}

	if !httpsURL.IsAbs() {
		if !strings.Contains(httpsURL.Path, "/") {
			user, err := getUserName()
			if err != nil {
				return err
			}
			httpsURL.Path = user + "/" + httpsURL.Path
		}
		httpsURL.Scheme = "https"
		httpsURL.Host = "github.com"
		if httpsURL.Path[0] != '/' {
			httpsURL.Path = "/" + httpsURL.Path
		}
	}
	s.https = httpsURL
	sshURL, err := url.Parse(fmt.Sprintf("ssh://git@%s%s", httpsURL.Host, httpsURL.Path))
	if err != nil {
		return err
	}
	s.ssh = sshURL
	return nil
}

func (s *Spec) SSHURL() *url.URL {
	u := *s.ssh
	return &u
}

func (s *Spec) URL() *url.URL {
	u := *s.https
	return &u
}

func (s Spec) String() string {
	return s.ref
}

// Convert SCP-like URL to SSH URL(e.g. [user@]host.xz:path/to/repo.git/)
// ref. http://git-scm.com/docs/git-fetch#_git_urls
// (golang hasn't supported Perl-like negative look-behind match)
var hasSchemePattern = regexp.MustCompile("^[^:]+://")
var scpLikeUrlPattern = regexp.MustCompile("^([^@]+@)?([^:]+):/?(.+)$")

func getUserName() (string, error) {
	user, err := git.GetOneConf("gogh.user")
	if err != nil {
		return "", err
	}
	if user != "" {
		return user, nil
	}
	if user := os.Getenv("GITHUB_USER"); user != "" {
		return user, nil
	}
	switch runtime.GOOS {
	case "windows":
		if user := os.Getenv("USERNAME"); user != "" {
			return user, nil
		}
	default:
		if user := os.Getenv("USER"); user != "" {
			return user, nil
		}
	}
	// Make the error if it does not match any pattern
	return "", fmt.Errorf("set gogh.user to your gitconfig")
}

// Specs is array of Spec
type Specs []Spec

func (s *Specs) Set(value string) error {
	spec := new(Spec)
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

func (s *Spec) Remote(ssh bool) (Remote, error) {
	url := s.URL()
	if ssh {
		url = s.SSHURL()
	}

	rmt, err := NewRepository(url)
	if err != nil {
		return nil, err
	}

	if rmt.IsValid() == false {
		return nil, fmt.Errorf("Not a valid repository: %s", url)
	}
	return rmt, nil
}
