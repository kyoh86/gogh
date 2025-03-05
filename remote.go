package gogh

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/kyoh86/gogh/v3/internal/github"
	"github.com/kyoh86/gogh/v3/internal/githubv4"
)

const DefaultHost = "github.com"

type RemoteController struct {
	adaptor github.Adaptor
}

func NewRemoteController(adaptor github.Adaptor) *RemoteController {
	return &RemoteController{
		adaptor: adaptor,
	}
}

func parseSpec(repo *github.Repository) (Spec, error) {
	rawURL := strings.TrimSuffix(repo.GetCloneURL(), ".git")
	u, err := url.Parse(rawURL)
	if err != nil {
		return Spec{}, fmt.Errorf("parse clone-url %q: %w", rawURL, err)
	}
	owner, name := path.Split(u.Path)

	return NewSpec(u.Host, strings.TrimLeft(strings.TrimRight(owner, "/"), "/"), name)
}

func ingestRepository(repo *github.Repository) (Repository, error) {
	var parentSpec *Spec
	if parent := repo.GetParent(); parent != nil {
		spec, err := parseSpec(parent)
		if err != nil {
			return Repository{}, fmt.Errorf("parse parent repository as local spec: %w", err)
		}
		parentSpec = &spec
	}
	spec, err := parseSpec(repo)
	if err != nil {
		return Repository{}, fmt.Errorf("parse repository as local spec: %w", err)
	}
	return Repository{
		Spec:        spec,
		URL:         strings.TrimSuffix(repo.GetCloneURL(), ".git"),
		Description: repo.GetDescription(),
		Homepage:    repo.GetHomepage(),
		Language:    repo.GetLanguage(),
		UpdatedAt:   repo.GetUpdatedAt().Time,
		Archived:    repo.GetArchived(),
		Private:     repo.GetPrivate(),
		IsTemplate:  repo.GetIsTemplate(),
		Fork:        repo.GetFork(),
		Parent:      parentSpec,
	}, nil
}

const (
	RepositoryListMaxLimitPerPage = 100
)

type RepositoryRelation string

const (
	RepositoryRelationOwner              = RepositoryRelation("owner")
	RepositoryRelationOrganizationMember = RepositoryRelation("organizationMember")
	RepositoryRelationCollaborator       = RepositoryRelation("collaborator")
)

var AllRepositoryRelation = []RepositoryRelation{
	RepositoryRelationOwner,
	RepositoryRelationOrganizationMember,
	RepositoryRelationCollaborator,
}

func (r RepositoryRelation) String() string {
	return string(r)
}

type RepositoryOrderField = githubv4.RepositoryOrderField

var AllRepositoryOrderField = []githubv4.RepositoryOrderField{
	githubv4.RepositoryOrderFieldCreatedAt,
	githubv4.RepositoryOrderFieldName,
	githubv4.RepositoryOrderFieldPushedAt,
	githubv4.RepositoryOrderFieldStargazers,
	githubv4.RepositoryOrderFieldUpdatedAt,
}

type OrderDirection = githubv4.OrderDirection

var AllOrderDirection = []githubv4.OrderDirection{
	githubv4.OrderDirectionAsc,
	githubv4.OrderDirectionDesc,
}

type RemoteListOption struct {
	Private    *bool
	IsFork     *bool
	IsArchived *bool
	Order      OrderDirection
	Sort       RepositoryOrderField
	Relation   []RepositoryRelation
	Limit      int
}

func (o *RemoteListOption) GetOptions() *github.RepositoryListOptions {
	owner := github.RepositoryAffiliationOwner
	if o == nil {
		return &github.RepositoryListOptions{
			OrderBy: github.RepositoryOrder{
				Field:     github.RepositoryOrderFieldUpdatedAt,
				Direction: githubv4.OrderDirectionDesc,
			},
			Limit:             RepositoryListMaxLimitPerPage,
			OwnerAffiliations: []github.RepositoryAffiliation{owner},
		}
	}
	opt := &github.RepositoryListOptions{
		IsFork:     o.IsFork,
		IsArchived: o.IsArchived,
	}

	if o.Sort == "" {
		opt.OrderBy = github.RepositoryOrder{
			Field:     github.RepositoryOrderFieldUpdatedAt,
			Direction: githubv4.OrderDirectionDesc,
		}
	} else {
		opt.OrderBy = github.RepositoryOrder{
			Field: o.Sort,
		}
		if o.Order == "" {
			if opt.OrderBy.Field == github.RepositoryOrderFieldName {
				opt.OrderBy.Direction = github.OrderDirectionAsc
			} else {
				opt.OrderBy.Direction = github.OrderDirectionDesc
			}
		} else {
			opt.OrderBy.Direction = o.Order
		}
	}
	if o.Limit == 0 {
		opt.Limit = RepositoryListMaxLimitPerPage
	} else {
		opt.Limit = o.Limit
	}

	if len(o.Relation) == 0 {
		opt.OwnerAffiliations = []github.RepositoryAffiliation{owner}
	} else {
		member := github.RepositoryAffiliationOrganizationMember
		collabo := github.RepositoryAffiliationCollaborator
		for _, r := range o.Relation {
			switch r {
			case RepositoryRelationOwner:
				opt.OwnerAffiliations = append(opt.OwnerAffiliations, owner)
			case RepositoryRelationOrganizationMember:
				opt.OwnerAffiliations = append(opt.OwnerAffiliations, member)
			case RepositoryRelationCollaborator:
				opt.OwnerAffiliations = append(opt.OwnerAffiliations, collabo)
			}
		}
	}

	if o.Private != nil {
		if *o.Private {
			private := github.RepositoryPrivacyPrivate
			opt.Privacy = private
		} else {
			public := github.RepositoryPrivacyPublic
			opt.Privacy = public
		}
	}
	return opt
}

func (c *RemoteController) Me(
	ctx context.Context,
) (string, error) {
	user, err := c.adaptor.GetMe(ctx)
	if err != nil {
		return "", fmt.Errorf("get a user: %w", err)
	}
	return user, nil
}

func (c *RemoteController) List(
	ctx context.Context,
	option *RemoteListOption,
) (allSpecs []Repository, _ error) {
	sch, ech := c.ListAsync(ctx, option)
	for {
		select {
		case spec, more := <-sch:
			if !more {
				return
			}
			allSpecs = append(allSpecs, spec)
		case err := <-ech:
			if err != nil {
				return nil, err
			}
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				return nil, err
			}
		}
	}
}

func ingestRepositoryFragment(
	host string,
	repo *github.RepositoryFragment,
) (ret Repository, _ error) {
	ret.URL = repo.Url
	ret.IsTemplate = repo.IsTemplate
	ret.Archived = repo.IsArchived
	ret.Private = repo.IsPrivate
	ret.Fork = repo.IsFork
	spec, err := NewSpec(host, repo.Owner.GetLogin(), repo.Name)
	if err != nil {
		return Repository{}, err
	}
	ret.Spec = spec
	ret.Description = repo.Description
	ret.Homepage = repo.HomepageUrl
	ret.Language = repo.PrimaryLanguage.Name
	ret.UpdatedAt = repo.UpdatedAt

	if repo.Parent.Owner != nil && repo.Parent.Name != "" {
		parent, err := NewSpec(host, repo.Parent.Owner.GetLogin(), repo.Parent.Name)
		if err != nil {
			return Repository{}, err
		}
		ret.Parent = &parent
	}
	return ret, nil
}

var errOverLimit = errors.New("over limit")

func (c *RemoteController) repoListSpecList(
	repos []*github.RepositoryFragment,
	count *int,
	limit int,
	ch chan<- Repository,
) error {
	for _, repo := range repos {
		if limit > 0 && limit <= *count {
			return errOverLimit
		}
		spec, err := ingestRepositoryFragment(c.adaptor.GetHost(), repo)
		if err != nil {
			return err
		}
		ch <- spec
		*count++
	}
	return nil
}

func (c *RemoteController) ListAsync(
	ctx context.Context,
	option *RemoteListOption,
) (<-chan Repository, <-chan error) {
	opt := option.GetOptions()
	sch := make(chan Repository, 1)
	ech := make(chan error, 1)
	go func() {
		defer close(sch)
		defer close(ech)

		var count int
		var limit int
		switch {
		case opt.Limit == 0:
			limit = 0
			opt.Limit = RepositoryListMaxLimitPerPage
		case opt.Limit > RepositoryListMaxLimitPerPage:
			limit = opt.Limit
			opt.Limit = RepositoryListMaxLimitPerPage
		default:
			limit = opt.Limit
		}
		for {
			repos, page, err := c.adaptor.RepositoryList(ctx, opt)
			if err != nil {
				ech <- err
				return
			}
			if err := c.repoListSpecList(repos, &count, limit, sch); err != nil {
				if errors.Is(err, errOverLimit) {
					return
				}
				ech <- err
				return
			}
			if !page.HasNextPage {
				return
			}
			opt.After = page.EndCursor
		}
	}()
	return sch, ech
}

type RemoteCreateOption struct {
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

func (o *RemoteCreateOption) buildRepository(name string) *github.Repository {
	if o == nil {
		return &github.Repository{Name: &name}
	}
	return &github.Repository{
		Name:                NilablePtr(name),
		Description:         NilablePtr(o.Description),
		Homepage:            NilablePtr(o.Homepage),
		Private:             NilablePtr(o.Private),
		HasIssues:           FalsePtr(o.DisableIssues),
		HasProjects:         FalsePtr(o.DisableProjects),
		HasWiki:             FalsePtr(o.DisableWiki),
		HasDownloads:        FalsePtr(o.DisableDownloads),
		IsTemplate:          NilablePtr(o.IsTemplate),
		TeamID:              NilablePtr(o.TeamID),
		AutoInit:            NilablePtr(o.AutoInit),
		GitignoreTemplate:   NilablePtr(o.GitignoreTemplate),
		LicenseTemplate:     NilablePtr(o.LicenseTemplate),
		AllowSquashMerge:    FalsePtr(o.PreventSquashMerge),
		AllowMergeCommit:    FalsePtr(o.PreventMergeCommit),
		AllowRebaseMerge:    FalsePtr(o.PreventRebaseMerge),
		DeleteBranchOnMerge: NilablePtr(o.DeleteBranchOnMerge),
	}
}

func (o *RemoteCreateOption) GetOrganization() string {
	if o == nil {
		return ""
	}
	return o.Organization
}

func (c *RemoteController) Create(
	ctx context.Context,
	name string,
	option *RemoteCreateOption,
) (Repository, error) {
	repo, _, err := c.adaptor.RepositoryCreate(
		ctx,
		option.GetOrganization(),
		option.buildRepository(name),
	)
	if err != nil {
		return Repository{}, fmt.Errorf("create a repository: %w", err)
	}
	return ingestRepository(repo)
}

type RemoteCreateFromTemplateOption struct {
	Owner       string `json:"owner,omitempty"`
	Description string `json:"description,omitempty"`
	Private     bool   `json:"private,omitempty"`
}

func (o *RemoteCreateFromTemplateOption) buildTemplateRepoRequest(
	name string,
) *github.TemplateRepoRequest {
	if o == nil {
		return &github.TemplateRepoRequest{Name: &name}
	}
	return &github.TemplateRepoRequest{
		Name:        NilablePtr(name),
		Owner:       NilablePtr(o.Owner),
		Description: NilablePtr(o.Description),
		Private:     FalsePtr(o.Private),
	}
}

func (c *RemoteController) CreateFromTemplate(
	ctx context.Context,
	templateOwner, templateName, name string,
	option *RemoteCreateFromTemplateOption,
) (Repository, error) {
	repo, _, err := c.adaptor.RepositoryCreateFromTemplate(
		ctx,
		templateOwner,
		templateName,
		option.buildTemplateRepoRequest(name),
	)
	if err != nil {
		return Repository{}, fmt.Errorf("create a repository from template: %w", err)
	}
	return ingestRepository(repo)
}

type RemoteForkOption struct {
	// Organization is the name of the organization that owns the repository.
	Organization string
}

func (o *RemoteForkOption) GetOptions() *github.RepositoryCreateForkOptions {
	if o == nil {
		return nil
	}
	return &github.RepositoryCreateForkOptions{
		Organization: o.Organization,
	}
}

func (c *RemoteController) Fork(
	ctx context.Context,
	owner string,
	name string,
	option *RemoteForkOption,
) (Repository, error) {
	repo, _, err := c.adaptor.RepositoryCreateFork(ctx, owner, name, option.GetOptions())
	if err != nil {
		var acc *github.AcceptedError
		if !errors.As(err, &acc) {
			return Repository{}, fmt.Errorf("fork a repository: %w", err)
		}
	}
	return ingestRepository(repo)
}

type RemoteGetOption struct{}

func (c *RemoteController) Get(
	ctx context.Context,
	owner string,
	name string,
	_ *RemoteGetOption,
) (Repository, error) {
	repo, _, err := c.adaptor.RepositoryGet(ctx, owner, name)
	if err != nil {
		return Repository{}, fmt.Errorf("get a repository: %w", err)
	}
	return ingestRepository(repo)
}

type RemoteDeleteOption struct{}

func (c *RemoteController) Delete(
	ctx context.Context,
	owner string,
	name string,
	_ *RemoteDeleteOption,
) error {
	if _, err := c.adaptor.RepositoryDelete(ctx, owner, name); err != nil {
		return fmt.Errorf("delete a repository: %w", err)
	}
	return nil
}

type MemberOfOption struct{}

func (c *RemoteController) MemberOf(
	ctx context.Context,
	orgName string,
	_ *MemberOfOption,
) (bool, error) {
	orgs, _, err := c.adaptor.OrganizationList(ctx)
	if err != nil {
		return false, err
	}
	for _, org := range orgs {
		if *org.Login == orgName {
			return true, nil
		}
	}
	return false, nil
}
