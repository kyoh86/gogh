package gogh

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/kyoh86/gogh/internal/context"
	"github.com/kyoh86/gogh/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatter(t *testing.T) {
	t.Run("dry run formatters", func(t *testing.T) {
		project, err := parseProject(&context.MockContext{MGitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/gogh")
		require.NoError(t, err)
		for _, formatter := range []ProjectListFormatter{
			ShortFormatter(),
			URLFormatter(),
			FullPathFormatter(),
			RelPathFormatter(),
		} {
			require.NoError(t, err)
			formatter.Add(project)
			require.NoError(t, formatter.PrintAll(ioutil.Discard, "\n"))
		}
	})

	t.Run("rel path formatter", func(t *testing.T) {
		project1, err := parseProject(&context.MockContext{MGitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		project2, err := parseProject(&context.MockContext{MGitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		formatter := RelPathFormatter()
		require.NoError(t, err)
		formatter.Add(project1)
		formatter.Add(project2)
		assert.Equal(t, 2, formatter.Len())
		var buf bytes.Buffer
		require.NoError(t, formatter.PrintAll(&buf, ":"))
		assert.Equal(t, `github.com/kyoh86/foo:github.com/kyoh86/bar:`, buf.String())
	})
	t.Run("writer error by rel path formatter", func(t *testing.T) {
		project, err := parseProject(&context.MockContext{MGitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		formatter := RelPathFormatter()
		require.NoError(t, err)
		formatter.Add(project)
		require.EqualError(t, formatter.PrintAll(testutil.DefaultErrorWriter, ""), "error writer")
	})

	t.Run("full path formatter", func(t *testing.T) {
		project1, err := parseProject(&context.MockContext{MGitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		project2, err := parseProject(&context.MockContext{MGitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		formatter := FullPathFormatter()
		require.NoError(t, err)
		formatter.Add(project1)
		formatter.Add(project2)
		assert.Equal(t, 2, formatter.Len())
	})
	t.Run("writer error by full path formatter", func(t *testing.T) {
		project, err := parseProject(&context.MockContext{MGitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		formatter := FullPathFormatter()
		require.NoError(t, err)
		formatter.Add(project)
		require.EqualError(t, formatter.PrintAll(testutil.DefaultErrorWriter, ""), "error writer")
	})

	t.Run("url formatter", func(t *testing.T) {
		project1, err := parseProject(&context.MockContext{MGitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		project2, err := parseProject(&context.MockContext{MGitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		formatter := URLFormatter()
		require.NoError(t, err)
		formatter.Add(project1)
		formatter.Add(project2)
		assert.Equal(t, 2, formatter.Len())
		var buf bytes.Buffer
		require.NoError(t, formatter.PrintAll(&buf, ":"))
		assert.Equal(t, `https://github.com/kyoh86/foo:https://github.com/kyoh86/bar:`, buf.String())
	})
	t.Run("writer error by url formatter", func(t *testing.T) {
		project, err := parseProject(&context.MockContext{MGitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		formatter := URLFormatter()
		require.NoError(t, err)
		formatter.Add(project)
		require.EqualError(t, formatter.PrintAll(testutil.DefaultErrorWriter, ""), "error writer")
	})

	t.Run("short formatter", func(t *testing.T) {
		project1, err := parseProject(&context.MockContext{MGitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		project2, err := parseProject(&context.MockContext{MGitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		project3, err := parseProject(&context.MockContext{MGitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh87/bar")
		require.NoError(t, err)
		project4, err := parseProject(&context.MockContext{MGitHubHost: "example.com"}, "/go/src", "/go/src/example.com/kyoh86/bar")
		require.NoError(t, err)
		project5, err := parseProject(&context.MockContext{MGitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/baz")
		require.NoError(t, err)
		project6, err := parseProject(&context.MockContext{MGitHubHost: "github.com"}, "/foo", "/foo/github.com/kyoh86/baz")
		require.NoError(t, err)
		formatter := ShortFormatter()
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
		project, err := parseProject(&context.MockContext{MGitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		formatter := ShortFormatter()
		require.NoError(t, err)
		formatter.Add(project)
		require.EqualError(t, formatter.PrintAll(testutil.DefaultErrorWriter, ""), "error writer")
	})

}
