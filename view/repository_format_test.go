package view_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/kyoh86/gogh/v2"
	testtarget "github.com/kyoh86/gogh/v2/view"
)

func TestRepositoryPrinters(t *testing.T) {
	spec, err := gogh.NewSpec("github.com", "kyoh86", "gogh")
	if err != nil {
		t.Fatal(err)
	}
	repo := gogh.Repository{
		Spec: spec,
		URL:  "https://github.com/kyoh86/gogh",
	}
	for _, testcase := range []struct {
		title   string
		printer func(io.Writer) testtarget.RepositoryPrinter
		want    string
	}{
		{
			title:   "spec",
			printer: testtarget.NewRepositorySpecPrinter,
			want:    "github.com/kyoh86/gogh\n",
		},
		{
			title:   "url",
			printer: testtarget.NewRepositoryURLPrinter,
			want:    "https://github.com/kyoh86/gogh\n",
		},
		{
			title:   "json",
			printer: testtarget.NewRepositoryJSONPrinter,
			want:    `{"spec":{"host":"github.com","owner":"kyoh86","name":"gogh"},"url":"https://github.com/kyoh86/gogh"}` + "\n",
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			var buf bytes.Buffer
			printer := testcase.printer(&buf)
			if err := printer.Print(repo); err != nil {
				t.Fatal(err)
			}
			if err := printer.Close(); err != nil {
				t.Fatal(err)
			}
			got := buf.String()
			if testcase.want != got {
				t.Errorf("result mismatched; want: %s; got: %s", testcase.want, got)
			}
		})
	}
}
