package repository

import (
	"github.com/kyoh86/gogh/v3/core/store"
)

// DefaultNameService provides access to default configuration settings
type DefaultNameService interface {
	// GetDefaultHost returns the current default host
	GetDefaultHost() string

	// GetMap returns a map of hosts and their corresponding owners
	GetMap() map[string]string

	// GetDefaultHostAndOwner returns the current default host and owner
	GetDefaultHostAndOwner() (host string, owner string)

	// GetDefaultOwnerFor returns the default owner for the specified host
	GetDefaultOwnerFor(host string) (string, error)

	// SetDefaultHost sets the default host
	SetDefaultHost(host string) error

	// SetDefaultOwnerFor sets the default owner for the specified host
	SetDefaultOwnerFor(host, owner string) error
}

// DefaultNameStore is a service for saving and loading tokens
type DefaultNameStore store.Store[DefaultNameService]
