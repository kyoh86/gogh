package view_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	testtarget "github.com/kyoh86/gogh/v3/ui/cli/view"
)

type repoRef struct {
	fullPath string
	path     string
	host     string
	owner    string
	name     string
}

// Host is a hostname (i.g.: "github.com")
func (r *repoRef) Host() string { return r.host }

// Owner is a owner name (i.g.: "kyoh86")
func (r *repoRef) Owner() string { return r.owner }

// Name of the repository (i.g.: "gogh")
func (r *repoRef) Name() string { return r.name }

// Path returns the path from root of the repository (i.g.: "github.com/kyoh86/gogh")
func (r *repoRef) Path() string { return r.path }

// FullPath returns the full path of the repository (i.g.: "/path/to/workspace/github.com/kyoh86/gogh")
func (r *repoRef) FullPath() string { return r.fullPath }

func TestLocalRepoFormat(t *testing.T) {
	repo := &repoRef{
		fullPath: "/path/to/workspace/github.com/kyoh86/gogh",
		path:     "github.com/kyoh86/gogh",
		host:     "github.com",
		owner:    "kyoh86",
		name:     "gogh",
	}

	// NOTE: When the path is checked, it should be passed with filepath.Clean.
	// Because windows uses '\' for path separator.
	for _, testcase := range []struct {
		title  string
		format testtarget.LocalRepoFormat
		expect string
	}{
		{
			title:  "FullPath",
			format: testtarget.LocalRepoFormatFullPath,
			expect: repo.fullPath,
		},
		{
			title:  "Path",
			format: testtarget.LocalRepoFormatPath,
			expect: repo.path,
		},
		{
			title:  "FieldsWithSpace",
			format: testtarget.LocalRepoFormatFields(" "),
			expect: strings.Join([]string{
				repo.fullPath,
				repo.path,
				repo.host,
				repo.owner,
				repo.name,
			}, " "),
		},
		{
			title:  "FieldsWithSpecial",
			format: testtarget.LocalRepoFormatFields("<<>>"),
			expect: strings.Join([]string{
				repo.fullPath,
				repo.path,
				repo.host,
				repo.owner,
				repo.name,
			}, "<<>>"),
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			actual, err := testcase.format.Format(repo)
			if err != nil {
				t.Fatalf("failed to format: %s", err)
			}
			if testcase.expect != actual {
				t.Errorf("expect %q but %q is gotten", testcase.expect, actual)
			}
		})
	}

	t.Run("JSON", func(t *testing.T) {
		formatted, err := testtarget.LocalRepoFormatJSON(repo)
		if err != nil {
			t.Fatalf("failed to format: %s", err)
		}
		var got map[string]any
		if err := json.Unmarshal([]byte(formatted), &got); err != nil {
			t.Fatalf("failed to unmarshal JSON formatted: %s", err)
		}
		want := map[string]any{
			"fullPath": repo.fullPath,
			"path":     repo.path,
			"host":     repo.host,
			"owner":    repo.owner,
			"name":     repo.name,
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("json obj mismatch (-want +got):\n%s", diff)
		}
	})
}
