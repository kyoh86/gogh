package remote

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/kyoh86/gogh/v3/domain/reporef"
	"github.com/kyoh86/gogh/v3/infra/github"
	"github.com/kyoh86/gogh/v3/infra/githubv4"
	"github.com/kyoh86/gogh/v3/util"
)

const DefaultHost = "github.com"

type Controller struct {
	adaptor github.Adaptor
}

func NewController(adaptor github.Adaptor) *Controller {
	return &Controller{
		adaptor: adaptor,
	}
}

func parseRepoRef(repo *github.Repository) (reporef.RepoRef, error) {
	rawURL := strings.TrimSuffix(repo.GetCloneURL(), ".git")
	u, err := url.Parse(rawURL)
	if err != nil {
		return reporef.RepoRef{}, fmt.Errorf("parse clone-url %q: %w", rawURL, err)
	}
	owner, name := path.Split(u.Path)

	return reporef.NewRepoRef(u.Host, strings.TrimLeft(strings.TrimRight(owner, "/"), "/"), name)
}

func ingestRepository(repo *github.Repository) (Repo, error) {
	var parentRepoRef *reporef.RepoRef
	if parent := repo.GetParent(); parent != nil {
		repoRef, err := parseRepoRef(parent)
		if err != nil {
			return Repo{}, fmt.Errorf("parse parent repository as local repo ref: %w", err)
		}
		parentRepoRef = &repoRef
	}
	repoRef, err := parseRepoRef(repo)
	if err != nil {
		return Repo{}, fmt.Errorf("parse repository as local repo ref: %w", err)
	}
	return Repo{
		Ref:         repoRef,
		URL:         strings.TrimSuffix(repo.GetCloneURL(), ".git"),
		Description: repo.GetDescription(),
		Homepage:    repo.GetHomepage(),
		Language:    repo.GetLanguage(),
		UpdatedAt:   repo.GetUpdatedAt().Time,
		Archived:    repo.GetArchived(),
		Private:     repo.GetPrivate(),
		IsTemplate:  repo.GetIsTemplate(),
		Fork:        repo.GetFork(),
		Parent:      parentRepoRef,
	}, nil
}

const (
	RepoListMaxLimitPerPage = 100
)

type RepoRelation string

const (
	RepoRelationOwner              = RepoRelation("owner")
	RepoRelationOrganizationMember = RepoRelation("organizationMember")
	RepoRelationCollaborator       = RepoRelation("collaborator")
)

var AllRepoRelation = []RepoRelation{
	RepoRelationOwner,
	RepoRelationOrganizationMember,
	RepoRelationCollaborator,
}

func (r RepoRelation) String() string {
	return string(r)
}

type RepoOrderField = githubv4.RepositoryOrderField

var AllRepoOrderField = []githubv4.RepositoryOrderField{
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

type ListOption struct {
	Private    *bool
	IsFork     *bool
	IsArchived *bool
	Order      OrderDirection
	Sort       RepoOrderField
	Relation   []RepoRelation
	Limit      int
}

func (o *ListOption) GetOptions() *github.RepositoryListOptions {
	owner := github.RepositoryAffiliationOwner
	if o == nil {
		return &github.RepositoryListOptions{
			OrderBy: github.RepositoryOrder{
				Field:     github.RepositoryOrderFieldUpdatedAt,
				Direction: githubv4.OrderDirectionDesc,
			},
			Limit:             RepoListMaxLimitPerPage,
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
		opt.Limit = RepoListMaxLimitPerPage
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
			case RepoRelationOwner:
				opt.OwnerAffiliations = append(opt.OwnerAffiliations, owner)
			case RepoRelationOrganizationMember:
				opt.OwnerAffiliations = append(opt.OwnerAffiliations, member)
			case RepoRelationCollaborator:
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

func (c *Controller) Me(
	ctx context.Context,
) (string, error) {
	user, err := c.adaptor.GetMe(ctx)
	if err != nil {
		return "", fmt.Errorf("get a user: %w", err)
	}
	return user, nil
}

func (c *Controller) List(
	ctx context.Context,
	option *ListOption,
) (allRepos []Repo, _ error) {
	sch, ech := c.ListAsync(ctx, option)
	for {
		select {
		case ref, more := <-sch:
			if !more {
				return
			}
			allRepos = append(allRepos, ref)
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
) (ret Repo, _ error) {
	ret.URL = repo.Url
	ret.IsTemplate = repo.IsTemplate
	ret.Archived = repo.IsArchived
	ret.Private = repo.IsPrivate
	ret.Fork = repo.IsFork
	ref, err := reporef.NewRepoRef(host, repo.Owner.GetLogin(), repo.Name)
	if err != nil {
		return Repo{}, err
	}
	ret.Ref = ref
	ret.Description = repo.Description
	ret.Homepage = repo.HomepageUrl
	ret.Language = repo.PrimaryLanguage.Name
	ret.UpdatedAt = repo.UpdatedAt

	if repo.Parent.Owner != nil && repo.Parent.Name != "" {
		parent, err := reporef.NewRepoRef(host, repo.Parent.Owner.GetLogin(), repo.Parent.Name)
		if err != nil {
			return Repo{}, err
		}
		ret.Parent = &parent
	}
	return ret, nil
}

var errOverLimit = errors.New("over limit")

func (c *Controller) repoList(
	repos []*github.RepositoryFragment,
	count *int,
	limit int,
	ch chan<- Repo,
) error {
	for _, repo := range repos {
		if limit > 0 && limit <= *count {
			return errOverLimit
		}
		ref, err := ingestRepositoryFragment(c.adaptor.GetHost(), repo)
		if err != nil {
			return err
		}
		ch <- ref
		*count++
	}
	return nil
}

func (c *Controller) ListAsync(
	ctx context.Context,
	option *ListOption,
) (<-chan Repo, <-chan error) {
	opt := option.GetOptions()
	sch := make(chan Repo, 1)
	ech := make(chan error, 1)
	go func() {
		defer close(sch)
		defer close(ech)

		var count int
		var limit int
		switch {
		case opt.Limit == 0:
			limit = 0
			opt.Limit = RepoListMaxLimitPerPage
		case opt.Limit > RepoListMaxLimitPerPage:
			limit = opt.Limit
			opt.Limit = RepoListMaxLimitPerPage
		default:
			limit = opt.Limit
		}
		for {
			repos, page, err := c.adaptor.RepositoryList(ctx, opt)
			if err != nil {
				ech <- err
				return
			}
			if err := c.repoList(repos, &count, limit, sch); err != nil {
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

type CreateOption struct {
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

func (o *CreateOption) buildRepository(name string) *github.Repository {
	if o == nil {
		return &github.Repository{Name: &name}
	}
	return &github.Repository{
		Name:                util.NilablePtr(name),
		Description:         util.NilablePtr(o.Description),
		Homepage:            util.NilablePtr(o.Homepage),
		Private:             util.NilablePtr(o.Private),
		HasIssues:           util.FalsePtr(o.DisableIssues),
		HasProjects:         util.FalsePtr(o.DisableProjects),
		HasWiki:             util.FalsePtr(o.DisableWiki),
		HasDownloads:        util.FalsePtr(o.DisableDownloads),
		IsTemplate:          util.NilablePtr(o.IsTemplate),
		TeamID:              util.NilablePtr(o.TeamID),
		AutoInit:            util.NilablePtr(o.AutoInit),
		GitignoreTemplate:   util.NilablePtr(o.GitignoreTemplate),
		LicenseTemplate:     util.NilablePtr(o.LicenseTemplate),
		AllowSquashMerge:    util.FalsePtr(o.PreventSquashMerge),
		AllowMergeCommit:    util.FalsePtr(o.PreventMergeCommit),
		AllowRebaseMerge:    util.FalsePtr(o.PreventRebaseMerge),
		DeleteBranchOnMerge: util.NilablePtr(o.DeleteBranchOnMerge),
	}
}

func (o *CreateOption) GetOrganization() string {
	if o == nil {
		return ""
	}
	return o.Organization
}

func (c *Controller) Create(
	ctx context.Context,
	name string,
	option *CreateOption,
) (Repo, error) {
	repo, _, err := c.adaptor.RepositoryCreate(
		ctx,
		option.GetOrganization(),
		option.buildRepository(name),
	)
	if err != nil {
		return Repo{}, fmt.Errorf("create a repository: %w", err)
	}
	return ingestRepository(repo)
}

type CreateFromTemplateOption struct {
	Owner       string `json:"owner,omitempty"`
	Description string `json:"description,omitempty"`
	Private     bool   `json:"private,omitempty"`
}

func (o *CreateFromTemplateOption) buildTemplateRepoRequest(
	name string,
) *github.TemplateRepoRequest {
	if o == nil {
		return &github.TemplateRepoRequest{Name: &name}
	}
	return &github.TemplateRepoRequest{
		Name:        util.NilablePtr(name),
		Owner:       util.NilablePtr(o.Owner),
		Description: util.NilablePtr(o.Description),
		Private:     util.FalsePtr(o.Private),
	}
}

func (c *Controller) CreateFromTemplate(
	ctx context.Context,
	templateOwner, templateName, name string,
	option *CreateFromTemplateOption,
) (Repo, error) {
	repo, _, err := c.adaptor.RepositoryCreateFromTemplate(
		ctx,
		templateOwner,
		templateName,
		option.buildTemplateRepoRequest(name),
	)
	if err != nil {
		return Repo{}, fmt.Errorf("create a repository from template: %w", err)
	}
	return ingestRepository(repo)
}

type ForkOption struct {
	// Organization is the name of the organization that owns the repository.
	Organization string
}

func (o *ForkOption) GetOptions() *github.RepositoryCreateForkOptions {
	if o == nil {
		return nil
	}
	return &github.RepositoryCreateForkOptions{
		Organization: o.Organization,
	}
}

func (c *Controller) Fork(
	ctx context.Context,
	owner string,
	name string,
	option *ForkOption,
) (Repo, error) {
	repo, _, err := c.adaptor.RepositoryCreateFork(ctx, owner, name, option.GetOptions())
	if err != nil {
		var acc *github.AcceptedError
		if !errors.As(err, &acc) {
			return Repo{}, fmt.Errorf("fork a repository: %w", err)
		}
	}
	return ingestRepository(repo)
}

type GetOption struct{}

func (c *Controller) Get(
	ctx context.Context,
	owner string,
	name string,
	_ *GetOption,
) (Repo, error) {
	repo, _, err := c.adaptor.RepositoryGet(ctx, owner, name)
	if err != nil {
		return Repo{}, fmt.Errorf("get a repository: %w", err)
	}
	return ingestRepository(repo)
}

type DeleteOption struct{}

func (c *Controller) Delete(
	ctx context.Context,
	owner string,
	name string,
	_ *DeleteOption,
) error {
	if _, err := c.adaptor.RepositoryDelete(ctx, owner, name); err != nil {
		return fmt.Errorf("delete a repository: %w", err)
	}
	return nil
}
