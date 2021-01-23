package gogh_test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	git "github.com/go-git/go-git/v5"
	testtarget "github.com/kyoh86/gogh/v2"
)

func description(t *testing.T, host, user, name string) testtarget.Description {
	t.Helper()
	d, err := testtarget.ValidateDescription(host, user, name)
	if err != nil {
		t.Fatalf("invalid description: %s", err)
	}
	return *d
}
func TestLocalController(t *testing.T) {
	ctx := context.Background()

	root := t.TempDir()
	local := testtarget.NewLocalController(ctx, root)

	t.Run("Create", func(t *testing.T) {
		d := description(t, "github.com", "kyoh86", "gogh")
		t.Run("First", func(t *testing.T) {
			project, err := local.Create(ctx, d, nil)
			if err != nil {
				t.Fatalf("failed to create a project: %s", err)
			}
			if root != project.Root() {
				t.Errorf("expect root %q but %q is gotten", root, project.Root())
			}
			if d.Host() != project.Host() {
				t.Errorf("expect host %q but %q is gotten", d.Host(), project.Host())
			}
			if d.User() != project.User() {
				t.Errorf("expect user %q but %q is gotten", d.User(), project.User())
			}
			if d.Name() != project.Name() {
				t.Errorf("expect name %q but %q is gotten", d.Name(), project.Name())
			}

			// check built properties
			expectRelPath := "github.com/kyoh86/gogh"
			if expectRelPath != project.RelPath() {
				t.Errorf("expect rel-path %q but %q is gotten", expectRelPath, project.RelPath())
			}
			expectURL := "https://github.com/kyoh86/gogh"
			if expectURL != project.URL() {
				t.Errorf("expect url %q but %q is gotten", expectURL, project.URL())
			}
			expectFullPath := filepath.Join(root, "github.com/kyoh86/gogh")
			if expectFullPath != project.FullPath() {
				t.Errorf("expect full-path %q but %q is gotten", expectFullPath, project.FullPath())
			}

			// check git remote
			repo, err := git.PlainOpen(expectFullPath)
			if err != nil {
				t.Fatalf("failed to open git repository in created project: %s", err)
			}
			remote, err := repo.Remote(git.DefaultRemoteName)
			if err != nil {
				t.Fatalf("failed to get remote %s: %s", git.DefaultRemoteName, err)
			}
			urls := remote.Config().URLs
			switch len(urls) {
			case 0:
				t.Fatal("created repository has no url")
			default:
				t.Fatalf("created repository has multiple urls: %+v", urls)
				fallthrough
			case 1:
				if expectURL != urls[0] {
					t.Errorf("expect the repository created for %q but %q actually", expectURL, urls[0])
				}
			}
		})

		t.Run("Duplicated", func(t *testing.T) {
			if _, err := local.Create(ctx, d, nil); err == nil {
				t.Fatalf("expect failure with creating a project that has already exist: %s", err)
			}
			if _, err := local.Clone(ctx, d, nil); err == nil {
				t.Fatalf("expect failure with cloning a project that has already exist: %s", err)
			}
		})
	})

	t.Run("List", func(t *testing.T) {
		// create noise
		// file
		if err := ioutil.WriteFile(filepath.Join(root, "github.com", "kyoh86", "file"), nil, 0644); err != nil {
			t.Fatalf("failed to create dummy file: %s", err)
		}
		// invalid name
		if err := os.MkdirAll(filepath.Join(root, "github.com", "kyoh86", "invalid name"), 0755); err != nil {
			t.Fatalf("failed to create dummy directory: %s", err)
		}

		expect := "github.com/kyoh86/gogh"

		// cases
		for _, testcase := range []struct {
			title string
			query string
		}{
			{
				title: "empty",
				query: "",
			},
			{
				title: "matched for name",
				query: "gogh",
			},
			{
				title: "matched for user",
				query: "kyoh86",
			},
			{
				title: "matched for user/name",
				query: "kyoh86/gogh",
			},
			{
				title: "matched for user/name",
				query: "kyoh86/gogh",
			},
		} {
			t.Run(testcase.title, func(t *testing.T) {
				actual, err := local.List(ctx, testcase.query)
				if err != nil {
					t.Fatalf("failed to get project list: %s", err)
				}
				if len(actual) != 1 {
					t.Fatalf("expect just one project is matched, but actual %d matched", len(actual))
				}
				for _, act := range actual {
					if expect != act.RelPath() {
						t.Errorf("expect that %q is matched but actual: %q", expect, act.RelPath())
					}
				}
			})
		}
	})

	t.Run("Remove", func(t *testing.T) {
		t.Run("Valid", func(t *testing.T) {
			if err := local.Remove(ctx, description(t, "github.com", "kyoh86", "gogh")); err != nil {
				t.Fatalf("failed to remove project: %s", err)
			}
			stat, err := os.Stat(filepath.Join(root, "github.com", "kyoh86", "gogh"))
			if !os.IsNotExist(err) {
				t.Fatalf("failed to remove instance: %+v", stat)
			}
		})
		// create noise
		// file
		if err := ioutil.WriteFile(filepath.Join(root, "github.com", "kyoh86", "file"), nil, 0644); err != nil {
			t.Fatalf("failed to create dummy file: %s", err)
		}
		// unmanaged root
		otherRoot := t.TempDir()
		if err := os.MkdirAll(filepath.Join(otherRoot, "github.com", "kyoh86", "gogh"), 0755); err != nil {
			t.Fatalf("failed to create unmanaged project: %s", err)
		}
		for _, testcase := range []struct {
			title       string
			description testtarget.Description
		}{
			{
				title:       "not exist",
				description: description(t, "github.com", "kyoh86", "not-exist"),
			},
			{
				title:       "instance is not a dir",
				description: description(t, "github.com", "kyoh86", "file"),
			},
		} {
			t.Run(testcase.title, func(t *testing.T) {
				actual := local.Remove(ctx, testcase.description)
				if actual == nil {
					t.Errorf("expect error when the description %s, but not", testcase.title)
				}
			})
		}
	})

	t.Run("Clone", func(t *testing.T) {
		t.Skip("skip remote access")
		d := description(t, "github.com", "kyoh86", "gogh")
		project, err := local.Clone(ctx, d, nil)
		if err != nil {
			t.Fatalf("failed to clone a project: %s", err)
		}
		if root != project.Root() {
			t.Errorf("expect root %q but %q is gotten", root, project.Root())
		}
		if d.Host() != project.Host() {
			t.Errorf("expect host %q but %q is gotten", d.Host(), project.Host())
		}
		if d.User() != project.User() {
			t.Errorf("expect user %q but %q is gotten", d.User(), project.User())
		}
		if d.Name() != project.Name() {
			t.Errorf("expect name %q but %q is gotten", d.Name(), project.Name())
		}

		// check built properties
		expectRelPath := "github.com/kyoh86/gogh"
		if expectRelPath != project.RelPath() {
			t.Errorf("expect rel-path %q but %q is gotten", expectRelPath, project.RelPath())
		}
		expectURL := "https://github.com/kyoh86/gogh"
		if expectURL != project.URL() {
			t.Errorf("expect url %q but %q is gotten", expectURL, project.URL())
		}
		expectFullPath := filepath.Join(root, "github.com/kyoh86/gogh")
		if expectFullPath != project.FullPath() {
			t.Errorf("expect full-path %q but %q is gotten", expectFullPath, project.FullPath())
		}

		// check git remote
		repo, err := git.PlainOpen(expectFullPath)
		if err != nil {
			t.Fatalf("failed to open git repository in cloned project: %s", err)
		}
		remote, err := repo.Remote(git.DefaultRemoteName)
		if err != nil {
			t.Fatalf("failed to get remote %s: %s", git.DefaultRemoteName, err)
		}
		urls := remote.Config().URLs
		switch len(urls) {
		case 0:
			t.Fatal("cloned repository has no url")
		default:
			t.Fatalf("cloned repository has multiple urls: %+v", urls)
			fallthrough
		case 1:
			if expectURL != urls[0] {
				t.Errorf("expect the repository cloned for %q but %q actually", expectURL, urls[0])
			}
		}
	})

	// UNDONE: func (f *ProjectFormatter) Format(ctx context.Context, project Project) (string, error)
	// UNDONE:  - fullpath, relpath, url, fields(fullpath url host user name), json
}
