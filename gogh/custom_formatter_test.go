package gogh

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kyoh86/gogh/internal/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomListFormatter(t *testing.T) {
	t.Run("null separator", func(t *testing.T) {
		project1, err := parseProject(&context.MockContext{MGitHubHost: "github.com"}, "/go/src", "/go/src/github.com/kyoh86/foo")
		require.NoError(t, err)
		formatter, err := CustomFormatter("{{short .}}{{null}}{{full .}}{{null}}{{relative .}}")
		require.NoError(t, err)
		formatter.Add(project1)
		assert.Equal(t, 1, formatter.Len())
		var buf bytes.Buffer
		require.NoError(t, formatter.PrintAll(&buf, "\r\n"))
		expected := "foo\x00/go/src/github.com/kyoh86/foo\x00github.com/kyoh86/foo\r\n"
		assert.Equal(t, expected, buf.String())
	})
	t.Run("normal separator", func(t *testing.T) {
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

		formatter, err := CustomFormatter("{{short .}};;{{full .}};;{{relative .}}")
		require.NoError(t, err)
		formatter.Add(project1)
		formatter.Add(project2)
		formatter.Add(project3)
		formatter.Add(project4)
		formatter.Add(project5)
		formatter.Add(project6)
		assert.Equal(t, 6, formatter.Len())
		var buf bytes.Buffer
		require.NoError(t, formatter.PrintAll(&buf, "\n"))
		expected := `
foo                            ;; /go/src/ github.com/kyoh86/foo ;;github.com/kyoh86/foo             
github.com/kyoh86/bar          ;; /go/src/ github.com/kyoh86/bar ;;github.com/kyoh86/bar                               
kyoh87/bar                     ;; /go/src/ github.com/kyoh87/bar ;;github.com/kyoh87/bar                    
example.com/kyoh86/bar         ;; /go/src/ example.com/kyoh86/bar;;example.com/kyoh86/bar                               
/go/src/github.com/kyoh86/baz  ;; /go/src/ github.com/kyoh86/baz ;;github.com/kyoh86/baz                                       
/foo/github.com/kyoh86/baz     ;; /foo/    github.com/kyoh86/baz ;;github.com/kyoh86/baz                                    
`
		expected = strings.Replace(expected, " ", "", -1)
		expected = strings.TrimLeft(expected, "\n")
		assert.Equal(t, expected, buf.String())
	})
}