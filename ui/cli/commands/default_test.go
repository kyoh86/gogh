package commands_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	testtarget "github.com/kyoh86/gogh/v3/ui/cli/commands"
)

type testCase[T comparable] struct {
	name         string
	value        T
	defaultValue T
	want         T
}

func TestDefaultValue_int(t *testing.T) {
	intCases := []testCase[int]{
		{"int", 0, 10, 10},
		{"int", 5, 10, 5},
	}
	for _, tt := range intCases {
		t.Run(tt.name, func(t *testing.T) {
			got := testtarget.DefaultValue(tt.value, tt.defaultValue)
			if got != tt.want {
				t.Errorf("DefaultValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultValue_string(t *testing.T) {
	intCases := []testCase[string]{
		{"string", "", "default", "default"},
		{"string", "value", "default", "value"},
	}
	for _, tt := range intCases {
		t.Run(tt.name, func(t *testing.T) {
			got := testtarget.DefaultValue(tt.value, tt.defaultValue)
			if got != tt.want {
				t.Errorf("DefaultValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultSlice(t *testing.T) {
	cases := []struct {
		name         string
		value        []string
		defaultValue []string
		want         []string
	}{
		{"string slice", nil, []string{"default"}, []string{"default"}},
		{"string slice", []string{}, []string{"default"}, []string{}},
		{"string slice", []string{"value"}, []string{"default"}, []string{"value"}},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := testtarget.DefaultSlice(tt.value, tt.defaultValue)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("result mismatch;\n-want, +got\n%s", diff)
			}
		})
	}
}
