package gogh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
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
