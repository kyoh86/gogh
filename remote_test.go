package gogh_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

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

func TestRemoteController_Get(t *testing.T) {
	ctx := context.Background()

	host := testtarget.DefaultHost
	user := "kyoh86"
	t.Run("Error", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		internalError := errors.New("test error")
		mock.EXPECT().RepositoryGet(ctx, user, "gogh").Return(nil, nil, internalError)
		if _, err := remote.Get(ctx, user, "gogh", nil); !errors.Is(err, internalError) {
			t.Errorf("expect passing internal error %q but actual %q", internalError, err)
		}
	})
	t.Run("Success", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryGet(ctx, user, "gogh").Return(&github.Repository{
			CloneURL: ptr.String("https://" + testtarget.DefaultHost + "/" + user + "/gogh"),
		}, nil, nil)

		actual, err := remote.Get(ctx, user, "gogh", nil)
		if err != nil {
			t.Fatalf("failed to get: %s", err)
		}
		if host != actual.Host() {
			t.Errorf("expect host %q but %q gotten", host, actual.Host())
		}
		if user != actual.Owner() {
			t.Errorf("expect user %q but %q gotten", user, actual.Owner())
		}
		if actual.Name() != "gogh" {
			t.Errorf("expect name %q but %q gotten", "gogh", actual.Name())
		}
	})
}

func TestRemoteListOption_GetQuery(t *testing.T) {
	for _, tc := range []struct {
		title string
		opt   *testtarget.RemoteListOption
		want  string
	}{
		{
			title: "nil",
			opt:   nil,
			want:  "user:@me",
		},
		{
			title: "empty",
			opt:   &testtarget.RemoteListOption{},
			want:  "user:@me",
		},
		{
			title: "single user",
			opt:   &testtarget.RemoteListOption{Users: []string{"kyoh86"}},
			want:  `user:"kyoh86"`,
		},
		{
			title: "archived",
			opt:   &testtarget.RemoteListOption{Archived: ptr.Bool(true)},
			want:  `user:@me archived:true`,
		},
		{
			title: "not archived",
			opt:   &testtarget.RemoteListOption{Archived: ptr.Bool(false)},
			want:  `user:@me archived:false`,
		},
		{
			title: "fork",
			opt:   &testtarget.RemoteListOption{IsFork: ptr.Bool(true)},
			want:  `user:@me fork:true`,
		},
		{
			title: "no fork",
			opt:   &testtarget.RemoteListOption{IsFork: ptr.Bool(false)},
			want:  `user:@me fork:false`,
		},
		{
			title: "private",
			opt:   &testtarget.RemoteListOption{IsPrivate: ptr.Bool(true)},
			want:  `user:@me is:private`,
		},
		{
			title: "public",
			opt:   &testtarget.RemoteListOption{IsPrivate: ptr.Bool(false)},
			want:  `user:@me is:public`,
		},
		{
			title: "languate",
			opt:   &testtarget.RemoteListOption{Language: "go"},
			want:  `user:@me language:"go"`,
		},
		{
			title: "all options",
			opt: &testtarget.RemoteListOption{
				Users:     []string{"kyoh86"},
				Archived:  ptr.Bool(false),
				IsFork:    ptr.Bool(true),
				IsPrivate: ptr.Bool(false),
				Language:  "go",
			},
			want: `user:"kyoh86" archived:false fork:true is:public language:"go"`,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got := tc.opt.GetQuery()
			if tc.want != got {
				t.Errorf("-want: %s\n+got:%s", tc.want, got)
			}
		})
	}
}

func TestRemoteController_List(t *testing.T) {
	ctx := context.Background()

	host := testtarget.DefaultHost
	user := "kyoh86"
	org := "kyoh86-tryouts"
	t.Run("Error", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		internalError := errors.New("test error")
		mock.EXPECT().SearchRepository(ctx, "user:@me", jsonMatcher{&github.SearchOptions{
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}}).Return(nil, nil, internalError)

		if _, err := remote.List(ctx, nil); !errors.Is(err, internalError) {
			t.Errorf("expect passing internal error %q but actual %q", internalError, err)
		}
	})

	t.Run("Timeout", func(t *testing.T) {
		timeout := 500 * time.Millisecond
		sleep := 2 * time.Second
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		internalError := context.DeadlineExceeded

		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().SearchRepository(ctx, "user:@me", jsonMatcher{&github.SearchOptions{
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}}).DoAndReturn(func(_ context.Context, _ string, _ *github.SearchOptions) ([]*github.Repository, *github.Response, error) {
			<-time.After(sleep)
			return []*github.Repository{{
				CloneURL: ptr.String("https://github.com/" + org + "/org-repo-1.git"),
			}, {
				CloneURL: ptr.String("https://github.com/" + org + "/org-repo-2.git"),
			}}, &github.Response{NextPage: 0}, nil
		})

		if _, err := remote.List(ctx, nil); !errors.Is(err, internalError) {
			t.Errorf("expect passing internal error %q but actual %q", internalError, err)
		}
	})

	t.Run("InvalidRepo", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().SearchRepository(ctx, "user:@me", jsonMatcher{&github.SearchOptions{
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}}).Return([]*github.Repository{{
			CloneURL: ptr.String("://invalid//url"),
		}}, &github.Response{NextPage: 1}, nil)

		if _, err := remote.List(ctx, nil); err == nil {
			t.Fatal("expect failure to listup but not")
		}
	})

	t.Run("NilOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().SearchRepository(ctx, "user:@me", jsonMatcher{&github.SearchOptions{
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}}).Return([]*github.Repository{{
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-1.git"),
		}, {
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-2.git"),
		}, {
			CloneURL: ptr.String("https://github.com/" + org + "/org-repo-1.git"),
		}, {
			CloneURL: ptr.String("https://github.com/" + org + "/org-repo-2.git"),
		}}, &github.Response{NextPage: 0}, nil)

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
			if expect.user != actual.Owner() {
				t.Errorf("expect user %q but %q gotten", expect.user, actual.Owner())
			}
			if expect.name != actual.Name() {
				t.Errorf("expect name %q but %q gotten", expect.name, actual.Name())
			}
		}
	})

	t.Run("Paging", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().SearchRepository(ctx, "user:@me", jsonMatcher{&github.SearchOptions{
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}}).Return([]*github.Repository{{
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-1.git"),
		}, {
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-2.git"),
		}}, &github.Response{NextPage: 1}, nil)
		mock.EXPECT().SearchRepository(ctx, "user:@me", jsonMatcher{&github.SearchOptions{
			ListOptions: github.ListOptions{
				Page:    1,
				PerPage: 100,
			},
		}}).Return([]*github.Repository{{
			CloneURL: ptr.String("https://github.com/" + org + "/org-repo-1.git"),
		}, {
			CloneURL: ptr.String("https://github.com/" + org + "/org-repo-2.git"),
		}}, &github.Response{NextPage: 0}, nil)

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
			if expect.user != actual.Owner() {
				t.Errorf("expect user %q but %q gotten", expect.user, actual.Owner())
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
		mock.EXPECT().SearchRepository(ctx, "user:@me", jsonMatcher{&github.SearchOptions{
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}}).Return([]*github.Repository{{
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-1.git"),
		}, {
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-2.git"),
		}, {
			CloneURL: ptr.String("https://github.com/" + org + "/org-repo-1.git"),
		}, {
			CloneURL: ptr.String("https://github.com/" + org + "/org-repo-2.git"),
		}}, &github.Response{NextPage: 0}, nil)

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
			if expect.user != actual.Owner() {
				t.Errorf("expect user %q but %q gotten", expect.user, actual.Owner())
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
		mock.EXPECT().SearchRepository(ctx, fmt.Sprintf("user:%q", user), jsonMatcher{&github.SearchOptions{
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}}).Return([]*github.Repository{{
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-1.git"),
		}, {
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-2.git"),
		}}, &github.Response{NextPage: 0}, nil)

		specs, err := remote.List(ctx, &testtarget.RemoteListOption{
			Users: []string{user},
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
			if expect.user != actual.Owner() {
				t.Errorf("expect user %q but %q gotten", expect.user, actual.Owner())
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
		mock.EXPECT().SearchRepository(ctx, "user:@me is:public", &github.SearchOptions{
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}).Return([]*github.Repository{{
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-1.git"),
		}, {
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-2.git"),
		}}, &github.Response{NextPage: 0}, nil)

		specs, err := remote.List(ctx, &testtarget.RemoteListOption{
			IsPrivate: ptr.Bool(false),
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
			if expect.user != actual.Owner() {
				t.Errorf("expect user %q but %q gotten", expect.user, actual.Owner())
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
		mock.EXPECT().SearchRepository(ctx, fmt.Sprintf("user:%q is:public", user), &github.SearchOptions{
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}).Return([]*github.Repository{{
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-1.git"),
		}, {
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-2.git"),
		}}, &github.Response{NextPage: 0}, nil)

		specs, err := remote.List(ctx, &testtarget.RemoteListOption{
			Users:     []string{user},
			IsPrivate: ptr.Bool(false),
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
			if expect.user != actual.Owner() {
				t.Errorf("expect user %q but %q gotten", expect.user, actual.Owner())
			}
			if expect.name != actual.Name() {
				t.Errorf("expect name %q but %q gotten", expect.name, actual.Name())
			}
		}
	})
}

func TestRemoteController_Create(t *testing.T) {
	ctx := context.Background()

	user := "kyoh86"
	org := "kyoh86-tryouts"
	t.Run("Error", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		internalError := errors.New("test error")

		mock.EXPECT().RepositoryCreate(ctx, "", jsonMatcher{&github.Repository{
			Name: ptr.String("gogh"),
		}}).Return(nil, nil, internalError)

		if _, err := remote.Create(ctx, "gogh", nil); !errors.Is(err, internalError) {
			t.Errorf("expect passing internal error %q but actual %q", internalError, err)
		}
	})

	t.Run("NilOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryCreate(ctx, "", jsonMatcher{&github.Repository{
			Name: ptr.String("gogh"),
		}}).Return(&github.Repository{
			CloneURL: ptr.String("https://github.com/" + user + "/gogh.git"),
		}, nil, nil)
		spec, err := remote.Create(ctx, "gogh", nil)
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != user {
			t.Errorf("expect that a spec be created with user %q, but actual %q", user, spec.Owner())
		}
		if spec.Name() != "gogh" {
			t.Errorf("expect that a spec be created with name %q, but actual %q", "gogh", spec.Name())
		}
	})

	t.Run("EmptyOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryCreate(ctx, "", jsonMatcher{&github.Repository{
			Name: ptr.String("user-repo-1"),
		}}).Return(&github.Repository{
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-1.git"),
		}, nil, nil)

		spec, err := remote.Create(ctx, "user-repo-1", &testtarget.RemoteCreateOption{})
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != user {
			t.Errorf("expect that a spec be created with user %q, but actual %q", user, spec.Owner())
		}
		if spec.Name() != "user-repo-1" {
			t.Errorf("expect that a spec be created with name %q, but actual %q", "user-repo-1", spec.Name())
		}
	})

	t.Run("WithOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryCreate(ctx, "", jsonMatcher{&github.Repository{
			Name:       ptr.String("user-repo-1"),
			Homepage:   ptr.String("https://kyoh86.dev"),
			TeamID:     ptr.Int64(3),
			HasIssues:  ptr.Bool(false),
			IsTemplate: ptr.Bool(true),
		}}).Return(&github.Repository{
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-1.git"),
		}, nil, nil)

		spec, err := remote.Create(ctx, "user-repo-1", &testtarget.RemoteCreateOption{
			Homepage:      "https://kyoh86.dev",
			TeamID:        3,
			DisableIssues: true,
			IsTemplate:    true,
		})
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != user {
			t.Errorf("expect that a spec be created with user %q, but actual %q", user, spec.Owner())
		}
		if spec.Name() != "user-repo-1" {
			t.Errorf("expect that a spec be created with name %q, but actual %q", "user-repo-1", spec.Name())
		}
	})

	t.Run("WithOrganization", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryCreate(ctx, org, &github.Repository{
			Name: ptr.String("org-repo-1"),
		}).Return(&github.Repository{
			CloneURL: ptr.String("https://github.com/" + org + "/org-repo-1.git"),
		}, nil, nil)

		spec, err := remote.Create(ctx, "org-repo-1", &testtarget.RemoteCreateOption{
			Organization: org,
		})
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != org {
			t.Errorf("expect that a spec be created with user %q, but actual %q", org, spec.Owner())
		}
		if spec.Name() != "org-repo-1" {
			t.Errorf("expect that a spec be created with name %q, but actual %q", "org-repo-1", spec.Name())
		}
	})

	t.Run("WithOrganizationAndOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryCreate(ctx, org, &github.Repository{
			Name:     ptr.String("org-repo-1"),
			Homepage: ptr.String("https://kyoh86.dev"),
		}).Return(&github.Repository{
			CloneURL: ptr.String("https://github.com/" + org + "/org-repo-1.git"),
		}, nil, nil)

		spec, err := remote.Create(ctx, "org-repo-1", &testtarget.RemoteCreateOption{
			Organization: org,
			Homepage:     "https://kyoh86.dev",
		})
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != org {
			t.Errorf("expect that a spec be created with user %q, but actual %q", org, spec.Owner())
		}
		if spec.Name() != "org-repo-1" {
			t.Errorf("expect that a spec be created with name %q, but actual %q", "org-repo-1", spec.Name())
		}
	})
}

func TestRemoteController_CreateFromTemplate(t *testing.T) {
	ctx := context.Background()

	user := "kyoh86"
	t.Run("Error", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		internalError := errors.New("test error")

		mock.EXPECT().RepositoryCreateFromTemplate(ctx, "temp-owner", "temp-name", jsonMatcher{&github.TemplateRepoRequest{
			Name: ptr.String("gogh"),
		}}).Return(nil, nil, internalError)

		if _, err := remote.CreateFromTemplate(ctx, "temp-owner", "temp-name", "gogh", nil); !errors.Is(err, internalError) {
			t.Errorf("expect passing internal error %q but actual %q", internalError, err)
		}
	})

	t.Run("NilOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryCreateFromTemplate(ctx, "temp-owner", "temp-name", jsonMatcher{&github.TemplateRepoRequest{
			Name: ptr.String("gogh"),
		}}).Return(&github.Repository{
			CloneURL: ptr.String("https://github.com/" + user + "/gogh.git"),
		}, nil, nil)
		spec, err := remote.CreateFromTemplate(ctx, "temp-owner", "temp-name", "gogh", nil)
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != user {
			t.Errorf("expect that a spec be created with user %q, but actual %q", user, spec.Owner())
		}
		if spec.Name() != "gogh" {
			t.Errorf("expect that a spec be created with name %q, but actual %q", "gogh", spec.Name())
		}
	})

	t.Run("EmptyOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryCreateFromTemplate(ctx, "temp-owner", "temp-name", jsonMatcher{&github.TemplateRepoRequest{
			Name: ptr.String("user-repo-1"),
		}}).Return(&github.Repository{
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-1.git"),
		}, nil, nil)

		spec, err := remote.CreateFromTemplate(ctx, "temp-owner", "temp-name", "user-repo-1", &testtarget.RemoteCreateFromTemplateOption{})
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != user {
			t.Errorf("expect that a spec be created with user %q, but actual %q", user, spec.Owner())
		}
		if spec.Name() != "user-repo-1" {
			t.Errorf("expect that a spec be created with name %q, but actual %q", "user-repo-1", spec.Name())
		}
	})

	t.Run("WithOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryCreateFromTemplate(ctx, "temp-owner", "temp-name", jsonMatcher{&github.TemplateRepoRequest{
			Name:  ptr.String("user-repo-1"),
			Owner: ptr.String("custom-user"),
		}}).Return(&github.Repository{
			CloneURL: ptr.String("https://github.com/custom-user/user-repo-1.git"),
		}, nil, nil)

		spec, err := remote.CreateFromTemplate(ctx, "temp-owner", "temp-name", "user-repo-1", &testtarget.RemoteCreateFromTemplateOption{
			Owner: "custom-user",
		})
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != "custom-user" {
			t.Errorf("expect that a spec be created with user %q, but actual %q", "custom-user", spec.Owner())
		}
		if spec.Name() != "user-repo-1" {
			t.Errorf("expect that a spec be created with name %q, but actual %q", "user-repo-1", spec.Name())
		}
	})
}

func TestRemoteController_Fork(t *testing.T) {
	ctx := context.Background()

	user := "kyoh86"
	org := "kyoh86-tryouts"
	t.Run("Error", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		internalError := errors.New("test error")

		mock.EXPECT().RepositoryCreateFork(ctx, user, "user-repo-1", nil).Return(nil, nil, internalError)

		if _, err := remote.Fork(ctx, user, "user-repo-1", nil); !errors.Is(err, internalError) {
			t.Errorf("expect passing internal error %q but actual %q", internalError, err)
		}
	})

	t.Run("NilOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryCreateFork(ctx, user, "user-repo-1", nil).Return(&github.Repository{
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-1.git"),
		}, nil, nil)
		spec, err := remote.Fork(ctx, user, "user-repo-1", nil)
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != user {
			t.Errorf("expect that a spec be created with user %q, but actual %q", user, spec.Owner())
		}
		if spec.Name() != "user-repo-1" {
			t.Errorf("expect that a spec be created with name %q, but actual %q", "user-repo-1", spec.Name())
		}
	})

	t.Run("EmptyOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryCreateFork(ctx, user, "user-repo-1", jsonMatcher{&github.RepositoryCreateForkOptions{}}).Return(&github.Repository{
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-1.git"),
		}, nil, nil)

		spec, err := remote.Fork(ctx, user, "user-repo-1", &testtarget.RemoteForkOption{})
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != user {
			t.Errorf("expect that a spec be created with user %q, but actual %q", user, spec.Owner())
		}
		if spec.Name() != "user-repo-1" {
			t.Errorf("expect that a spec be created with name %q, but actual %q", "user-repo-1", spec.Name())
		}
	})

	t.Run("WithOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryCreateFork(ctx, user, "user-repo-1", jsonMatcher{&github.RepositoryCreateForkOptions{
			Organization: org,
		}}).Return(&github.Repository{
			CloneURL: ptr.String("https://github.com/" + user + "/user-repo-1.git"),
		}, nil, nil)

		spec, err := remote.Fork(ctx, user, "user-repo-1", &testtarget.RemoteForkOption{
			Organization: org,
		})
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != user {
			t.Errorf("expect that a spec be created with user %q, but actual %q", user, spec.Owner())
		}
		if spec.Name() != "user-repo-1" {
			t.Errorf("expect that a spec be created with name %q, but actual %q", "user-repo-1", spec.Name())
		}
	})
}

func TestRemoteController_GetSource(t *testing.T) {
	ctx := context.Background()

	host := testtarget.DefaultHost
	user := "kyoh86"
	org := "kyoh86-tryouts"
	t.Run("Error", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		internalError := errors.New("test error")
		mock.EXPECT().RepositoryGet(ctx, user, "gogh").Return(nil, nil, internalError)
		if _, err := remote.GetSource(ctx, user, "gogh", nil); !errors.Is(err, internalError) {
			t.Errorf("expect passing internal error %q but actual %q", internalError, err)
		}
	})
	t.Run("NotForked", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryGet(ctx, user, "gogh").Return(&github.Repository{
			CloneURL: ptr.String("https://" + testtarget.DefaultHost + "/" + user + "/gogh"),
		}, nil, nil)
		actual, err := remote.GetSource(ctx, user, "gogh", nil)
		if err != nil {
			t.Fatalf("failed to get: %s", err)
		}
		if host != actual.Host() {
			t.Errorf("expect host %q but %q gotten", host, actual.Host())
		}
		if user != actual.Owner() {
			t.Errorf("expect user %q but %q gotten", user, actual.Owner())
		}
		if actual.Name() != "gogh" {
			t.Errorf("expect name %q but %q gotten", "gogh", actual.Name())
		}
	})
	t.Run("Success", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryGet(ctx, user, "gogh").Return(&github.Repository{
			CloneURL: ptr.String("https://" + testtarget.DefaultHost + "/" + user + "/gogh"),
			Source: &github.Repository{
				CloneURL: ptr.String("https://" + testtarget.DefaultHost + "/" + org + "/gogh"),
			},
		}, nil, nil)
		actual, err := remote.GetSource(ctx, user, "gogh", nil)
		if err != nil {
			t.Fatalf("failed to get: %s", err)
		}
		if host != actual.Host() {
			t.Errorf("expect host %q but %q gotten", host, actual.Host())
		}
		if org != actual.Owner() {
			t.Errorf("expect org %q but %q gotten", org, actual.Owner())
		}
		if actual.Name() != "gogh" {
			t.Errorf("expect name %q but %q gotten", "gogh", actual.Name())
		}
	})
}

func TestRemoteController_GetParent(t *testing.T) {
	ctx := context.Background()

	host := testtarget.DefaultHost
	user := "kyoh86"
	org := "kyoh86-tryouts"
	t.Run("Error", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		internalError := errors.New("test error")
		mock.EXPECT().RepositoryGet(ctx, user, "gogh").Return(nil, nil, internalError)
		if _, err := remote.GetParent(ctx, user, "gogh", nil); !errors.Is(err, internalError) {
			t.Errorf("expect passing internal error %q but actual %q", internalError, err)
		}
	})
	t.Run("NotForked", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryGet(ctx, user, "gogh").Return(&github.Repository{
			CloneURL: ptr.String("https://" + testtarget.DefaultHost + "/" + user + "/gogh"),
		}, nil, nil)
		actual, err := remote.GetParent(ctx, user, "gogh", nil)
		if err != nil {
			t.Fatalf("failed to get: %s", err)
		}
		if host != actual.Host() {
			t.Errorf("expect host %q but %q gotten", host, actual.Host())
		}
		if user != actual.Owner() {
			t.Errorf("expect user %q but %q gotten", user, actual.Owner())
		}
		if actual.Name() != "gogh" {
			t.Errorf("expect name %q but %q gotten", "gogh", actual.Name())
		}
	})
	t.Run("Success", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryGet(ctx, user, "gogh").Return(&github.Repository{
			CloneURL: ptr.String("https://" + testtarget.DefaultHost + "/" + user + "/gogh"),
			Parent: &github.Repository{
				CloneURL: ptr.String("https://" + testtarget.DefaultHost + "/" + org + "/gogh"),
			},
		}, nil, nil)
		actual, err := remote.GetParent(ctx, user, "gogh", nil)
		if err != nil {
			t.Fatalf("failed to get: %s", err)
		}
		if host != actual.Host() {
			t.Errorf("expect host %q but %q gotten", host, actual.Host())
		}
		if org != actual.Owner() {
			t.Errorf("expect org %q but %q gotten", org, actual.Owner())
		}
		if actual.Name() != "gogh" {
			t.Errorf("expect name %q but %q gotten", "gogh", actual.Name())
		}
	})
}

func TestRemoteController_Delete(t *testing.T) {
	ctx := context.Background()

	user := "kyoh86"
	t.Run("Error", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		internalError := errors.New("test error")
		mock.EXPECT().RepositoryDelete(ctx, user, "gogh").Return(nil, internalError)

		if err := remote.Delete(ctx, user, "gogh", nil); !errors.Is(err, internalError) {
			t.Errorf("expect passing internal error %q but actual %q", internalError, err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryDelete(ctx, user, "gogh").Return(nil, nil)

		if err := remote.Delete(ctx, user, "gogh", nil); err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
	})
}
