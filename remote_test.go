package gogh_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	github "github.com/google/go-github/v33/github"
	testtarget "github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/internal/mock_github"
	"github.com/wacul/ptr"
)

func TestRemoteController(t *testing.T) {
	ctx := context.Background()

	// TODO: list -> github.com/kyoh86/gogh, github.com/kyoh86/vim-gogh, github.kyoh86.dev/kyoh86/dotfiles
	// TODO: list --host github.com --owner kyoh86 -> github.com/kyoh86/gogh, github.com/kyoh86/vim-gogh

	// TODO: authorized server
	// TODO: create github.com/kyoh86/gogh -> github.com/kyoh86/gogh created
	// TODO: remove github.com/kyoh86/gogh -> github.com/kyoh86/gogh removed

	// NOTE: parsing gogh, kyoh86/gogh or github.com/kyoh86/gogh is the respoonsibilities of the "Descriptor" ->

	// TODO: multiple server
	t.Run("Unauthorized", func(t *testing.T) {
		host := testtarget.DefaultHost
		user := "kyoh86"
		org := "kyoh86-tryouts"

		server, err := testtarget.NewServerFor(host, user)
		if err != nil {
			t.Fatalf("failed to create a new server: %s", err)
		}
		remote := testtarget.NewRemoteController(server)

		t.Run("List", func(t *testing.T) {
			t.Run("NilOption", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mock := mock_github.NewMockAdaptor(ctrl)
				remote.SetAdaptor(mock)
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
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mock := mock_github.NewMockAdaptor(ctrl)
				remote.SetAdaptor(mock)
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

			t.Run("ByOrganization", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mock := mock_github.NewMockAdaptor(ctrl)
				remote.SetAdaptor(mock)
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

				projects, err := remote.List(ctx, &testtarget.RemoteListOption{
					Organization: org,
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

			t.Run("ByUser", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mock := mock_github.NewMockAdaptor(ctrl)
				remote.SetAdaptor(mock)
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

			t.Run("Options", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mock := mock_github.NewMockAdaptor(ctrl)
				remote.SetAdaptor(mock)
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

			t.Run("ByOrganizationWithOptions", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mock := mock_github.NewMockAdaptor(ctrl)
				remote.SetAdaptor(mock)
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

				projects, err := remote.List(ctx, &testtarget.RemoteListOption{
					Organization: org,
					Options: &github.RepositoryListOptions{
						Visibility: "public",
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

			t.Run("ByUserWithOptions", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mock := mock_github.NewMockAdaptor(ctrl)
				remote.SetAdaptor(mock)
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
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mock := mock_github.NewMockAdaptor(ctrl)
				remote.SetAdaptor(mock)
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
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mock := mock_github.NewMockAdaptor(ctrl)
				remote.SetAdaptor(mock)
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

			t.Run("OrganizationAndQuery", func(t *testing.T) {
				// TODO:
			})

			t.Run("OrganizationAndQueryNoMatch", func(t *testing.T) {
				// TODO:
			})

			t.Run("UserAndQuery", func(t *testing.T) {
				// TODO:
			})

			t.Run("UserAndQueryNoMatch", func(t *testing.T) {
				// TODO:
			})
		})
	})

	// func (r *RemoteController) Create(ctx context.Context, description Description) (Project, error)
	// func (r *RemoteController) Remove(ctx context.Context, description Description) error
}
