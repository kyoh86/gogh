package config

import (
	"context"
	"io"
	"log"
	"os"
)

// Config holds configuration file values.
type Config struct {
	context.Context
	Log    LogConfig         `yaml:"log"`
	VRoot  StringArrayConfig `yaml:"root,omitempty" env:"GOGH_ROOT"`
	GitHub GitHubConfig      `yaml:"github,omitempty"`
}

type LogConfig struct {
	Level        string     `yaml:"level,omitempty" env:"GOGH_LOG_LEVEL"`
	Date         BoolConfig `yaml:"date,omitempty" env:"GOGH_LOG_DATE"`                 // the date in the local time zone: 2009/01/23
	Time         BoolConfig `yaml:"time,omitempty" env:"GOGH_LOG_TIME"`                 // the time in the local time zone: 01:23:23
	MicroSeconds BoolConfig `yaml:"microseconds,omitempty" env:"GOGH_LOG_MICROSECONDS"` // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	LongFile     BoolConfig `yaml:"longfile,omitempty" env:"GOGH_LOG_LONGFILE"`         // full file name and line number: /a/b/c/d.go:23
	ShortFile    BoolConfig `yaml:"shortfile,omitempty" env:"GOGH_LOG_SHORTFILE"`       // final file name element and line number: d.go:23. overrides Llongfile
	UTC          BoolConfig `yaml:"utc,omitempty" env:"GOGH_LOG_UTC"`                   // if Ldate or Ltime is set, use UTC rather than the local time zone
}

type GitHubConfig struct {
	Token string `yaml:"token,omitempty" env:"GOGH_GITHUB_TOKEN"`
	User  string `yaml:"user,omitempty" env:"GOGH_GITHUB_USER"`
	Host  string `yaml:"host,omitempty" env:"GOGH_GITHUB_HOST"`
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

func (c *Config) Root() []string {
	return c.VRoot
}

func (c *Config) PrimaryRoot() string {
	return c.VRoot[0]
}
