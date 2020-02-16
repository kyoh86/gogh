// +build remote_test

package gogh_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// You can kick this test with `go test -tags remote_test

func TestRepoIsPublic(t *testing.T) {
	t.Run("public repo", func(t *testing.T) {
		r, err := ParseRepo("kyoh86/gogh")
		require.NoError(t, err)

		ctx := context.MockContext{MGitHubHost: "github.com", MGitHubUser: "kyoh86"}
		is, err := r.IsPublic(&ctx)
		require.NoError(t, err)
		assert.True(t, is)
	})

	t.Run("private repo", func(t *testing.T) {
		r, err := ParseRepo("kyoh86/unknown")
		require.NoError(t, err)

		ctx := context.MockContext{MGitHubHost: "github.com", MGitHubUser: "kyoh86"}
		is, err := r.IsPublic(&ctx)
		require.NoError(t, err)
		assert.False(t, is)
	})
}
