package config

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/joeshaw/envdecode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func TestPathListOption(t *testing.T) {

	type testStruct struct {
		PathList PathListOption `env:"PATH_LIST" yaml:"paths,omitempty"`
	}

	t.Run("encode to yaml", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, yaml.NewEncoder(&buf).Encode(testStruct{}))
		assert.Equal(t, "{}", strings.TrimSpace(buf.String()))

		buf.Reset()
		require.NoError(t, yaml.NewEncoder(&buf).Encode(testStruct{PathList: PathListOption{}}))
		assert.Equal(t, "{}", strings.TrimSpace(buf.String()))

		buf.Reset()
		require.NoError(t, yaml.NewEncoder(&buf).Encode(testStruct{PathList: PathListOption{"/foo", "/bar"}}))
		assert.Equal(t, "paths:\n- /foo\n- /bar", strings.TrimSpace(buf.String()))
	})
	t.Run("decode from YAML", func(t *testing.T) {
		var testValue testStruct
		require.NoError(t, yaml.Unmarshal([]byte(`{}`), &testValue))
		assert.Equal(t, PathListOption(nil), testValue.PathList)

		require.NoError(t, yaml.Unmarshal([]byte(`paths: []`), &testValue))
		assert.Equal(t, PathListOption{}, testValue.PathList)

		require.NoError(t, yaml.Unmarshal([]byte(`paths: ["/foo", "/bar"]`), &testValue))
		assert.Equal(t, PathListOption{"/foo", "/bar"}, testValue.PathList)

		assert.Error(t, yaml.Unmarshal([]byte(`paths: invalid`), &testValue))
	})
	t.Run("get from envar", func(t *testing.T) {
		var testValue testStruct
		resetEnv(t)
		require.NoError(t, os.Setenv("PATH_LIST", "/foo:/bar"))
		require.NoError(t, envdecode.Decode(&testValue))
		assert.Equal(t, PathListOption{"/foo", "/bar"}, testValue.PathList)
	})
}
