package gogh

import (
	"io/ioutil"
	"testing"

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
	//TODO: check output
	//TODO: error case
}
