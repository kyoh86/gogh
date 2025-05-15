package bundle_dump

import (
	"context"
	"iter"
	"net/url"

	"github.com/kyoh86/gogh/v3/core/git"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/workspace"
)

// UseCase defines the use case for listing repositories
type UseCase struct {
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
	hostingService   hosting.HostingService
	gitService       git.GitService
}

// NewUseCase creates a new instance of UseCase
func NewUseCase(
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
	hostingService hosting.HostingService,
	gitService git.GitService,
) *UseCase {
	return &UseCase{
		workspaceService: workspaceService,
		finderService:    finderService,
		hostingService:   hostingService,
		gitService:       gitService,
	}
}

// BundleEntry represents a repository entry in the bundle
type BundleEntry struct {
	Name  string
	Alias *string
}

type Options = workspace.ListOptions

// Execute retrieves a list of repositories under the specified workspace roots
func (u *UseCase) Execute(ctx context.Context, opts Options) iter.Seq2[*BundleEntry, error] {
	return func(yield func(*BundleEntry, error) bool) {
		for repo, err := range u.finderService.ListAllRepository(ctx, u.workspaceService, opts) {
			if err != nil {
				yield(nil, err)
				return
			}
			if repo == nil {
				continue
			}
			name := repo.Path()
			remotes, err := u.gitService.GetDefaultRemotes(ctx, repo.FullPath())
			if err != nil {
				yield(nil, err)
				return
			}

			for _, remote := range remotes {
				uobj, err := url.Parse(remote)
				if err != nil {
					yield(nil, err)
					return
				}
				if uobj.Host != repo.Host() {
					continue
				}
				ref, err := u.hostingService.ParseURL(uobj)
				if err != nil {
					yield(nil, err)
					return
				}
				remoteName := ref.String()
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
