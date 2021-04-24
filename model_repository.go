package gogh

import "time"

type Repository struct {
	PushedAt    time.Time
	Parent      *Spec
	spec        Spec
	URL         string
	Description string
	Homepage    string
	Language    string
	Topics      []string
	Archived    bool
	Private     bool
	IsTemplate  bool
	Fork        bool
}

func (r Repository) Spec() Spec    { return r.spec }
func (r Repository) Host() string  { return r.spec.Host() }
func (r Repository) Owner() string { return r.spec.Owner() }
func (r Repository) Name() string  { return r.spec.Name() }
