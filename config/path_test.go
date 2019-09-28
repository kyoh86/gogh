package config

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandPath(t *testing.T) {
	user, err := user.Current()
	require.NoError(t, err)
	expHome := user.HomeDir

	t.Run("empty to empty", func(t *testing.T) {
		act := expandPath("")
		assert.Equal(t, "", act)
	})

	t.Run("success to expand `~` to homedir", func(t *testing.T) {
		act := expandPath("~")
		assert.Equal(t, expHome, act)
	})

	t.Run("success to expand `~/foo` to `foo` under the homedir", func(t *testing.T) {
		act := expandPath("~/foo")
		assert.Equal(t, filepath.Join(expHome, "foo"), act)
	})

	t.Run("should not to expand abs path", func(t *testing.T) {
		act := expandPath("/foo/bar")
		assert.Equal(t, "/foo/bar", act)
	})

	t.Run("should not to expand tilde-prefixed-name", func(t *testing.T) {
		act := expandPath("~foo")
		assert.Equal(t, "~foo", act)
	})

	t.Run("success to expand env", func(t *testing.T) {
		os.Setenv("FOO", "bar")
		act := expandPath("~/$FOO")
		assert.Equal(t, filepath.Join(expHome, "bar"), act)
	})
}
