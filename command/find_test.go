package command_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFind(t *testing.T) {
	svc := initTest(t)
	defer svc.teardown(t)

	proj1 := filepath.Join(svc.root1, "github.com", "kyoh86", "vim-gogh", ".git")
	require.NoError(t, os.MkdirAll(proj1, 0755))
	proj2 := filepath.Join(svc.root2, "github.com", "kyoh86", "gogh", ".git")
	require.NoError(t, os.MkdirAll(proj2, 0755))
	proj3 := filepath.Join(svc.root2, "github.com", "kyoh85", "test", ".git")
	require.NoError(t, os.MkdirAll(proj3, 0755))

	svc.ev.EXPECT().GithubUser().Return("kyoh86")
	assert.EqualError(t, command.Find(svc.ev, true, mustParseRepoSpec(t, "gogh")), "project not found")

	svc.ev.EXPECT().GithubUser().Return("kyoh86")
	assert.NoError(t, command.Find(svc.ev, false, mustParseRepoSpec(t, "gogh")))

	svc.ev.EXPECT().GithubUser().Return("kyoh86")
	assert.NoError(t, command.Find(svc.ev, false, mustParseRepoSpec(t, "kyoh85/test")))

	assert.NoError(t, command.Find(svc.ev, true, mustParseRepoSpec(t, "vim-gogh")))
}
