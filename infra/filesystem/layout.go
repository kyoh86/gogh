package filesystem

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"github.com/kyoh86/gogh/v3/util"
)

// Layout is a filesystem-based standard repository layout implementation
type Layout struct{}

// NewLayout creates a new instance of Layout
func NewLayout() *Layout {
	return &Layout{}
}

// Match returns the reference corresponding to the given path
func (l *Layout) Match(root workspace.Root, path string) (*repository.Reference, error) {
	// ルートからの相対パスを取得
	relPath, err := filepath.Rel(root, path)
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

// PathFor は与えられたリファレンスに対応するパスを返す
func (l *Layout) PathFor(root workspace.Root, ref *repository.Reference) string {
	return filepath.Join(root, ref.Host(), ref.Owner(), ref.Name())
}

// Ensure Layout implements workspace.Layout
var _ workspace.Layout = (*Layout)(nil)
