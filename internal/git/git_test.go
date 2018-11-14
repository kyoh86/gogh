package git

import (
	"io/ioutil"
	"net/url"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/internal/run"
	. "github.com/onsi/gomega"
)

func TestGitBackend(t *testing.T) {
	RegisterTestingT(t)

	tempDir, err := ioutil.TempDir("", "gogh-test")
	if err != nil {
		t.Fatal(err)
	}

	localDir := filepath.Join(tempDir, "repo")

	remoteURL, err := url.Parse("https://example.com/git/repo")
	if err != nil {
		t.Fatal(err)
	}

	commands := []*exec.Cmd{}
	lastCommand := func() *exec.Cmd { return commands[len(commands)-1] }
	run.CommandRunner = func(cmd *exec.Cmd) error {
		commands = append(commands, cmd)
		return nil
	}

	err = Clone(remoteURL, localDir, false)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(1))
	Expect(lastCommand().Args).To(Equal([]string{
		"git", "clone", remoteURL.String(), localDir,
	}))

	err = Clone(remoteURL, localDir, true)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(2))
	Expect(lastCommand().Args).To(Equal([]string{
		"git", "clone", "--depth", "1", remoteURL.String(), localDir,
	}))

	err = Update(localDir)

	Expect(err).NotTo(HaveOccurred())
	Expect(commands).To(HaveLen(3))
	Expect(lastCommand().Args).To(Equal([]string{
		"git", "pull", "--ff-only",
	}))
	Expect(lastCommand().Dir).To(Equal(localDir))
}
