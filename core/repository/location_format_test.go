package repository_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	testtarget "github.com/kyoh86/gogh/v3/core/repository"
)

func TestLocationFormat(t *testing.T) {
	loc := testtarget.NewLocation(
		"/path/to/workspace/github.com/kyoh86/gogh",
		"github.com",
		"kyoh86",
		"gogh",
	)

	// NOTE: When the path is checked, it should be passed with filepath.Clean.
	// Because windows uses '\' for path separator.
	for _, testcase := range []struct {
		title  string
		format testtarget.LocationFormat
		expect string
	}{
		{
			title:  "FullPath",
			format: testtarget.LocationFormatFullPath,
			expect: loc.FullPath(),
		},
		{
			title:  "Path",
			format: testtarget.LocationFormatPath,
			expect: loc.Path(),
		},
		{
			title:  "FieldsWithSpace",
			format: testtarget.LocationFormatFields(" "),
			expect: strings.Join([]string{
				loc.FullPath(),
				loc.Path(),
				loc.Host(),
				loc.Owner(),
				loc.Name(),
			}, " "),
		},
		{
			title:  "FieldsWithSpecial",
			format: testtarget.LocationFormatFields("<<>>"),
			expect: strings.Join([]string{
				loc.FullPath(),
				loc.Path(),
				loc.Host(),
				loc.Owner(),
				loc.Name(),
			}, "<<>>"),
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			actual, err := testcase.format.Format(*loc)
			if err != nil {
				t.Fatalf("failed to format: %s", err)
			}
			if testcase.expect != actual {
				t.Errorf("expect %q but %q is gotten", testcase.expect, actual)
			}
		})
	}

	t.Run("JSON", func(t *testing.T) {
		formatted, err := testtarget.LocationFormatJSON(*loc)
		if err != nil {
			t.Fatalf("failed to format: %s", err)
		}
		var got map[string]any
		if err := json.Unmarshal([]byte(formatted), &got); err != nil {
			t.Fatalf("failed to unmarshal JSON formatted: %s", err)
		}
		want := map[string]any{
			"fullPath": loc.FullPath(),
			"path":     loc.Path(),
			"host":     loc.Host(),
			"owner":    loc.Owner(),
			"name":     loc.Name(),
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("json obj mismatch (-want +got):\n%s", diff)
		}
	})
}
