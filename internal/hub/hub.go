package hub

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/kyoh86/gogh/gogh"
	"golang.org/x/oauth2"
)

// New builds GitHub Client with GitHub API token that is configured.
func New(env gogh.Env) (*Client, error) {
	if host := env.GithubHost(); host != "" && host != "github.com" {
		url := fmt.Sprintf("https://%s/api/v3", host)
		client, err := github.NewEnterpriseClient(url, url, oauth2Client(env))
		if err != nil {
			return nil, err
		}
		return &Client{client}, nil
	}
	return &Client{github.NewClient(oauth2Client(env))}, nil
}

func authenticated(env gogh.Env) bool {
	return env.GithubToken() != ""
}

func oauth2Client(env gogh.Env) *http.Client {
	if !authenticated(env) {
		return nil
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: env.GithubToken()})
	return oauth2.NewClient(context.TODO(), ts)
}

type Client struct {
	client *github.Client
}

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
func (i *Client) Repos(env gogh.Env, user string, own, collaborate, member bool, visibility, sort, direction string) ([]string, error) {
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
	if user == "" && !authenticated(env) {
		user = "kyoh86" // TODO: cache.GithubUser() or the get 'me' with Github token
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

		repos, res, err := i.client.Repositories.List(context.TODO(), user, opts)
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

// Fork will fork a repository for yours (or for the organization).
func (i *Client) Fork(
	env gogh.Env,
	repo *gogh.Repo,
	organization string,
) (result *gogh.Repo, retErr error) {
	/*
		Build GitHub requests.
		See: https://developer.github.com/v3/repos/forks/#parameters-1
			- organization string
				Optional parameter to specify the organization name if forking into an organization.
	*/
	// If the context has no authentication token, specifies context user name for "me".
	if organization == "" && !authenticated(env) {
		organization = "kyoh86" // TODO: cache.GithubUser() or the get 'me' with Github token
	}

	opts := &github.RepositoryCreateForkOptions{
		Organization: organization,
	}

	newRepo, _, err := i.client.Repositories.CreateFork(context.TODO(), repo.Owner(), repo.Name(), opts)
	if newRepo != nil {
		result, retErr = gogh.ParseRepo(env, newRepo.GetHTMLURL())
	}
	if err != nil {
		retErr = fmt.Errorf("creating fork: %w", err)
	}
	return result, retErr
}

// Create new repository.
func (i *Client) Create(
	env gogh.Env,
	repo *gogh.Repo,
	description string,
	homepage *url.URL,
	private bool,
) (newRepo *github.Repository, retErr error) {
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
	name := repo.Name()
	newRepo = &github.Repository{
		Name:        &name,
		Description: &description,
		Private:     &private,
	}
	if homepage != nil {
		page := homepage.String()
		newRepo.Homepage = &page
	}

	organization := repo.Owner()
	if organization == "kyoh86" { // TODO: cache.GithubUser() or the get 'me' with Github token
		organization = ""
	}
	newRepo, _, err := i.client.Repositories.Create(context.TODO(), organization, newRepo)
	return newRepo, err
}
