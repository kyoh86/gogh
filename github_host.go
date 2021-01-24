package gogh

import (
	"encoding/json"

	yaml "github.com/goccy/go-yaml"
)

const (
	DefaultGithubHostName = "github.com"
)

var (
	DefaultGithubHost GithubHost
)

func init() {
	h, err := NewGithubHost(DefaultGithubHostName)
	if err != nil {
		panic(err)
	}
	DefaultGithubHost = h
}

type GithubHost struct {
	t taggedGithubHost
}

func NewGithubHost(host string) (GithubHost, error) {
	if err := ValidateHost(host); err != nil {
		return GithubHost{}, err
	}
	return GithubHost{t: taggedGithubHost{Name: host}}, nil
}

func (h GithubHost) Name() string  { return h.t.Name }
func (h GithubHost) User() string  { return h.t.User }
func (h GithubHost) Token() string { return h.t.Token }

func (h *GithubHost) SetName(v string) error {
	if err := ValidateHost(v); err != nil {
		return err
	}
	h.t.Name = v
	return nil
}
func (h *GithubHost) SetUser(v string) error  { h.t.User = v; return nil }
func (h *GithubHost) SetToken(v string) error { h.t.Token = v; return nil }

type taggedGithubHost struct {
	Name  string `json:"name" yaml:"name"`
	User  string `json:"user,omitempty" yaml:"user,omitempty"`
	Token string `json:"token,omitempty" yaml:"token,omitempty"`
}

func (h *GithubHost) UnmarshalJSON(b []byte) error {
	var t taggedGithubHost
	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}
	if err := ValidateHost(t.Name); err != nil {
		return err
	}
	h.t = t
	return nil
}

func (h GithubHost) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.t)
}

func (h *GithubHost) UnmarshalYAML(unmarshaler func(interface{}) error) error {
	var t taggedGithubHost
	if err := unmarshaler(&t); err != nil {
		return err
	}
	if err := ValidateHost(t.Name); err != nil {
		return err
	}
	h.t = t
	return nil
}

func (h GithubHost) MarshalYAML() (interface{}, error) {
	return h.t, nil
}

var _ json.Unmarshaler = (*GithubHost)(nil)
var _ yaml.InterfaceUnmarshaler = (*GithubHost)(nil)

var _ json.Marshaler = (*GithubHost)(nil)
var _ yaml.InterfaceMarshaler = (*GithubHost)(nil)
