package gogh

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateName(t *testing.T) {
	assert.EqualError(t, ValidateName(""), "empty project name", "empty project name is invalid")
	assert.EqualError(t, ValidateName("."), "'.' or '..' is reserved name", "'dot' conflicts with 'current directory'")
	assert.EqualError(t, ValidateName(".."), "'.' or '..' is reserved name", "'dot' conflicts with 'parent directory'")
	assert.EqualError(t, ValidateName("kyoh86/gogh"), "project name may only contain alphanumeric characters, dots or hyphens", "slashes must not be contained in project name")
	assert.NoError(t, ValidateName("----..--.."), "hyphens and dots are usable in project name")
}

func TestValidateOwner(t *testing.T) {
	expect := "owner name may only contain alphanumeric characters or single hyphens, and cannot begin or end with a hyphen"
	assert.EqualError(t, ValidateOwner(""), expect, "fail when empty owner is given")
	assert.EqualError(t, ValidateOwner("kyoh_86"), expect, "fail when owner name contains invalid charactor")
	assert.EqualError(t, ValidateOwner("-kyoh86"), expect, "fail when owner name starts with hyphen")
	assert.EqualError(t, ValidateOwner("kyoh86-"), expect, "fail when owner name ends with hyphen")
	assert.NoError(t, ValidateOwner("kyoh86"), "success")
}

func TestValidateRoot(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)
	assert.EqualError(t, ValidateRoot([]string{}), "no root", "fail when no path in root")
	assert.NoError(t, ValidateRoot([]string{"/path/to/not/existing", tmp}))
}

func TestValidateContext(t *testing.T) {
	ctx := &Config{
		VRoot: RootConfig{"/path/to/not/existing"},
	}
	ctx.GitHub.User = ""
	assert.Error(t, ValidateContext(ctx), "fail when empty owner is given")
	ctx.GitHub.User = "kyoh86"
	assert.NoError(t, ValidateContext(ctx), "success")
}
