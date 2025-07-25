package filesystem

import (
	"context"
	"errors"
	"iter"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

type FinderService struct{}

func NewFinderService() *FinderService {
	return &FinderService{}
}

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
func (f *FinderService) FindByReference(ctx context.Context, ws workspace.WorkspaceService, ref repository.Reference) (*repository.Location, error) {
	for _, root := range ws.GetRoots() {
		abs := filepath.Join(root, ref.String())
		isDir, err := f.isDir(abs)
		if err != nil {
			return nil, err
		}
		if isDir {
			return repository.NewLocation(
				abs,
				ref.Host(),
				ref.Owner(),
				ref.Name(),
			), nil
		}
	}
	return nil, workspace.ErrNotMatched
}

// FindByPath implements workspace.FinderService.
func (f *FinderService) FindByPath(ctx context.Context, ws workspace.WorkspaceService, path string) (*repository.Location, error) {
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
		ref, err := layout.Match(abs)
		switch {
		case err == nil:
			isDir, err := f.isDir(abs)
			if err != nil {
				return nil, err
			}
			if isDir {
				return repository.NewLocation(
					abs,
					ref.Host(),
					ref.Owner(),
					ref.Name(),
				), nil
			}
		case errors.Is(err, workspace.ErrNotMatched):
			// Ignore directories that do not match the layout
		default:
			return nil, err
		}
	}
	return nil, workspace.ErrNotMatched
}

// ListAllRepository implements workspace.FinderService.
func (f *FinderService) ListAllRepository(ctx context.Context, ws workspace.WorkspaceService, opts workspace.ListOptions) iter.Seq2[*repository.Location, error] {
	return func(yield func(*repository.Location, error) bool) {
		var i int
		for _, root := range ws.GetRoots() {
			layout := ws.GetLayoutFor(root)
			for ref, err := range f.ListRepositoryInRoot(ctx, layout, opts) {
				if err != nil {
					yield(nil, err)
					return
				}
				if !yield(ref, nil) {
					return
				}
				i++
				if opts.Limit > 0 && i >= opts.Limit {
					return
				}
			}
		}
	}
}

// ListRepositoryInRoot implements workspace.FinderService.
func (f *FinderService) ListRepositoryInRoot(ctx context.Context, l workspace.LayoutService, opts workspace.ListOptions) iter.Seq2[*repository.Location, error] {
	var i int
	return func(yield func(*repository.Location, error) bool) {
		if err := filepath.Walk(l.GetRoot(), func(p string, info os.FileInfo, err error) error {
			if os.IsNotExist(err) {
				return nil
			}
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return nil
			}
			ref, err := l.ExactMatch(p)
			switch {
			case errors.Is(err, workspace.ErrNotMatched):
				// Ignore directories that do not match the layout
			case err == nil:
				if ref == nil {
					return filepath.SkipDir
				}
				match, err := matchPattern(opts.Patterns, *ref)
				if err != nil {
					return err
				}
				if !match {
					return filepath.SkipDir
				}
				if !yield(repository.NewLocation(
					p,
					ref.Host(),
					ref.Owner(),
					ref.Name(),
				), nil) {
					return filepath.SkipAll
				}
				return filepath.SkipDir
			default:
				return err
			}
			i++
			if opts.Limit > 0 && i >= opts.Limit {
				return filepath.SkipAll
			}
			return nil
		}); err != nil {
			yield(nil, err)
			return
		}
	}
}

func matchPattern(patterns []string, ref repository.Reference) (bool, error) {
	if len(patterns) == 0 {
		return true, nil
	}
	for _, pattern := range patterns {
		if match, err := doublestar.Match(pattern, ref.String()); err != nil {
			return false, err
		} else if match {
			return true, nil
		}
	}
	return false, nil
}

var _ workspace.FinderService = (*FinderService)(nil)
