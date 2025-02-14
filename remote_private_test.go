package gogh

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	exgithub "github.com/google/go-github/v69/github"
	"github.com/kyoh86/gogh/v3/internal/github"
	"github.com/kyoh86/gogh/v3/internal/githubv4"
)

func TestParseSpec(t *testing.T) {
	t.Run("valid http URL", func(t *testing.T) {
		want := Spec{host: "github.com", owner: "kyoh86", name: "gogh"}
		got, err := parseSpec(&github.Repository{
			CloneURL: Ptr("https://github.com/kyoh86/gogh"),
		})
		if err != nil {
			t.Fatalf("unexpected error: %q", err)
		}
		if diff := cmp.Diff(want, got, cmp.AllowUnexported(want)); diff != "" {
			t.Errorf("result mistmatch\n-want, +got:\n%s", diff)
		}
	})

	t.Run("unsupported (ssh) URL", func(t *testing.T) {
		_, err := parseSpec(&github.Repository{
			CloneURL: Ptr("git@github.com:kyoh86/gogh.git"),
		})
		if err == nil {
			t.Fatal("expected error, but not")
		}
	})
}

func TestIngestRepository(t *testing.T) {
	t.Run("valid repository", func(t *testing.T) {
		tim, _ := time.Parse("2006-01-02", "2021-01-01")
		want := Repository{
			URL:         "https://github.com/kyoh86/gogh",
			Description: "valid description",
			Homepage:    "valid homepage",
			Language:    "valid language",
			UpdatedAt:   tim,
			Archived:    true,
			Private:     false,
			IsTemplate:  true,
			Fork:        false,
			Spec:        Spec{host: "github.com", owner: "kyoh86", name: "gogh"},
			Parent:      &Spec{host: "github.com", owner: "kyoh86-tryouts", name: "test"},
		}
		got, err := ingestRepository(&github.Repository{
			Description: Ptr("valid description"),
			Homepage:    Ptr("valid homepage"),
			UpdatedAt:   &exgithub.Timestamp{Time: tim},
			CloneURL:    Ptr("https://github.com/kyoh86/gogh"),
			Language:    Ptr("valid language"),
			Fork:        Ptr(false),
			Parent: &github.Repository{
				CloneURL: Ptr("https://github.com/kyoh86-tryouts/test"),
			},
			Archived:   Ptr(true),
			Private:    Ptr(false),
			IsTemplate: Ptr(true),
		})
		if err != nil {
			t.Fatalf("unexpected error: %q", err)
		}
		if diff := cmp.Diff(want, got, cmp.AllowUnexported(want), cmp.AllowUnexported(want.Spec)); diff != "" {
			t.Errorf("result mistmatch\n-want, +got:\n%s", diff)
		}
	})

	t.Run("unsupported (ssh) URL", func(t *testing.T) {
		tim, _ := time.Parse("2006-01-02", "2021-01-01")
		_, err := ingestRepository(&github.Repository{
			Description: Ptr("valid description"),
			Homepage:    Ptr("valid homepage"),
			UpdatedAt:   &exgithub.Timestamp{Time: tim},
			CloneURL:    Ptr("git@github.com:kyoh86/gogh.git"),
			Language:    Ptr("valid language"),
			Fork:        Ptr(false),
			Parent: &github.Repository{
				CloneURL: Ptr("https://github.com/kyoh86-tryouts/test"),
			},
			Archived:   Ptr(true),
			Private:    Ptr(false),
			IsTemplate: Ptr(true),
		})
		if err == nil {
			t.Fatal("expected error, but not")
		}
	})

	t.Run("unsupported (ssh) parent URL", func(t *testing.T) {
		tim, _ := time.Parse("2006-01-02", "2021-01-01")
		_, err := ingestRepository(&github.Repository{
			Description: Ptr("valid description"),
			Homepage:    Ptr("valid homepage"),
			UpdatedAt:   &exgithub.Timestamp{Time: tim},
			CloneURL:    Ptr("https://github.com/kyoh86/gogh"),
			Language:    Ptr("valid language"),
			Fork:        Ptr(false),
			Parent: &github.Repository{
				CloneURL: Ptr("git@github.com:kyoh86-tryouts/test"),
			},
			Archived:   Ptr(true),
			Private:    Ptr(false),
			IsTemplate: Ptr(true),
		})
		if err == nil {
			t.Fatal("expected error, but not")
		}
	})
}

func TestIngestRepositoryFragment(t *testing.T) {
	tim, _ := time.Parse("2006-01-02", "2021-01-01")
	t.Run("valid repository", func(t *testing.T) {
		want := Repository{
			URL:         "https://github.com/kyoh86/gogh",
			Description: "valid description",
			Homepage:    "https://example.com",
			Language:    "valid language",
			UpdatedAt:   tim,
			Archived:    true,
			Private:     false,
			IsTemplate:  true,
			Fork:        false,
			Spec:        Spec{host: "github.com", owner: "kyoh86", name: "gogh"},
			Parent:      &Spec{host: "github.com", owner: "kyoh86-tryouts", name: "test"},
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
		if diff := cmp.Diff(want, got, cmp.AllowUnexported(want), cmp.AllowUnexported(want.Spec)); diff != "" {
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
