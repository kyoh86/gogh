package gogh

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)
	require.NoError(t, os.Setenv("GOPATH", tmp))
	require.NoError(t, os.MkdirAll(filepath.Join(tmp, "src"), 0755))

	t.Run("plain", func(t *testing.T) {
		_, _, err := capture(func() {
			require.NoError(t, List(false, false, false, false, ""))
		})
		require.NoError(t, err)
	})

	t.Run("short", func(t *testing.T) {
		_, _, err := capture(func() {
			require.NoError(t, List(false, false, true, false, ""))
		})
		require.NoError(t, err)
	})
}
