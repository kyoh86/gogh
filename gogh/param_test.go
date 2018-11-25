package gogh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParam(t *testing.T) {
	var param []string
	param = appendIf(param, "-a", false)
	param = appendIf(param, "-b", true)
	param = appendIf(param, "-c", false)
	param = appendIf(param, "-d", true)
	param = appendIfFilled(param, "-e", "")
	param = appendIfFilled(param, "-f", "file1")
	param = appendIfFilled(param, "-g", "")
	param = appendIfFilled(param, "-h", "file2")
	assert.Equal(t, []string{"-b", "-d", "-f", "file1", "-h", "file2"}, param)
}
