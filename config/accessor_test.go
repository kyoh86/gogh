package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessor(t *testing.T) {
	const (
		dummyToken = "token1"
		dummyHost  = "hostx1"
		dummyUser  = "kyoh86"
		dummyLevel = "trace"
	)
	t.Run("getting", func(t *testing.T) {
		mustOption := func(acc *OptionAccessor, err error) *OptionAccessor {
			t.Helper()
			require.NoError(t, err)
			return acc
		}
		var cfg Config
		cfg.GitHub.Token = dummyToken
		cfg.GitHub.Host = dummyHost
		cfg.GitHub.User = dummyUser
		cfg.VRoot = []string{"/foo", "/bar"}

		_, err := Option("invalid name")
		assert.EqualError(t, err, "invalid option name")
		assert.Equal(t, "*****", mustOption(Option("github.token")).Get(&cfg))
		assert.Equal(t, dummyHost, mustOption(Option("github.host")).Get(&cfg))
		assert.Equal(t, dummyUser, mustOption(Option("github.user")).Get(&cfg))
		assert.Equal(t, "/foo:/bar", mustOption(Option("root")).Get(&cfg))
	})
	t.Run("putting", func(t *testing.T) {
		mustOption := func(acc *OptionAccessor, err error) *OptionAccessor {
			t.Helper()
			require.NoError(t, err)
			return acc
		}
		var cfg Config
		assert.NoError(t, mustOption(Option("github.host")).Put(&cfg, dummyHost))
		assert.NoError(t, mustOption(Option("github.user")).Put(&cfg, dummyUser))
		assert.NoError(t, mustOption(Option("root")).Put(&cfg, "/foo:/bar"))

		assert.Equal(t, "", cfg.GitHub.Token)
		assert.Equal(t, dummyHost, cfg.GitHub.Host)
		assert.Equal(t, dummyUser, cfg.GitHub.User)
		assert.Equal(t, PathListOption{"/foo", "/bar"}, cfg.VRoot)
	})
	t.Run("putting error", func(t *testing.T) {
		mustOption := func(acc *OptionAccessor, err error) *OptionAccessor {
			t.Helper()
			require.NoError(t, err)
			return acc
		}
		var cfg Config
		assert.EqualError(t, mustOption(Option("github.token")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("github.host")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("github.user")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("root")).Put(&cfg, ""), "empty value")

		assert.Error(t, mustOption(Option("github.user")).Put(&cfg, "-kyoh86"), "invalid github username")
		assert.Error(t, mustOption(Option("root")).Put(&cfg, "\x00"), "invalid value")

		assert.Equal(t, "", cfg.GitHub.Token)
		assert.Equal(t, "", cfg.GitHub.Host)
		assert.Equal(t, "", cfg.GitHub.User)
		assert.Empty(t, cfg.VRoot)
	})

	t.Run("unsetting", func(t *testing.T) {
		var cfg Config
		cfg.GitHub.Token = dummyToken
		cfg.GitHub.Host = dummyHost
		cfg.GitHub.User = dummyUser
		cfg.VRoot = []string{"/foo", "/bar"}

		_, err := Option("invalid name")
		assert.EqualError(t, err, "invalid option name")
		for _, name := range OptionNames() {
			if name == gitHubTokenOptionAccessor.optionName {
				continue
			}
			acc, err := Option(name)
			require.NoError(t, err)
			assert.NoError(t, acc.Unset(&cfg), name)
		}
		assert.Equal(t, "", cfg.GitHub.Host)
		assert.Equal(t, "", cfg.GitHub.User)
		assert.Empty(t, cfg.VRoot)
	})
}
