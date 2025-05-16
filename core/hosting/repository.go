package hosting

import (
	"time"

	"github.com/kyoh86/gogh/v4/core/repository"
)

// Repository represents a repository on a remote source
type Repository struct {
	// Ref is a reference of the repository
	Ref repository.Reference
	// URL is a full URL for the repository (e.g.: "https://github.com/kyoh86/gogh")
	URL string

	// CloneURL is a clone URL for the repository (e.g.: "https://github.com/kyoh86/gogh.git")
	CloneURL string

	// UpdatedAt is the last updated time of the repository
	UpdatedAt time.Time
	// Parent is the parent repository if it is a fork
	Parent *ParentRepository

	// Description is a description of the repository (e.g.: "Gogh is a collection of themes for Gnome Terminal and Pantheon Terminal")
	Description string
	// Homepage is a homepage of the repository (e.g.: "https://example.com")
	Homepage string
	// Language is a primary language of the repository (e.g.: "Go")
	Language string
	// Archived is if the repository is archived
	Archived bool
	// Private is if the repository is private
	Private bool
	// IsTemplate is if the repository is a template
	IsTemplate bool
	// Fork is if the repository is a fork
	Fork bool
}

// ParentRepository represents a parent repository of a fork
type ParentRepository struct {
	// Ref is a reference of the parent repository
	Ref repository.Reference
	// CloneURL is a clone URL for the parent repository
	CloneURL string
}
