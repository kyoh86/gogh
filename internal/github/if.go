package github

import (
	"context"

	github "github.com/google/go-github/v69/github"
	"github.com/kyoh86/gogh/v2/internal/githubv4"
)

type (
	Repository                  = github.Repository
	RepositoryCreateForkOptions = github.RepositoryCreateForkOptions
	TemplateRepoRequest         = github.TemplateRepoRequest
	ListOptions                 = github.ListOptions
	Response                    = github.Response
	User                        = github.User
	RepositoryPrivacy           = githubv4.RepositoryPrivacy
	RepositoryOrder             = githubv4.RepositoryOrder
	RepositoryOrderField        = githubv4.RepositoryOrderField
	OrderDirection              = githubv4.OrderDirection
	RepositoryAffiliation       = githubv4.RepositoryAffiliation
	PageInfoFragment            = githubv4.PageInfoFragment
	LanguageFragment            = githubv4.LanguageFragment
	OwnerFragment               = githubv4.OwnerFragment
	ParentRepositoryFragment    = githubv4.ParentRepositoryFragment
	RepositoryFragment          = githubv4.RepositoryFragment
)

const (
	OrderDirectionAsc  OrderDirection = githubv4.OrderDirectionAsc
	OrderDirectionDesc OrderDirection = githubv4.OrderDirectionDesc
)

const (
	RepositoryPrivacyPublic  RepositoryPrivacy = githubv4.RepositoryPrivacyPublic
	RepositoryPrivacyPrivate RepositoryPrivacy = githubv4.RepositoryPrivacyPrivate
)

const (
	RepositoryOrderFieldCreatedAt  RepositoryOrderField = githubv4.RepositoryOrderFieldCreatedAt
	RepositoryOrderFieldUpdatedAt  RepositoryOrderField = githubv4.RepositoryOrderFieldUpdatedAt
	RepositoryOrderFieldPushedAt   RepositoryOrderField = githubv4.RepositoryOrderFieldPushedAt
	RepositoryOrderFieldName       RepositoryOrderField = githubv4.RepositoryOrderFieldName
	RepositoryOrderFieldStargazers RepositoryOrderField = githubv4.RepositoryOrderFieldStargazers
)

const (
	RepositoryAffiliationOwner              RepositoryAffiliation = githubv4.RepositoryAffiliationOwner
	RepositoryAffiliationCollaborator       RepositoryAffiliation = githubv4.RepositoryAffiliationCollaborator
	RepositoryAffiliationOrganizationMember RepositoryAffiliation = githubv4.RepositoryAffiliationOrganizationMember
)

type RepositoryListOptions struct {
	OrderBy           RepositoryOrder
	After             string
	Privacy           RepositoryPrivacy
	OwnerAffiliations []RepositoryAffiliation
	Limit             int
	IsFork            *bool
	IsArchived        *bool
}

type Adaptor interface {
	GetHost() string
	GetMe(ctx context.Context) (string, error)
	GetAuthenticatedUser(ctx context.Context) (*User, *Response, error)
	UserGet(ctx context.Context, user string) (*User, *Response, error)
	RepositoryList(
		ctx context.Context,
		opts *RepositoryListOptions,
	) ([]*RepositoryFragment, PageInfoFragment, error)
	RepositoryCreate(
		ctx context.Context,
		org string,
		repo *Repository,
	) (*Repository, *Response, error)
	RepositoryCreateFork(
		ctx context.Context,
		owner string,
		repo string,
		opts *RepositoryCreateForkOptions,
	) (*Repository, *Response, error)
	RepositoryCreateFromTemplate(
		ctx context.Context,
		templateOwner, templateRepo string,
		templateRepoReq *TemplateRepoRequest,
	) (*Repository, *Response, error)
	RepositoryDelete(ctx context.Context, owner string, repo string) (*Response, error)
	RepositoryGet(ctx context.Context, owner string, repo string) (*Repository, *Response, error)
}
