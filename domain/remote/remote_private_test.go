package remote

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	exgithub "github.com/google/go-github/v69/github"
	"github.com/kyoh86/gogh/v3/domain/reporef"
	"github.com/kyoh86/gogh/v3/infra/github"
	"github.com/kyoh86/gogh/v3/infra/githubv4"
	"github.com/kyoh86/gogh/v3/util"
)

func TestParseRepoRef(t *testing.T) {
	t.Run("valid http URL", func(t *testing.T) {
		want, err := reporef.NewRepoRef("github.com", "kyoh86", "gogh")
		if err != nil {
			t.Fatalf("unexpected error: %q", err)
		}
		got, err := parseRepoRef(&github.Repository{
			CloneURL: util.Ptr("https://github.com/kyoh86/gogh"),
		})
		if err != nil {
			t.Fatalf("unexpected error: %q", err)
		}
		if diff := cmp.Diff(want, got, cmp.AllowUnexported(want)); diff != "" {
			t.Errorf("result mistmatch\n-want, +got:\n%s", diff)
		}
	})

	t.Run("unsupported (ssh) URL", func(t *testing.T) {
		_, err := parseRepoRef(&github.Repository{
			CloneURL: util.Ptr("git@github.com:kyoh86/gogh.git"),
		})
		if err == nil {
			t.Fatal("expected error, but not")
		}
	})
}

func TestIngestRepository(t *testing.T) {
	t.Run("valid repository", func(t *testing.T) {
		tim, _ := time.Parse("2006-01-02", "2021-01-01")
		wantRef, err := reporef.NewRepoRef("github.com", "kyoh86", "gogh")
		if err != nil {
			t.Fatalf("unexpected error: %q", err)
		}
		wantParent, err := reporef.NewRepoRef("github.com", "kyoh86-tryouts", "test")
		if err != nil {
			t.Fatalf("unexpected error: %q", err)
		}
		want := RemoteRepo{
			URL:         "https://github.com/kyoh86/gogh",
			Description: "valid description",
			Homepage:    "valid homepage",
			Language:    "valid language",
			UpdatedAt:   tim,
			Archived:    true,
			Private:     false,
			IsTemplate:  true,
			Fork:        false,
			Ref:         wantRef,
			Parent:      &wantParent,
		}
		got, err := ingestRepository(&github.Repository{
			Description: util.Ptr("valid description"),
			Homepage:    util.Ptr("valid homepage"),
			UpdatedAt:   &exgithub.Timestamp{Time: tim},
			CloneURL:    util.Ptr("https://github.com/kyoh86/gogh"),
			Language:    util.Ptr("valid language"),
			Fork:        util.Ptr(false),
			Parent: &github.Repository{
				CloneURL: util.Ptr("https://github.com/kyoh86-tryouts/test"),
			},
			Archived:   util.Ptr(true),
			Private:    util.Ptr(false),
			IsTemplate: util.Ptr(true),
		})
		if err != nil {
			t.Fatalf("unexpected error: %q", err)
		}
		if diff := cmp.Diff(want, got, cmp.AllowUnexported(want), cmp.AllowUnexported(want.Ref)); diff != "" {
			t.Errorf("result mistmatch\n-want, +got:\n%s", diff)
		}
	})

	t.Run("unsupported (ssh) URL", func(t *testing.T) {
		tim, _ := time.Parse("2006-01-02", "2021-01-01")
		_, err := ingestRepository(&github.Repository{
			Description: util.Ptr("valid description"),
			Homepage:    util.Ptr("valid homepage"),
			UpdatedAt:   &exgithub.Timestamp{Time: tim},
			CloneURL:    util.Ptr("git@github.com:kyoh86/gogh.git"),
			Language:    util.Ptr("valid language"),
			Fork:        util.Ptr(false),
			Parent: &github.Repository{
				CloneURL: util.Ptr("https://github.com/kyoh86-tryouts/test"),
			},
			Archived:   util.Ptr(true),
			Private:    util.Ptr(false),
			IsTemplate: util.Ptr(true),
		})
		if err == nil {
			t.Fatal("expected error, but not")
		}
	})

	t.Run("unsupported (ssh) parent URL", func(t *testing.T) {
		tim, _ := time.Parse("2006-01-02", "2021-01-01")
		_, err := ingestRepository(&github.Repository{
			Description: util.Ptr("valid description"),
			Homepage:    util.Ptr("valid homepage"),
			UpdatedAt:   &exgithub.Timestamp{Time: tim},
			CloneURL:    util.Ptr("https://github.com/kyoh86/gogh"),
			Language:    util.Ptr("valid language"),
			Fork:        util.Ptr(false),
			Parent: &github.Repository{
				CloneURL: util.Ptr("git@github.com:kyoh86-tryouts/test"),
			},
			Archived:   util.Ptr(true),
			Private:    util.Ptr(false),
			IsTemplate: util.Ptr(true),
		})
		if err == nil {
			t.Fatal("expected error, but not")
		}
	})
}

func TestIngestRepositoryFragment(t *testing.T) {
	tim, _ := time.Parse("2006-01-02", "2021-01-01")
	t.Run("valid repository", func(t *testing.T) {
		wantRef, err := reporef.NewRepoRef("github.com", "kyoh86", "gogh")
		if err != nil {
			t.Fatalf("unexpected error: %q", err)
		}
		wantParent, err := reporef.NewRepoRef("github.com", "kyoh86-tryouts", "test")
		if err != nil {
			t.Fatalf("unexpected error: %q", err)
		}
		want := RemoteRepo{
			URL:         "https://github.com/kyoh86/gogh",
			Description: "valid description",
			Homepage:    "https://example.com",
			Language:    "valid language",
			UpdatedAt:   tim,
			Archived:    true,
			Private:     false,
			IsTemplate:  true,
			Fork:        false,
			Ref:         wantRef,
			Parent:      &wantParent,
		}
		got, err := ingestRepositoryFragment("github.com", &github.RepositoryFragment{
			Owner:           &githubv4.RepositoryFragmentOwnerUser{OwnerFragmentUser: githubv4.OwnerFragmentUser{Login: "kyoh86"}},
			Name:            "gogh",
			Description:     "valid description",
			HomepageUrl:     "https://example.com",
			UpdatedAt:       tim,
			Url:             "https://github.com/kyoh86/gogh",
			PrimaryLanguage: githubv4.RepositoryFragmentPrimaryLanguage{LanguageFragment: githubv4.LanguageFragment{Name: "valid language"}},
			IsFork:          false,
			Parent: githubv4.RepositoryFragmentParentRepository{
				ParentRepositoryFragment: githubv4.ParentRepositoryFragment{
					Owner: &githubv4.ParentRepositoryFragmentOwnerOrganization{OwnerFragmentOrganization: githubv4.OwnerFragmentOrganization{Login: "kyoh86-tryouts"}},
					Name:  "test",
				},
			},
			IsArchived: true,
			IsPrivate:  false,
			IsTemplate: true,
		})
		if err != nil {
			t.Fatalf("unexpected error: %q", err)
		}
		if diff := cmp.Diff(want, got, cmp.AllowUnexported(want), cmp.AllowUnexported(want.Ref)); diff != "" {
			t.Errorf("result mistmatch\n-want, +got:\n%s", diff)
		}
	})

	t.Run("invalid owner", func(t *testing.T) {
		_, err := ingestRepositoryFragment("github.com", &github.RepositoryFragment{
			Owner:      &githubv4.RepositoryFragmentOwnerUser{OwnerFragmentUser: githubv4.OwnerFragmentUser{Login: ".."}},
			Name:       "gogh",
			Url:        "https://github.com/kyoh86/gogh",
			UpdatedAt:  tim,
			IsFork:     false,
			IsArchived: true,
			IsPrivate:  false,
			IsTemplate: true,
		})
		if err == nil {
			t.Fatal("expected error, but not")
		}
	})

	t.Run("invalid parent", func(t *testing.T) {
		_, err := ingestRepositoryFragment("github.com", &github.RepositoryFragment{
			Owner:     &githubv4.RepositoryFragmentOwnerUser{OwnerFragmentUser: githubv4.OwnerFragmentUser{Login: "kyoh86"}},
			Name:      "gogh",
			Url:       "https://github.com/kyoh86/gogh",
			UpdatedAt: tim,
			Parent: githubv4.RepositoryFragmentParentRepository{
				ParentRepositoryFragment: githubv4.ParentRepositoryFragment{
					Owner: &githubv4.ParentRepositoryFragmentOwnerOrganization{OwnerFragmentOrganization: githubv4.OwnerFragmentOrganization{Login: ".."}},
					Name:  "..",
				},
			},
			IsFork:     false,
			IsArchived: true,
			IsPrivate:  false,
			IsTemplate: true,
		})
		if err == nil {
			t.Fatal("expected error, but not")
		}
	})
}
