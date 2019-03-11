package remote

import (
	"strings"

	"github.com/google/go-github/v24/github"
	"github.com/kyoh86/gogh/gogh"
)

// Repo will get a list of repositories for a user.
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
func Repo(ctx gogh.Context, user string, own, collaborate, member bool, visibility, sort, direction string) ([]string, error) {
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
		user = ctx.UserName()
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
