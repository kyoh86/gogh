package view_test

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kyoh86/gogh/v3"
	testtarget "github.com/kyoh86/gogh/v3/view"
)

func TestProjectFormat(t *testing.T) {
	tempDir := t.TempDir()
	spec, err := gogh.NewSpec("github.com", "kyoh86", "gogh")
	if err != nil {
		t.Fatalf("failed to init Spec: %s", err)
	}
	project := gogh.NewProject(tempDir, spec)
	if err != nil {
		t.Fatalf("failed to get project from Spec: %s", err)
	}
	if err := gogh.CreateLocalProject(context.Background(), project, spec.URL(), nil); err != nil {
		t.Fatalf("failed to prepare local project from Spec: %s", err)
	}

	wantPath := filepath.Join(tempDir, "github.com/kyoh86/gogh")

	// NOTE: When the path is checked, it should be passed with filepath.Clean.
	// Because windows uses '\' for path separator.
	for _, testcase := range []struct {
		title  string
		format testtarget.ProjectFormat
		expect string
	}{
		{
			title:  "FullFilePath",
			format: testtarget.ProjectFormatFullFilePath,
			expect: wantPath,
		},
		{
			title:  "RelPath",
			format: testtarget.ProjectFormatRelPath,
			expect: "github.com/kyoh86/gogh",
		},
		{
			title:  "RelFilePath",
			format: testtarget.ProjectFormatRelFilePath,
			expect: filepath.Clean("github.com/kyoh86/gogh"),
		},
		{
			title:  "URL",
			format: testtarget.ProjectFormatURL,
			expect: "https://github.com/kyoh86/gogh",
		},
		{
			title:  "FieldsWithSpace",
			format: testtarget.ProjectFormatFields(" "),
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
			format: testtarget.ProjectFormatFields("<<>>"),
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
			actual, err := testcase.format.Format(project)
			if err != nil {
				t.Fatalf("failed to format: %s", err)
			}
			if testcase.expect != actual {
				t.Errorf("expect %q but %q is gotten", testcase.expect, actual)
			}
		})
	}

	t.Run("JSON", func(t *testing.T) {
		formatted, err := testtarget.ProjectFormatJSON(project)
		if err != nil {
			t.Fatalf("failed to format: %s", err)
		}
		var got map[string]interface{}
		if err := json.Unmarshal([]byte(formatted), &got); err != nil {
			t.Fatalf("failed to unmarshal JSON formatted: %s", err)
		}
		want := map[string]interface{}{
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
