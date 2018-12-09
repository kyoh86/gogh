package gogh

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestLocalName(t *testing.T) {
	t.Run("valid name", func(t *testing.T) {
		var name LocalName
		require.NoError(t, name.Set("kyoh86/gogh"))
		assert.Equal(t, "kyoh86/gogh", name.String())
		assert.Equal(t, "kyoh86", name.User())
		assert.Equal(t, "gogh", name.Name())
	})
	t.Run("invalid name (single term)", func(t *testing.T) {
		var name LocalName
		assert.NotNil(t, name.Set("gogh"))
	})
	t.Run("invalid name (triple term)", func(t *testing.T) {
		var name LocalName
		assert.NotNil(t, name.Set("github.com/kyoh86/gogh"))
	})
}

func TestRemoteName(t *testing.T) {
	t.Run("full HTTPS URL", func(t *testing.T) {
		ctx := &implContext{}
		spec, err := ParseRemoteName("https://github.com/kyoh86/pusheen-explorer")
		require.NoError(t, err)
		assert.Equal(t, "https://github.com/kyoh86/pusheen-explorer", spec.String())
		assert.Equal(t, "https://github.com/kyoh86/pusheen-explorer", spec.URL(ctx, false).String())
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer", spec.URL(ctx, true).String())
	})

	t.Run("scp like URL 1", func(t *testing.T) {
		ctx := &implContext{}
		spec, err := ParseRemoteName("git@github.com:kyoh86/pusheen-explorer.git")
		require.NoError(t, err)
		assert.Equal(t, "git@github.com:kyoh86/pusheen-explorer.git", spec.String())
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer.git", spec.URL(ctx, false).String())
	})

	t.Run("scp like URL 2", func(t *testing.T) {
		ctx := &implContext{}
		spec, err := ParseRemoteName("git@github.com:/kyoh86/pusheen-explorer.git")
		require.NoError(t, err)
		assert.Equal(t, "git@github.com:/kyoh86/pusheen-explorer.git", spec.String())
		assert.Equal(t, "ssh://git@github.com/kyoh86/pusheen-explorer.git", spec.URL(ctx, false).String())
	})

	t.Run("scp like URL 3", func(t *testing.T) {
		ctx := &implContext{}
		spec, err := ParseRemoteName("github.com:kyoh86/pusheen-explorer.git")
		require.NoError(t, err)
		assert.Equal(t, "github.com:kyoh86/pusheen-explorer.git", spec.String())
		assert.Equal(t, "ssh://github.com/kyoh86/pusheen-explorer.git", spec.URL(ctx, false).String())
	})

	t.Run("owner/name spec", func(t *testing.T) {
		ctx := &implContext{}
		spec, err := ParseRemoteName("kyoh86/gogh")
		require.NoError(t, err)
		assert.Equal(t, "kyoh86/gogh", spec.String())
		assert.Equal(t, "https://github.com/kyoh86/gogh", spec.URL(ctx, false).String())
	})

	t.Run("name only spec", func(t *testing.T) {
		ctx := &implContext{userName: "kyoh86"}
		spec, err := ParseRemoteName("gogh")
		require.NoError(t, err)
		assert.Equal(t, "gogh", spec.String())
		assert.Equal(t, "https://github.com/kyoh86/gogh", spec.URL(ctx, false).String())
	})

	t.Run("fail when invalid url given", func(t *testing.T) {
		_, err := ParseRemoteName("://////")
		assert.NotNil(t, err)
	})
}

func TestRemoteNames(t *testing.T) {
	var specs RemoteNames
	require.True(t, specs.IsCumulative())
	assert.Empty(t, specs.String())
	assert.NoError(t, specs.Set("kyoh86/gogh"), "owner/name spec")
	assert.NoError(t, specs.Set("gogh"), "name only spec")
	assert.NotNil(t, specs.Set("://////"), "fail when invalid url given")
	assert.Len(t, specs, 2)
	assert.Equal(t, "kyoh86/gogh", specs[0].String())
	assert.Equal(t, "gogh", specs[1].String())
	assert.Equal(t, "kyoh86/gogh,gogh", specs.String())
}
