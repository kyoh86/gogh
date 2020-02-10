package config

import (
	"context"
	"io"
	"log"
	"os"
	"time"
)

// Config holds configuration file values.
type Config struct {
	Log    LogConfig      `yaml:"log,omitempty"`
	VRoot  PathListOption `yaml:"root,omitempty" env:"GOGH_ROOT"`
	GitHub GitHubConfig   `yaml:"github,omitempty"`
}

var _ context.Context = (*Config)(nil)

type LogConfig struct {
	Level        string     `yaml:"level,omitempty" env:"GOGH_LOG_LEVEL"`
	Date         BoolOption `yaml:"date" env:"GOGH_LOG_DATE"`                 // the date in the local time zone: 2009/01/23
	Time         BoolOption `yaml:"time" env:"GOGH_LOG_TIME"`                 // the time in the local time zone: 01:23:23
	MicroSeconds BoolOption `yaml:"microseconds" env:"GOGH_LOG_MICROSECONDS"` // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	LongFile     BoolOption `yaml:"longfile" env:"GOGH_LOG_LONGFILE"`         // full file name and line number: /a/b/c/d.go:23
	ShortFile    BoolOption `yaml:"shortfile" env:"GOGH_LOG_SHORTFILE"`       // final file name element and line number: d.go:23. overrides Llongfile
	UTC          BoolOption `yaml:"utc" env:"GOGH_LOG_UTC"`                   // if Ldate or Ltime is set, use UTC rather than the local time zone
}

type GitHubConfig struct {
	Token string `yaml:"-" env:"GOGH_GITHUB_TOKEN"`
	User  string `yaml:"user,omitempty" env:"GOGH_GITHUB_USER"`
	Host  string `yaml:"host,omitempty" env:"GOGH_GITHUB_HOST"`
}

// Deadline : empty context.Context
func (*Config) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done : empty context.Context
func (*Config) Done() <-chan struct{} {
	return nil
}

// Err : empty context.Context
func (*Config) Err() error {
	return nil
}

// Value : empty context.Context
func (*Config) Value(key interface{}) interface{} {
	return nil
}

func (*Config) String() string {
	return "gogh/config.Config"
}

func (c *Config) Stdin() io.Reader {
	return os.Stdin
}

func (c *Config) Stdout() io.Writer {
	return os.Stdout
}

func (c *Config) Stderr() io.Writer {
	return os.Stderr
}

func (c *Config) GitHubUser() string {
	return c.GitHub.User
}

func (c *Config) GitHubToken() string {
	return c.GitHub.Token
}

func (c *Config) GitHubHost() string {
	return c.GitHub.Host
}

func (c *Config) LogLevel() string {
	return c.Log.Level
}

func (c *Config) LogFlags() int {
	var f int
	if c.Log.Date.Bool() {
		f |= log.Ldate
	}
	if c.Log.Time.Bool() {
		f |= log.Ltime
	}
	if c.Log.MicroSeconds.Bool() {
		f |= log.Lmicroseconds
	}
	if c.Log.LongFile.Bool() {
		f |= log.Llongfile
	}
	if c.Log.ShortFile.Bool() {
		f |= log.Lshortfile
	}
	if c.Log.UTC.Bool() {
		f |= log.LUTC
	}
	return f
}

func (c *Config) LogDate() bool         { return c.Log.Date.Bool() }
func (c *Config) LogTime() bool         { return c.Log.Time.Bool() }
func (c *Config) LogMicroSeconds() bool { return c.Log.MicroSeconds.Bool() }
func (c *Config) LogLongFile() bool     { return c.Log.LongFile.Bool() }
func (c *Config) LogShortFile() bool    { return c.Log.ShortFile.Bool() }
func (c *Config) LogUTC() bool          { return c.Log.UTC.Bool() }

func (c *Config) Root() []string {
	ret := make([]string, 0, len(c.VRoot))
	for _, p := range c.VRoot {
		ret = append(ret, expandPath(p))
	}
	return ret
}

func (c *Config) PrimaryRoot() string {
	return expandPath(c.VRoot[0])
}
