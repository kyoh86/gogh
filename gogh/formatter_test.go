package gogh

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatter(t *testing.T) {
	t.Run("dry run formatters", func(t *testing.T) {
		ctx := &implContext{
			roots: []string{"/go/src"},
		}
		repo, err := FromFullPath(ctx, "/go/src/github.com/kyoh86/gogh")
		require.NoError(t, err)
		for _, f := range RepoListFormats() {
			formatter, err := RepoListFormat(f).Formatter()
			require.NoError(t, err)
			formatter.Add(repo)
			require.NoError(t, formatter.PrintAll(ioutil.Discard, "\n"))
		}
	})

	t.Run("rel path formatter", func(t *testing.T) {
		ctx := &implContext{
			roots: []string{"/go/src"},
		}
		repo1, err := FromFullPath(ctx, "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		repo2, err := FromFullPath(ctx, "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		formatter, err := RepoListFormatRelPath.Formatter()
		require.NoError(t, err)
		formatter.Add(repo1)
		formatter.Add(repo2)
		var buf bytes.Buffer
		require.NoError(t, formatter.PrintAll(&buf, ":"))
		assert.Equal(t, `github.com/kyoh86/foo:github.com/kyoh86/bar:`, buf.String())
	})

	t.Run("full path formatter", func(t *testing.T) {
		ctx := &implContext{
			roots: []string{"/go/src"},
		}
		repo1, err := FromFullPath(ctx, "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		repo2, err := FromFullPath(ctx, "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		formatter, err := RepoListFormatFullPath.Formatter()
		require.NoError(t, err)
		formatter.Add(repo1)
		formatter.Add(repo2)
		var buf bytes.Buffer
		require.NoError(t, formatter.PrintAll(&buf, ":"))
		assert.Equal(t, `/go/src/github.com/kyoh86/foo:/go/src/github.com/kyoh86/bar:`, buf.String())
	})

	t.Run("short formatter", func(t *testing.T) {
		ctx := &implContext{
			roots: []string{
				"/go/src",
				"/foo",
			},
		}
		repo1, err := FromFullPath(ctx, "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		repo2, err := FromFullPath(ctx, "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		repo3, err := FromFullPath(ctx, "/go/src/github.com/kyoh87/bar")
		require.NoError(t, err)
		repo4, err := FromFullPath(ctx, "/go/src/example.com/kyoh86/bar")
		require.NoError(t, err)
		repo5, err := FromFullPath(ctx, "/go/src/github.com/kyoh86/baz")
		require.NoError(t, err)
		repo6, err := FromFullPath(ctx, "/foo/github.com/kyoh86/baz")
		require.NoError(t, err)
		formatter, err := RepoListFormatShort.Formatter()
		require.NoError(t, err)
		formatter.Add(repo1)
		formatter.Add(repo2)
		formatter.Add(repo3)
		formatter.Add(repo4)
		formatter.Add(repo5)
		formatter.Add(repo6)
		var buf bytes.Buffer
		require.NoError(t, formatter.PrintAll(&buf, ":"))
		assert.Equal(t, `foo:github.com/kyoh86/bar:kyoh87/bar:example.com/kyoh86/bar:/go/src/github.com/kyoh86/baz:/foo/github.com/kyoh86/baz:`, buf.String())
	})

	t.Run("invalid formatter", func(t *testing.T) {
		_, err := RepoListFormat("dummy").Formatter()
		assert.Errorf(t, err, "%q is invalid repo format", "dummy")
	})

	t.Run("writer error", func(t *testing.T) {
		ctx := &implContext{
			roots: []string{"/go/src"},
		}
		repo, err := FromFullPath(ctx, "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		formatter, err := RepoListFormatShort.Formatter()
		require.NoError(t, err)
		formatter.Add(repo)
		require.Error(t, formatter.PrintAll(&invalidWriter{}, ""), "invalid writer")
	})
}

type invalidWriter struct {
}

func (w *invalidWriter) Write([]byte) (int, error) {
	return 0, errors.New("invalid writer")
}
