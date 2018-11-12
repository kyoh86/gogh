package git

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAll(t *testing.T) {
	all, err := GetAllConf("gogh.non.existent.key")
	assert.NoError(t, err)
	assert.Empty(t, all)
}

func WithGitconfigFile(t *testing.T, configContent string) func() {
	t.Helper()
	tmpdir, err := ioutil.TempDir("", "gogh-test")
	require.NoError(t, err)

	tmpGitconfigFile := filepath.Join(tmpdir, "gitconfig")

	ioutil.WriteFile(
		tmpGitconfigFile,
		[]byte(configContent),
		0777,
	)

	prevGitConfigEnv := os.Getenv("GIT_CONFIG")
	os.Setenv("GIT_CONFIG", tmpGitconfigFile)

	return func() {
		os.Setenv("GIT_CONFIG", prevGitConfigEnv)
	}
}
