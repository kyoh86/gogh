package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessor(t *testing.T) {
	t.Run("getting", func(t *testing.T) {
		mustOption := func(acc *OptionAccessor, err error) *OptionAccessor {
			t.Helper()
			require.NoError(t, err)
			return acc
		}
		var cfg Config
		cfg.GitHub.Token = "token1"
		cfg.GitHub.Host = "hostx1"
		cfg.GitHub.User = "kyoh86"
		cfg.Log.Level = "trace"
		cfg.Log.Date = TrueOption
		cfg.Log.Time = FalseOption
		cfg.Log.LongFile = TrueOption
		cfg.Log.ShortFile = TrueOption
		cfg.Log.UTC = TrueOption
		cfg.VRoot = []string{"/foo", "/bar"}

		_, err := Option("invalid name")
		assert.EqualError(t, err, "invalid option name")
		assert.Equal(t, "token1", mustOption(Option("github.token")).Get(&cfg))
		assert.Equal(t, "hostx1", mustOption(Option("github.host")).Get(&cfg))
		assert.Equal(t, "kyoh86", mustOption(Option("github.user")).Get(&cfg))
		assert.Equal(t, "trace", mustOption(Option("log.level")).Get(&cfg))
		assert.Equal(t, "yes", mustOption(Option("log.date")).Get(&cfg))
		assert.Equal(t, "no", mustOption(Option("log.time")).Get(&cfg))
		assert.Equal(t, "yes", mustOption(Option("log.longfile")).Get(&cfg))
		assert.Equal(t, "yes", mustOption(Option("log.shortfile")).Get(&cfg))
		assert.Equal(t, "yes", mustOption(Option("log.utc")).Get(&cfg))
		assert.Equal(t, "/foo:/bar", mustOption(Option("root")).Get(&cfg))
	})
	t.Run("putting", func(t *testing.T) {
		mustOption := func(acc *OptionAccessor, err error) *OptionAccessor {
			t.Helper()
			require.NoError(t, err)
			return acc
		}
		var cfg Config
		assert.NoError(t, mustOption(Option("github.host")).Put(&cfg, "hostx1"))
		assert.NoError(t, mustOption(Option("github.user")).Put(&cfg, "kyoh86"))
		assert.NoError(t, mustOption(Option("log.level")).Put(&cfg, "trace"))
		assert.NoError(t, mustOption(Option("log.date")).Put(&cfg, "yes"))
		assert.NoError(t, mustOption(Option("log.time")).Put(&cfg, "no"))
		assert.NoError(t, mustOption(Option("log.longfile")).Put(&cfg, "yes"))
		assert.NoError(t, mustOption(Option("log.shortfile")).Put(&cfg, "yes"))
		assert.NoError(t, mustOption(Option("log.utc")).Put(&cfg, "yes"))
		assert.NoError(t, mustOption(Option("root")).Put(&cfg, "/foo:/bar"))

		assert.Equal(t, "", cfg.GitHub.Token)
		assert.Equal(t, "hostx1", cfg.GitHub.Host)
		assert.Equal(t, "kyoh86", cfg.GitHub.User)
		assert.Equal(t, "trace", cfg.Log.Level)
		assert.Equal(t, TrueOption, cfg.Log.Date)
		assert.True(t, cfg.LogDate())
		assert.Equal(t, FalseOption, cfg.Log.Time)
		assert.False(t, cfg.LogTime())
		assert.Equal(t, TrueOption, cfg.Log.LongFile)
		assert.True(t, cfg.LogLongFile())
		assert.Equal(t, TrueOption, cfg.Log.ShortFile)
		assert.True(t, cfg.LogShortFile())
		assert.Equal(t, TrueOption, cfg.Log.UTC)
		assert.True(t, cfg.LogUTC())
		assert.Equal(t, PathListOption{"/foo", "/bar"}, cfg.VRoot)
	})
	t.Run("putting error", func(t *testing.T) {
		mustOption := func(acc *OptionAccessor, err error) *OptionAccessor {
			t.Helper()
			require.NoError(t, err)
			return acc
		}
		var cfg Config
		assert.EqualError(t, mustOption(Option("github.token")).Put(&cfg, "token1"), "token must not save")

		assert.EqualError(t, mustOption(Option("github.host")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("github.user")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("log.level")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("log.date")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("log.time")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("log.longfile")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("log.shortfile")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("log.utc")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("root")).Put(&cfg, ""), "empty value")

		assert.Error(t, mustOption(Option("github.user")).Put(&cfg, "-kyoh86"), "invalid github username")
		assert.Error(t, mustOption(Option("log.level")).Put(&cfg, "foobar"), "invalid log level")
		assert.Error(t, mustOption(Option("log.date")).Put(&cfg, "invalid value"), "invalid value")
		assert.Error(t, mustOption(Option("log.time")).Put(&cfg, "invalid value"), "invalid value")
		assert.Error(t, mustOption(Option("log.longfile")).Put(&cfg, "invalid value"), "invalid value")
		assert.Error(t, mustOption(Option("log.shortfile")).Put(&cfg, "invalid value"), "invalid value")
		assert.Error(t, mustOption(Option("log.utc")).Put(&cfg, "invalid value"), "invalid value")
		assert.Error(t, mustOption(Option("root")).Put(&cfg, "\x00"), "invalid value")

		assert.Equal(t, "", cfg.GitHub.Token)
		assert.Equal(t, "", cfg.GitHub.Host)
		assert.Equal(t, "", cfg.GitHub.User)
		assert.Equal(t, "", cfg.Log.Level)
		assert.Equal(t, EmptyBoolOption, cfg.Log.Date)
		assert.Equal(t, EmptyBoolOption, cfg.Log.Time)
		assert.Equal(t, EmptyBoolOption, cfg.Log.LongFile)
		assert.Equal(t, EmptyBoolOption, cfg.Log.ShortFile)
		assert.Equal(t, EmptyBoolOption, cfg.Log.UTC)
		assert.Empty(t, cfg.VRoot)

	})
	t.Run("unsetting", func(t *testing.T) {
		var cfg Config
		cfg.GitHub.Token = "token1"
		cfg.GitHub.Host = "hostx1"
		cfg.GitHub.User = "kyoh86"
		cfg.Log.Level = "trace"
		cfg.Log.Date = TrueOption
		cfg.Log.Time = FalseOption
		cfg.Log.LongFile = TrueOption
		cfg.Log.ShortFile = TrueOption
		cfg.Log.UTC = TrueOption
		cfg.VRoot = []string{"/foo", "/bar"}

		_, err := Option("invalid name")
		assert.EqualError(t, err, "invalid option name")
		for _, name := range OptionNames() {
			acc, err := Option(name)
			require.NoError(t, err)
			assert.NoError(t, acc.Unset(&cfg), name)
		}
		assert.Equal(t, "", cfg.GitHub.Token)
		assert.Equal(t, "", cfg.GitHub.Host)
		assert.Equal(t, "", cfg.GitHub.User)
		assert.Equal(t, "", cfg.Log.Level)
		assert.Equal(t, EmptyBoolOption, cfg.Log.Date)
		assert.Equal(t, EmptyBoolOption, cfg.Log.Time)
		assert.Equal(t, EmptyBoolOption, cfg.Log.LongFile)
		assert.Equal(t, EmptyBoolOption, cfg.Log.ShortFile)
		assert.Equal(t, EmptyBoolOption, cfg.Log.UTC)
		assert.Empty(t, cfg.VRoot)
	})
}
