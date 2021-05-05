package gogh

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	exgithub "github.com/google/go-github/v35/github"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/kyoh86/gogh/v2/internal/githubv4"
	"github.com/wacul/ptr"
)

func TestParseSpec(t *testing.T) {
	t.Run("valid http URL", func(t *testing.T) {
		want := Spec{host: "github.com", owner: "kyoh86", name: "gogh"}
		got, err := parseSpec(&github.Repository{
			CloneURL: ptr.String("https://github.com/kyoh86/gogh"),
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
			CloneURL: ptr.String("git@github.com:kyoh86/gogh.git"),
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
			PushedAt:    tim,
			Archived:    true,
			Private:     false,
			IsTemplate:  true,
			Fork:        false,
			Spec:        Spec{host: "github.com", owner: "kyoh86", name: "gogh"},
			Parent:      &Spec{host: "github.com", owner: "kyoh86-tryouts", name: "test"},
		}
		got, err := ingestRepository(&github.Repository{
			Description: ptr.String("valid description"),
			Homepage:    ptr.String("valid homepage"),
			PushedAt:    &exgithub.Timestamp{Time: tim},
			CloneURL:    ptr.String("https://github.com/kyoh86/gogh"),
			Language:    ptr.String("valid language"),
			Fork:        ptr.Bool(false),
			Parent: &github.Repository{
				CloneURL: ptr.String("https://github.com/kyoh86-tryouts/test"),
			},
			Archived:   ptr.Bool(true),
			Private:    ptr.Bool(false),
			IsTemplate: ptr.Bool(true),
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
			Description: ptr.String("valid description"),
			Homepage:    ptr.String("valid homepage"),
			PushedAt:    &exgithub.Timestamp{Time: tim},
			CloneURL:    ptr.String("git@github.com:kyoh86/gogh.git"),
			Language:    ptr.String("valid language"),
			Fork:        ptr.Bool(false),
			Parent: &github.Repository{
				CloneURL: ptr.String("https://github.com/kyoh86-tryouts/test"),
			},
			Archived:   ptr.Bool(true),
			Private:    ptr.Bool(false),
			IsTemplate: ptr.Bool(true),
		})
		if err == nil {
			t.Fatal("expected error, but not")
		}
	})

	t.Run("unsupported (ssh) parent URL", func(t *testing.T) {
		tim, _ := time.Parse("2006-01-02", "2021-01-01")
		_, err := ingestRepository(&github.Repository{
			Description: ptr.String("valid description"),
			Homepage:    ptr.String("valid homepage"),
			PushedAt:    &exgithub.Timestamp{Time: tim},
			CloneURL:    ptr.String("https://github.com/kyoh86/gogh"),
			Language:    ptr.String("valid language"),
			Fork:        ptr.Bool(false),
			Parent: &github.Repository{
				CloneURL: ptr.String("git@github.com:kyoh86-tryouts/test"),
			},
			Archived:   ptr.Bool(true),
			Private:    ptr.Bool(false),
			IsTemplate: ptr.Bool(true),
		})
		if err == nil {
			t.Fatal("expected error, but not")
		}
	})
}

func TestIngestRepositoryFragment(t *testing.T) {
	t.Run("valid repository", func(t *testing.T) {
		tim, _ := time.Parse("2006-01-02", "2021-01-01")
		want := Repository{
			URL:         "https://github.com/kyoh86/gogh",
			Description: "valid description",
			Homepage:    "valid homepage",
			Language:    "valid language",
			PushedAt:    tim,
			Archived:    true,
			Private:     false,
			IsTemplate:  true,
			Fork:        false,
			Spec:        Spec{host: "github.com", owner: "kyoh86", name: "gogh"},
			Parent:      &Spec{host: "github.com", owner: "kyoh86-tryouts", name: "test"},
		}
		got, err := ingestRepositoryFragment("github.com", &github.RepositoryFragment{
			Owner:           githubv4.OwnerFragment{Login: "kyoh86"},
			Name:            "gogh",
			Description:     ptr.String("valid description"),
			HomepageURL:     ptr.String("valid homepage"),
			PushedAt:        ptr.String(tim.Format(time.RFC3339)),
			URL:             "https://github.com/kyoh86/gogh",
			PrimaryLanguage: &githubv4.LanguageFragment{Name: "valid language"},
			IsFork:          false,
			Parent: &github.ParentRepositoryFragment{
				Owner: githubv4.OwnerFragment{Login: "kyoh86-tryouts"},
				Name:  "test",
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

	t.Run("invalid time", func(t *testing.T) {
		_, err := ingestRepositoryFragment("github.com", &github.RepositoryFragment{
			Owner:      githubv4.OwnerFragment{Login: "kyoh86"},
			Name:       "gogh",
			URL:        "https://github.com/kyoh86/gogh",
			PushedAt:   ptr.String("invalid time"),
			IsFork:     false,
			IsArchived: true,
			IsPrivate:  false,
			IsTemplate: true,
		})
		if err == nil {
			t.Fatal("expected error, but not")
		}
	})

	t.Run("invalid owner", func(t *testing.T) {
		_, err := ingestRepositoryFragment("github.com", &github.RepositoryFragment{
			Owner:      githubv4.OwnerFragment{Login: ".."},
			Name:       "gogh",
			URL:        "https://github.com/kyoh86/gogh",
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
			Owner: githubv4.OwnerFragment{Login: "kyoh86"},
			Name:  "gogh",
			URL:   "https://github.com/kyoh86/gogh",
			Parent: &github.ParentRepositoryFragment{
				Owner: githubv4.OwnerFragment{Login: ".."},
				Name:  "..",
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
