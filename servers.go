package gogh

import (
	"errors"
)

var (
	ErrNoServer          = errors.New("no server registered")
	ErrServerNotFound    = errors.New("server not found")
	ErrUnremovableServer = errors.New("default server is not able to be removed")
)

type Servers struct {
	defaultServer *Server
	serverMap     map[string]Server
}

// NewServers will build Spec with a default server and alternative servers.
func NewServers(defaultServer Server, alternatives ...Server) Servers {
	h := defaultServer.Host()
	serverMap := map[string]Server{
		h: defaultServer,
	}
	for _, s := range alternatives {
		h := s.Host()
		if _, ok := serverMap[h]; ok {
			continue
		}
		serverMap[h] = s
	}
	return Servers{
		serverMap:     serverMap,
		defaultServer: &defaultServer,
	}
}
func (s *Servers) init() {
	if s.serverMap == nil {
		s.serverMap = map[string]Server{}
	}
}

func (s *Servers) Set(host, user, token string) error {
	s.init()

	server, err := NewServerFor(host, user, token)
	if err != nil {
		return err
	}
	if s.defaultServer == nil || s.defaultServer.Host() == host {
		s.defaultServer = &server
	}
	s.serverMap[host] = server
	return nil
}

func (s *Servers) Find(host string) (Server, error) {
	if len(s.serverMap) == 0 {
		return Server{}, ErrNoServer
	}
	server, ok := s.serverMap[host]
	if !ok {
		return Server{}, ErrServerNotFound
	}
	return server, nil
}

func (s *Servers) SetDefault(host string) error {
	s.init()

	server, ok := s.serverMap[host]
	if !ok {
		return ErrServerNotFound
	}
	s.defaultServer = &server
	return nil
}

func (s *Servers) Default() (Server, error) {
	if server := s.defaultServer; server != nil {
		return *server, nil
	}
	return Server{}, ErrNoServer
}

func (s *Servers) List() (list []Server, _ error) {
	d := s.defaultServer
	if d == nil {
		return nil, nil
	}
	list = append(list, *d)

	for _, server := range s.serverMap {
		if server.Host() == d.Host() {
			continue
		}
		list = append(list, server)
	}
	return list, nil
}

func (s *Servers) Remove(host string) error {
	s.init()

	if _, ok := s.serverMap[host]; !ok {
		return ErrServerNotFound
	}
	if len(s.serverMap) == 1 {
		s.serverMap = map[string]Server{}
		s.defaultServer = nil
	}
	if s.defaultServer != nil && s.defaultServer.Host() == host {
		return ErrUnremovableServer
	}
	delete(s.serverMap, host)
	return nil
}
