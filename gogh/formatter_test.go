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
		project, err := parseProject(&implContext{gitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/gogh")
		require.NoError(t, err)
		for _, f := range ProjectListFormats() {
			formatter, err := ProjectListFormat(f).Formatter()
			require.NoError(t, err)
			formatter.Add(project)
			require.NoError(t, formatter.PrintAll(ioutil.Discard, "\n"))
		}
	})

	t.Run("rel path formatter", func(t *testing.T) {
		project1, err := parseProject(&implContext{gitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		project2, err := parseProject(&implContext{gitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		formatter, err := ProjectListFormatRelPath.Formatter()
		require.NoError(t, err)
		formatter.Add(project1)
		formatter.Add(project2)
		assert.Equal(t, 2, formatter.Len())
		var buf bytes.Buffer
		require.NoError(t, formatter.PrintAll(&buf, ":"))
		assert.Equal(t, `github.com/kyoh86/foo:github.com/kyoh86/bar:`, buf.String())
	})
	t.Run("writer error by rel path formatter", func(t *testing.T) {
		project, err := parseProject(&implContext{gitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		formatter, err := ProjectListFormatRelPath.Formatter()
		require.NoError(t, err)
		formatter.Add(project)
		require.EqualError(t, formatter.PrintAll(&invalidWriter{}, ""), "invalid writer")
	})

	t.Run("full path formatter", func(t *testing.T) {
		project1, err := parseProject(&implContext{gitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		project2, err := parseProject(&implContext{gitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		formatter, err := ProjectListFormatFullPath.Formatter()
		require.NoError(t, err)
		formatter.Add(project1)
		formatter.Add(project2)
		assert.Equal(t, 2, formatter.Len())
		var buf bytes.Buffer
		require.NoError(t, formatter.PrintAll(&buf, ":"))
		assert.Equal(t, `/go/src/github.com/kyoh86/foo:/go/src/github.com/kyoh86/bar:`, buf.String())
	})
	t.Run("writer error by full path formatter", func(t *testing.T) {
		project, err := parseProject(&implContext{gitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		formatter, err := ProjectListFormatFullPath.Formatter()
		require.NoError(t, err)
		formatter.Add(project)
		require.EqualError(t, formatter.PrintAll(&invalidWriter{}, ""), "invalid writer")
	})

	t.Run("url formatter", func(t *testing.T) {
		project1, err := parseProject(&implContext{gitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		project2, err := parseProject(&implContext{gitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		formatter, err := ProjectListFormatURL.Formatter()
		require.NoError(t, err)
		formatter.Add(project1)
		formatter.Add(project2)
		assert.Equal(t, 2, formatter.Len())
		var buf bytes.Buffer
		require.NoError(t, formatter.PrintAll(&buf, ":"))
		assert.Equal(t, `https://github.com/kyoh86/foo:https://github.com/kyoh86/bar:`, buf.String())
	})
	t.Run("writer error by url formatter", func(t *testing.T) {
		project, err := parseProject(&implContext{gitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		formatter, err := ProjectListFormatURL.Formatter()
		require.NoError(t, err)
		formatter.Add(project)
		require.EqualError(t, formatter.PrintAll(&invalidWriter{}, ""), "invalid writer")
	})

	t.Run("short formatter", func(t *testing.T) {
		project1, err := parseProject(&implContext{gitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		project2, err := parseProject(&implContext{gitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		project3, err := parseProject(&implContext{gitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh87/bar")
		require.NoError(t, err)
		project4, err := parseProject(&implContext{gitHubHost: "example.com"}, "/go/src", "/go/src/example.com/kyoh86/bar")
		require.NoError(t, err)
		project5, err := parseProject(&implContext{gitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/baz")
		require.NoError(t, err)
		project6, err := parseProject(&implContext{gitHubHost: "github.com"}, "/foo", "/foo/github.com/kyoh86/baz")
		require.NoError(t, err)
		formatter, err := ProjectListFormatShort.Formatter()
		require.NoError(t, err)
		formatter.Add(project1)
		formatter.Add(project2)
		formatter.Add(project3)
		formatter.Add(project4)
		formatter.Add(project5)
		formatter.Add(project6)
		assert.Equal(t, 6, formatter.Len())
		var buf bytes.Buffer
		require.NoError(t, formatter.PrintAll(&buf, ":"))
		assert.Equal(t, `foo:github.com/kyoh86/bar:kyoh87/bar:example.com/kyoh86/bar:/go/src/github.com/kyoh86/baz:/foo/github.com/kyoh86/baz:`, buf.String())
	})
	t.Run("writer error by short formatter", func(t *testing.T) {
		project, err := parseProject(&implContext{gitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		formatter, err := ProjectListFormatShort.Formatter()
		require.NoError(t, err)
		formatter.Add(project)
		require.EqualError(t, formatter.PrintAll(&invalidWriter{}, ""), "invalid writer")
	})

	t.Run("invalid formatter", func(t *testing.T) {
		_, err := ProjectListFormat("dummy").Formatter()
		assert.Errorf(t, err, "%q is invalid project format", "dummy")
	})

}

type invalidWriter struct {
}

func (w *invalidWriter) Write([]byte) (int, error) {
	return 0, errors.New("invalid writer")
}
