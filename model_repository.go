package gogh

import "time"

type Repository struct {
	PushedAt    *time.Time `json:"pushedAt,omitempty"`
	Parent      *Spec      `json:"parent,omitempty"`
	Spec        Spec       `json:"spec"`
	URL         string     `json:"url"`
	Description string     `json:"description,omitempty"`
	Homepage    string     `json:"homepage,omitempty"`
	Language    string     `json:"language,omitempty"`
	Archived    bool       `json:"archived,omitempty"`
	Private     bool       `json:"private,omitempty"`
	IsTemplate  bool       `json:"isTemplate,omitempty"`
	Fork        bool       `json:"fork,omitempty"`
}

func (r Repository) Host() string  { return r.Spec.Host() }
func (r Repository) Owner() string { return r.Spec.Owner() }
func (r Repository) Name() string  { return r.Spec.Name() }
