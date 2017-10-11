package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckScopes(t *testing.T) {
	{
		assert.NoError(t, checkScopes([]string{"public_repo", "repo", "user"}))
	}
	{
		err := checkScopes([]string{}).(*scopeError)
		assert.Error(t, err)
		assert.Contains(t, err.required, "repo")
		assert.Contains(t, err.required, "user")
	}
	{
		err := checkScopes(nil).(*scopeError)
		assert.Error(t, err)
		assert.Contains(t, err.required, "repo")
		assert.Contains(t, err.required, "user")
	}
	{
		err := checkScopes([]string{"user", "admin"}).(*scopeError)
		assert.Error(t, err)
		assert.Contains(t, err.required, "repo")
		assert.NotContains(t, err.required, "user")
		assert.NotContains(t, err.required, "admin")
	}
}
