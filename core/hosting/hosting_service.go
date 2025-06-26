package hosting

import (
	"context"
	"iter"
	"net/url"

	"github.com/kyoh86/gogh/v4/core/auth"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/typ"
)

// HostingService provides access to remote repositories
type HostingService interface {
	// GetURLOf converts a repository reference to its remote URL
	GetURLOf(ref repository.Reference) (*url.URL, error)
	// ParseURL converts a remote URL to a repository reference
	ParseURL(u *url.URL) (*repository.Reference, error)
	// GetTokenFor retrieves an authentication token and user for a specific repository reference
	GetTokenFor(ctx context.Context, host, owner string) (string, auth.Token, error)
	// GetRepository retrieves repository information from a remote source
	GetRepository(context.Context, repository.Reference) (*Repository, error)
	// ListRepository retrieves a list of repositories from a remote source
	ListRepository(context.Context, ListRepositoryOptions) iter.Seq2[*Repository, error]
	// DeleteRepository deletes a repository from a remote source
	DeleteRepository(context.Context, repository.Reference) error
	// CreateRepository creates a new repository on the remote hosting service
	CreateRepository(
		ctx context.Context,
		ref repository.Reference,
		opts CreateRepositoryOptions,
	) (*Repository, error)
	// CreateRepositoryFromTemplate creates a new repository from an existing template repository
	CreateRepositoryFromTemplate(
		ctx context.Context,
		ref repository.Reference,
		tmp repository.Reference,
		opts CreateRepositoryFromTemplateOptions,
	) (*Repository, error)
	// ForkRepository creates a fork of a repository on the remote hosting service
	ForkRepository(
		ctx context.Context,
		ref repository.Reference,
		target repository.Reference,
		opts ForkRepositoryOptions,
	) (*Repository, error)
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
	IsFork typ.Tristate
	// IsArchived filters for repositories that are archived
	IsArchived typ.Tristate
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
	// None
	RepositoryPrivacyNone RepositoryPrivacy = iota
	// Private
	RepositoryPrivacyPrivate
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

type CreateRepositoryOptions struct {
	Description         string
	Homepage            string
	Organization        string
	LicenseTemplate     string
	GitignoreTemplate   string
	TeamID              int64
	DisableDownloads    bool
	IsTemplate          bool
	Private             bool
	DisableWiki         bool
	AutoInit            bool
	DisableProjects     bool
	DisableIssues       bool
	PreventSquashMerge  bool
	PreventMergeCommit  bool
	PreventRebaseMerge  bool
	DeleteBranchOnMerge bool
}

type CreateRepositoryFromTemplateOptions struct {
	Description        string
	IncludeAllBranches bool
	Private            bool
}

type ForkRepositoryOptions struct {
	DefaultBranchOnly bool
}
