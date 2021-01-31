package gogh

import (
	"errors"
	"sync"
)

var (
	ErrNoServer          = errors.New("no server registered")
	ErrServerNotFound    = errors.New("server not found")
	ErrUnremovableServer = errors.New("default server is not able to be removed")
)

type Servers struct {
	mutex         sync.Mutex
	defaultServer *Server
	serverMap     map[string]Server
}

func (s *Servers) init() {
	if s.serverMap == nil {
		s.serverMap = map[string]Server{}
	}
}

func (s *Servers) Set(host, user, token string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
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
	s.mutex.Lock()
	defer s.mutex.Unlock()

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
	s.mutex.Lock()
	defer s.mutex.Unlock()
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

func (s *Servers) Remove(host string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
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
