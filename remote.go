package gogh

import (
	"context"
	"net/http"

	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/wacul/ptr"
)

type RemoteController struct {
	adaptor github.Adaptor
}

func GithubAdaptor(ctx context.Context, server Server) (github.Adaptor, error) {
	var client *http.Client
	if server.Token() != "" {
		client = github.NewAuthClient(ctx, server.Token())
	}
	//UNDONE: support Enterprise with server.baseURL and server.uploadURL
	return github.NewAdaptor(client), nil
}

func NewRemoteController(adaptor github.Adaptor) *RemoteController {
	return &RemoteController{
		adaptor: adaptor,
	}
}

type RemoteListOption struct {
	User    string
	Query   string
	Options *github.RepositoryListOptions
}

func (c *RemoteController) List(ctx context.Context, option *RemoteListOption) ([]Project, error) {
	// UNDONE: implement
	return nil, nil
}

type RemoteListByOrgOption struct {
	Query   string
	Options *github.RepositoryListByOrgOptions
}

func (c *RemoteController) ListByOrg(ctx context.Context, org string, option *RemoteListByOrgOption) ([]Project, error) {
	// UNDONE: implement
	return nil, nil
}

type RemoteCreateOption struct {
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

// falsePtr converts bool (default: false) to *bool (default: nil as true)
func falsePtr(b bool) *bool {
	if b {
		return ptr.Bool(false)
	}
	return nil // == ptr.Bool(true)
}

func (o *RemoteCreateOption) buildRepository(name string) *github.Repository {
	if o == nil {
		return &github.Repository{Name: &name}
	}
	return &github.Repository{
		Name:                &name,
		Description:         &o.Description,
		Homepage:            &o.Homepage,
		Visibility:          &o.Visibility,
		HasIssues:           falsePtr(o.DisableIssues),
		HasProjects:         falsePtr(o.DisableProjects),
		HasWiki:             falsePtr(o.DisableWiki),
		HasDownloads:        falsePtr(o.DisableDownloads),
		IsTemplate:          &o.IsTemplate,
		TeamID:              &o.TeamID,
		AutoInit:            &o.AutoInit,
		GitignoreTemplate:   &o.GitignoreTemplate,
		LicenseTemplate:     &o.LicenseTemplate,
		AllowSquashMerge:    falsePtr(o.PreventSquashMerge),
		AllowMergeCommit:    falsePtr(o.PreventMergeCommit),
		AllowRebaseMerge:    falsePtr(o.PreventRebaseMerge),
		DeleteBranchOnMerge: &o.DeleteBranchOnMerge,
	}
}

func (c *RemoteController) Create(ctx context.Context, description Description, option *RemoteCreateOption) (*Project, error) {
	// UNDONE: implement
	return nil, nil
}

type RemoteDeleteOption struct{}

func (c *RemoteController) Delete(ctx context.Context, description Description, _ *RemoteDeleteOption) (*Project, error) {
	// UNDONE: implement
	return nil, nil
}
