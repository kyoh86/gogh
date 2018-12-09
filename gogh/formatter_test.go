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
		local, err := parseLocal("/go/src", "/go/src/github.com/kyoh86/gogh")
		require.NoError(t, err)
		for _, f := range RepoListFormats() {
			formatter, err := RepoListFormat(f).Formatter()
			require.NoError(t, err)
			formatter.Add(local)
			require.NoError(t, formatter.PrintAll(ioutil.Discard, "\n"))
		}
	})

	t.Run("rel path formatter", func(t *testing.T) {
		local1, err := parseLocal("/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		local2, err := parseLocal("/go/src", "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		formatter, err := RepoListFormatRelPath.Formatter()
		require.NoError(t, err)
		formatter.Add(local1)
		formatter.Add(local2)
		var buf bytes.Buffer
		require.NoError(t, formatter.PrintAll(&buf, ":"))
		assert.Equal(t, `github.com/kyoh86/foo:github.com/kyoh86/bar:`, buf.String())
	})

	t.Run("full path formatter", func(t *testing.T) {
		local1, err := parseLocal("/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		local2, err := parseLocal("/go/src", "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		formatter, err := RepoListFormatFullPath.Formatter()
		require.NoError(t, err)
		formatter.Add(local1)
		formatter.Add(local2)
		var buf bytes.Buffer
		require.NoError(t, formatter.PrintAll(&buf, ":"))
		assert.Equal(t, `/go/src/github.com/kyoh86/foo:/go/src/github.com/kyoh86/bar:`, buf.String())
	})

	t.Run("short formatter", func(t *testing.T) {
		local1, err := parseLocal("/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		local2, err := parseLocal("/go/src", "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		local3, err := parseLocal("/go/src", "/go/src/github.com/kyoh87/bar")
		require.NoError(t, err)
		local4, err := parseLocal("/go/src", "/go/src/example.com/kyoh86/bar")
		require.NoError(t, err)
		local5, err := parseLocal("/go/src", "/go/src/github.com/kyoh86/baz")
		require.NoError(t, err)
		local6, err := parseLocal("/foo", "/foo/github.com/kyoh86/baz")
		require.NoError(t, err)
		formatter, err := RepoListFormatShort.Formatter()
		require.NoError(t, err)
		formatter.Add(local1)
		formatter.Add(local2)
		formatter.Add(local3)
		formatter.Add(local4)
		formatter.Add(local5)
		formatter.Add(local6)
		var buf bytes.Buffer
		require.NoError(t, formatter.PrintAll(&buf, ":"))
		assert.Equal(t, `foo:github.com/kyoh86/bar:kyoh87/bar:example.com/kyoh86/bar:/go/src/github.com/kyoh86/baz:/foo/github.com/kyoh86/baz:`, buf.String())
	})

	t.Run("invalid formatter", func(t *testing.T) {
		_, err := RepoListFormat("dummy").Formatter()
		assert.Errorf(t, err, "%q is invalid repo format", "dummy")
	})

	t.Run("writer error", func(t *testing.T) {
		local, err := parseLocal("/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		formatter, err := RepoListFormatShort.Formatter()
		require.NoError(t, err)
		formatter.Add(local)
		require.Error(t, formatter.PrintAll(&invalidWriter{}, ""), "invalid writer")
	})
}

type invalidWriter struct {
}

func (w *invalidWriter) Write([]byte) (int, error) {
	return 0, errors.New("invalid writer")
}
