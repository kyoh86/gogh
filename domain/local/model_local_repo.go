package local

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/kyoh86/gogh/v3/domain/reporef"
)

// LocalRepo is the location of a repository in the local.
// It is a valid location, that never means "exist".
type LocalRepo struct {
	root string
	ref  reporef.RepoRef
}

func (l LocalRepo) Root() string {
	return l.root
}

func (l LocalRepo) Host() string  { return l.ref.Host() }
func (l LocalRepo) Owner() string { return l.ref.Owner() }
func (l LocalRepo) Name() string  { return l.ref.Name() }

func (l LocalRepo) FullLevels() []string {
	return []string{l.root, l.ref.Host(), l.ref.Owner(), l.ref.Name()}
}

func (l LocalRepo) RelLevels() []string {
	return l.ref.RelLevels()
}

func (l LocalRepo) FullFilePath() string {
	return filepath.Join(l.FullLevels()...)
}

func (l LocalRepo) RelFilePath() string {
	return filepath.Join(l.RelLevels()...)
}

func (l LocalRepo) RelPath() string {
	return path.Join(l.RelLevels()...)
}

// CheckEntity checks if the local repository is exist in the local file-system.
func (l LocalRepo) CheckEntity() error {
	path := l.FullFilePath()
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return errors.New("local repository is not dir")
	}
	return nil
}

// UNDONE: CheckEntityInFileSystem() support fs.FS

func NewLocalRepo(root string, ref reporef.RepoRef) LocalRepo {
	return LocalRepo{root: root, ref: ref}
}
