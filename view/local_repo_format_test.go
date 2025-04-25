package view_test

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kyoh86/gogh/v3"
	"github.com/kyoh86/gogh/v3/domain/reporef"
	testtarget "github.com/kyoh86/gogh/v3/view"
)

func TestLocalRepoFormat(t *testing.T) {
	tempDir := t.TempDir()
	ref, err := reporef.NewRepoRef("github.com", "kyoh86", "gogh")
	if err != nil {
		t.Fatalf("failed to init Ref: %s", err)
	}
	locRepo := gogh.NewLocalRepo(tempDir, ref)
	if err != nil {
		t.Fatalf("failed to get a local repository from Ref: %s", err)
	}
	if err := gogh.CreateLocalRepo(context.Background(), locRepo, ref.URL(), nil); err != nil {
		t.Fatalf("failed to prepare local repository from Ref: %s", err)
	}

	wantPath := filepath.Join(tempDir, "github.com/kyoh86/gogh")

	// NOTE: When the path is checked, it should be passed with filepath.Clean.
	// Because windows uses '\' for path separator.
	for _, testcase := range []struct {
		title  string
		format testtarget.LocalRepoFormat
		expect string
	}{
		{
			title:  "FullFilePath",
			format: testtarget.LocalRepoFormatFullFilePath,
			expect: wantPath,
		},
		{
			title:  "RelPath",
			format: testtarget.LocalRepoFormatRelPath,
			expect: "github.com/kyoh86/gogh",
		},
		{
			title:  "RelFilePath",
			format: testtarget.LocalRepoFormatRelFilePath,
			expect: filepath.Clean("github.com/kyoh86/gogh"),
		},
		{
			title:  "URL",
			format: testtarget.LocalRepoFormatURL,
			expect: "https://github.com/kyoh86/gogh",
		},
		{
			title:  "FieldsWithSpace",
			format: testtarget.LocalRepoFormatFields(" "),
			expect: strings.Join([]string{
				wantPath,
				filepath.Clean("github.com/kyoh86/gogh"),
				"https://github.com/kyoh86/gogh",
				"github.com/kyoh86/gogh",
				"github.com",
				"kyoh86",
				"gogh",
			}, " "),
		},
		{
			title:  "FieldsWithSpecial",
			format: testtarget.LocalRepoFormatFields("<<>>"),
			expect: strings.Join([]string{
				wantPath,
				filepath.Clean("github.com/kyoh86/gogh"),
				"https://github.com/kyoh86/gogh",
				"github.com/kyoh86/gogh",
				"github.com",
				"kyoh86",
				"gogh",
			}, "<<>>"),
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			actual, err := testcase.format.Format(locRepo)
			if err != nil {
				t.Fatalf("failed to format: %s", err)
			}
			if testcase.expect != actual {
				t.Errorf("expect %q but %q is gotten", testcase.expect, actual)
			}
		})
	}

	t.Run("JSON", func(t *testing.T) {
		formatted, err := testtarget.LocalRepoFormatJSON(locRepo)
		if err != nil {
			t.Fatalf("failed to format: %s", err)
		}
		var got map[string]any
		if err := json.Unmarshal([]byte(formatted), &got); err != nil {
			t.Fatalf("failed to unmarshal JSON formatted: %s", err)
		}
		want := map[string]any{
			"fullFilePath": wantPath,
			"relFilePath":  filepath.Clean("github.com/kyoh86/gogh"),
			"url":          "https://github.com/kyoh86/gogh",
			"relPath":      "github.com/kyoh86/gogh",
			"host":         "github.com",
			"owner":        "kyoh86",
			"name":         "gogh",
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("json obj mismatch (-want +got):\n%s", diff)
		}
	})
}
