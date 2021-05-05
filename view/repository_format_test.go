package view_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/kyoh86/gogh/v2"
	testtarget "github.com/kyoh86/gogh/v2/view"
)

func TestRepositoryPrinters(t *testing.T) {
	uat, err := time.Parse(time.RFC3339, "2021-05-01T01:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	spec, err := gogh.NewSpec("github.com", "kyoh86", "gogh")
	if err != nil {
		t.Fatal(err)
	}
	repo := gogh.Repository{
		UpdatedAt: uat,
		Spec:      spec,
		URL:       "https://github.com/kyoh86/gogh",
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
			want:    `{"updatedAt":"2021-05-01T01:00:00Z","spec":{"host":"github.com","owner":"kyoh86","name":"gogh"},"url":"https://github.com/kyoh86/gogh"}` + "\n",
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
