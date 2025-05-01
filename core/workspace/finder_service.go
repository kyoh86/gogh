package workspace

import (
	"context"

	"github.com/kyoh86/gogh/v3/core/repository"
)

// FindOptions はリポジトリ検索のオプション
type FindOptions struct {
	// Pattern は検索パターン (glob形式、例: **/go-*)
	Pattern string
}

// RepoInfo はリポジトリの情報
type RepoInfo struct {
	// Path はリポジトリの絶対パス
	Path string

	// Reference はリポジトリ参照
	Reference repository.Reference

	// Root はこのリポジトリが含まれるルート
	Root Root

	// IsValid は有効なリポジトリかどうか
	IsValid bool
}

// FinderService はリポジトリ検索サービス
// 実装は外部システム層で行われる
type FinderService interface {
	// FindByReference は参照に一致するリポジトリを検索
	FindByReference(ctx context.Context, ws WorkspaceService, layout Layout, ref repository.Reference) (*RepoInfo, error)

	// FindInRoot は特定のルート下でリポジトリを検索
	FindInRoot(ctx context.Context, root Root, layout Layout, opts FindOptions) ([]RepoInfo, error)

	// FindAll は全ルート下でリポジトリを検索
	FindAll(ctx context.Context, ws WorkspaceService, layout Layout, opts FindOptions) ([]RepoInfo, error)
}
