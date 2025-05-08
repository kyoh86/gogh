package bundle_dump

import (
	"context"
	"iter"
	"net/url"
	"strings"

	"github.com/kyoh86/gogh/v3/core/workspace"
	"github.com/kyoh86/gogh/v3/infra/git"
	"github.com/kyoh86/gogh/v3/util"
)

// UseCase defines the use case for listing repositories
type UseCase struct {
	workspaceService workspace.WorkspaceService
}

// NewUseCase creates a new instance of UseCase
func NewUseCase(
	workspaceService workspace.WorkspaceService,
) *UseCase {
	return &UseCase{
		workspaceService: workspaceService,
	}
}

// BundleEntry represents a repository entry in the bundle
type BundleEntry struct {
	Name  string
	Alias *string
}

// Execute retrieves a list of repositories under the specified workspace roots
func (u *UseCase) Execute(ctx context.Context) iter.Seq2[*BundleEntry, error] {
	ws := u.workspaceService
	return func(yield func(*BundleEntry, error) bool) {
		gitService := git.NewService()
		for _, root := range ws.GetRoots() {
			layout := ws.GetLayoutFor(root)
			for repo, err := range layout.ListRepository(ctx, 0) {
				if err != nil {
					yield(nil, err)
					return
				}
				if repo == nil {
					continue
				}
				name := repo.Path()
				remotes, err := gitService.GetDefaultRemotes(ctx, repo.FullPath())
				if err != nil {
					yield(nil, err)
					return
				}

				for _, remote := range remotes {
					//TODO: Parse URL with hosting service?
					uobj, err := url.Parse(remote)
					if err != nil {
						yield(nil, err)
						return
					}
					if uobj.Host != repo.Host() {
						continue
					}
					remoteName := strings.Join([]string{uobj.Host, strings.TrimPrefix(strings.TrimSuffix(uobj.Path, ".git"), "/")}, "/")
					entry := &BundleEntry{
						Name: name,
					}
					if remoteName != name {
						entry.Alias = &remoteName
					}
					if !yield(entry, nil) {
						return
					}
					break
				}
			}
		}
	}
}
