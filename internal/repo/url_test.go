package repo

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindGitHubURL(t *testing.T) {
	got, err := findGitHubURL([]string{
		"https://google.com/?hl=ja",
		"http://github.com/me/mine",
		"https://github.com/me/mine.git",
		"https://github.com/you/yours.git",
	})
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "https", got.Scheme)
	assert.Equal(t, "github.com", got.Host)
	assert.Equal(t, "/me/mine.git", got.Path)
}

func TestParseIdentifier(t *testing.T) {
	id, err := parseIdentifier(&url.URL{
		Scheme: "TestScheme",
		Host:   "TestHost",
		Path:   "/me/mine.git",
	})
	require.NoError(t, err)
	require.NotNil(t, id)
	assert.Equal(t, "TestScheme", id.Scheme)
	assert.Equal(t, "TestHost", id.Host)
	assert.Equal(t, "me", id.Owner)
	assert.Equal(t, "mine", id.Name)
}
