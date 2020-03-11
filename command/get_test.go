package command_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/gogh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	svc := initTest(t)
	defer svc.teardown(t)
	gitCtrl := gomock.NewController(t)

	// Assert that expected methods is invoked.
	defer gitCtrl.Finish()

	m := NewMockGitClient(gitCtrl)

	gomock.InOrder(
		m.EXPECT().Clone(gomock.Eq(filepath.Join(svc.root1, "github.com/kyoh86/gogh")), gomock.Any(), gomock.Eq(false)),
		m.EXPECT().Clone(gomock.Eq(filepath.Join(svc.root1, "github.com/kyoh86/vim-gogh")), gomock.Any(), gomock.Eq(false)),
	)
	assert.NoError(t, command.GetAll(svc.env, m, false, false, false, []gogh.Repo{
		*mustParseRepo(t, svc.env, "kyoh86/gogh"),
		*mustParseRepo(t, svc.env, "kyoh86/vim-gogh"),
	}))

	m.EXPECT().Clone(gomock.Eq(filepath.Join(svc.root1, "github.com/kyoh86/gogh")), gomock.Any(), gomock.Eq(false))
	assert.NoError(t, command.Get(svc.env, m, false, false, false, mustParseRepo(t, svc.env, "kyoh86/gogh")), "success getting one")

	require.NoError(t, os.MkdirAll(filepath.Join(svc.root1, "github.com", "kyoh86", "gogh", ".git"), 0755))
	assert.NoError(t, command.Get(svc.env, m, false, false, false, mustParseRepo(t, svc.env, "kyoh86/gogh")), "success getting one that is already exist")

	m.EXPECT().Update(gomock.Eq(filepath.Join(svc.root1, "github.com/kyoh86/gogh")))
	assert.NoError(t, command.Get(svc.env, m, true, false, false, mustParseRepo(t, svc.env, "kyoh86/gogh")), "success updating one that is already exist")
}
