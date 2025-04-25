package local

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/kyoh86/gogh/v3/domain/reporef"
)

// Repo is the location of a repository in the local.
// It is a valid location, that never means "exist".
type Repo struct {
	root string
	ref  reporef.RepoRef
}

func (l Repo) Root() string {
	return l.root
}

func (l Repo) Host() string  { return l.ref.Host() }
func (l Repo) Owner() string { return l.ref.Owner() }
func (l Repo) Name() string  { return l.ref.Name() }

func (l Repo) FullLevels() []string {
	return []string{l.root, l.ref.Host(), l.ref.Owner(), l.ref.Name()}
}

func (l Repo) RelLevels() []string {
	return l.ref.RelLevels()
}

func (l Repo) FullFilePath() string {
	return filepath.Join(l.FullLevels()...)
}

func (l Repo) RelFilePath() string {
	return filepath.Join(l.RelLevels()...)
}

func (l Repo) RelPath() string {
	return path.Join(l.RelLevels()...)
}

// CheckEntity checks if the local repository is exist in the local file-system.
func (l Repo) CheckEntity() error {
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

func NewRepo(root string, ref reporef.RepoRef) Repo {
	return Repo{root: root, ref: ref}
}
