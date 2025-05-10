package github

import (
	"testing"
)

func TestConvertSshUrlToHttps(t *testing.T) {
	// convertSshUrlToHttps
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "git@github.com:kyoh86/gogh.git",
			expected: "https://github.com/kyoh86/gogh",
		},
	}
	for _, c := range cases {
		actual := convertSSHToHTTPS(c.input)
		if actual != c.expected {
			t.Errorf("convertSshUrlToHttps(%s) = %s; expected %s", c.input, actual, c.expected)
		}
	}
}
