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

func TestBoolOption(t *testing.T) {
	assert.True(t, TrueOption.Bool())
	assert.False(t, FalseOption.Bool())

	type testStruct struct {
		Bool BoolOption `env:"BOOL" yaml:"bool,omitempty"`
	}

	t.Run("encode to yaml", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, yaml.NewEncoder(&buf).Encode(testStruct{}))
		assert.Equal(t, "{}", strings.TrimSpace(buf.String()))

		buf.Reset()
		require.NoError(t, yaml.NewEncoder(&buf).Encode(testStruct{Bool: FalseOption}))
		assert.Equal(t, `bool: "no"`, strings.TrimSpace(buf.String()))

		buf.Reset()
		require.NoError(t, yaml.NewEncoder(&buf).Encode(testStruct{Bool: TrueOption}))
		assert.Equal(t, `bool: "yes"`, strings.TrimSpace(buf.String()))
	})
	t.Run("decode from YAML", func(t *testing.T) {
		var testValue testStruct
		require.NoError(t, yaml.Unmarshal([]byte(`{}`), &testValue))
		assert.Equal(t, EmptyBoolOption, testValue.Bool)

		require.NoError(t, yaml.Unmarshal([]byte(`bool: no`), &testValue))
		assert.Equal(t, FalseOption, testValue.Bool)

		require.NoError(t, yaml.Unmarshal([]byte(`bool: yes`), &testValue))
		assert.Equal(t, TrueOption, testValue.Bool)

		require.NoError(t, yaml.Unmarshal([]byte(`bool: ""`), &testValue))
		assert.Equal(t, EmptyBoolOption, testValue.Bool)

		assert.Error(t, yaml.Unmarshal([]byte(`bool: invalid`), &testValue))
	})
	t.Run("get from envar", func(t *testing.T) {
		var testValue testStruct
		resetEnv(t)
		require.NoError(t, os.Setenv("BOOL", "no"))
		require.NoError(t, envdecode.Decode(&testValue))
		assert.Equal(t, FalseOption, testValue.Bool)

		testValue = testStruct{}
		resetEnv(t)
		require.NoError(t, os.Setenv("BOOL", "yes"))
		require.NoError(t, envdecode.Decode(&testValue))
		assert.Equal(t, TrueOption, testValue.Bool)
	})
}
