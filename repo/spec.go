package repo

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/kyoh86/gogh/internal/git"
)

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
