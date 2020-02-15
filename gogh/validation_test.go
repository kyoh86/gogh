package gogh

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/kyoh86/gogh/internal/context"
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
	assert.EqualError(t, ValidateOwner("kyoh_86"), expect, "fail when owner name contains invalid character")
	assert.EqualError(t, ValidateOwner("-kyoh86"), expect, "fail when owner name starts with hyphen")
	assert.EqualError(t, ValidateOwner("kyoh86-"), expect, "fail when owner name ends with hyphen")
	assert.NoError(t, ValidateOwner("kyoh86"), "success")
}

func TestValidateRoots(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)
	assert.EqualError(t, ValidateRoots([]string{}), "no root", "fail when no path in root")
	assert.NoError(t, ValidateRoots([]string{"/path/to/not/existing", tmp}))
	assert.Error(t, ValidateRoots([]string{"\x00", tmp}))
}

func TestValidateContext(t *testing.T) {
	t.Run("invalid root", func(t *testing.T) {
		ctx := &context.MockContext{
			MRoot:       []string{"/\x00"},
			MGitHubUser: "kyoh86",
		}
		assert.Error(t, ValidateContext(ctx))
	})
	t.Run("invalid owner", func(t *testing.T) {
		ctx := &context.MockContext{
			MRoot:       []string{"/path/to/not/existing"},
			MGitHubUser: "",
		}
		assert.Error(t, ValidateContext(ctx))
	})
	t.Run("valid context", func(t *testing.T) {
		ctx := &context.MockContext{
			MRoot:       []string{"/path/to/not/existing"},
			MGitHubUser: "kyoh86",
		}
		assert.NoError(t, ValidateContext(ctx))
	})
}
