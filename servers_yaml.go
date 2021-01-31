package gogh

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

func (s *Servers) UnmarshalYAML(unmarshaler func(interface{}) error) error {
	var slice yaml.MapSlice
	if err := unmarshaler(&slice); err != nil {
		return err
	}

	var d *Server
	m := map[string]Server{}
	for _, item := range slice {
		host, ok := item.Key.(string)
		if !ok {
			return ErrInvalidHost(fmt.Sprintf("invalid host: %v", item.Key))
		}
		info, ok := item.Value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid value: %v", item.Value)
		}

		user, ok := info["user"].(string)
		if !ok {
			return fmt.Errorf("invalid user: %v", info["user"])
		}
		token, ok := info["token"].(string)
		if !ok {
			return fmt.Errorf("invalid token: %v", info["token"])
		}
		server, err := NewServerFor(host, user, token)
		if err != nil {
			return err
		}
		if d == nil {
			d = &server
		}
		m[host] = server
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.defaultServer = d
	s.serverMap = m

	return nil
}

func (s *Servers) MarshalYAML() (interface{}, error) {
	if s.defaultServer == nil {
		return nil, nil
	}

	var slice yaml.MapSlice
	for _, server := range s.serverMap {
		slice = append(slice, yaml.MapItem{
			Key: server.Host(),
			Value: yaml.MapSlice{
				{Key: "user", Value: server.User()},
				{Key: "token", Value: server.Token()},
			},
		})
	}
	return slice, nil
}

var _ yaml.InterfaceUnmarshaler = (*Servers)(nil)
var _ yaml.InterfaceMarshaler = (*Servers)(nil)
