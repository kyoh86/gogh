package repo

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	git "gopkg.in/src-d/go-git.v4"
)

type Repository struct {
	wd         string
	root       string
	repository *git.Repository
}

func Open(wd string) (*Repository, error) {
	root, err := rootDir(wd)
	if err != nil {
		return nil, err
	}
	r, err := git.PlainOpen(root)
	if err != nil {
		return nil, err
	}
	return &Repository{wd: wd, root: root, repository: r}, nil
}

// Root will get a git directory
func (r *Repository) Root() string {
	return r.root
}

func rootDir(wd string) (string, error) {
	var needle = wd
	for {
		parent, name := filepath.Split(needle)
		if parent == needle {
			break
		}
		parent = strings.TrimRight(parent, string([]rune{filepath.Separator}))
		if name == ".git" {
			return parent, nil
		}

		_, err := os.Stat(filepath.Join(needle, ".git"))
		if os.IsNotExist(err) {
			needle = parent
			continue
		}
		if err != nil {
			return "", errors.Wrap(err, "stat current directory")
		}
		return needle, nil
	}
	return "", errors.New("not a git repository (or any of the parent directories)")
}
