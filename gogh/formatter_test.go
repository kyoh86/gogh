package gogh_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := NewMockContext(ctrl)
	ctx.EXPECT().GithubHost().AnyTimes().Return("github.com")

	t.Run("dry run formatters", func(t *testing.T) {
		project, err := gogh.ParseProject(ctx, "/go/src", "/go/src/github.com/kyoh86/gogh")
		require.NoError(t, err)
		for _, formatter := range []gogh.ProjectListFormatter{
			gogh.ShortFormatter(),
			gogh.URLFormatter(),
			gogh.FullPathFormatter(),
			gogh.RelPathFormatter(),
		} {
			require.NoError(t, err)
			formatter.Add(project)
			require.NoError(t, formatter.PrintAll(ioutil.Discard, "\n"))
		}
	})

	t.Run("rel path formatter", func(t *testing.T) {
		project1, err := gogh.ParseProject(ctx, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		project2, err := gogh.ParseProject(ctx, "/go/src", "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		formatter := gogh.RelPathFormatter()
		require.NoError(t, err)
		formatter.Add(project1)
		formatter.Add(project2)
		assert.Equal(t, 2, formatter.Len())
		var buf bytes.Buffer
		require.NoError(t, formatter.PrintAll(&buf, ":"))
		assert.Equal(t, `github.com/kyoh86/foo:github.com/kyoh86/bar:`, buf.String())
	})
	t.Run("writer error by rel path formatter", func(t *testing.T) {
		project, err := gogh.ParseProject(ctx, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		formatter := gogh.RelPathFormatter()
		require.NoError(t, err)
		formatter.Add(project)
		require.EqualError(t, formatter.PrintAll(testutil.DefaultErrorWriter, ""), "error writer")
	})

	t.Run("full path formatter", func(t *testing.T) {
		project1, err := gogh.ParseProject(ctx, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		project2, err := gogh.ParseProject(ctx, "/go/src", "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		formatter := gogh.FullPathFormatter()
		require.NoError(t, err)
		formatter.Add(project1)
		formatter.Add(project2)
		assert.Equal(t, 2, formatter.Len())
	})
	t.Run("writer error by full path formatter", func(t *testing.T) {
		project, err := gogh.ParseProject(ctx, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		formatter := gogh.FullPathFormatter()
		require.NoError(t, err)
		formatter.Add(project)
		require.EqualError(t, formatter.PrintAll(testutil.DefaultErrorWriter, ""), "error writer")
	})

	t.Run("url formatter", func(t *testing.T) {
		project1, err := gogh.ParseProject(ctx, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		project2, err := gogh.ParseProject(ctx, "/go/src", "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		formatter := gogh.URLFormatter()
		require.NoError(t, err)
		formatter.Add(project1)
		formatter.Add(project2)
		assert.Equal(t, 2, formatter.Len())
		var buf bytes.Buffer
		require.NoError(t, formatter.PrintAll(&buf, ":"))
		assert.Equal(t, `https://github.com/kyoh86/foo:https://github.com/kyoh86/bar:`, buf.String())
	})
	t.Run("writer error by url formatter", func(t *testing.T) {
		project, err := gogh.ParseProject(ctx, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		formatter := gogh.URLFormatter()
		require.NoError(t, err)
		formatter.Add(project)
		require.EqualError(t, formatter.PrintAll(testutil.DefaultErrorWriter, ""), "error writer")
	})

	t.Run("short formatter", func(t *testing.T) {
		project1, err := gogh.ParseProject(ctx, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		project2, err := gogh.ParseProject(ctx, "/go/src", "/go/src/github.com/kyoh86/bar")
		require.NoError(t, err)
		project3, err := gogh.ParseProject(ctx, "/go/src", "/go/src/github.com/kyoh87/bar")
		require.NoError(t, err)

		expCtrl := gomock.NewController(t)
		defer expCtrl.Finish()
		expCtx := NewMockContext(expCtrl)
		expCtx.EXPECT().GithubHost().AnyTimes().Return("example.com")

		project4, err := gogh.ParseProject(expCtx, "/go/src", "/go/src/example.com/kyoh86/bar")
		require.NoError(t, err)
		project5, err := gogh.ParseProject(ctx, "/go/src", "/go/src/github.com/kyoh86/baz")
		require.NoError(t, err)
		project6, err := gogh.ParseProject(ctx, "/foo", "/foo/github.com/kyoh86/baz")
		require.NoError(t, err)
		formatter := gogh.ShortFormatter()
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
		project, err := gogh.ParseProject(ctx, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		formatter := gogh.ShortFormatter()
		require.NoError(t, err)
		formatter.Add(project)
		require.EqualError(t, formatter.PrintAll(testutil.DefaultErrorWriter, ""), "error writer")
	})
}
