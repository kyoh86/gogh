package view_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	testtarget "github.com/kyoh86/gogh/v3/ui/cli/view"
)

func TestRemoteRepoPrinters(t *testing.T) {
	uat, err := time.Parse(time.RFC3339, "2021-05-01T01:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	ref := repository.NewReference("github.com", "kyoh86", "gogh")
	repo := hosting.Repository{
		UpdatedAt: uat,
		Ref:       ref,
		URL:       "https://github.com/kyoh86/gogh",
	}
	for _, testcase := range []struct {
		title   string
		printer func(io.Writer) testtarget.RemoteRepoPrinter
		want    string
	}{
		{
			title:   "ref",
			printer: testtarget.NewRemoteRepoRefPrinter,
			want:    "github.com/kyoh86/gogh\n",
		},
		{
			title:   "url",
			printer: testtarget.NewRemoteRepoURLPrinter,
			want:    "https://github.com/kyoh86/gogh\n",
		},
		{
			title:   "json",
			printer: testtarget.NewRemoteRepoJSONPrinter,
			want:    `{"ref":{"host":"github.com","owner":"kyoh86","name":"gogh"},"url":"https://github.com/kyoh86/gogh","updatedAt":"2021-05-01T01:00:00Z"}` + "\n",
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
