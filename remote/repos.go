package remote

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/v24/github"
	"github.com/kyoh86/gogh/gogh"
)

// Repos will get a list of repositories for a user.
// Parameters:
//   * user:        Who has the repositories. Empty means the "me" (authenticated user, or GOGH_GITHUB_USER).
//   * own:         Include repositories that are owned by the user
//   * collaborate: Include repositories that the user has been added to as a collaborator
//   * member:      Include repositories that the user has access to through being a member of an organization. This includes every repository on every team that the user is on
//   * visibility:  Can be one of all, public, or private
//   * sort:        Can be one of created, updated, pushed, full_name
//   * direction:   Can be one of asc or desc default. Default means asc when using full_name, otherwise desc
// Returns:
//   List of the url for repoisitories
func Repos(ctx gogh.Context, user string, own, collaborate, member bool, visibility, sort, direction string) ([]string, error) {
	client, err := NewClient(ctx)
	if err != nil {
		return nil, err
	}

	/*
		Build GitHub requests.
		See: https://developer.github.com/v3/repos/#parameters
		- affiliation string  Comma-separated list of values. Can include:
				- owner: Repositories that are owned by the authenticated user.
				- collaborator: Repositories that the user has been added to as a collaborator.
				- organization_member: Repositories that the user has access to through being a member of an organization. This includes every repository on every team that the user is on.
	*/
	var affs []string
	if own {
		affs = append(affs, "owner")
	}
	if collaborate {
		affs = append(affs, "collaborator")
	}
	if member {
		affs = append(affs, "organization_member")
	}
	// If the context has no authentication token, specifies context user name for "me".
	if user == "" && !authenticated(ctx) {
		user = ctx.GitHubUser()
	}

	opts := &github.RepositoryListOptions{
		Visibility:  visibility,
		Affiliation: strings.Join(affs, ","),
		Sort:        sort,
		Direction:   direction,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	var list []string
	last := 1
	for page := 1; page <= last; page++ {
		opts.ListOptions.Page = page

		repos, res, err := client.Repositories.List(ctx, user, opts)
		if err != nil {
			return nil, err
		}

		last = res.LastPage
		for _, repo := range repos {
			list = append(list, repo.GetHTMLURL())
		}
	}
	return list, nil
}

func Fork(
	ctx gogh.Context,
	repo *gogh.Repo,
	organization string,
) (result *gogh.Repo, retErr error) {
	client, err := NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating new client: %w", err)
	}

	/*
		Build GitHub requests.
		See: https://developer.github.com/v3/repos/forks/#parameters-1
			- organization string
				Optional parameter to specify the organization name if forking into an organization.
	*/
	// If the context has no authentication token, specifies context user name for "me".
	if organization == "" && !authenticated(ctx) {
		organization = ctx.GitHubUser()
	}

	opts := &github.RepositoryCreateForkOptions{
		Organization: organization,
	}

	newRepo, _, err := client.Repositories.CreateFork(ctx, repo.Owner(ctx), repo.Name(ctx), opts)
	if newRepo != nil {
		result, retErr = gogh.ParseRepo(newRepo.GetHTMLURL())
	}
	if err != nil {
		retErr = fmt.Errorf("creating fork: %w", err)
	}
	return result, retErr
}

// repository will be created under that org. If the empty string is
// specified, it will be created for the authenticated user.
func Create(
	ctx gogh.Context,
	repo *gogh.Repo,
	description string,
	homepage *url.URL,
	private bool,
) (newRepo *github.Repository, retErr error) {
	client, err := NewClient(ctx)
	if err != nil {
		return nil, err
	}

	// Build request parameters.
	// See: https://developer.github.com/v3/repos/#create
	// Parameters
	// - name	string
	//		Required. The name of the repository.
	// - description	string
	//		A short description of the repository.
	// - homepage	string
	//		A URL with more information about the repository.
	// - private	boolean
	//		Either true to create a private repository or false to create a public one. Creating private repositories requires a paid GitHub account. Default: false
	// - visibility	string
	//		Can be public or private. If your organization is associated with an enterprise account using GitHub Enterprise Cloud, visibility can also be internal. For more information, see "Creating an internal repository" in the GitHub Help documentation.
	//		The visibility parameter overrides the private parameter when you use both parameters with the nebula-preview preview header.
	name := repo.Name(ctx)
	newRepo = &github.Repository{
		Name:        &name,
		Description: &description,
		Private:     &private,
	}
	if homepage != nil {
		page := homepage.String()
		newRepo.Homepage = &page
	}

	owner := repo.ExplicitOwner(ctx)
	if owner == ctx.GitHubUser() {
		owner = ""
	}
	newRepo, _, err = client.Repositories.Create(ctx, owner, newRepo)
	return newRepo, err
}
