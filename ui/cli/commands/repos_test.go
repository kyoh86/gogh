package commands

import (
	"testing"
)

func TestQuoteEnums(t *testing.T) {
	for _, testcase := range []struct {
		title  string
		want   string
		source []string
	}{
		{
			title:  "empty",
			source: []string{},
			want:   "",
		},
		{
			title:  "nil",
			source: nil,
			want:   "",
		},
		{
			title:  "minimal",
			source: []string{"a", "b"},
			want:   `"a" or "b"`,
		},
		{
			title:  "not minimal",
			source: []string{"a", "b", "c"},
			want:   `"a", "b" or "c"`,
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			got := quoteEnums(testcase.source)
			if testcase.want != got {
				t.Errorf("quoteEnums(%v) = %q, want %q", testcase.source, got, testcase.want)
			}
		})
	}
}
