package gogh

import (
	"fmt"

	"github.com/kyoh86/gogh/v2/internal/github"
)

const (
	DefaultHost = github.DefaultHost
)

type Server struct {
	host  string
	user  string
	token string
}

func (s Server) String() string {
	return fmt.Sprintf("%s@%s", s.user, s.host)
}

func NewServerFor(host, user, token string) (Server, error) {
	if err := ValidateHost(host); err != nil {
		return Server{}, err
	}
	if err := ValidateOwner(user); err != nil {
		return Server{}, err
	}
	return Server{host: host, user: user, token: token}, nil
}

func NewServer(user, token string) (Server, error) {
	return NewServerFor(DefaultHost, user, token)
}

func (s Server) Host() string  { return s.host }
func (s Server) User() string  { return s.user }
func (s Server) Token() string { return s.token }
