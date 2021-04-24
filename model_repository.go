package gogh

import "time"

type Repository struct {
	PushedAt    time.Time `json:"pushedAt"`
	Parent      *Spec     `json:"parent"`
	Spec        Spec      `json:"spec"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	Homepage    string    `json:"homepage"`
	Language    string    `json:"language"`
	Topics      []string  `json:"topics"`
	Archived    bool      `json:"archived"`
	Private     bool      `json:"private"`
	IsTemplate  bool      `json:"isTemplate"`
	Fork        bool      `json:"fork"`
}

func (r Repository) Host() string  { return r.Spec.Host() }
func (r Repository) Owner() string { return r.Spec.Owner() }
func (r Repository) Name() string  { return r.Spec.Name() }
