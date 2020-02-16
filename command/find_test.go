package command_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kyoh86/gogh/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFind(t *testing.T) {
	root1, err := ioutil.TempDir(os.TempDir(), "gogh-test1")
	require.NoError(t, err)
	defer os.RemoveAll(root1)
	root2, err := ioutil.TempDir(os.TempDir(), "gogh-test2")
	require.NoError(t, err)
	defer os.RemoveAll(root2)

	proj1 := filepath.Join(root1, "github.com", "kyoh86", "vim-gogh", ".git")
	require.NoError(t, os.MkdirAll(proj1, 0755))
	proj2 := filepath.Join(root2, "github.com", "kyoh86", "gogh", ".git")
	require.NoError(t, os.MkdirAll(proj2, 0755))
	proj3 := filepath.Join(root2, "github.com", "kyoh85", "test", ".git")
	require.NoError(t, os.MkdirAll(proj3, 0755))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := NewMockContext(ctrl)

	ctx.EXPECT().Root().AnyTimes().Return([]string{root1, root2})
	ctx.EXPECT().PrimaryRoot().AnyTimes().Return(root1)
	ctx.EXPECT().GitHubHost().AnyTimes().Return("github.com")
	ctx.EXPECT().GitHubUser().AnyTimes().Return("kyoh86")
	ctx.EXPECT().Done().AnyTimes()

	assert.EqualError(t, command.Find(ctx, true, mustParseRepo(t, ctx, "gogh")), "project not found")

	assert.NoError(t, command.Find(ctx, false, mustParseRepo(t, ctx, "gogh")))

	assert.NoError(t, command.Find(ctx, false, mustParseRepo(t, ctx, "kyoh85/test")))

	assert.NoError(t, command.Find(ctx, true, mustParseRepo(t, ctx, "vim-gogh")))
}
