package hub

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/kyoh86/gogh/env"
	"github.com/kyoh86/gogh/gogh"
	keyring "github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

// New builds GitHub Client with GitHub API token that is configured.
func New(authContext context.Context, ev gogh.Env) (*Client, error) {
	if host := ev.GithubHost(); host != "" && host != "github.com" {
		url := fmt.Sprintf("https://%s/api/v3", host)
		httpClient, err := oauth2Client(authContext, ev)
		if err != nil {
			return nil, err
		}
		client, err := github.NewEnterpriseClient(url, url, httpClient)
		if err != nil {
			return nil, err
		}
		return &Client{client}, nil
	}
	httpClient, err := oauth2Client(authContext, ev)
	if err != nil {
		return nil, err
	}
	return &Client{github.NewClient(httpClient)}, nil
}

func getToken(ev gogh.Env) (string, error) {
	if ev.GithubUser() == "" {
		return "", errors.New("github.user is necessary to access GitHub")
	}
	envar := os.Getenv("GOGH_GITHUB_TOKEN")
	if envar != "" {
		return envar, nil
	}
	return keyring.Get(strings.Join([]string{ev.GithubHost(), env.KeyringService}, "."), ev.GithubUser())
}

func oauth2Client(authContext context.Context, ev gogh.Env) (*http.Client, error) {
	token, err := getToken(ev)
	if err != nil {
		return nil, err
	}
	if token == "" {
		return nil, errors.New("github.token is necessary to access GitHub")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return oauth2.NewClient(authContext, ts), nil
}

type Client struct {
	client *github.Client
}

// Repos will get a list of repositories for a user.
// Parameters:
//   * user:        Who has the repositories. Empty means the token user
//   * own:         Include repositories that are owned by the user
//   * collaborate: Include repositories that the user has been added to as a collaborator
//   * member:      Include repositories that the user has access to through being a member of an organization. This includes every repository on every team that the user is on
//   * visibility:  Can be one of all, public, or private
//   * sort:        Can be one of created, updated, pushed, full_name
//   * direction:   Can be one of asc or desc default. Default means asc when using full_name, otherwise desc
// Returns:
//   List of the url for repoisitories
func (i *Client) Repos(ctx context.Context, ev gogh.Env, user string, own, collaborate, member bool, visibility, sort, direction string) ([]string, error) {
	/*
		Build GitHub requests.
		See: https://developer.github.com/v3/repos/#parameters
		- affiliation string  Comma-separated list of values. Can include:
				- owner: Repositories that are owned by the token user.
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

		repos, res, err := i.client.Repositories.List(ctx, user, opts)
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
	ctx context.Context,
	ev gogh.Env,
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
	if organization == "" {
		token, err := getToken(ev)
		if err == nil && token != "" {
			organization = ev.GithubUser()
		}
	}

	opts := &github.RepositoryCreateForkOptions{
		Organization: organization,
	}

	newRepo, _, err := i.client.Repositories.CreateFork(ctx, repo.Owner(), repo.Name(), opts)
	if newRepo != nil {
		result, retErr = gogh.ParseRepo(ev, newRepo.GetHTMLURL())
	}
	if err != nil {
		retErr = fmt.Errorf("creating fork: %w", err)
	}
	return result, retErr
}

// Create new repository.
func (i *Client) Create(
	ctx context.Context,
	ev gogh.Env,
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
	if organization == ev.GithubUser() {
		organization = ""
	}
	newRepo, _, err := i.client.Repositories.Create(ctx, organization, newRepo)
	return newRepo, err
}
