package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	all, err := GetAllConf("gogh.non.existent.key")
	assert.NoError(t, err)
	assert.Empty(t, all)
}
