package gogh

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandPath(t *testing.T) {
	sep := filepath.Separator

	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, expandPath(""))
	})

	t.Run("fullpath", func(t *testing.T) {
		path := filepath.Join(string(sep), "foo", "bar")
		assert.Equal(t, path, expandPath(path))
	})

	t.Run("homedir", func(t *testing.T) {
		assert.NotEqual(t, "~", expandPath("~"))
		assert.NotContains(t, expandPath("~"), "~")
	})

	t.Run("under the homedir", func(t *testing.T) {
		path := filepath.Join("~", "foo")
		assert.NotEqual(t, path, expandPath(path))
		assert.NotContains(t, expandPath(path), "~")
	})

	t.Run("tilda tilda", func(t *testing.T) {
		assert.Equal(t, "~~", expandPath("~~"))
	})
}
