package filesystem

import (
	"context"
	"errors"
	"iter"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kyoh86/gogh/v3/core/workspace"
)

type FinderService struct {
}

func NewFinderService() *FinderService {
	return &FinderService{}
}

// FindByPath implements workspace.FinderService.
func (f *FinderService) FindByPath(ctx context.Context, ws workspace.WorkspaceService, path string) (workspace.RepoInfo, error) {
	pre := func(root string) bool {
		return true
	}
	if filepath.IsAbs(path) {
		slashPath := filepath.ToSlash(path)
		pre = func(root string) bool {
			return strings.HasPrefix(slashPath, filepath.ToSlash(root)+"/")
		}
	}
	for _, root := range ws.GetRoots() {
		if !pre(root) {
			continue
		}
		layout := ws.GetLayoutFor(root)
		ref, err := layout.Match(path)
		if err != nil {
			return nil, err
		}
		return &repoRef{
			fullPath: path,
			path:     strings.Join([]string{ref.Host(), ref.Owner(), ref.Name()}, "/"),
			host:     ref.Host(),
			owner:    ref.Owner(),
			name:     ref.Name(),
		}, nil
	}
	return nil, workspace.ErrNotMatched

}

// ListAllRepository implements workspace.FinderService.
func (f *FinderService) ListAllRepository(ctx context.Context, ws workspace.WorkspaceService, opt workspace.ListOptions) iter.Seq2[workspace.RepoInfo, error] {
	return func(yield func(workspace.RepoInfo, error) bool) {
		var i int
		for _, root := range ws.GetRoots() {
			layout := ws.GetLayoutFor(root)
			for ref, err := range f.ListRepositoryInRoot(ctx, layout, opt) {
				if err != nil {
					yield(nil, err)
					return
				}
				if ref == nil {
					continue
				}
				if !yield(ref, nil) {
					return
				}
				i++
				if opt.Limit > 0 && i >= opt.Limit {
					return
				}
			}
		}
	}
}

// ListRepositoryInRoot implements workspace.FinderService.
func (f *FinderService) ListRepositoryInRoot(ctx context.Context, l workspace.LayoutService, opt workspace.ListOptions) iter.Seq2[workspace.RepoInfo, error] {
	var i int
	return func(yield func(workspace.RepoInfo, error) bool) {
		if err := filepath.Walk(l.GetRoot(), func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return nil
			}
			ref, err := l.Match(p)
			switch {
			case errors.Is(err, workspace.ErrNotMatched):
				// Ignore directories that do not match the layout
			case err == nil:
				if !yield(&repoRef{
					fullPath: p,
					path:     path.Join(ref.Host(), ref.Owner(), ref.Name()),
					host:     ref.Host(),
					owner:    ref.Owner(),
					name:     ref.Name(),
				}, nil) {
					return filepath.SkipAll
				}
			default:
				return err
			}
			i++
			if opt.Limit > 0 && i >= opt.Limit {
				return filepath.SkipAll
			}
			return nil
		}); err != nil {
			yield(nil, err)
		}
	}
}

var _ workspace.FinderService = (*FinderService)(nil)
