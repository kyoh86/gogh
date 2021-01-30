package gogh

import (
	"context"
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

func (c *RemoteController) repoSpec(repo *github.Repository) (Spec, error) {
	rawURL := strings.TrimSuffix(repo.GetCloneURL(), ".git")
	u, err := url.Parse(rawURL)
	if err != nil {
		return Spec{}, fmt.Errorf("parse clone-url %q: %w", rawURL, err)
	}
	user, name := path.Split(u.Path)
	return NewSpec(u.Host, strings.TrimLeft(strings.TrimRight(user, "/"), "/"), name)
}

func (c *RemoteController) repoListSpecList(query string, repos []*github.Repository) (specs []Spec, _ error) {
	for _, repo := range repos {
		spec, err := c.repoSpec(repo)
		if err != nil {
			return nil, err
		}
		if !strings.Contains(spec.String(), query) {
			continue
		}
		specs = append(specs, spec)
	}
	return
}

type RemoteListOption struct {
	User    string
	Query   string
	Options *github.RepositoryListOptions
}

func (o *RemoteListOption) GetUser() string {
	if o == nil {
		return ""
	}
	return o.User
}

func (o *RemoteListOption) GetQuery() string {
	if o == nil {
		return ""
	}
	return o.Query
}

func (o *RemoteListOption) GetOptions() *github.RepositoryListOptions {
	if o == nil {
		return nil
	}
	return o.Options
}

func (c *RemoteController) List(ctx context.Context, option *RemoteListOption) ([]Spec, error) {
	repos, _, err := c.adaptor.RepositoryList(ctx, option.GetUser(), option.GetOptions())
	if err != nil {
		return nil, err
	}
	return c.repoListSpecList(option.GetQuery(), repos)
}

type RemoteListByOrgOption struct {
	Query   string
	Options *github.RepositoryListByOrgOptions
}

func (o *RemoteListByOrgOption) GetQuery() string {
	if o == nil {
		return ""
	}
	return o.Query
}

func (o *RemoteListByOrgOption) GetOptions() *github.RepositoryListByOrgOptions {
	if o == nil {
		return nil
	}
	return o.Options
}

func (c *RemoteController) ListByOrg(ctx context.Context, org string, option *RemoteListByOrgOption) ([]Spec, error) {
	repos, _, err := c.adaptor.RepositoryListByOrg(ctx, org, option.GetOptions())
	if err != nil {
		return nil, err
	}
	return c.repoListSpecList(option.GetQuery(), repos)
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

func (c *RemoteController) Create(ctx context.Context, name string, option *RemoteCreateOption) (Spec, error) {
	repo, _, err := c.adaptor.RepositoryCreate(ctx, option.GetOrganization(), option.buildRepository(name))
	if err != nil {
		return Spec{}, fmt.Errorf("create a repository: %w", err)
	}
	return c.repoSpec(repo)
}

type RemoteDeleteOption struct{}

func (c *RemoteController) Delete(ctx context.Context, name string, _ *RemoteDeleteOption) error {
	// UNDONE: implement
	return nil
}
