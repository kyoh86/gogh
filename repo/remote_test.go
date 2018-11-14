package repo

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"
	"testing"

	"github.com/kyoh86/gogh/internal/run"
	. "github.com/onsi/gomega"
)

func parseURL(urlString string) *url.URL {
	u, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}
	return u
}

func TestNewRepositoryGitHub(t *testing.T) {
	RegisterTestingT(t)

	var (
		repo Remote
		err  error
	)

	repo, err = NewRepository(parseURL("https://github.com/motemen/pusheen-explorer"))
	Expect(err).To(BeNil())
	Expect(repo.IsValid()).To(Equal(true))

	repo, err = NewRepository(parseURL("https://github.com/motemen/pusheen-explorer/"))
	Expect(err).To(BeNil())
	Expect(repo.IsValid()).To(Equal(true))

	repo, err = NewRepository(parseURL("https://github.com/motemen/pusheen-explorer/blob/master/README.md"))
	Expect(err).To(BeNil())
	Expect(repo.IsValid()).To(Equal(false))

	repo, err = NewRepository(parseURL("https://example.com/motemen/pusheen-explorer"))
	Expect(err).To(BeNil())
	Expect(repo.IsValid()).To(Equal(true))
}

func newFakeRunner(dispatch map[string]error) run.RunFunc {
	return func(cmd *exec.Cmd) error {
		cmdString := strings.Join(cmd.Args, " ")
		for cmdPrefix, err := range dispatch {
			if strings.Index(cmdString, cmdPrefix) == 0 {
				return err
			}
		}
		panic(fmt.Sprintf("No fake dispatch found for: %s", cmdString))
	}
}

func TestNewRepositoryGitHubGist(t *testing.T) {
	RegisterTestingT(t)

	var (
		repo Remote
		err  error
	)

	repo, err = NewRepository(parseURL("https://gist.github.com/motemen/9733745"))
	Expect(err).To(BeNil())
	Expect(repo.IsValid()).To(Equal(true))
}

func TestNewRepositoryGoogleCode(t *testing.T) {
	RegisterTestingT(t)

	var (
		repo Remote
		err  error
	)

	repo, err = NewRepository(parseURL("https://code.google.com/p/git-core"))
	Expect(err).To(BeNil())
	Expect(repo.IsValid()).To(Equal(true))
	run.CommandRunner = newFakeRunner(map[string]error{
		"git ls-remote": nil,
	})
}
