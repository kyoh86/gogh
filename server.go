package gogh

import (
	"encoding/json"

	yaml "github.com/goccy/go-yaml"
)

const (
	DefaultHost = "github.com"
)

var (
	DefaultServer Server
)

func init() {
	s, err := NewServer(DefaultHost)
	if err != nil {
		panic(err)
	}
	DefaultServer = s
}

type Server struct {
	t taggedServer
}

func NewServer(host string) (Server, error) {
	if err := ValidateHost(host); err != nil {
		return Server{}, err
	}
	return Server{t: taggedServer{Host: host}}, nil
}

func (s Server) Host() string  { return s.t.Host }
func (s Server) User() string  { return s.t.User }
func (s Server) Token() string { return s.t.Token }

func (s *Server) SetHost(v string) error {
	if err := ValidateHost(v); err != nil {
		return err
	}
	s.t.Host = v
	return nil
}
func (s *Server) SetUser(v string) error {
	if err := ValidateUser(v); err != nil {
		return err
	}
	s.t.User = v
	return nil
}
func (s *Server) SetToken(v string) error { s.t.Token = v; return nil }

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
