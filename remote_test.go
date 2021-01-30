package gogh_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	testtarget "github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/kyoh86/gogh/v2/internal/github_mock"
	"github.com/wacul/ptr"
)

func MockAdaptor(t *testing.T) (*github_mock.MockAdaptor, func()) {
	ctrl := gomock.NewController(t)
	mock := github_mock.NewMockAdaptor(ctrl)
	return mock, ctrl.Finish
}

func TestRemoteController(t *testing.T) {
	ctx := context.Background()

	// TODO: authorized server
	// TODO: create github.com/kyoh86/gogh -> github.com/kyoh86/gogh created
	// TODO: remove github.com/kyoh86/gogh -> github.com/kyoh86/gogh removed

	host := testtarget.DefaultHost
	user := "kyoh86"
	org := "kyoh86-tryouts"

	t.Run("ListByOrg", func(t *testing.T) {
		t.Run("NilOption", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryListByOrg(ctx, org, nil).Return([]*github.Repository{{
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-1"),
			}, {
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-2"),
			}}, nil, nil)

			projects, err := remote.ListByOrg(ctx, org, nil)
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			expects := []string{
				host + "/" + org + "/org-repo-1",
				host + "/" + org + "/org-repo-2",
			}
			for i, expect := range expects {
				actual := projects[i].RelPath()
				if expect != actual {
					t.Errorf("expect project %q at %d but %q is gotten", expect, i, actual)
				}
			}
		})

		t.Run("EmptyOption", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryListByOrg(ctx, org, nil).Return([]*github.Repository{{
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-1"),
			}, {
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-2"),
			}}, nil, nil)

			projects, err := remote.ListByOrg(ctx, org, &testtarget.RemoteListByOrgOption{})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			expects := []string{
				host + "/" + org + "/org-repo-1",
				host + "/" + org + "/org-repo-2",
			}
			for i, expect := range expects {
				actual := projects[i].RelPath()
				if expect != actual {
					t.Errorf("expect project %q at %d but %q is gotten", expect, i, actual)
				}
			}
		})

		t.Run("WithOptions", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryListByOrg(ctx, org, &github.RepositoryListOptions{
				Visibility: "public",
			}).Return([]*github.Repository{{
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-1"),
			}, {
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-2"),
			}}, nil, nil)

			projects, err := remote.ListByOrg(ctx, org, &testtarget.RemoteListByOrgOption{
				Options: &github.RepositoryListByOrgOptions{
					Type: "private",
				},
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			expects := []string{
				host + "/" + org + "/org-repo-1",
				host + "/" + org + "/org-repo-2",
			}
			for i, expect := range expects {
				actual := projects[i].RelPath()
				if expect != actual {
					t.Errorf("expect project %q at %d but %q is gotten", expect, i, actual)
				}
			}
		})

		t.Run("WithQuery", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryListByOrg(ctx, org, nil).Return([]*github.Repository{{
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-1"),
			}, {
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-2"),
			}}, nil, nil)

			projects, err := remote.ListByOrg(ctx, org, &testtarget.RemoteListByOrgOption{
				Query: "repo-1",
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			expects := []string{
				host + "/" + org + "/org-repo-1",
			}
			for i, expect := range expects {
				actual := projects[i].RelPath()
				if expect != actual {
					t.Errorf("expect project %q at %d but %q is gotten", expect, i, actual)
				}
			}
		})

		t.Run("WithQueryNoMatch", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryListByOrg(ctx, org, nil).Return([]*github.Repository{{
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-1"),
			}, {
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-2"),
			}}, nil, nil)

			projects, err := remote.ListByOrg(ctx, org, &testtarget.RemoteListByOrgOption{
				Query: "no-match",
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(projects) > 0 {
				t.Errorf("expect no project is found but %d projects are found", len(projects))
			}
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("NilOption", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryList(ctx, "", nil).Return([]github.Repository{{
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-1"),
			}, {
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-2"),
			}, {
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-1"),
			}, {
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-2"),
			}}, nil, nil)

			projects, err := remote.List(ctx, nil)
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(projects) != 2 {
				t.Fatalf("expect some projects, but %d is gotten", len(projects))
			}
			expects := []string{
				host + "/" + user + "/user-repo-1",
				host + "/" + user + "/user-repo-2",
				host + "/" + org + "/org-repo-1",
				host + "/" + org + "/org-repo-2",
			}
			for i, expect := range expects {
				actual := projects[i].RelPath()
				if expect != actual {
					t.Errorf("expect project %q at %d but %q is gotten", expect, i, actual)
				}
			}
		})

		t.Run("EmptyOption", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryList(ctx, "", nil).Return([]*github.Repository{{
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-1"),
			}, {
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-2"),
			}, {
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-1"),
			}, {
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-2"),
			}}, nil, nil)

			projects, err := remote.List(ctx, &testtarget.RemoteListOption{})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			expects := []string{
				host + "/" + user + "/user-repo-1",
				host + "/" + user + "/user-repo-2",
				host + "/" + org + "/org-repo-1",
				host + "/" + org + "/org-repo-2",
			}
			for i, expect := range expects {
				actual := projects[i].RelPath()
				if expect != actual {
					t.Errorf("expect project %q at %d but %q is gotten", expect, i, actual)
				}
			}
		})

		t.Run("ByUser", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryList(ctx, user, nil).Return([]*github.Repository{{
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-1"),
			}, {
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-2"),
			}}, nil, nil)

			projects, err := remote.List(ctx, &testtarget.RemoteListOption{
				User: user,
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			expects := []string{
				host + "/" + user + "/user-repo-1",
				host + "/" + user + "/user-repo-2",
			}
			for i, expect := range expects {
				actual := projects[i].RelPath()
				if expect != actual {
					t.Errorf("expect project %q at %d but %q is gotten", expect, i, actual)
				}
			}
		})

		t.Run("WithOptions", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryList(ctx, "", &github.RepositoryListOptions{
				Visibility: "public",
			}).Return([]*github.Repository{{
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-1"),
			}, {
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-2"),
			}, {
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-1"),
			}, {
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-2"),
			}}, nil, nil)

			projects, err := remote.List(ctx, &testtarget.RemoteListOption{
				Options: &github.RepositoryListOptions{
					Visibility: "public",
				},
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			expects := []string{
				host + "/" + user + "/user-repo-1",
				host + "/" + user + "/user-repo-2",
				host + "/" + org + "/org-repo-1",
				host + "/" + org + "/org-repo-2",
			}
			for i, expect := range expects {
				actual := projects[i].RelPath()
				if expect != actual {
					t.Errorf("expect project %q at %d but %q is gotten", expect, i, actual)
				}
			}
		})

		t.Run("ByUserWithOptions", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryList(ctx, user, &github.RepositoryListOptions{
				Visibility: "public",
			}).Return([]*github.Repository{{
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-1"),
			}, {
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-2"),
			}}, nil, nil)

			projects, err := remote.List(ctx, &testtarget.RemoteListOption{
				User: user,
				Options: &github.RepositoryListOptions{
					Visibility: "public",
				},
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			expects := []string{
				host + "/" + user + "/user-repo-1",
				host + "/" + user + "/user-repo-2",
			}
			for i, expect := range expects {
				actual := projects[i].RelPath()
				if expect != actual {
					t.Errorf("expect project %q at %d but %q is gotten", expect, i, actual)
				}
			}
		})

		t.Run("Query", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryList(ctx, "", nil).Return([]*github.Repository{{
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-1"),
			}, {
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-2"),
			}, {
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-1"),
			}, {
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-2"),
			}}, nil, nil)

			projects, err := remote.List(ctx, &testtarget.RemoteListOption{
				Query: "repo-1",
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			expects := []string{
				host + "/" + user + "/user-repo-1",
				host + "/" + org + "/org-repo-1",
			}
			for i, expect := range expects {
				actual := projects[i].RelPath()
				if expect != actual {
					t.Errorf("expect project %q at %d but %q is gotten", expect, i, actual)
				}
			}
		})

		t.Run("QueryNoMatch", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryList(ctx, "", nil).Return([]*github.Repository{{
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-1"),
			}, {
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-2"),
			}, {
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-1"),
			}, {
				Organization: &github.Organization{
					Login: &org,
				},
				Name: ptr.String("org-repo-2"),
			}}, nil, nil)

			projects, err := remote.List(ctx, &testtarget.RemoteListOption{
				Query: "no-match",
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(projects) == 1 {
				t.Fatalf("expect one project, but %d is gotten", len(projects))
			}
		})

		t.Run("ByUserWithQuery", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryList(ctx, user, nil).Return([]*github.Repository{{
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-1"),
			}, {
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-2"),
			}}, nil, nil)

			projects, err := remote.List(ctx, &testtarget.RemoteListOption{
				User:  user,
				Query: "repo-1",
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			expects := []string{
				host + "/" + user + "/user-repo-1",
			}
			for i, expect := range expects {
				actual := projects[i].RelPath()
				if expect != actual {
					t.Errorf("expect project %q at %d but %q is gotten", expect, i, actual)
				}
			}
		})

		t.Run("ByUserWithQueryNoMatch", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryList(ctx, user, nil).Return([]*github.Repository{{
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-1"),
			}, {
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-2"),
			}}, nil, nil)

			projects, err := remote.List(ctx, &testtarget.RemoteListOption{
				User:  user,
				Query: "no-match",
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(projects) > 0 {
				t.Errorf("expect no project is found but %d projects are found", len(projects))
			}
		})
	})

	/*
		t.Run("Create", func(t *testing.T) {
			t.Run("NilOption", func(t *testing.T) {
				mock, teardown := MockAdaptor(t)
				defer teardown()
				remote := testtarget.NewRemoteController(mock)
				mock.EXPECT().RepositoryCreate(ctx, "", &github.Repository{
					Name: ptr.String("gogh"),
				}).Return(&github.Repository{
					Owner: &github.User{
						Login: &user,
					},
					Name: ptr.String("gogh"),
				}, nil, nil)
				project, err := remote.Create(ctx, "gogh", nil)
				if err != nil {
					t.Fatalf("failed to listup: %s", err)
				}
				if project.Name() != "gogh" {
					t.Errorf("expect that a project be created with name %q, but actual %q", "gogh", project.Name())
				}
			})
			t.Run("EmptyOption", func(t *testing.T) {
				mock, teardown := MockAdaptor(t)
				defer teardown()
				remote := testtarget.NewRemoteController(mock)
				mock.EXPECT().RepositoryCreate(ctx, "", &github.Repository{}).Return([]*github.Repository{{
					Owner: &github.User{
						Login: &user,
					},
					Name: ptr.String("user-repo-1"),
				}, {
					Owner: &github.User{
						Login: &user,
					},
					Name: ptr.String("user-repo-2"),
				}}, nil, nil)

				projects, err := remote.Create(ctx, &testtarget.RemoteCreateOption{
					User:  user,
					Query: "no-match",
				})
				if err != nil {
					t.Fatalf("failed to listup: %s", err)
				}
			})
		})
		t.Run("Remove", func(t *testing.T) {
			// TODO: remove github.com/kyoh86/gogh -> github.com/kyoh86/gogh removed
		})
	*/
}
