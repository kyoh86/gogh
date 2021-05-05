package gogh_test

import (
	"context"
	"errors"
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
	mock.EXPECT().GetHost().AnyTimes().Return("github.com")
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
		mock.EXPECT().RepositoryList(ctx, jsonMatcher{&github.RepositoryListOptions{Limit: ptr.Int64(100)}}).Return(nil, github.PageInfoFragment{}, internalError)

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
		mock.EXPECT().RepositoryList(ctx, jsonMatcher{&github.RepositoryListOptions{Limit: ptr.Int64(100)}}).DoAndReturn(func(_ context.Context, _ *github.RepositoryListOptions) ([]*github.RepositoryFragment, github.PageInfoFragment, error) {
			<-time.After(sleep)
			return []*github.RepositoryFragment{{
				Owner: github.OwnerFragment{Login: org},
				Name:  "org-repo-1.git",
			}, {
				Owner: github.OwnerFragment{Login: org},
				Name:  "org-repo-2.git",
			}}, github.PageInfoFragment{}, nil
		})

		if _, err := remote.List(ctx, nil); !errors.Is(err, internalError) {
			t.Errorf("expect passing internal error %q but actual %q", internalError, err)
		}
	})

	t.Run("InvalidRepo", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryList(ctx, jsonMatcher{&github.RepositoryListOptions{Limit: ptr.Int64(100)}}).Return([]*github.RepositoryFragment{{
			Owner: github.OwnerFragment{Login: "."}, Name: ".",
		}}, github.PageInfoFragment{}, nil)

		if _, err := remote.List(ctx, nil); err == nil {
			t.Fatal("expect failure to listup but not")
		}
	})

	t.Run("NilOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryList(ctx, jsonMatcher{&github.RepositoryListOptions{Limit: ptr.Int64(100)}}).Return([]*github.RepositoryFragment{{
			Owner: github.OwnerFragment{Login: user}, Name: "user-repo-1",
		}, {
			Owner: github.OwnerFragment{Login: user}, Name: "user-repo-2",
		}, {
			Owner: github.OwnerFragment{Login: org}, Name: "org-repo-1",
		}, {
			Owner: github.OwnerFragment{Login: org}, Name: "org-repo-2",
		}}, github.PageInfoFragment{}, nil)

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
		mock.EXPECT().RepositoryList(ctx, jsonMatcher{&github.RepositoryListOptions{Limit: ptr.Int64(100)}}).Return([]*github.RepositoryFragment{{
			Owner: github.OwnerFragment{Login: user}, Name: "user-repo-1",
		}, {
			Owner: github.OwnerFragment{Login: user}, Name: "user-repo-2",
		}}, github.PageInfoFragment{HasNextPage: true, EndCursor: ptr.String("next-page")}, nil)
		mock.EXPECT().RepositoryList(ctx, jsonMatcher{&github.RepositoryListOptions{
			After: ptr.String("next-page"),
			Limit: ptr.Int64(100),
		}}).Return([]*github.RepositoryFragment{{
			Owner: github.OwnerFragment{Login: org}, Name: "org-repo-1",
		}, {
			Owner: github.OwnerFragment{Login: org}, Name: "org-repo-2",
		}}, github.PageInfoFragment{}, nil)

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

	owner := github.RepositoryAffiliationOwner
	emptyOption := &github.RepositoryListOptions{
		OwnerAffiliations: []*github.RepositoryAffiliation{&owner},
		Limit:             ptr.Int64(100),
	}
	t.Run("EmptyOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryList(ctx, jsonMatcher{emptyOption}).Return([]*github.RepositoryFragment{{
			Owner: github.OwnerFragment{Login: user}, Name: "user-repo-1",
		}, {
			Owner: github.OwnerFragment{Login: user}, Name: "user-repo-2",
		}, {
			Owner: github.OwnerFragment{Login: org}, Name: "org-repo-1",
		}, {
			Owner: github.OwnerFragment{Login: org}, Name: "org-repo-2",
		}}, github.PageInfoFragment{}, nil)

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
		mock.EXPECT().RepositoryList(ctx, jsonMatcher{emptyOption}).Return([]*github.RepositoryFragment{{
			Owner: github.OwnerFragment{Login: user}, Name: "user-repo-1",
		}, {
			Owner: github.OwnerFragment{Login: user}, Name: "user-repo-2",
		}}, github.PageInfoFragment{}, nil)
		specs, err := remote.List(ctx, &testtarget.RemoteListOption{})
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
		public := github.RepositoryPrivacyPublic
		mock.EXPECT().RepositoryList(ctx, jsonMatcher{&github.RepositoryListOptions{
			Privacy:           &public,
			OwnerAffiliations: emptyOption.OwnerAffiliations,
			Limit:             emptyOption.Limit,
		}}).Return([]*github.RepositoryFragment{{
			Owner: github.OwnerFragment{Login: user}, Name: "user-repo-1",
		}, {
			Owner: github.OwnerFragment{Login: user}, Name: "user-repo-2",
		}}, github.PageInfoFragment{}, nil)

		specs, err := remote.List(ctx, &testtarget.RemoteListOption{
			Private: ptr.Bool(false),
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
		public := github.RepositoryPrivacyPublic
		mock.EXPECT().RepositoryList(ctx, jsonMatcher{github.RepositoryListOptions{
			Privacy:           &public,
			OwnerAffiliations: emptyOption.OwnerAffiliations,
			Limit:             emptyOption.Limit,
		}}).Return([]*github.RepositoryFragment{{
			Owner: github.OwnerFragment{Login: user}, Name: "user-repo-1",
		}, {
			Owner: github.OwnerFragment{Login: user}, Name: "user-repo-2",
		}}, github.PageInfoFragment{}, nil)

		specs, err := remote.List(ctx, &testtarget.RemoteListOption{
			Private: ptr.Bool(false),
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
