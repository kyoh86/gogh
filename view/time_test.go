package view_test

import (
	"testing"
	"time"

	testtarget "github.com/kyoh86/gogh/v2/view"
)

func TestFuzzyAgoAbbr(t *testing.T) {
	now, err := time.Parse(time.RFC3339, "2021-05-01T01:00:00.000Z")
	if err != nil {
		t.Fatal(err)
	}
	for _, testcase := range []struct {
		title string
		now   time.Time
		at    time.Time
		want  string
	}{
		{
			title: "now",
			now:   now,
			at:    now,
			want:  "0m",
		},
		{
			title: "59m59.999s",
			now:   now,
			at:    now.Add(-59*time.Minute - 59*time.Second - 999*time.Millisecond),
			want:  "59m",
		},
		{
			title: "23h59m59.999s",
			now:   now,
			at:    now.Add(-23*time.Hour - 59*time.Minute - 59*time.Second - 999*time.Millisecond),
			want:  "23h",
		},
		{
			title: "29d23h59m59.999s",
			now:   now,
			at:    now.Add(-29*24*time.Hour - 23*time.Hour - 59*time.Minute - 59*time.Second - 999*time.Millisecond),
			want:  "29d",
		},
		{
			title: "30d",
			now:   now,
			at:    now.Add(-30 * 24 * time.Hour),
			want:  "2021-04-01",
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			got := testtarget.FuzzyAgoAbbr(testcase.now, testcase.at)
			if got != testcase.want {
				t.Errorf("want: %s, got: %s", testcase.want, got)
			}
		})
	}
}
