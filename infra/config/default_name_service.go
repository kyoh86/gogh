package config

import (
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/typ"
	"github.com/kyoh86/gogh/v3/infra/github"
)

const DefaultHost = github.GlobalHost

//TODO: move to core
//TODO: type DefaultNameService interface {
//TODO: type defaultNameServiceImpl struct {
//TODO: func NewDefaultNameService() DefaultNameService {

// DefaultNameService implements the repository.DefaultNameService interface
type DefaultNameService struct {
	hosts       typ.Map[string, string]
	defaultHost string
	changed     bool
}

// NewDefaultNameService creates a new DefaultNameService instance
func NewDefaultNameService() repository.DefaultNameService {
	return &DefaultNameService{
		hosts:       typ.Map[string, string]{},
		defaultHost: DefaultHost,
	}
}

// GetMap implements auth.DefaultsService
func (d DefaultNameService) GetMap() map[string]string {
	if d.hosts == nil {
		return nil
	}
	return d.hosts
}

// GetDefaultHost implements auth.DefaultsService
func (d DefaultNameService) GetDefaultHost() string {
	if d.defaultHost == "" {
		return DefaultHost
	}
	return d.defaultHost
}

// GetDefaultHostAndOwner implements auth.DefaultsService
func (d DefaultNameService) GetDefaultHostAndOwner() (host string, owner string) {
	hostName := d.GetDefaultHost()
	ownerName, _ := d.hosts.TryGet(hostName)
	return hostName, ownerName
}

// GetDefaultOwnerFor implements auth.DefaultsService
func (d DefaultNameService) GetDefaultOwnerFor(host string) (string, error) {
	owner, _ := d.hosts.TryGet(host)
	return owner, nil
}

// SetDefaultHost implements auth.DefaultsService
func (d *DefaultNameService) SetDefaultHost(host string) error {
	if err := repository.ValidateHost(host); err != nil {
		return err
	}
	d.defaultHost = host
	d.changed = true
	return nil
}

// SetDefaultOwnerFor implements auth.DefaultsService
func (d *DefaultNameService) SetDefaultOwnerFor(host, owner string) error {
	if err := repository.ValidateHost(host); err != nil {
		return err
	}
	if err := repository.ValidateOwner(owner); err != nil {
		return err
	}
	d.hosts.Set(host, owner)
	d.changed = true
	return nil
}

// HasChanges implements repository.DefaultNameService.
func (d *DefaultNameService) HasChanges() bool {
	return d.changed
}

// MarkSaved implements repository.DefaultNameService.
func (d *DefaultNameService) MarkSaved() {
	d.changed = false
}

var _ repository.DefaultNameService = (*DefaultNameService)(nil)
