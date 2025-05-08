package hosting

import (
	"time"

	"github.com/kyoh86/gogh/v3/core/repository"
)

// Repository represents a repository on a remote source
type Repository struct {
	// Ref is a reference of the repository
	Ref repository.Reference `json:"ref"`
	// URL is a full URL for the repository (i.g.: "https://github.com/kyoh86/gogh")
	URL string `json:"url"`

	// CloneURL is a clone URL for the repository (i.g.: "
	CloneURL string `json:"cloneUrl,omitempty"`

	// UpdatedAt is the last updated time of the repository
	UpdatedAt time.Time `json:"updatedAt"`
	// Parent is the parent repository if it is a fork
	Parent *ParentRepository `json:"parent,omitempty"`

	// Description is a description of the repository (i.g.: "Gogh is a collection of themes for Gnome Terminal and Pantheon Terminal")
	Description string `json:"description,omitempty"`
	// Homepage is a homepage of the repository (i.g.: "https://example.com")
	Homepage string `json:"homepage,omitempty"`
	// Language is a primary language of the repository (i.g.: "Go")
	Language string `json:"language,omitempty"`
	// Archived is if the repository is archived
	Archived bool `json:"language,omitempty"`
	// Private is if the repository is private
	Private bool `json:"private,omitempty"`
	// IsTemplate is if the repository is a template
	IsTemplate bool `json:"isTemplate,omitempty"`
	// Fork is if the repository is a fork
	Fork bool `json:"fork,omitempty"`
}

// ParentRepository represents a parent repository of a fork
type ParentRepository struct {
	// Ref is a reference of the parent repository
	Ref repository.Reference
	// CloneURL is a clone URL for the parent repository
	CloneURL string
}
