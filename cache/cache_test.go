package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	c := new(Cache)
	c.SetGithubUser("valid-key", "kyoh86")
	assert.Equal(t, "kyoh86", c.GetGithubUser("valid-key"))
	assert.Empty(t, c.GetGithubUser("invalid-key"))
}
