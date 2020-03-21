package gogh_test

import (
	"testing"

	"github.com/kyoh86/gogh/gogh"
	"github.com/stretchr/testify/assert"
)

func TestValidateName(t *testing.T) {
	assert.EqualError(t, gogh.ValidateName(""), "project name is empty", "empty project name")
	assert.EqualError(t, gogh.ValidateName("."), "'.' is reserved name", "'dot' conflicts with 'current directory'")
	assert.EqualError(t, gogh.ValidateName(".."), "'..' is reserved name", "'dot' conflicts with 'parent directory'")
	assert.EqualError(t, gogh.ValidateName("kyoh86/gogh"), "invalid project name", "slashes must not be contained in project name")
	assert.NoError(t, gogh.ValidateName("----..--.."), "hyphens and dots are usable in project name")
}

func TestValidateOwner(t *testing.T) {
	expect := "invalid owner name"
	assert.EqualError(t, gogh.ValidateOwner(""), expect, "fail when empty owner is given")
	assert.EqualError(t, gogh.ValidateOwner("kyoh_86"), expect, "fail when owner name contains invalid character")
	assert.EqualError(t, gogh.ValidateOwner("-kyoh86"), expect, "fail when owner name starts with hyphen")
	assert.EqualError(t, gogh.ValidateOwner("kyoh86-"), expect, "fail when owner name ends with hyphen")
	assert.NoError(t, gogh.ValidateOwner("kyoh86"), "success")
}
