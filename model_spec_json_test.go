package gogh_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	testtarget "github.com/kyoh86/gogh/v3"
)

func TestSpecJSON(t *testing.T) {
	spec, err := testtarget.NewSpec("github.com", "kyoh86", "gogh")
	if err != nil {
		t.Fatal(err)
	}

	for _, testcase := range []struct {
		title string
		input interface{}
		want  string
	}{
		{
			title: "bared",
			input: spec,
			want:  `{"host":"github.com","owner":"kyoh86","name":"gogh"}`,
		},
		{
			title: "pointer",
			input: &spec,
			want:  `{"host":"github.com","owner":"kyoh86","name":"gogh"}`,
		},
		{
			title: "wrap",
			input: struct {
				Spec testtarget.Spec
			}{Spec: spec},
			want: `{"Spec":{"host":"github.com","owner":"kyoh86","name":"gogh"}}`,
		},
		{
			title: "wrap pointer",
			input: struct {
				Spec *testtarget.Spec
			}{Spec: &spec},
			want: `{"Spec":{"host":"github.com","owner":"kyoh86","name":"gogh"}}`,
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
		buf, err := json.Marshal(spec)
		if err != nil {
			t.Fatal(err)
		}
		var got testtarget.Spec
		if err := json.Unmarshal(buf, &got); err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(spec, got, cmp.AllowUnexported(spec)); diff != "" {
			t.Errorf("result mismatch;\n-want, +got\n%s", diff)
		}
	})

	t.Run("Unmarshal invalid input", func(t *testing.T) {
		var got testtarget.Spec
		if err := json.Unmarshal([]byte(`{"host":42}`), &got); err == nil {
			t.Error("expected error, but got nil")
		}
	})
}
