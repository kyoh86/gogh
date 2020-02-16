package command_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	svc := initTest(t)
	defer svc.teardown(t)
	teardown := testutil.Stubin(t, []byte("y\n"))
	defer teardown()

	proj1 := filepath.Join(svc.root1, "github.com", "kyoh86", "gogh-test-1", ".git")
	require.NoError(t, os.MkdirAll(proj1, 0755))
	proj2 := filepath.Join(svc.root2, "github.com", "kyoh86", "gogh-test-2", ".git")
	require.NoError(t, os.MkdirAll(proj2, 0755))
	proj3 := filepath.Join(svc.root2, "github.com", "kyoh85", "gogh-test-3", ".git")
	require.NoError(t, os.MkdirAll(proj3, 0755))

	t.Run("delete proj2 explicitly", func(t *testing.T) {
		assert.NoError(t, command.Delete(svc.ctx, false, "gogh-test-2"))
		var err error
		_, err = os.Stat(proj1)
		assert.NoError(t, err)
		_, err = os.Stat(proj2)
		assert.True(t, os.IsNotExist(err))
		require.NoError(t, os.MkdirAll(proj2, 0755))
		_, err = os.Stat(proj3)
		assert.NoError(t, err)
	})

	t.Run("delete proj3 with fuzzy", func(t *testing.T) {
		teardown := testutil.Stubin(t, []byte("y\n"))
		defer teardown()

		assert.NoError(t, command.Delete(svc.ctx, false, "3"))
		var err error
		_, err = os.Stat(proj1)
		assert.NoError(t, err)
		_, err = os.Stat(proj2)
		assert.NoError(t, err)
		_, err = os.Stat(proj3)
		assert.True(t, os.IsNotExist(err))
		require.NoError(t, os.MkdirAll(proj3, 0755))
	})

	t.Run("delete proj1 in primary", func(t *testing.T) {
		teardown := testutil.Stubin(t, []byte("y\n"))
		defer teardown()

		assert.NoError(t, command.Delete(svc.ctx, true, "test"))
		var err error
		_, err = os.Stat(proj1)
		assert.True(t, os.IsNotExist(err))
		require.NoError(t, os.MkdirAll(proj1, 0755))
		_, err = os.Stat(proj2)
		assert.NoError(t, err)
		_, err = os.Stat(proj3)
		assert.NoError(t, err)
	})

	t.Run("did not match", func(t *testing.T) {
		assert.EqualError(t, command.Delete(svc.ctx, true, "foobar"), "any projects did not matched for \"foobar\"")
		var err error
		_, err = os.Stat(proj1)
		assert.NoError(t, err)
		_, err = os.Stat(proj2)
		assert.NoError(t, err)
		_, err = os.Stat(proj3)
		assert.NoError(t, err)
	})
}
