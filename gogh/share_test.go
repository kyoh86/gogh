package gogh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepoShared(t *testing.T) {
	t.Run("valid shareds", func(t *testing.T) {
		var shared RepoShared
		assert.NoError(t, shared.Set("false"))
		assert.Equal(t, "false", shared.String())
		assert.NoError(t, shared.Set("true"))
		assert.Equal(t, "true", shared.String())
		assert.NoError(t, shared.Set("umask"))
		assert.Equal(t, "umask", shared.String())
		assert.NoError(t, shared.Set("group"))
		assert.Equal(t, "group", shared.String())
		assert.NoError(t, shared.Set("all"))
		assert.Equal(t, "all", shared.String())
		assert.NoError(t, shared.Set("world"))
		assert.Equal(t, "world", shared.String())
		assert.NoError(t, shared.Set("everybody"))
		assert.Equal(t, "everybody", shared.String())
		assert.NoError(t, shared.Set("0777"))
		assert.Equal(t, "0777", shared.String())
		assert.NoError(t, shared.Set("777"))
		assert.Equal(t, "777", shared.String())
	})
	t.Run("invalid shared", func(t *testing.T) {
		var shared RepoShared
		assert.NotNil(t, shared.Set("gogh"))
		assert.NotNil(t, shared.Set("800"))
	})
	t.Run("invalid shared (triple term)", func(t *testing.T) {
		var shared RepoShared
		assert.NotNil(t, shared.Set("github.com/kyoh86/gogh"))
	})
}
