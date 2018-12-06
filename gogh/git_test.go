package gogh

import (
	"io/ioutil"
	"net/url"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitBackend(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "gogh-test")
	require.NoError(t, err)

	localDir := filepath.Join(tempDir, "repo")

	remoteURL, err := url.Parse("https://example.com/git/repo")
	require.NoError(t, err)

	commands := []*exec.Cmd{}
	lastCommand := func() *exec.Cmd { return commands[len(commands)-1] }
	commandRunner = func(cmd *exec.Cmd) error {
		commands = append(commands, cmd)
		return nil
	}

	err = gitClone(remoteURL, localDir, false)
	require.NoError(t, err)
	assert.Len(t, commands, 1)
	assert.Equal(t, []string{
		"git", "clone", remoteURL.String(), localDir,
	}, lastCommand().Args)

	err = gitClone(remoteURL, localDir, true)
	require.NoError(t, err)
	assert.Len(t, commands, 2)
	assert.Equal(t, []string{
		"git", "clone", "--depth", "1", remoteURL.String(), localDir,
	}, lastCommand().Args)

	err = gitUpdate(localDir)
	require.NoError(t, err)
	assert.Len(t, commands, 3)
	assert.Equal(t, []string{
		"git", "pull", "--ff-only",
	}, lastCommand().Args)
	assert.Equal(t, localDir, lastCommand().Dir)
}
