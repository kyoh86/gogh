package repository_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	testtarget "github.com/kyoh86/gogh/v4/core/repository"
)

func TestReferenceJSON(t *testing.T) {
	ref := testtarget.NewReference("github.com", "kyoh86", "gogh")

	for _, testcase := range []struct {
		title string
		input any
		want  string
	}{
		{
			title: "bared",
			input: ref,
			want:  `{"host":"github.com","owner":"kyoh86","name":"gogh"}`,
		},
		{
			title: "pointer",
			input: &ref,
			want:  `{"host":"github.com","owner":"kyoh86","name":"gogh"}`,
		},
		{
			title: "wrap",
			input: struct {
				Reference testtarget.Reference
			}{Reference: ref},
			want: `{"Reference":{"host":"github.com","owner":"kyoh86","name":"gogh"}}`,
		},
		{
			title: "wrap pointer",
			input: struct {
				Reference *testtarget.Reference
			}{Reference: &ref},
			want: `{"Reference":{"host":"github.com","owner":"kyoh86","name":"gogh"}}`,
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			buf, err := json.Marshal(testcase.input)
			if err != nil {
				t.Fatal(err)
			}
			got := string(buf)
			if testcase.want != got {
				t.Errorf("result mismatch; want: %s, got: %s", testcase.want, got)
			}
		})
	}

	t.Run("Marshal & Unmarshal", func(t *testing.T) {
		buf, err := json.Marshal(ref)
		if err != nil {
			t.Fatal(err)
		}
		var got testtarget.Reference
		if err := json.Unmarshal(buf, &got); err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(ref, got, cmp.AllowUnexported(ref)); diff != "" {
			t.Errorf("result mismatch;\n-want, +got\n%s", diff)
		}
	})

	t.Run("Unmarshal invalid input", func(t *testing.T) {
		var got testtarget.Reference
		if err := json.Unmarshal([]byte(`{"host":42}`), &got); err == nil {
			t.Error("expected error, but got nil")
		}
	})
}
