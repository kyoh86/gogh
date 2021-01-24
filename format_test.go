package gogh_test

import (
	"path/filepath"
	"strings"
	"testing"

	testtarget "github.com/kyoh86/gogh/v2"
)

func TestFormat(t *testing.T) {
	description, err := testtarget.NewDescription("github.com", "kyoh86", "gogh")
	if err != nil {
		t.Fatalf("failed to init Description: %s", err)
	}
	project := testtarget.NewProject("/tmp", description)
	if err != nil {
		t.Fatalf("failed to get project from Description: %s", err)
	}

	// NOTE: When the path is checked, it should be passed with filepath.Clean.
	// Because windows uses '\' for path separator.
	for _, testcase := range []struct {
		title  string
		format testtarget.Format
		expect string
	}{
		{
			title:  "FullFilePath",
			format: testtarget.FormatFullFilePath,
			expect: filepath.Clean("/tmp/github.com/kyoh86/gogh"),
		},
		{
			title:  "RelPath",
			format: testtarget.FormatRelPath,
			expect: "github.com/kyoh86/gogh",
		},
		{
			title:  "RelFilePath",
			format: testtarget.FormatRelFilePath,
			expect: filepath.Clean("github.com/kyoh86/gogh"),
		},
		{
			title:  "URL",
			format: testtarget.FormatURL,
			expect: "https://github.com/kyoh86/gogh",
		},
		{
			title:  "FieldsWithSpace",
			format: testtarget.FormatFields(" "),
			expect: strings.Join([]string{
				filepath.Clean("/tmp/github.com/kyoh86/gogh"),
				"github.com/kyoh86/gogh",
				"github.com",
				"kyoh86",
				"gogh",
			}, " "),
		},
		{
			title:  "FieldsWithSpecial",
			format: testtarget.FormatFields("<<>>"),
			expect: strings.Join([]string{
				filepath.Clean("/tmp/github.com/kyoh86/gogh"),
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
	// UNDONE: Test JSON Formatter: it should be checked with assert.JSONEq
}
