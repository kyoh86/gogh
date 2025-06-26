package repository_print_test

import (
	"bytes"
	"context"
	"iter"
	"testing"
	"time"

	testtarget "github.com/kyoh86/gogh/v4/app/repository_print"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/repository"
)

func sliceToIter2[T any](slices []T) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for _, v := range slices {
			if !yield(v, nil) {
				return
			}
		}
	}
}

func TestRepositoryPrinter(t *testing.T) {
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
		title  string
		format string
		want   string
	}{
		{
			title:  "ref",
			format: "ref",
			want:   "github.com/kyoh86/gogh\n",
		},
		{
			title:  "url",
			format: "url",
			want:   "https://github.com/kyoh86/gogh\n",
		},
		{
			title:  "json",
			format: "json",
			want:   `{"ref":{"host":"github.com","owner":"kyoh86","name":"gogh"},"url":"https://github.com/kyoh86/gogh","updatedAt":"2021-05-01T01:00:00Z"}` + "\n",
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			var buf bytes.Buffer
			printer := testtarget.NewUsecase(&buf, testcase.format)
			if err := printer.Execute(context.Background(), sliceToIter2([]*hosting.Repository{&repo})); err != nil {
				t.Fatal(err)
			}
			got := buf.String()
			if testcase.want != got {
				t.Errorf("result mismatched; want: %s; got: %s", testcase.want, got)
			}
		})
	}
}
