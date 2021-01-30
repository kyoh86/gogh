package gogh

import (
	"encoding/json"

	yaml "github.com/goccy/go-yaml"
	"github.com/kyoh86/gogh/v2/internal/github"
)

const (
	DefaultHost = github.DefaultHost
)

type Server struct {
	t taggedServer
}

func NewServerFor(host, user, token string) (Server, error) {
	if err := ValidateHost(host); err != nil {
		return Server{}, err
	}
	if err := ValidateUser(user); err != nil {
		return Server{}, err
	}
	return Server{t: taggedServer{Host: host, User: user, Token: token}}, nil
}

func NewServer(user, token string) (Server, error) {
	return NewServerFor(DefaultHost, user, token)
}

func (s Server) Host() string { return s.t.Host }
func (s Server) User() string { return s.t.User }

func (s Server) Token() string { return s.t.Token }

type taggedServer struct {
	Host  string `json:"host" yaml:"host"`
	User  string `json:"user,omitempty" yaml:"user,omitempty"`
	Token string `json:"token,omitempty" yaml:"token,omitempty"`
}

func (s *Server) UnmarshalJSON(b []byte) error {
	return s.UnmarshalYAML(func(obj interface{}) error {
		return json.Unmarshal(b, obj)
	})
}

func (s Server) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.t)
}

func (s *Server) UnmarshalYAML(unmarshaler func(interface{}) error) error {
	var t taggedServer
	if err := unmarshaler(&t); err != nil {
		return err
	}
	if err := ValidateHost(t.Host); err != nil {
		return err
	}
	if err := ValidateUser(t.User); err != nil {
		return err
	}
	s.t = t
	return nil
}

func (s Server) MarshalYAML() (interface{}, error) {
	return s.t, nil
}

var _ json.Unmarshaler = (*Server)(nil)
var _ yaml.InterfaceUnmarshaler = (*Server)(nil)

var _ json.Marshaler = (*Server)(nil)
var _ yaml.InterfaceMarshaler = (*Server)(nil)
