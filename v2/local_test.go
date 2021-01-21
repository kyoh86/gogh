package gogh_test

import (
	"context"
	"path/filepath"
	"testing"

	git "github.com/go-git/go-git/v5"
	testtarget "github.com/kyoh86/gogh/v2"
)

func TestLocal(t *testing.T) {
	ctx := context.Background()

	t.Run("Empty", func(t *testing.T) {
		root := t.TempDir()
		local, err := testtarget.NewLocal(ctx)
		if err != nil {
			t.Fatalf("failed to create new local: %q", err)
		}
		local.SetRoot(root)

		t.Run("Roots", func(t *testing.T) {
			expect := root
			roots := local.Roots()
			l := len(roots)
			switch l {
			case 0:
				t.Fatal("expect not empty roots")
			default:
				t.Errorf("expect single root but %d roots are gotten", l)
				fallthrough
			case 1:
				if actual := roots[0]; actual != expect {
					t.Errorf("expect root %q but %q are gotten", expect, actual)
				}
			}
		})

		t.Run("Create", func(t *testing.T) {
			d := testtarget.Description{
				Host: "github.com",
				User: "kyoh86",
				Name: "gogh",
			}
			t.Run("First", func(t *testing.T) {
				project, err := local.Create(ctx, d)
				if err != nil {
					t.Fatalf("failed to create a project: %s", err)
				}
				if root != project.Root {
					t.Errorf("expect root %q but %q is gotten", root, project.Root)
				}
				if d.Host != project.Host {
					t.Errorf("expect host %q but %q is gotten", d.Host, project.Host)
				}
				if d.User != project.User {
					t.Errorf("expect user %q but %q is gotten", d.User, project.User)
				}
				if d.Name != project.Name {
					t.Errorf("expect name %q but %q is gotten", d.Name, project.Name)
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
				_, err := local.Create(ctx, d)
				if err == nil {
					t.Fatalf("expect failure with creating a project that has already exist: %s", err)
				}
			})
		})

		// UNDONE: func (l *Local) List(ctx context.Context, params *LocalListParam) ([]Project, error)
		// UNDONE: func (l *Local) LocalWalk(ctx context.Context, walker LocalWalkFunc) error
		// UNDONE: func (l *Local) LocalWalkInPrimary(ctx context.Context, walker LocalWalkFunc) error
		// UNDONE: func (l *Local) Clone(ctx context.Context, description Description) (Project, error)
		// UNDONE: func (l *Local) Remove(ctx context.Context, project Project) error
	})
}
