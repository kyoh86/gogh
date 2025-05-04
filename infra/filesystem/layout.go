package filesystem

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"github.com/kyoh86/gogh/v3/util"
)

// Layout is a filesystem-based standard repository layout implementation
type Layout struct {
	root workspace.Root
}

// NewLayout creates a new instance of Layout
func NewLayout(root workspace.Root) *Layout {
	return &Layout{root: root}
}

// Match returns the reference corresponding to the given path
func (l *Layout) Match(path string) (*repository.Reference, error) {
	// ルートからの相対パスを取得
	relPath, err := filepath.Rel(l.root, path)
	if err != nil {
		return nil, err
	}

	// パスコンポーネントを分解
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	if len(parts) < 3 {
		return nil, errors.New("path does not match repository layout")
	}

	// host/owner/nameの形式でリファレンスを作成
	return util.Ptr(repository.NewReference(parts[0], parts[1], parts[2])), nil
}

func (l *Layout) PathFor(ref repository.Reference) string {
	return filepath.Join(l.root, ref.Host(), ref.Owner(), ref.Name())
}

func (l *Layout) CreateRepositoryFolder(ref repository.Reference) (string, error) {
	path := l.PathFor(ref)
	return path, os.MkdirAll(path, 0755)
}

func (l *Layout) DeleteRepository(ref repository.Reference) error {
	path := l.PathFor(ref)
	return os.RemoveAll(path)
}

// Ensure Layout implements workspace.Layout
var _ workspace.Layout = (*Layout)(nil)
