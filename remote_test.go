package gogh_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	testtarget "github.com/kyoh86/gogh/v3"
	"github.com/kyoh86/gogh/v3/infra/github"
	"github.com/kyoh86/gogh/v3/infra/github_mock"
	"github.com/kyoh86/gogh/v3/infra/githubv4"
	"go.uber.org/mock/gomock"
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
			CloneURL: testtarget.Ptr("https://" + testtarget.DefaultHost + "/" + user + "/gogh"),
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

func TestRemoteListOption_GetOptions(t *testing.T) {
	private := github.RepositoryPrivacyPrivate
	public := github.RepositoryPrivacyPublic
	owner := github.RepositoryAffiliationOwner
	member := github.RepositoryAffiliationOrganizationMember
	collabo := github.RepositoryAffiliationCollaborator
	for _, testcase := range []struct {
		title string
		base  *testtarget.RemoteListOption
		want  *github.RepositoryListOptions
	}{
		{
			title: "nil",
			base:  nil,
			want: &github.RepositoryListOptions{
				OrderBy: github.RepositoryOrder{
					Field:     githubv4.RepositoryOrderFieldUpdatedAt,
					Direction: githubv4.OrderDirectionDesc,
				},
				Limit:             testtarget.RepositoryListMaxLimitPerPage,
				OwnerAffiliations: []github.RepositoryAffiliation{owner},
			},
		},
		{
			title: "empty",
			base:  &testtarget.RemoteListOption{},
			want: &github.RepositoryListOptions{
				OrderBy: github.RepositoryOrder{
					Field:     githubv4.RepositoryOrderFieldUpdatedAt,
					Direction: githubv4.OrderDirectionDesc,
				},
				Limit:             testtarget.RepositoryListMaxLimitPerPage,
				OwnerAffiliations: []github.RepositoryAffiliation{owner},
			},
		},
		{
			title: "private",
			base:  &testtarget.RemoteListOption{Private: testtarget.Ptr(true)},
			want: &github.RepositoryListOptions{
				OrderBy: github.RepositoryOrder{
					Field:     githubv4.RepositoryOrderFieldUpdatedAt,
					Direction: githubv4.OrderDirectionDesc,
				},
				Limit:             testtarget.RepositoryListMaxLimitPerPage,
				OwnerAffiliations: []github.RepositoryAffiliation{owner},
				Privacy:           private,
			},
		},
		{
			title: "public",
			base:  &testtarget.RemoteListOption{Private: testtarget.Ptr(false)},
			want: &github.RepositoryListOptions{
				OrderBy: github.RepositoryOrder{
					Field:     githubv4.RepositoryOrderFieldUpdatedAt,
					Direction: githubv4.OrderDirectionDesc,
				},
				Limit:             testtarget.RepositoryListMaxLimitPerPage,
				OwnerAffiliations: []github.RepositoryAffiliation{owner},
				Privacy:           public,
			},
		},
		{
			title: "limit",
			base:  &testtarget.RemoteListOption{Limit: 1},
			want: &github.RepositoryListOptions{
				OrderBy: github.RepositoryOrder{
					Field:     githubv4.RepositoryOrderFieldUpdatedAt,
					Direction: githubv4.OrderDirectionDesc,
				},
				Limit:             1,
				OwnerAffiliations: []github.RepositoryAffiliation{owner},
			},
		},
		{
			title: "relations",
			base: &testtarget.RemoteListOption{Relation: []testtarget.RepositoryRelation{
				testtarget.RepositoryRelationOwner,
				testtarget.RepositoryRelationOrganizationMember,
				testtarget.RepositoryRelationCollaborator,
			}},
			want: &github.RepositoryListOptions{
				OrderBy: github.RepositoryOrder{
					Field:     githubv4.RepositoryOrderFieldUpdatedAt,
					Direction: githubv4.OrderDirectionDesc,
				},
				Limit:             testtarget.RepositoryListMaxLimitPerPage,
				OwnerAffiliations: []github.RepositoryAffiliation{owner, member, collabo},
			},
		},
		{
			title: "fork",
			base: &testtarget.RemoteListOption{
				IsFork: testtarget.Ptr(true),
			},
			want: &github.RepositoryListOptions{
				OrderBy: github.RepositoryOrder{
					Field:     githubv4.RepositoryOrderFieldUpdatedAt,
					Direction: githubv4.OrderDirectionDesc,
				},
				Limit:             testtarget.RepositoryListMaxLimitPerPage,
				OwnerAffiliations: []github.RepositoryAffiliation{owner},
				IsFork:            testtarget.Ptr(true),
			},
		},
		{
			title: "sort by name",
			base: &testtarget.RemoteListOption{
				Sort: github.RepositoryOrderFieldName,
			},
			want: &github.RepositoryListOptions{
				Limit:             testtarget.RepositoryListMaxLimitPerPage,
				OwnerAffiliations: []github.RepositoryAffiliation{owner},
				OrderBy: github.RepositoryOrder{
					Field:     github.RepositoryOrderFieldName,
					Direction: github.OrderDirectionAsc,
				},
			},
		},
		{
			title: "sort by stargazers",
			base: &testtarget.RemoteListOption{
				Sort: github.RepositoryOrderFieldStargazers,
			},
			want: &github.RepositoryListOptions{
				Limit:             testtarget.RepositoryListMaxLimitPerPage,
				OwnerAffiliations: []github.RepositoryAffiliation{owner},
				OrderBy: github.RepositoryOrder{
					Field:     github.RepositoryOrderFieldStargazers,
					Direction: github.OrderDirectionDesc,
				},
			},
		},
		{
			title: "sort by stargazers ascending",
			base: &testtarget.RemoteListOption{
				Sort:  github.RepositoryOrderFieldStargazers,
				Order: github.OrderDirectionAsc,
			},
			want: &github.RepositoryListOptions{
				Limit:             testtarget.RepositoryListMaxLimitPerPage,
				OwnerAffiliations: []github.RepositoryAffiliation{owner},
				OrderBy: github.RepositoryOrder{
					Field:     github.RepositoryOrderFieldStargazers,
					Direction: github.OrderDirectionAsc,
				},
			},
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			got := testcase.base.GetOptions()
			if diff := cmp.Diff(testcase.want, got); diff != "" {
				t.Errorf("result mismatched\n-want, +got:\n%s", diff)
			}
		})
	}
}

func TestRemoteController_List(t *testing.T) {
	ctx := context.Background()

	owner := github.RepositoryAffiliationOwner
	emptyOption := &github.RepositoryListOptions{
		OrderBy: github.RepositoryOrder{
			Field:     githubv4.RepositoryOrderFieldUpdatedAt,
			Direction: githubv4.OrderDirectionDesc,
		},
		Limit:             testtarget.RepositoryListMaxLimitPerPage,
		OwnerAffiliations: []github.RepositoryAffiliation{owner},
	}
	host := testtarget.DefaultHost
	user := "kyoh86"
	org := "kyoh86-tryouts"
	t.Run("Error", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		internalError := errors.New("test error")
		mock.EXPECT().
			RepositoryList(ctx, jsonMatcher{emptyOption}).
			Return(nil, github.PageInfoFragment{}, internalError)

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
		mock.EXPECT().
			RepositoryList(ctx, jsonMatcher{emptyOption}).
			DoAndReturn(func(_ context.Context, _ *github.RepositoryListOptions) ([]*github.RepositoryFragment, github.PageInfoFragment, error) {
				<-time.After(sleep)
				return []*github.RepositoryFragment{{
					Owner: &githubv4.RepositoryFragmentOwnerOrganization{OwnerFragmentOrganization: githubv4.OwnerFragmentOrganization{Login: org}},
					Name:  "org-repo-1.git",
				}, {
					Owner: &githubv4.RepositoryFragmentOwnerOrganization{OwnerFragmentOrganization: githubv4.OwnerFragmentOrganization{Login: org}},
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
		mock.EXPECT().
			RepositoryList(ctx, jsonMatcher{emptyOption}).
			Return([]*github.RepositoryFragment{{
				Owner: &githubv4.RepositoryFragmentOwnerUser{OwnerFragmentUser: githubv4.OwnerFragmentUser{Login: "."}},
			}}, github.PageInfoFragment{}, nil)

		if _, err := remote.List(ctx, nil); err == nil {
			t.Fatal("expect failure to listup but not")
		}
	})

	t.Run("NilOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().
			RepositoryList(ctx, jsonMatcher{emptyOption}).
			Return([]*github.RepositoryFragment{{
				UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
				Owner:     &githubv4.RepositoryFragmentOwnerUser{OwnerFragmentUser: githubv4.OwnerFragmentUser{Login: user}},
				Name:      "user-repo-1",
			}, {
				UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
				Owner:     &githubv4.RepositoryFragmentOwnerUser{OwnerFragmentUser: githubv4.OwnerFragmentUser{Login: user}},
				Name:      "user-repo-2",
			}, {
				UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
				Owner:     &githubv4.RepositoryFragmentOwnerOrganization{OwnerFragmentOrganization: githubv4.OwnerFragmentOrganization{Login: org}},
				Name:      "org-repo-1",
			}, {
				UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
				Owner:     &githubv4.RepositoryFragmentOwnerOrganization{OwnerFragmentOrganization: githubv4.OwnerFragmentOrganization{Login: org}},
				Name:      "org-repo-2",
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
		mock.EXPECT().
			RepositoryList(ctx, jsonMatcher{emptyOption}).
			Return([]*github.RepositoryFragment{{
				UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
				Owner:     &githubv4.RepositoryFragmentOwnerUser{OwnerFragmentUser: githubv4.OwnerFragmentUser{Login: user}},
				Name:      "user-repo-1",
			}, {
				UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
				Owner:     &githubv4.RepositoryFragmentOwnerUser{OwnerFragmentUser: githubv4.OwnerFragmentUser{Login: user}},
				Name:      "user-repo-2",
			}}, github.PageInfoFragment{HasNextPage: true, EndCursor: "next-page"}, nil)
		mock.EXPECT().RepositoryList(ctx, jsonMatcher{&github.RepositoryListOptions{
			OrderBy: github.RepositoryOrder{
				Field:     githubv4.RepositoryOrderFieldUpdatedAt,
				Direction: githubv4.OrderDirectionDesc,
			},
			After:             "next-page",
			Limit:             testtarget.RepositoryListMaxLimitPerPage,
			OwnerAffiliations: []github.RepositoryAffiliation{owner},
		}}).Return([]*github.RepositoryFragment{{
			UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
			Owner:     &githubv4.RepositoryFragmentOwnerOrganization{OwnerFragmentOrganization: githubv4.OwnerFragmentOrganization{Login: org}}, Name: "org-repo-1",
		}, {
			UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
			Owner:     &githubv4.RepositoryFragmentOwnerOrganization{OwnerFragmentOrganization: githubv4.OwnerFragmentOrganization{Login: org}}, Name: "org-repo-2",
		}}, github.PageInfoFragment{}, nil)
		specs, err := remote.List(ctx, &testtarget.RemoteListOption{Limit: 0})
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

	t.Run("Paging Over Max Per Page", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)

		const (
			perPage = testtarget.RepositoryListMaxLimitPerPage
			limit   = perPage + 2
		)
		var responses []*github.RepositoryFragment
		for i := int64(0); i < limit+1; i++ { // getting limit +1 items but ignore it.
			responses = append(responses, &github.RepositoryFragment{
				UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
				Owner:     &githubv4.RepositoryFragmentOwnerUser{OwnerFragmentUser: githubv4.OwnerFragmentUser{Login: user}},
				Name:      fmt.Sprintf("user-repo-%03d", i+1),
			})
		}
		mock.EXPECT().RepositoryList(ctx, jsonMatcher{&github.RepositoryListOptions{
			OrderBy: github.RepositoryOrder{
				Field:     githubv4.RepositoryOrderFieldUpdatedAt,
				Direction: githubv4.OrderDirectionDesc,
			},
			Limit:             perPage,
			OwnerAffiliations: []github.RepositoryAffiliation{owner},
		}}).Return(responses[0:perPage], github.PageInfoFragment{HasNextPage: true, EndCursor: "next-page"}, nil)
		mock.EXPECT().RepositoryList(ctx, jsonMatcher{&github.RepositoryListOptions{
			After: "next-page",
			OrderBy: github.RepositoryOrder{
				Field:     githubv4.RepositoryOrderFieldUpdatedAt,
				Direction: githubv4.OrderDirectionDesc,
			},
			Limit:             perPage,
			OwnerAffiliations: []github.RepositoryAffiliation{owner},
		}}).Return(responses[perPage:], github.PageInfoFragment{}, nil)
		specs, err := remote.List(ctx, &testtarget.RemoteListOption{Limit: limit})
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if int64(len(specs)) != limit {
			t.Fatalf("expect some specs, but %d is gotten", len(specs))
		}
		for i := int64(0); i < limit; i++ {
			got := specs[i]
			want := fmt.Sprintf("user-repo-%03d", i+1)
			if want != got.Name() {
				t.Errorf("expect name %q but %q gotten", want, got.Name())
			}
		}
	})

	t.Run("EmptyOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().
			RepositoryList(ctx, jsonMatcher{emptyOption}).
			Return([]*github.RepositoryFragment{{
				UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
				Owner:     &githubv4.RepositoryFragmentOwnerUser{OwnerFragmentUser: githubv4.OwnerFragmentUser{Login: user}}, Name: "user-repo-1",
			}, {
				UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
				Owner:     &githubv4.RepositoryFragmentOwnerUser{OwnerFragmentUser: githubv4.OwnerFragmentUser{Login: user}}, Name: "user-repo-2",
			}, {
				UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
				Owner:     &githubv4.RepositoryFragmentOwnerOrganization{OwnerFragmentOrganization: githubv4.OwnerFragmentOrganization{Login: org}}, Name: "org-repo-1",
			}, {
				UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
				Owner:     &githubv4.RepositoryFragmentOwnerOrganization{OwnerFragmentOrganization: githubv4.OwnerFragmentOrganization{Login: org}}, Name: "org-repo-2",
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
		mock.EXPECT().
			RepositoryList(ctx, jsonMatcher{emptyOption}).
			Return([]*github.RepositoryFragment{{
				UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
				Owner:     &githubv4.RepositoryFragmentOwnerUser{OwnerFragmentUser: githubv4.OwnerFragmentUser{Login: user}}, Name: "user-repo-1",
			}, {
				UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
				Owner:     &githubv4.RepositoryFragmentOwnerUser{OwnerFragmentUser: githubv4.OwnerFragmentUser{Login: user}}, Name: "user-repo-2",
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
			OrderBy: github.RepositoryOrder{
				Field:     githubv4.RepositoryOrderFieldUpdatedAt,
				Direction: githubv4.OrderDirectionDesc,
			},
			Privacy:           public,
			OwnerAffiliations: emptyOption.OwnerAffiliations,
			Limit:             emptyOption.Limit,
		}}).Return([]*github.RepositoryFragment{{
			UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
			Owner:     &githubv4.RepositoryFragmentOwnerUser{OwnerFragmentUser: githubv4.OwnerFragmentUser{Login: user}}, Name: "user-repo-1",
		}, {
			UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
			Owner:     &githubv4.RepositoryFragmentOwnerUser{OwnerFragmentUser: githubv4.OwnerFragmentUser{Login: user}}, Name: "user-repo-2",
		}}, github.PageInfoFragment{}, nil)

		specs, err := remote.List(ctx, &testtarget.RemoteListOption{
			Private: testtarget.Ptr(false),
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
			OrderBy: github.RepositoryOrder{
				Field:     githubv4.RepositoryOrderFieldUpdatedAt,
				Direction: githubv4.OrderDirectionDesc,
			},
			Privacy:           public,
			OwnerAffiliations: emptyOption.OwnerAffiliations,
			Limit:             emptyOption.Limit,
		}}).Return([]*github.RepositoryFragment{{
			UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
			Owner:     &githubv4.RepositoryFragmentOwnerUser{OwnerFragmentUser: githubv4.OwnerFragmentUser{Login: user}},
			Name:      "user-repo-1",
		}, {
			UpdatedAt: time.Date(2021, time.May, 1, 1, 0, 0, 0, time.UTC),
			Owner:     &githubv4.RepositoryFragmentOwnerUser{OwnerFragmentUser: githubv4.OwnerFragmentUser{Login: user}},
			Name:      "user-repo-2",
		}}, github.PageInfoFragment{}, nil)

		specs, err := remote.List(ctx, &testtarget.RemoteListOption{
			Private: testtarget.Ptr(false),
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
			Name: testtarget.Ptr("gogh"),
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
			Name: testtarget.Ptr("gogh"),
		}}).Return(&github.Repository{
			CloneURL: testtarget.Ptr("https://github.com/" + user + "/gogh.git"),
		}, nil, nil)
		spec, err := remote.Create(ctx, "gogh", nil)
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != user {
			t.Errorf(
				"expect that a spec be created with user %q, but actual %q",
				user,
				spec.Owner(),
			)
		}
		if spec.Name() != "gogh" {
			t.Errorf(
				"expect that a spec be created with name %q, but actual %q",
				"gogh",
				spec.Name(),
			)
		}
	})

	t.Run("EmptyOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryCreate(ctx, "", jsonMatcher{&github.Repository{
			Name: testtarget.Ptr("user-repo-1"),
		}}).Return(&github.Repository{
			CloneURL: testtarget.Ptr("https://github.com/" + user + "/user-repo-1.git"),
		}, nil, nil)

		spec, err := remote.Create(ctx, "user-repo-1", &testtarget.RemoteCreateOption{})
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != user {
			t.Errorf(
				"expect that a spec be created with user %q, but actual %q",
				user,
				spec.Owner(),
			)
		}
		if spec.Name() != "user-repo-1" {
			t.Errorf(
				"expect that a spec be created with name %q, but actual %q",
				"user-repo-1",
				spec.Name(),
			)
		}
	})

	t.Run("WithOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryCreate(ctx, "", jsonMatcher{&github.Repository{
			Name:       testtarget.Ptr("user-repo-1"),
			Homepage:   testtarget.Ptr("https://kyoh86.dev"),
			TeamID:     testtarget.Ptr(int64(3)),
			HasIssues:  testtarget.Ptr(false),
			IsTemplate: testtarget.Ptr(true),
		}}).Return(&github.Repository{
			CloneURL: testtarget.Ptr("https://github.com/" + user + "/user-repo-1.git"),
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
			t.Errorf(
				"expect that a spec be created with user %q, but actual %q",
				user,
				spec.Owner(),
			)
		}
		if spec.Name() != "user-repo-1" {
			t.Errorf(
				"expect that a spec be created with name %q, but actual %q",
				"user-repo-1",
				spec.Name(),
			)
		}
	})

	t.Run("WithOrganization", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryCreate(ctx, org, &github.Repository{
			Name: testtarget.Ptr("org-repo-1"),
		}).Return(&github.Repository{
			CloneURL: testtarget.Ptr("https://github.com/" + org + "/org-repo-1.git"),
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
			t.Errorf(
				"expect that a spec be created with name %q, but actual %q",
				"org-repo-1",
				spec.Name(),
			)
		}
	})

	t.Run("WithOrganizationAndOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryCreate(ctx, org, &github.Repository{
			Name:     testtarget.Ptr("org-repo-1"),
			Homepage: testtarget.Ptr("https://kyoh86.dev"),
		}).Return(&github.Repository{
			CloneURL: testtarget.Ptr("https://github.com/" + org + "/org-repo-1.git"),
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
			t.Errorf(
				"expect that a spec be created with name %q, but actual %q",
				"org-repo-1",
				spec.Name(),
			)
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

		mock.EXPECT().
			RepositoryCreateFromTemplate(ctx, "temp-owner", "temp-name", jsonMatcher{&github.TemplateRepoRequest{
				Name: testtarget.Ptr("gogh"),
			}}).
			Return(nil, nil, internalError)

		if _, err := remote.CreateFromTemplate(ctx, "temp-owner", "temp-name", "gogh", nil); !errors.Is(
			err,
			internalError,
		) {
			t.Errorf("expect passing internal error %q but actual %q", internalError, err)
		}
	})

	t.Run("NilOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().
			RepositoryCreateFromTemplate(ctx, "temp-owner", "temp-name", jsonMatcher{&github.TemplateRepoRequest{
				Name: testtarget.Ptr("gogh"),
			}}).
			Return(&github.Repository{
				CloneURL: testtarget.Ptr("https://github.com/" + user + "/gogh.git"),
			}, nil, nil)
		spec, err := remote.CreateFromTemplate(ctx, "temp-owner", "temp-name", "gogh", nil)
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != user {
			t.Errorf(
				"expect that a spec be created with user %q, but actual %q",
				user,
				spec.Owner(),
			)
		}
		if spec.Name() != "gogh" {
			t.Errorf(
				"expect that a spec be created with name %q, but actual %q",
				"gogh",
				spec.Name(),
			)
		}
	})

	t.Run("EmptyOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().
			RepositoryCreateFromTemplate(ctx, "temp-owner", "temp-name", jsonMatcher{&github.TemplateRepoRequest{
				Name: testtarget.Ptr("user-repo-1"),
			}}).
			Return(&github.Repository{
				CloneURL: testtarget.Ptr("https://github.com/" + user + "/user-repo-1.git"),
			}, nil, nil)

		spec, err := remote.CreateFromTemplate(
			ctx,
			"temp-owner",
			"temp-name",
			"user-repo-1",
			&testtarget.RemoteCreateFromTemplateOption{},
		)
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != user {
			t.Errorf(
				"expect that a spec be created with user %q, but actual %q",
				user,
				spec.Owner(),
			)
		}
		if spec.Name() != "user-repo-1" {
			t.Errorf(
				"expect that a spec be created with name %q, but actual %q",
				"user-repo-1",
				spec.Name(),
			)
		}
	})

	t.Run("WithOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().
			RepositoryCreateFromTemplate(ctx, "temp-owner", "temp-name", jsonMatcher{&github.TemplateRepoRequest{
				Name:  testtarget.Ptr("user-repo-1"),
				Owner: testtarget.Ptr("custom-user"),
			}}).
			Return(&github.Repository{
				CloneURL: testtarget.Ptr("https://github.com/custom-user/user-repo-1.git"),
			}, nil, nil)

		spec, err := remote.CreateFromTemplate(
			ctx,
			"temp-owner",
			"temp-name",
			"user-repo-1",
			&testtarget.RemoteCreateFromTemplateOption{
				Owner: "custom-user",
			},
		)
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != "custom-user" {
			t.Errorf(
				"expect that a spec be created with user %q, but actual %q",
				"custom-user",
				spec.Owner(),
			)
		}
		if spec.Name() != "user-repo-1" {
			t.Errorf(
				"expect that a spec be created with name %q, but actual %q",
				"user-repo-1",
				spec.Name(),
			)
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

		mock.EXPECT().
			RepositoryCreateFork(ctx, user, "user-repo-1", nil).
			Return(nil, nil, internalError)

		if _, err := remote.Fork(ctx, user, "user-repo-1", nil); !errors.Is(err, internalError) {
			t.Errorf("expect passing internal error %q but actual %q", internalError, err)
		}
	})

	t.Run("NilOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().RepositoryCreateFork(ctx, user, "user-repo-1", nil).Return(&github.Repository{
			CloneURL: testtarget.Ptr("https://github.com/" + user + "/user-repo-1.git"),
		}, nil, nil)
		spec, err := remote.Fork(ctx, user, "user-repo-1", nil)
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != user {
			t.Errorf(
				"expect that a spec be created with user %q, but actual %q",
				user,
				spec.Owner(),
			)
		}
		if spec.Name() != "user-repo-1" {
			t.Errorf(
				"expect that a spec be created with name %q, but actual %q",
				"user-repo-1",
				spec.Name(),
			)
		}
	})

	t.Run("EmptyOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().
			RepositoryCreateFork(ctx, user, "user-repo-1", jsonMatcher{&github.RepositoryCreateForkOptions{}}).
			Return(&github.Repository{
				CloneURL: testtarget.Ptr("https://github.com/" + user + "/user-repo-1.git"),
			}, nil, nil)

		spec, err := remote.Fork(ctx, user, "user-repo-1", &testtarget.RemoteForkOption{})
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != user {
			t.Errorf(
				"expect that a spec be created with user %q, but actual %q",
				user,
				spec.Owner(),
			)
		}
		if spec.Name() != "user-repo-1" {
			t.Errorf(
				"expect that a spec be created with name %q, but actual %q",
				"user-repo-1",
				spec.Name(),
			)
		}
	})

	t.Run("WithOption", func(t *testing.T) {
		mock, teardown := MockAdaptor(t)
		defer teardown()
		remote := testtarget.NewRemoteController(mock)
		mock.EXPECT().
			RepositoryCreateFork(ctx, user, "user-repo-1", jsonMatcher{&github.RepositoryCreateForkOptions{
				Organization: org,
			}}).
			Return(&github.Repository{
				CloneURL: testtarget.Ptr("https://github.com/" + user + "/user-repo-1.git"),
			}, nil, nil)

		spec, err := remote.Fork(ctx, user, "user-repo-1", &testtarget.RemoteForkOption{
			Organization: org,
		})
		if err != nil {
			t.Fatalf("failed to listup: %s", err)
		}
		if spec.Owner() != user {
			t.Errorf(
				"expect that a spec be created with user %q, but actual %q",
				user,
				spec.Owner(),
			)
		}
		if spec.Name() != "user-repo-1" {
			t.Errorf(
				"expect that a spec be created with name %q, but actual %q",
				"user-repo-1",
				spec.Name(),
			)
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
