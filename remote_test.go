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

			specs, err := remote.ListByOrg(ctx, org, nil)
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(specs) != 2 {
				t.Fatalf("expect 2 specs, but %d is gotten", len(specs))
			}
			for i, expect := range []struct {
				host string
				user string
				name string
			}{{
				host: host,
				user: org,
				name: "org-repo-1",
			}, {
				host: host,
				user: org,
				name: "org-repo-2",
			}} {
				actual := specs[i]
				if expect.host != actual.Host() {
					t.Errorf("expect host %q but %q gotten", expect.host, actual.Host())
				}
				if expect.user != actual.User() {
					t.Errorf("expect user %q but %q gotten", expect.user, actual.User())
				}
				if expect.name != actual.Name() {
					t.Errorf("expect name %q but %q gotten", expect.name, actual.Name())
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

			specs, err := remote.ListByOrg(ctx, org, &testtarget.RemoteListByOrgOption{})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(specs) != 2 {
				t.Fatalf("expect 2 specs, but %d is gotten", len(specs))
			}
			for i, expect := range []struct {
				host string
				user string
				name string
			}{{
				host: host,
				user: org,
				name: "org-repo-1",
			}, {
				host: host,
				user: org,
				name: "org-repo-2",
			}} {
				actual := specs[i]
				if expect.host != actual.Host() {
					t.Errorf("expect host %q but %q gotten", expect.host, actual.Host())
				}
				if expect.user != actual.User() {
					t.Errorf("expect user %q but %q gotten", expect.user, actual.User())
				}
				if expect.name != actual.Name() {
					t.Errorf("expect name %q but %q gotten", expect.name, actual.Name())
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

			specs, err := remote.ListByOrg(ctx, org, &testtarget.RemoteListByOrgOption{
				Options: &github.RepositoryListByOrgOptions{
					Type: "private",
				},
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(specs) != 2 {
				t.Fatalf("expect 2 specs, but %d is gotten", len(specs))
			}
			for i, expect := range []struct {
				host string
				user string
				name string
			}{{
				host: host,
				user: org,
				name: "org-repo-1",
			}, {
				host: host,
				user: org,
				name: "org-repo-2",
			}} {
				actual := specs[i]
				if expect.host != actual.Host() {
					t.Errorf("expect host %q but %q gotten", expect.host, actual.Host())
				}
				if expect.user != actual.User() {
					t.Errorf("expect user %q but %q gotten", expect.user, actual.User())
				}
				if expect.name != actual.Name() {
					t.Errorf("expect name %q but %q gotten", expect.name, actual.Name())
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

			specs, err := remote.ListByOrg(ctx, org, &testtarget.RemoteListByOrgOption{
				Query: "repo-1",
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(specs) != 1 {
				t.Fatalf("expect a spec, but %d is gotten", len(specs))
			}
			for i, expect := range []struct {
				host string
				user string
				name string
			}{{
				host: host,
				user: org,
				name: "org-repo-1",
			}} {
				actual := specs[i]
				if expect.host != actual.Host() {
					t.Errorf("expect host %q but %q gotten", expect.host, actual.Host())
				}
				if expect.user != actual.User() {
					t.Errorf("expect user %q but %q gotten", expect.user, actual.User())
				}
				if expect.name != actual.Name() {
					t.Errorf("expect name %q but %q gotten", expect.name, actual.Name())
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

			specs, err := remote.ListByOrg(ctx, org, &testtarget.RemoteListByOrgOption{
				Query: "no-match",
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(specs) > 0 {
				t.Errorf("expect no spec is found but %d specs are found", len(specs))
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

			specs, err := remote.List(ctx, nil)
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(specs) != 4 {
				t.Fatalf("expect some specs, but %d is gotten", len(specs))
			}
			for i, expect := range []struct {
				host string
				user string
				name string
			}{{
				host: host,
				user: user,
				name: "user-repo-1",
			}, {
				host: host,
				user: user,
				name: "user-repo-2",
			}, {
				host: host,
				user: org,
				name: "org-repo-1",
			}, {
				host: host,
				user: org,
				name: "org-repo-2",
			}} {
				actual := specs[i]
				if expect.host != actual.Host() {
					t.Errorf("expect host %q but %q gotten", expect.host, actual.Host())
				}
				if expect.user != actual.User() {
					t.Errorf("expect user %q but %q gotten", expect.user, actual.User())
				}
				if expect.name != actual.Name() {
					t.Errorf("expect name %q but %q gotten", expect.name, actual.Name())
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

			specs, err := remote.List(ctx, &testtarget.RemoteListOption{})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(specs) != 4 {
				t.Fatalf("expect some specs, but %d is gotten", len(specs))
			}
			for i, expect := range []struct {
				host string
				user string
				name string
			}{{
				host: host,
				user: user,
				name: "user-repo-1",
			}, {
				host: host,
				user: user,
				name: "user-repo-2",
			}, {
				host: host,
				user: org,
				name: "org-repo-1",
			}, {
				host: host,
				user: org,
				name: "org-repo-2",
			}} {
				actual := specs[i]
				if expect.host != actual.Host() {
					t.Errorf("expect host %q but %q gotten", expect.host, actual.Host())
				}
				if expect.user != actual.User() {
					t.Errorf("expect user %q but %q gotten", expect.user, actual.User())
				}
				if expect.name != actual.Name() {
					t.Errorf("expect name %q but %q gotten", expect.name, actual.Name())
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

			specs, err := remote.List(ctx, &testtarget.RemoteListOption{
				User: user,
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(specs) != 2 {
				t.Fatalf("expect some specs, but %d is gotten", len(specs))
			}
			for i, expect := range []struct {
				host string
				user string
				name string
			}{{
				host: host,
				user: user,
				name: "user-repo-1",
			}, {
				host: host,
				user: user,
				name: "user-repo-2",
			}} {
				actual := specs[i]
				if expect.host != actual.Host() {
					t.Errorf("expect host %q but %q gotten", expect.host, actual.Host())
				}
				if expect.user != actual.User() {
					t.Errorf("expect user %q but %q gotten", expect.user, actual.User())
				}
				if expect.name != actual.Name() {
					t.Errorf("expect name %q but %q gotten", expect.name, actual.Name())
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

			specs, err := remote.List(ctx, &testtarget.RemoteListOption{
				Options: &github.RepositoryListOptions{
					Visibility: "public",
				},
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(specs) != 2 {
				t.Fatalf("expect some specs, but %d is gotten", len(specs))
			}
			for i, expect := range []struct {
				host string
				user string
				name string
			}{{
				host: host,
				user: user,
				name: "user-repo-1",
			}, {
				host: host,
				user: user,
				name: "user-repo-2",
			}} {
				actual := specs[i]
				if expect.host != actual.Host() {
					t.Errorf("expect host %q but %q gotten", expect.host, actual.Host())
				}
				if expect.user != actual.User() {
					t.Errorf("expect user %q but %q gotten", expect.user, actual.User())
				}
				if expect.name != actual.Name() {
					t.Errorf("expect name %q but %q gotten", expect.name, actual.Name())
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

			specs, err := remote.List(ctx, &testtarget.RemoteListOption{
				User: user,
				Options: &github.RepositoryListOptions{
					Visibility: "public",
				},
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(specs) != 2 {
				t.Fatalf("expect some specs, but %d is gotten", len(specs))
			}
			for i, expect := range []struct {
				host string
				user string
				name string
			}{{
				host: host,
				user: user,
				name: "user-repo-1",
			}, {
				host: host,
				user: user,
				name: "user-repo-2",
			}} {
				actual := specs[i]
				if expect.host != actual.Host() {
					t.Errorf("expect host %q but %q gotten", expect.host, actual.Host())
				}
				if expect.user != actual.User() {
					t.Errorf("expect user %q but %q gotten", expect.user, actual.User())
				}
				if expect.name != actual.Name() {
					t.Errorf("expect name %q but %q gotten", expect.name, actual.Name())
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

			specs, err := remote.List(ctx, &testtarget.RemoteListOption{
				Query: "repo-1",
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(specs) != 2 {
				t.Fatalf("expect some specs, but %d is gotten", len(specs))
			}
			for i, expect := range []struct {
				host string
				user string
				name string
			}{{
				host: host,
				user: user,
				name: "user-repo-1",
			}, {
				host: host,
				user: org,
				name: "org-repo-1",
			}} {
				actual := specs[i]
				if expect.host != actual.Host() {
					t.Errorf("expect host %q but %q gotten", expect.host, actual.Host())
				}
				if expect.user != actual.User() {
					t.Errorf("expect user %q but %q gotten", expect.user, actual.User())
				}
				if expect.name != actual.Name() {
					t.Errorf("expect name %q but %q gotten", expect.name, actual.Name())
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

			specs, err := remote.List(ctx, &testtarget.RemoteListOption{
				Query: "no-match",
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(specs) > 1 {
				t.Fatalf("expect no spec is matched, but %d is gotten", len(specs))
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

			specs, err := remote.List(ctx, &testtarget.RemoteListOption{
				User:  user,
				Query: "repo-1",
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(specs) != 1 {
				t.Fatalf("expect one spec, but %d is gotten", len(specs))
			}
			for i, expect := range []struct {
				host string
				user string
				name string
			}{{
				host: host,
				user: user,
				name: "user-repo-1",
			}} {
				actual := specs[i]
				if expect.host != actual.Host() {
					t.Errorf("expect host %q but %q gotten", expect.host, actual.Host())
				}
				if expect.user != actual.User() {
					t.Errorf("expect user %q but %q gotten", expect.user, actual.User())
				}
				if expect.name != actual.Name() {
					t.Errorf("expect name %q but %q gotten", expect.name, actual.Name())
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

			specs, err := remote.List(ctx, &testtarget.RemoteListOption{
				User:  user,
				Query: "no-match",
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if len(specs) > 0 {
				t.Errorf("expect no spec is found but %d specs are found", len(specs))
			}
		})
	})

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
			spec, err := remote.Create(ctx, "gogh", nil)
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if spec.User() != user {
				t.Errorf("expect that a spec be created with user %q, but actual %q", user, spec.User())
			}
			if spec.Name() != "gogh" {
				t.Errorf("expect that a spec be created with name %q, but actual %q", "gogh", spec.Name())
			}
		})
		t.Run("EmptyOption", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryCreate(ctx, "", &github.Repository{
				Name: ptr.String("gogh"),
			}).Return(&github.Repository{
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-1"),
			}, nil, nil)

			spec, err := remote.Create(ctx, "gogh", &testtarget.RemoteCreateOption{})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if spec.User() != user {
				t.Errorf("expect that a spec be created with user %q, but actual %q", user, spec.User())
			}
			if spec.Name() != "gogh" {
				t.Errorf("expect that a spec be created with name %q, but actual %q", "gogh", spec.Name())
			}
		})

		t.Run("WithOption", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryCreate(ctx, "", &github.Repository{
				Name:     ptr.String("gogh"),
				Homepage: ptr.String("https://kyoh86.dev"),
			}).Return(&github.Repository{
				Owner: &github.User{
					Login: &user,
				},
				Name: ptr.String("user-repo-1"),
			}, nil, nil)

			spec, err := remote.Create(ctx, "gogh", &testtarget.RemoteCreateOption{
				Homepage: "https://kyoh86.dev",
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if spec.User() != user {
				t.Errorf("expect that a spec be created with user %q, but actual %q", user, spec.User())
			}
			if spec.Name() != "gogh" {
				t.Errorf("expect that a spec be created with name %q, but actual %q", "gogh", spec.Name())
			}
		})

		t.Run("WithOrganization", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryCreate(ctx, "kyoh86-tryouts", &github.Repository{
				Name: ptr.String("gogh"),
			}).Return(&github.Repository{
				Organization: &github.Organization{
					Login: ptr.String("kyoh86-tryouts"),
				},
				Name: ptr.String("user-repo-1"),
			}, nil, nil)

			spec, err := remote.Create(ctx, "gogh", &testtarget.RemoteCreateOption{
				Organization: "kyoh86-tryouts",
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if spec.User() != "kyoh86-tryouts" {
				t.Errorf("expect that a spec be created with user %q, but actual %q", "kyoh86-tryouts", spec.User())
			}
			if spec.Name() != "gogh" {
				t.Errorf("expect that a spec be created with name %q, but actual %q", "gogh", spec.Name())
			}
		})

		t.Run("WithOrganizationAndOption", func(t *testing.T) {
			mock, teardown := MockAdaptor(t)
			defer teardown()
			remote := testtarget.NewRemoteController(mock)
			mock.EXPECT().RepositoryCreate(ctx, "kyoh86-tryouts", &github.Repository{
				Name:     ptr.String("gogh"),
				Homepage: ptr.String("https://kyoh86.dev"),
			}).Return(&github.Repository{
				Organization: &github.Organization{
					Login: ptr.String("kyoh86-tryouts"),
				},
				Name: ptr.String("user-repo-1"),
			}, nil, nil)

			spec, err := remote.Create(ctx, "gogh", &testtarget.RemoteCreateOption{
				Organization: "kyoh86-tryouts",
				Homepage:     "https://kyoh86.dev",
			})
			if err != nil {
				t.Fatalf("failed to listup: %s", err)
			}
			if spec.User() != "kyoh86-tryouts" {
				t.Errorf("expect that a spec be created with user %q, but actual %q", "kyoh86-tryouts", spec.User())
			}
			if spec.Name() != "gogh" {
				t.Errorf("expect that a spec be created with name %q, but actual %q", "gogh", spec.Name())
			}
		})
	})
	t.Run("Remove", func(t *testing.T) {
		// TODO: remove github.com/kyoh86/gogh -> github.com/kyoh86/gogh removed
	})
}
