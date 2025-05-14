package repository

import (
	"github.com/kyoh86/gogh/v3/core/gogh"
	"github.com/kyoh86/gogh/v3/core/store"
	"github.com/kyoh86/gogh/v3/core/typ"
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

	store.Content
}

// defaultNameServiceImpl implements the DefaultNameService interface
type defaultNameServiceImpl struct {
	hosts       typ.Map[string, string]
	defaultHost string
	changed     bool
}

// NewDefaultNameService creates a new DefaultNameService instance
func NewDefaultNameService() DefaultNameService {
	return &defaultNameServiceImpl{
		hosts:       typ.Map[string, string]{},
		defaultHost: gogh.DefaultHost,
	}
}

// GetMap implements DefaultNameService
func (d defaultNameServiceImpl) GetMap() map[string]string {
	if d.hosts == nil {
		return nil
	}
	return d.hosts
}

// GetDefaultHost implements DefaultNameService
func (d defaultNameServiceImpl) GetDefaultHost() string {
	if d.defaultHost == "" {
		return gogh.DefaultHost
	}
	return d.defaultHost
}

// GetDefaultHostAndOwner implements DefaultNameService
func (d defaultNameServiceImpl) GetDefaultHostAndOwner() (host string, owner string) {
	hostName := d.GetDefaultHost()
	ownerName, _ := d.hosts.TryGet(hostName)
	return hostName, ownerName
}

// GetDefaultOwnerFor implements DefaultNameService
func (d defaultNameServiceImpl) GetDefaultOwnerFor(host string) (string, error) {
	owner, _ := d.hosts.TryGet(host)
	return owner, nil
}

// SetDefaultHost implements DefaultNameService
func (d *defaultNameServiceImpl) SetDefaultHost(host string) error {
	if err := ValidateHost(host); err != nil {
		return err
	}
	d.defaultHost = host
	d.changed = true
	return nil
}

// SetDefaultOwnerFor implements DefaultNameService
func (d *defaultNameServiceImpl) SetDefaultOwnerFor(host, owner string) error {
	if err := ValidateHost(host); err != nil {
		return err
	}
	if err := ValidateOwner(owner); err != nil {
		return err
	}
	d.hosts.Set(host, owner)
	d.changed = true
	return nil
}

// HasChanges implements DefaultNameService.
func (d *defaultNameServiceImpl) HasChanges() bool {
	return d.changed
}

// MarkSaved implements DefaultNameService.
func (d *defaultNameServiceImpl) MarkSaved() {
	d.changed = false
}

var _ DefaultNameService = (*defaultNameServiceImpl)(nil)
