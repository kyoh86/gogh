package migrate

import (
	"fmt"

	"github.com/kyoh86/gogh/v2"
)

// Server represents a GitHub server with user and token.
// Deprecated: use Token and TokenManager instead.
type Server struct {
	host  string
	user  string
	token string
}

func (s Server) String() string {
	return fmt.Sprintf("%s@%s", s.user, s.host)
}

func NewServerFor(host, user, token string) (Server, error) {
	if err := gogh.ValidateHost(host); err != nil {
		return Server{}, err
	}
	if err := gogh.ValidateOwner(user); err != nil {
		return Server{}, err
	}
	return Server{host: host, user: user, token: token}, nil
}

func NewServer(user, token string) (Server, error) {
	return NewServerFor(gogh.DefaultHost, user, token)
}

func (s Server) Host() string  { return s.host }
func (s Server) User() string  { return s.user }
func (s Server) Token() string { return s.token }
