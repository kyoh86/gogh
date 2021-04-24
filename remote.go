package gogh

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/kyoh86/gogh/v2/internal/github"
)

type RemoteController struct {
	adaptor github.Adaptor
}

func NewRemoteController(adaptor github.Adaptor) *RemoteController {
	return &RemoteController{
		adaptor: adaptor,
	}
}

func (c *RemoteController) parseSpec(repo *github.Repository) (Spec, error) {
	rawURL := strings.TrimSuffix(repo.GetCloneURL(), ".git")
	u, err := url.Parse(rawURL)
	if err != nil {
		return Spec{}, fmt.Errorf("parse clone-url %q: %w", rawURL, err)
	}
	owner, name := path.Split(u.Path)

	return NewSpec(u.Host, strings.TrimLeft(strings.TrimRight(owner, "/"), "/"), name)
}

func (c *RemoteController) ingest(repo *github.Repository) (Repository, error) {
	var parentSpec *Spec
	if parent := repo.GetParent(); parent != nil {
		spec, err := c.parseSpec(parent)
		if err != nil {
			return Repository{}, fmt.Errorf("parse parent repository as local spec: %w", err)
		}
		parentSpec = &spec
	}
	spec, err := c.parseSpec(repo)
	if err != nil {
		return Repository{}, fmt.Errorf("parse repository as local spec: %w", err)
	}
	return Repository{
		spec:        spec,
		description: repo.GetDescription(),
		homepage:    repo.GetHomepage(),
		language:    repo.GetLanguage(),
		topics:      repo.Topics,
		pushedAt:    repo.GetPushedAt().Time,
		archived:    repo.GetArchived(),
		private:     repo.GetPrivate(),
		isTemplate:  repo.GetIsTemplate(),
		fork:        repo.GetFork(),
		parent:      parentSpec,
	}, nil
}

func (c *RemoteController) repoListSpecList(repos []*github.Repository, ch chan<- Repository) error {
	for _, repo := range repos {
		spec, err := c.ingest(repo)
		if err != nil {
			return err
		}
		ch <- spec
	}
	return nil
}

type RemoteListOption struct {
	Users []string
	Query string

	// How to sort the search results. Possible values are `stars`,
	// `fork` and `updated`.
	Sort string

	// Sort order if sort parameter is provided. Possible values are: asc,
	// desc. Default is desc.
	Order string

	// If non-nil, filters repositories according to whether they have been archived.
	Archived *bool

	// If non-null, filters repositories according to whether they are forks of another repository.
	IsFork *bool

	// If non-null, filters repositories according to whether they are private.
	IsPrivate *bool

	// Filter by primary coding language.
	Language string

	// If non-nill, limit a number of repositories to list.
	Limit *int
}

func (o *RemoteListOption) GetQuery() string {
	if o == nil {
		return "user:@me"
	}
	terms := make([]string, 0, 10)
	if len(o.Users) == 0 {
		terms = append(terms, "user:@me")
	} else {
		for _, u := range o.Users {
			terms = append(terms, fmt.Sprintf("user:%q", u))
		}
	}

	if o.Archived != nil {
		terms = append(terms, fmt.Sprintf("archived:%v", *o.Archived))
	}
	if o.IsFork != nil {
		terms = append(terms, fmt.Sprintf("fork:%v", *o.IsFork))
	}
	if o.IsPrivate != nil {
		if *o.IsPrivate {
			terms = append(terms, "is:private")
		} else {
			terms = append(terms, "is:public")
		}
	}
	if o.Language != "" {
		terms = append(terms, fmt.Sprintf("language:%q", o.Language))
	}
	return strings.Join(terms, " ")
}

func (o *RemoteListOption) GetOptions() *github.SearchOptions {
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	if o == nil {
		return opt
	}
	opt.Sort = o.Sort
	opt.Order = o.Order
	return opt
}

func (c *RemoteController) List(ctx context.Context, option *RemoteListOption) (allSpecs []Repository, _ error) {
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

func (c *RemoteController) ListAsync(ctx context.Context, option *RemoteListOption) (<-chan Repository, <-chan error) {
	opt := option.GetOptions()
	sch := make(chan Repository, 1)
	ech := make(chan error, 1)
	go func() {
		defer close(sch)
		defer close(ech)
		for {
			repos, resp, err := c.adaptor.SearchRepository(ctx, option.GetQuery(), opt)
			if err != nil {
				ech <- err
				return
			}
			if err := c.repoListSpecList(repos, sch); err != nil {
				ech <- err
				return
			}
			if resp.NextPage == 0 {
				return
			}
			opt.Page = resp.NextPage
		}
	}()
	return sch, ech
}

type RemoteCreateOption struct {
	// Organization is the name of the organization that owns the repository.
	Organization string

	// Description is a short description of the repository.
	Description string

	// Homepage is a URL with more information about the repository.
	Homepage string

	// Visibility can be public or private. If your organization is associated with an enterprise account using GitHub
	// Enterprise Cloud or GitHub Enterprise Server 2.20+, visibility can also be internal. For more information, see
	// "Creating an internal repository" in the GitHub Help documentation.  The visibility parameter overrides the private
	// parameter when you use both parameters with the nebula-preview preview header.
	Visibility string

	// DisableIssues is either false to enable issues for this repository or true to disable them.
	DisableIssues bool

	// DisableProjects is either false to enable projects for this repository or true to disable them. Note: If you're
	// creating a repository in an organization that has disabled repository projects, the default is false, and if you
	// pass true, the API returns an error.
	DisableProjects bool

	// DisableWiki is either false to enable the wiki for this repository or true to disable it.
	DisableWiki bool

	// DisableDownloads is either false to enable the downloads or true to disable it.
	DisableDownloads bool

	// IsTemplate is either true to make this repo available as a template repository or false to prevent it.
	IsTemplate bool

	// TeamID is the id of the team that will be granted access to this repository. This is only valid when creating a
	// repository in an organization.
	TeamID int64

	// AutoInit is pass true to create an initial commit with empty README.
	AutoInit bool

	// GitignoreTemplate is the desired language or platform .gitignore template to apply. Use the name of the template
	// without the extension. For example, "Haskell".
	GitignoreTemplate string

	// LicenseTemplate is an open source license template, and then use the license keyword as the licenseTemplate string.
	// For example, "mit" or "mpl-2.0".
	LicenseTemplate string

	// PreventSquashMerge is either false to allow squash-merging pull requests, or true to prevent squash-merging.
	PreventSquashMerge bool

	// PreventMergeCommit is either false to allow merging pull requests with a merge commit, or true to prevent merging
	// pull requests with merge commits.
	PreventMergeCommit bool

	// PreventRebaseMerge is either false to allow rebase-merging pull requests, or true to prevent rebase-merging.
	PreventRebaseMerge bool

	// DeleteBranchOnMerge is either true to allow automatically deleting head branches when pull requests are merged, or
	// false to prevent automatic deletion.
	DeleteBranchOnMerge bool
}

func (o *RemoteCreateOption) buildRepository(name string) *github.Repository {
	if o == nil {
		return &github.Repository{Name: &name}
	}
	return &github.Repository{
		Name:                stringPtr(name),
		Description:         stringPtr(o.Description),
		Homepage:            stringPtr(o.Homepage),
		Visibility:          stringPtr(o.Visibility),
		HasIssues:           falsePtr(o.DisableIssues),
		HasProjects:         falsePtr(o.DisableProjects),
		HasWiki:             falsePtr(o.DisableWiki),
		HasDownloads:        falsePtr(o.DisableDownloads),
		IsTemplate:          boolPtr(o.IsTemplate),
		TeamID:              int64Ptr(o.TeamID),
		AutoInit:            boolPtr(o.AutoInit),
		GitignoreTemplate:   stringPtr(o.GitignoreTemplate),
		LicenseTemplate:     stringPtr(o.LicenseTemplate),
		AllowSquashMerge:    falsePtr(o.PreventSquashMerge),
		AllowMergeCommit:    falsePtr(o.PreventMergeCommit),
		AllowRebaseMerge:    falsePtr(o.PreventRebaseMerge),
		DeleteBranchOnMerge: boolPtr(o.DeleteBranchOnMerge),
	}
}

func (o *RemoteCreateOption) GetOrganization() string {
	if o == nil {
		return ""
	}
	return o.Organization
}

func (c *RemoteController) Create(ctx context.Context, name string, option *RemoteCreateOption) (Repository, error) {
	repo, _, err := c.adaptor.RepositoryCreate(ctx, option.GetOrganization(), option.buildRepository(name))
	if err != nil {
		return Repository{}, fmt.Errorf("create a repository: %w", err)
	}
	return c.ingest(repo)
}

type RemoteCreateFromTemplateOption struct {
	Owner       string `json:"owner,omitempty"`
	Description string `json:"description,omitempty"`
	Private     bool   `json:"private,omitempty"`
}

func (o *RemoteCreateFromTemplateOption) buildTemplateRepoRequest(name string) *github.TemplateRepoRequest {
	if o == nil {
		return &github.TemplateRepoRequest{Name: &name}
	}
	return &github.TemplateRepoRequest{
		Name:        stringPtr(name),
		Owner:       stringPtr(o.Owner),
		Description: stringPtr(o.Description),
		Private:     falsePtr(o.Private),
	}
}

func (c *RemoteController) CreateFromTemplate(ctx context.Context, templateOwner, templateName, name string, option *RemoteCreateFromTemplateOption) (Repository, error) {
	repo, _, err := c.adaptor.RepositoryCreateFromTemplate(ctx, templateOwner, templateName, option.buildTemplateRepoRequest(name))
	if err != nil {
		return Repository{}, fmt.Errorf("create a repository from template: %w", err)
	}
	return c.ingest(repo)
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

func (c *RemoteController) Fork(ctx context.Context, owner string, name string, option *RemoteForkOption) (Repository, error) {
	repo, _, err := c.adaptor.RepositoryCreateFork(ctx, owner, name, option.GetOptions())
	if err != nil {
		var acc *github.AcceptedError
		if !errors.As(err, &acc) {
			return Repository{}, fmt.Errorf("fork a repository: %w", err)
		}
	}
	return c.ingest(repo)
}

type RemoteGetOption struct{}

func (c *RemoteController) Get(ctx context.Context, owner string, name string, _ *RemoteGetOption) (Repository, error) {
	repo, _, err := c.adaptor.RepositoryGet(ctx, owner, name)
	if err != nil {
		return Repository{}, fmt.Errorf("get a repository: %w", err)
	}
	return c.ingest(repo)
}

type RemoteSourceOption struct{}

func (c *RemoteController) GetSource(ctx context.Context, owner string, name string, _ *RemoteSourceOption) (Repository, error) {
	repo, _, err := c.adaptor.RepositoryGet(ctx, owner, name)
	if err != nil {
		return Repository{}, fmt.Errorf("get a repository: %w", err)
	}
	if source := repo.GetSource(); source != nil {
		return c.ingest(source)
	}
	return c.ingest(repo)
}

type RemoteParentOption struct{}

func (c *RemoteController) GetParent(ctx context.Context, owner string, name string, _ *RemoteParentOption) (Repository, error) {
	repo, _, err := c.adaptor.RepositoryGet(ctx, owner, name)
	if err != nil {
		return Repository{}, fmt.Errorf("get a repository: %w", err)
	}
	if parent := repo.GetParent(); parent != nil {
		return c.ingest(parent)
	}
	return c.ingest(repo)
}

type RemoteDeleteOption struct{}

func (c *RemoteController) Delete(ctx context.Context, owner string, name string, _ *RemoteDeleteOption) error {
	if _, err := c.adaptor.RepositoryDelete(ctx, owner, name); err != nil {
		return fmt.Errorf("delete a repository: %w", err)
	}
	return nil
}
