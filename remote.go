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
		Spec:        spec,
		URL:         strings.TrimSuffix(repo.GetCloneURL(), ".git"),
		Description: repo.GetDescription(),
		Homepage:    repo.GetHomepage(),
		Language:    repo.GetLanguage(),
		Topics:      repo.Topics,
		PushedAt:    repo.GetPushedAt().Time,
		Archived:    repo.GetArchived(),
		Private:     repo.GetPrivate(),
		IsTemplate:  repo.GetIsTemplate(),
		Fork:        repo.GetFork(),
		Parent:      parentSpec,
	}, nil
}

type RemoteListOption struct {
	IsPrivate *bool
	Limit     *int
	Archived  *bool
	IsFork    *bool
	Query     string
	Order     string
	Language  string
	Sort      string
	Users     []string
}

func (o *RemoteListOption) GetQuery() string {
	if o == nil {
		return "user:@me fork:true sort:updated"
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
	if o.IsFork == nil {
		terms = append(terms, "fork:true")
	} else if *o.IsFork {
		terms = append(terms, "fork:only")
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
	if o.Query != "" {
		terms = append(terms, o.Query)
	}
	if o.Sort == "" {
		if o.Query == "" {
			terms = append(terms, "sort:updated")
		}
	} else {
		terms = append(terms, fmt.Sprintf("sort:%q", o.Sort))
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
	Organization        string
	Description         string
	Homepage            string
	Visibility          string
	LicenseTemplate     string
	GitignoreTemplate   string
	TeamID              int64
	IsTemplate          bool
	DisableDownloads    bool
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

type RemoteDeleteOption struct{}

func (c *RemoteController) Delete(ctx context.Context, owner string, name string, _ *RemoteDeleteOption) error {
	if _, err := c.adaptor.RepositoryDelete(ctx, owner, name); err != nil {
		return fmt.Errorf("delete a repository: %w", err)
	}
	return nil
}

type RemoteListOrganizationsOption struct{}

func (o *RemoteListOrganizationsOption) GetOptions() *github.ListOptions {
	opt := &github.ListOptions{
		PerPage: 100,
	}
	if o == nil {
		return opt
	}
	return opt
}

func (c *RemoteController) ListOrganizations(ctx context.Context, option *RemoteListOrganizationsOption) (allNames []string, _ error) {
	nch, ech := c.ListOrganizationsAsync(ctx, option)
	for {
		select {
		case name, more := <-nch:
			if !more {
				return
			}
			allNames = append(allNames, name)
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

func (c *RemoteController) ListOrganizationsAsync(ctx context.Context, option *RemoteListOrganizationsOption) (<-chan string, <-chan error) {
	opt := option.GetOptions()
	nch := make(chan string, 1)
	ech := make(chan error, 1)
	go func() {
		defer close(nch)
		defer close(ech)
		for {
			orgs, resp, err := c.adaptor.OrganizationsList(ctx, opt)
			if err != nil {
				ech <- err
				return
			}
			for _, org := range orgs {
				nch <- org.GetLogin()
			}
			if resp.NextPage == 0 {
				return
			}
			opt.Page = resp.NextPage
		}
	}()
	return nch, ech
}
