package gogh

import "time"

type Repository struct {
	spec        Spec
	description string
	homepage    string
	language    string
	topics      []string
	pushedAt    time.Time
	parent      *Spec
	archived    bool
	private     bool
	isTemplate  bool
	fork        bool
}

func (r Repository) Spec() Spec    { return r.spec }
func (r Repository) Host() string  { return r.spec.Host() }
func (r Repository) Owner() string { return r.spec.Owner() }
func (r Repository) Name() string  { return r.spec.Name() }

func (r Repository) Description() string {
	return r.description
}

func (r Repository) Homepage() string {
	return r.homepage
}
func (r Repository) Language() string    { return r.language }
func (r Repository) Topics() []string    { return r.topics }
func (r Repository) PushedAt() time.Time { return r.pushedAt }
func (r Repository) Parent() *Spec       { return r.parent }
func (r Repository) IsArchived() bool    { return r.archived }
func (r Repository) IsPrivate() bool     { return r.private }
func (r Repository) IsTemplate() bool    { return r.isTemplate }
func (r Repository) IsFork() bool        { return r.fork }
