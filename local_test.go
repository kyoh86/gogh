package gogh_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	git "github.com/go-git/go-git/v5"
	"github.com/google/go-cmp/cmp"
	testtarget "github.com/kyoh86/gogh/v3"
)

func mustSpec(t *testing.T, host, owner, name string) testtarget.Spec {
	t.Helper()
	spec, err := testtarget.NewSpec(host, owner, name)
	if err != nil {
		t.Fatalf("invalid spec: %s", err)
	}
	return spec
}

func TestLocalController(t *testing.T) {
	ctx := context.Background()

	root := t.TempDir()
	local := testtarget.NewLocalController(root)

	t.Run("Create", func(t *testing.T) {
		spec := mustSpec(t, "github.com", "kyoh86", "gogh")
		t.Run("Exist", func(t *testing.T) {
			e, err := local.Exist(ctx, spec, nil)
			if err != nil {
				t.Fatalf("failed to create a project: %s", err)
			}
			if e {
				t.Errorf("%q exists", spec)
			}
		})

		t.Run("First", func(t *testing.T) {
			project, err := local.Create(ctx, spec, nil)
			if err != nil {
				t.Fatalf("failed to create a project: %s", err)
			}
			if root != project.Root() {
				t.Errorf("expect root %q but %q is gotten", root, project.Root())
			}
			if spec.Host() != project.Host() {
				t.Errorf("expect host %q but %q is gotten", spec.Host(), project.Host())
			}
			if spec.Owner() != project.Owner() {
				t.Errorf("expect owner %q but %q is gotten", spec.Owner(), project.Owner())
			}
			if spec.Name() != project.Name() {
				t.Errorf("expect name %q but %q is gotten", spec.Name(), project.Name())
			}

			// check built properties
			expectRelPath := "github.com/kyoh86/gogh"
			if expectRelPath != project.RelPath() {
				t.Errorf("expect rel-path %q but %q is gotten", expectRelPath, project.RelPath())
			}
			expectRelFilePath := filepath.Clean("github.com/kyoh86/gogh")
			if expectRelFilePath != project.RelFilePath() {
				t.Errorf(
					"expect rel-path %q but %q is gotten",
					expectRelFilePath,
					project.RelFilePath(),
				)
			}
			expectFullPath := filepath.Join(root, "github.com", "kyoh86", "gogh")
			if expectFullPath != project.FullFilePath() {
				t.Errorf(
					"expect full-path %q but %q is gotten",
					expectFullPath,
					project.FullFilePath(),
				)
			}

			// check git remote
			got, err := local.GetRemoteURLs(ctx, spec, git.DefaultRemoteName)
			if err != nil {
				t.Fatalf("failed to get remote urls from created project: %s", err)
			}
			want := []string{"https://github.com/kyoh86/gogh"}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("remote urls mismatch (-want +got):\n%s", diff)
			}
		})

		t.Run("NotExist", func(t *testing.T) {
			e, err := local.Exist(ctx, spec, nil)
			if err != nil {
				t.Fatalf("failed to create a project: %s", err)
			}
			if !e {
				t.Errorf("%q does not exist", spec)
			}
		})

		t.Run("Duplicated", func(t *testing.T) {
			if _, err := local.Create(ctx, spec, nil); err != git.ErrRepositoryAlreadyExists {
				t.Fatalf(
					"error mismatch: -want +got\n -%v\n +%v",
					git.ErrRepositoryAlreadyExists,
					err,
				)
			}
			if _, err := local.Clone(ctx, spec, "", nil); err != git.ErrRepositoryAlreadyExists {
				t.Fatalf(
					"error mismatch: -want +got\n -%v\n +%v",
					git.ErrRepositoryAlreadyExists,
					err,
				)
			}
		})
	})

	t.Run("PassWalkFnError", func(t *testing.T) {
		expect := errors.New("error for test")
		called := false
		actual := local.Walk(ctx, nil, func(p testtarget.Project) error {
			called = true
			return expect
		})
		if !called {
			t.Fatal("expect that walkFn is called, but not")
		}
		if !errors.Is(actual, expect) {
			t.Fatalf("expect passed through error %v from walkFn, but %v gotten", expect, actual)
		}
	})

	t.Run("SetRemotes", func(t *testing.T) {
		spec := mustSpec(t, "github.com", "kyoh86", "gogh")
		name := "upstream"
		t.Run("First", func(t *testing.T) {
			newSpec := mustSpec(t, "github.com", "kyoh86", "gogh-upstream")
			url := "https://github.com/kyoh86/gogh-upstream"

			if err := local.SetRemoteSpecs(ctx, spec, map[string][]testtarget.Spec{name: {newSpec}}); err != nil {
				t.Fatalf("failed to set remotes: %s", err)
			}
			// check git remote
			got, err := local.GetRemoteURLs(ctx, spec, name)
			if err != nil {
				t.Fatalf("failed to get remote urls from a project which is set remote: %s", err)
			}
			want := []string{url}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("remote urls mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("Overwrite", func(t *testing.T) {
			newSpec := mustSpec(t, "github.com", "kyoh86", "gogh-overwrite")
			url := "https://github.com/kyoh86/gogh-overwrite"

			if err := local.SetRemoteSpecs(ctx, spec, map[string][]testtarget.Spec{name: {newSpec}}); err != nil {
				t.Fatalf("failed to set remotes: %s", err)
			}
			// check git remote
			got, err := local.GetRemoteURLs(ctx, spec, name)
			if err != nil {
				t.Fatalf("failed to get remote urls from a project which is set remote: %s", err)
			}
			want := []string{url}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("remote urls mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("NotFound", func(t *testing.T) {
			if err := local.SetRemoteURLs(ctx, mustSpec(t, "github.com", "kyoh86", "unknown"), nil); err == nil {
				t.Error("expect that SetRemoteURLs is failed, but not")
			}
			if _, err := local.GetRemoteURLs(ctx, mustSpec(t, "github.com", "kyoh86", "unknown"), git.DefaultRemoteName); err == nil {
				t.Error("expect that GetRemoteURLs is failed, but not")
			}
			if _, err := local.GetRemoteURLs(ctx, spec, "unknown"); err == nil {
				t.Error("expect that GetRemoteURLs is failed, but not")
			}
		})
	})

	t.Run("List", func(t *testing.T) {
		// create noise
		// file
		if err := os.WriteFile(filepath.Join(root, "github.com", "kyoh86", "file"), nil, 0644); err != nil {
			t.Fatalf("failed to create dummy file: %s", err)
		}
		// invalid name
		invalidPath := filepath.Join(root, "github.com", "kyoh86", "invalid name")
		if err := os.MkdirAll(invalidPath, 0755); err != nil {
			t.Fatalf("failed to create dummy directory: %s", err)
		}
		_, err := git.PlainInit(invalidPath, false)
		if err != nil {
			t.Fatalf("failed to init dummy repository: %s", err)
		}

		expect := "github.com/kyoh86/gogh"

		// match cases
		for _, testcase := range []struct {
			title  string
			option *testtarget.LocalListOption
		}{
			{
				title:  "nil",
				option: nil,
			},
			{
				title:  "empty",
				option: &testtarget.LocalListOption{Query: ""},
			},
			{
				title:  "matched for name",
				option: &testtarget.LocalListOption{Query: "gogh"},
			},
			{
				title:  "matched for owner",
				option: &testtarget.LocalListOption{Query: "kyoh86"},
			},
			{
				title:  "matched for owner/name",
				option: &testtarget.LocalListOption{Query: "kyoh86/gogh"},
			},
			{
				title:  "matched for owner/name",
				option: &testtarget.LocalListOption{Query: "kyoh86/gogh"},
			},
		} {
			t.Run(testcase.title, func(t *testing.T) {
				actual, err := local.List(ctx, testcase.option)
				if err != nil {
					t.Fatalf("failed to get project list: %s", err)
				}
				if len(actual) != 1 {
					t.Fatalf(
						"expect just one project is matched, but actual %d matched",
						len(actual),
					)
				}
				for _, act := range actual {
					if expect != act.RelPath() {
						t.Errorf("expect that %q is matched but actual: %q", expect, act.RelPath())
					}
				}
			})
		}

		// unmatch case
		t.Run("Unmatch", func(t *testing.T) {
			actual, err := local.List(ctx, &testtarget.LocalListOption{Query: "dummy"})
			if err != nil {
				t.Fatalf("failed to get project list: %s", err)
			}
			if len(actual) != 0 {
				t.Errorf("expect that no project matched, but %d projects are gotten", len(actual))
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("Valid", func(t *testing.T) {
			if err := local.Delete(ctx, mustSpec(t, "github.com", "kyoh86", "gogh"), nil); err != nil {
				t.Fatalf("failed to remove project: %s", err)
			}
			stat, err := os.Stat(filepath.Join(root, "github.com", "kyoh86", "gogh"))
			if !os.IsNotExist(err) {
				t.Fatalf("failed to remove instance: %+v", stat)
			}
		})
		// create noise
		// file
		if err := os.WriteFile(filepath.Join(root, "github.com", "kyoh86", "file"), nil, 0644); err != nil {
			t.Fatalf("failed to create dummy file: %s", err)
		}
		// unmanaged root
		otherRoot := t.TempDir()
		if err := os.MkdirAll(filepath.Join(otherRoot, "github.com", "kyoh86", "gogh"), 0755); err != nil {
			t.Fatalf("failed to create unmanaged project: %s", err)
		}
		for _, testcase := range []struct {
			title string
			spec  testtarget.Spec
		}{
			{
				title: "not exist",
				spec:  mustSpec(t, "github.com", "kyoh86", "not-exist"),
			},
			{
				title: "instance is not a dir",
				spec:  mustSpec(t, "github.com", "kyoh86", "file"),
			},
		} {
			t.Run(testcase.title, func(t *testing.T) {
				actual := local.Delete(ctx, testcase.spec, nil)
				if actual == nil {
					t.Errorf("expect error when the spec %s, but not", testcase.title)
				}
			})
		}
	})

	t.Run("Clone", func(t *testing.T) {
		spec := mustSpec(t, "github.com", "kyoh86-tryouts", "bare")
		project, err := local.Clone(ctx, spec, "", nil)
		if err != nil {
			t.Fatalf("failed to clone a project: %s", err)
		}
		if root != project.Root() {
			t.Errorf("expect root %q but %q is gotten", root, project.Root())
		}
		if spec.Host() != project.Host() {
			t.Errorf("expect host %q but %q is gotten", spec.Host(), project.Host())
		}
		if spec.Owner() != project.Owner() {
			t.Errorf("expect owner %q but %q is gotten", spec.Owner(), project.Owner())
		}
		if spec.Name() != project.Name() {
			t.Errorf("expect name %q but %q is gotten", spec.Name(), project.Name())
		}

		// check built properties
		expectRelPath := "github.com/kyoh86-tryouts/bare"
		if expectRelPath != project.RelPath() {
			t.Errorf("expect rel-path %q but %q is gotten", expectRelPath, project.RelPath())
		}
		expectRelFilePath := filepath.Clean("github.com/kyoh86-tryouts/bare")
		if expectRelFilePath != project.RelFilePath() {
			t.Errorf(
				"expect rel-path %q but %q is gotten",
				expectRelFilePath,
				project.RelFilePath(),
			)
		}
		expectFullPath := filepath.Join(root, "github.com/kyoh86-tryouts/bare")
		if expectFullPath != project.FullFilePath() {
			t.Errorf("expect full-path %q but %q is gotten", expectFullPath, project.FullFilePath())
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
		default:
			t.Fatalf("cloned repository has multiple urls: %+v", urls)
			fallthrough
		case 1:
			expectURL := "https://github.com/kyoh86-tryouts/bare"
			if expectURL != urls[0] {
				t.Errorf("expect the repository cloned for %q but %q actually", expectURL, urls[0])
			}
		case 0:
			t.Fatal("cloned repository has no url")
		}
	})

	t.Run("Alias", func(t *testing.T) {
		spec := mustSpec(t, "github.com", "kyoh86-tryouts", "bare")
		alias := mustSpec(t, "example.com", "kyoh86", "alias")
		project, err := local.Clone(ctx, spec, "", &testtarget.LocalCloneOption{
			Alias: &alias,
		})
		if err != nil {
			t.Fatalf("failed to clone a project: %s", err)
		}
		if root != project.Root() {
			t.Errorf("expect root %q but %q is gotten", root, project.Root())
		}
		if spec.Host() != project.Host() {
			t.Errorf("expect host %q but %q is gotten", spec.Host(), project.Host())
		}
		if alias.Owner() != project.Owner() {
			t.Errorf("expect owner %q but %q is gotten", alias.Owner(), project.Owner())
		}
		if alias.Name() != project.Name() {
			t.Errorf("expect name %q but %q is gotten", alias.Name(), project.Name())
		}

		// check built properties
		wantRelPath := "github.com/kyoh86/alias"
		if wantRelPath != project.RelPath() {
			t.Errorf("want rel-path %q but %q is gotten", wantRelPath, project.RelPath())
		}
		wantRelFilePath := filepath.Clean("github.com/kyoh86/alias")
		if wantRelFilePath != project.RelFilePath() {
			t.Errorf("want rel-path %q but %q is gotten", wantRelFilePath, project.RelFilePath())
		}
		wantFullPath := filepath.Join(root, "github.com/kyoh86/alias")
		if wantFullPath != project.FullFilePath() {
			t.Errorf("want full-path %q but %q is gotten", wantFullPath, project.FullFilePath())
		}

		// check git remote
		repo, err := git.PlainOpen(wantFullPath)
		if err != nil {
			t.Fatalf("failed to open git repository in cloned project: %s", err)
		}
		remote, err := repo.Remote(git.DefaultRemoteName)
		if err != nil {
			t.Fatalf("failed to get remote %s: %s", git.DefaultRemoteName, err)
		}
		urls := remote.Config().URLs
		switch len(urls) {
		default:
			t.Fatalf("cloned repository has multiple urls: %+v", urls)
			fallthrough
		case 1:
			sourceURL := "https://github.com/kyoh86-tryouts/bare"
			if sourceURL != urls[0] {
				t.Errorf("expect the repository cloned for %q but %q actually", sourceURL, urls[0])
			}
		case 0:
			t.Fatal("cloned repository has no url")
		}
	})

	t.Run("CloneFailureWithInvalidToken", func(t *testing.T) {
		spec := mustSpec(t, "github.com", "kyoh86", "gogh")
		if _, err := local.Clone(ctx, spec, "invalid-token", nil); err == nil {
			t.Fatalf("expect failure to clone a project: %s", err)
		}
	})
}

func TestLocalControllerWithUnaccessableRoot(t *testing.T) {
	ctx := context.Background()

	tmp := t.TempDir()
	root := filepath.Join(tmp, "root")

	spec := mustSpec(t, "example.com", "kyoh86", "gogh")
	local := testtarget.NewLocalController(root)

	t.Run("NotExit", func(t *testing.T) {
		if _, err := local.List(ctx, nil); err != nil {
			t.Fatalf("failed to list not found root: %s", err)
		}
	})

	t.Run("NotWritable", func(t *testing.T) {
		// prepare a file for the root of the test
		if err := os.WriteFile(root, nil, 0644); err != nil {
			t.Fatalf("failed to prepare dummy file: %s", err)
		}

		if _, err := local.Create(ctx, spec, nil); err == nil {
			t.Errorf("expect failure to create")
		}
		if _, err := local.Clone(ctx, spec, "", nil); err == nil {
			t.Errorf("expect failure to clone")
		}
		if err := local.Delete(ctx, spec, nil); err == nil {
			t.Errorf("expect failure to remove")
		}
	})
}
