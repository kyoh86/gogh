package remote

import (
	"time"

	"github.com/kyoh86/gogh/v3/domain/reporef"
)

type Repo struct {
	UpdatedAt   time.Time        `json:"updatedAt"`
	Parent      *reporef.RepoRef `json:"parent,omitempty"`
	Ref         reporef.RepoRef  `json:"ref"`
	URL         string           `json:"url"`
	Description string           `json:"description,omitempty"`
	Homepage    string           `json:"homepage,omitempty"`
	Language    string           `json:"language,omitempty"`
	Archived    bool             `json:"archived,omitempty"`
	Private     bool             `json:"private,omitempty"`
	IsTemplate  bool             `json:"isTemplate,omitempty"`
	Fork        bool             `json:"fork,omitempty"`
}

func (r Repo) Host() string  { return r.Ref.Host() }
func (r Repo) Owner() string { return r.Ref.Owner() }
func (r Repo) Name() string  { return r.Ref.Name() }
