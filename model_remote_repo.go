package gogh

import "time"

type RemoteRepo struct {
	UpdatedAt   time.Time `json:"updatedAt"`
	Parent      *RepoRef  `json:"parent,omitempty"`
	Ref         RepoRef   `json:"ref"`
	URL         string    `json:"url"`
	Description string    `json:"description,omitempty"`
	Homepage    string    `json:"homepage,omitempty"`
	Language    string    `json:"language,omitempty"`
	Archived    bool      `json:"archived,omitempty"`
	Private     bool      `json:"private,omitempty"`
	IsTemplate  bool      `json:"isTemplate,omitempty"`
	Fork        bool      `json:"fork,omitempty"`
}

func (r RemoteRepo) Host() string  { return r.Ref.Host() }
func (r RemoteRepo) Owner() string { return r.Ref.Owner() }
func (r RemoteRepo) Name() string  { return r.Ref.Name() }
