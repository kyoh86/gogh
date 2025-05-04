package hosting

import (
	"context"
	"iter"

	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/util"
)

// HostingService provides access to remote repositories
type HostingService interface {
	// GetTokenFor retrieves an authentication token and user for a specific repository reference
	GetTokenFor(ctx context.Context, reference repository.Reference) (string, auth.Token, error)
	// GetRepository retrieves repository information from a remote source
	GetRepository(context.Context, repository.Reference) (*Repository, error)
	// ListRepository retrieves a list of repositories from a remote source
	ListRepository(context.Context, *ListRepositoryOptions) iter.Seq2[*Repository, error]
	// DeleteRepository deletes a repository from a remote source
	DeleteRepository(context.Context, repository.Reference) error
}

// BooleanFilter represents a filter state for boolean repository attributes
type BooleanFilter int

const (
	// BooleanFilterNone indicates no filtering should be applied
	BooleanFilterNone BooleanFilter = iota
	// BooleanFilterTrue filters for repositories where the attribute is true
	BooleanFilterTrue
	// BooleanFilterFalse filters for repositories where the attribute is false
	BooleanFilterFalse
)

// AsBoolPtr converts the BooleanFilter to a pointer to a boolean value
func (f BooleanFilter) AsBoolPtr() *bool {
	switch f {
	case BooleanFilterNone:
		return nil
	case BooleanFilterTrue:
		return util.Ptr(true)
	case BooleanFilterFalse:
		return util.Ptr(false)
	default:
		panic("invalid boolean filter")
	}
}

// ListRepositoryOptions represents options for listing repositories
type ListRepositoryOptions struct {
	// OrderBy specifies the ordering of the repositories
	OrderBy RepositoryOrder
	// Privacy specifies the privacy level of the repositories
	Privacy RepositoryPrivacy
	// OwnerAffiliations specifies the affiliations of the user to the repositories
	OwnerAffiliations []RepositoryAffiliation
	// Limit specifies the maximum number of repositories to return
	Limit int
	// IsFork filters for repositories that are forks
	IsFork BooleanFilter
	// IsArchived filters for repositories that are archived
	IsArchived BooleanFilter
}

// Ordering options for repository connections
type RepositoryOrder struct {
	// The ordering direction.
	Direction OrderDirection `json:"direction"`
	// The field to order repositories by.
	Field RepositoryOrderField `json:"field"`
}

// Possible directions in which to order a list of items when provided an `orderBy` argument.
type OrderDirection int

const (
	// Specifies an ascending order for a given `orderBy` argument.
	OrderDirectionAsc OrderDirection = 1 + iota
	// Specifies a descending order for a given `orderBy` argument.
	OrderDirectionDesc
)

// Properties by which repository connections can be ordered.
type RepositoryOrderField int

const (
	// Order repositories by creation time
	RepositoryOrderFieldCreatedAt RepositoryOrderField = iota
	// Order repositories by name
	RepositoryOrderFieldName
	// Order repositories by push time
	RepositoryOrderFieldPushedAt
	// Order repositories by number of stargazers
	RepositoryOrderFieldStargazers
	// Order repositories by update time
	RepositoryOrderFieldUpdatedAt
)

// The privacy of a repository
type RepositoryPrivacy int

const (
	// Private
	RepositoryPrivacyPrivate RepositoryPrivacy = iota
	// Public
	RepositoryPrivacyPublic
)

// The affiliation of a user to a repository
type RepositoryAffiliation int

const (
	// Repositories that the user has been added to as a collaborator.
	RepositoryAffiliationCollaborator RepositoryAffiliation = iota
	// Repositories that the user has access to through being a member of an
	// organization. This includes every repository on every team that the user is on.
	RepositoryAffiliationOrganizationMember
	// Repositories that are owned by the authenticated user.
	RepositoryAffiliationOwner
)
