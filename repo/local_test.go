package repo

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewRepository(t *testing.T) {
	RegisterTestingT(t)

	roots = []string{"/repos"}

	r, err := FromFullPath("/repos/github.com/kyoh86/gogh")
	Expect(err).To(BeNil())
	Expect(r.NonHostPath()).To(Equal("kyoh86/gogh"))
	Expect(r.Subpaths()).To(Equal([]string{"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"}))

	r, err = FromFullPath("/repos/stash.com/scm/kyoh86/gogh")
	Expect(err).To(BeNil())
	Expect(r.NonHostPath()).To(Equal("scm/kyoh86/gogh"))
	Expect(r.Subpaths()).To(Equal([]string{"gogh", "kyoh86/gogh", "scm/kyoh86/gogh", "stash.com/scm/kyoh86/gogh"}))

	githubURL, _ := url.Parse("ssh://git@github.com/kyoh86/gogh.git")
	r, err = FromURL(githubURL)
	Expect(err).To(BeNil())
	Expect(r.FullPath).To(Equal("/repos/github.com/kyoh86/gogh"))

	stashURL, _ := url.Parse("ssh://git@stash.com/scm/kyoh86/gogh")
	r, err = FromURL(stashURL)
	Expect(err).To(BeNil())
	Expect(r.FullPath).To(Equal("/repos/stash.com/scm/kyoh86/gogh"))

	svnSourceforgeURL, _ := url.Parse("http://svn.code.sf.net/p/gogh/code/trunk")
	r, err = FromURL(svnSourceforgeURL)
	Expect(err).To(BeNil())
	Expect(r.FullPath).To(Equal("/repos/svn.code.sf.net/p/gogh/code/trunk"))

	gitSourceforgeURL, _ := url.Parse("http://git.code.sf.net/p/gogh/code")
	r, err = FromURL(gitSourceforgeURL)
	Expect(err).To(BeNil())
	Expect(r.FullPath).To(Equal("/repos/git.code.sf.net/p/gogh/code"))

	svnSourceforgeJpURL, _ := url.Parse("http://scm.sourceforge.jp/svnroot/gogh/")
	r, err = FromURL(svnSourceforgeJpURL)
	Expect(err).To(BeNil())
	Expect(r.FullPath).To(Equal("/repos/scm.sourceforge.jp/svnroot/gogh"))

	gitSourceforgeJpURL, _ := url.Parse("http://scm.sourceforge.jp/gitroot/gogh/gogh.git")
	r, err = FromURL(gitSourceforgeJpURL)
	Expect(err).To(BeNil())
	Expect(r.FullPath).To(Equal("/repos/scm.sourceforge.jp/gitroot/gogh/gogh"))

	svnAssemblaURL, _ := url.Parse("https://subversion.assembla.com/svn/gogh/")
	r, err = FromURL(svnAssemblaURL)
	Expect(err).To(BeNil())
	Expect(r.FullPath).To(Equal("/repos/subversion.assembla.com/svn/gogh"))

	gitAssemblaURL, _ := url.Parse("https://git.assembla.com/gogh.git")
	r, err = FromURL(gitAssemblaURL)
	Expect(err).To(BeNil())
	Expect(r.FullPath).To(Equal("/repos/git.assembla.com/gogh"))
}

// https://gist.github.com/kyanny/c231f48e5d08b98ff2c3
func TestList_Symlink(t *testing.T) {
	RegisterTestingT(t)

	root, err := ioutil.TempDir("", "")
	Expect(err).To(BeNil())

	symDir, err := ioutil.TempDir("", "")
	Expect(err).To(BeNil())

	roots = []string{root}

	err = os.MkdirAll(filepath.Join(root, "github.com", "atom", "atom", ".git"), 0777)
	Expect(err).To(BeNil())

	err = os.MkdirAll(filepath.Join(root, "github.com", "zabbix", "zabbix", ".git"), 0777)
	Expect(err).To(BeNil())

	err = os.Symlink(symDir, filepath.Join(root, "github.com", "gogh"))
	Expect(err).To(BeNil())

	paths := []string{}
	Walk(func(repo *Local) error {
		paths = append(paths, repo.RelPath)
		return nil
	})

	Expect(paths).To(HaveLen(2))
}
