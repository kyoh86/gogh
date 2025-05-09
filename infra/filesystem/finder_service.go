package filesystem

import (
	"context"
	"errors"
	"iter"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
)

type FinderService struct {
}

func NewFinderService() *FinderService {
	return &FinderService{}
}

type repoRef struct {
	fullPath string
	path     string
	host     string
	owner    string
	name     string
	exists   bool
}

// Exists returns true if the repository exists
func (r *repoRef) Exists() bool { return r.exists }

// Host is a hostname (i.g.: "github.com")
func (r *repoRef) Host() string { return r.host }

// Owner is a owner name (i.g.: "kyoh86")
func (r *repoRef) Owner() string { return r.owner }

// Name of the repository (i.g.: "gogh")
func (r *repoRef) Name() string { return r.name }

// Path returns the path from root of the repository (i.g.: "github.com/kyoh86/gogh")
func (r *repoRef) Path() string { return r.path }

// FullPath returns the full path of the repository (i.g.: "/path/to/workspace/github.com/kyoh86/gogh")
func (r *repoRef) FullPath() string { return r.fullPath }

func (f *FinderService) isDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

// FindByReference implements workspace.FinderService.
func (f *FinderService) FindByReference(ctx context.Context, ws workspace.WorkspaceService, reference repository.Reference) (workspace.RepoInfo, error) {
	for _, root := range ws.GetRoots() {
		layout := ws.GetLayoutFor(root)
		ref, err := layout.Match(reference.String())
		switch {
		case err == nil:
			root := layout.GetRoot()
			path := reference.String()
			abs := filepath.Join(root, path)
			exists, err := f.isDir(abs)
			if err != nil {
				return nil, err
			}
			return &repoRef{
				fullPath: abs,
				path:     path,
				host:     ref.Host(),
				owner:    ref.Owner(),
				name:     ref.Name(),
				exists:   exists,
			}, nil
		case errors.Is(err, workspace.ErrNotMatched):
			// Ignore directories that do not match the layout
		default:
			return nil, err
		}
	}
	return nil, workspace.ErrNotMatched
}

// FindByPath implements workspace.FinderService.
func (f *FinderService) FindByPath(ctx context.Context, ws workspace.WorkspaceService, path string) (workspace.RepoInfo, error) {
	pre := func(root string) (string, bool) {
		return filepath.Join(root, path), true
	}
	if filepath.IsAbs(path) {
		slashPath := filepath.ToSlash(path)
		pre = func(root string) (string, bool) {
			return path, strings.HasPrefix(slashPath, filepath.ToSlash(root)+"/")
		}
	}
	for _, root := range ws.GetRoots() {
		abs, ok := pre(root)
		if !ok {
			continue
		}
		layout := ws.GetLayoutFor(root)
		ref, err := layout.Match(path)
		switch {
		case err == nil:
			exists, err := f.isDir(abs)
			if err != nil {
				return nil, err
			}
			return &repoRef{
				fullPath: abs,
				path:     strings.Join([]string{ref.Host(), ref.Owner(), ref.Name()}, "/"),
				host:     ref.Host(),
				owner:    ref.Owner(),
				name:     ref.Name(),
				exists:   exists,
			}, nil
		case errors.Is(err, workspace.ErrNotMatched):
			// Ignore directories that do not match the layout
		default:
			return nil, err
		}
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
					exists:   true,
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
