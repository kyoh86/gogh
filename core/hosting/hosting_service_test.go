package hosting_test

import (
	"context"
	"errors"
	"iter"
	"net/url"
	"testing"
	"time"

	"github.com/kyoh86/gogh/v4/core/auth"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/typ"
	"golang.org/x/oauth2"
)

// mockHostingService implements HostingService for testing
type mockHostingService struct {
	getURLOfFunc                     func(repository.Reference) (*url.URL, error)
	parseURLFunc                     func(*url.URL) (*repository.Reference, error)
	getTokenForFunc                  func(context.Context, string, string) (string, auth.Token, error)
	getRepositoryFunc                func(context.Context, repository.Reference) (*hosting.Repository, error)
	listRepositoryFunc               func(context.Context, hosting.ListRepositoryOptions) iter.Seq2[*hosting.Repository, error]
	deleteRepositoryFunc             func(context.Context, repository.Reference) error
	createRepositoryFunc             func(context.Context, repository.Reference, hosting.CreateRepositoryOptions) (*hosting.Repository, error)
	createRepositoryFromTemplateFunc func(context.Context, repository.Reference, repository.Reference, hosting.CreateRepositoryFromTemplateOptions) (*hosting.Repository, error)
	forkRepositoryFunc               func(context.Context, repository.Reference, repository.Reference, hosting.ForkRepositoryOptions) (*hosting.Repository, error)
}

func (m *mockHostingService) GetURLOf(ref repository.Reference) (*url.URL, error) {
	if m.getURLOfFunc != nil {
		return m.getURLOfFunc(ref)
	}
	return url.Parse("https://" + ref.Host() + "/" + ref.Owner() + "/" + ref.Name())
}

func (m *mockHostingService) ParseURL(u *url.URL) (*repository.Reference, error) {
	if m.parseURLFunc != nil {
		return m.parseURLFunc(u)
	}
	return nil, errors.New("not implemented")
}

func (m *mockHostingService) GetTokenFor(ctx context.Context, host, owner string) (string, auth.Token, error) {
	if m.getTokenForFunc != nil {
		return m.getTokenForFunc(ctx, host, owner)
	}
	return "user", oauth2.Token{AccessToken: "token"}, nil
}

func (m *mockHostingService) GetRepository(ctx context.Context, ref repository.Reference) (*hosting.Repository, error) {
	if m.getRepositoryFunc != nil {
		return m.getRepositoryFunc(ctx, ref)
	}
	return &hosting.Repository{
		Ref:       ref,
		URL:       "https://" + ref.Host() + "/" + ref.Owner() + "/" + ref.Name(),
		CloneURL:  "https://" + ref.Host() + "/" + ref.Owner() + "/" + ref.Name() + ".git",
		UpdatedAt: time.Now(),
	}, nil
}

func (m *mockHostingService) ListRepository(ctx context.Context, opts hosting.ListRepositoryOptions) iter.Seq2[*hosting.Repository, error] {
	if m.listRepositoryFunc != nil {
		return m.listRepositoryFunc(ctx, opts)
	}
	return func(yield func(*hosting.Repository, error) bool) {}
}

func (m *mockHostingService) DeleteRepository(ctx context.Context, ref repository.Reference) error {
	if m.deleteRepositoryFunc != nil {
		return m.deleteRepositoryFunc(ctx, ref)
	}
	return nil
}

func (m *mockHostingService) CreateRepository(ctx context.Context, ref repository.Reference, opts hosting.CreateRepositoryOptions) (*hosting.Repository, error) {
	if m.createRepositoryFunc != nil {
		return m.createRepositoryFunc(ctx, ref, opts)
	}
	return &hosting.Repository{
		Ref:         ref,
		URL:         "https://" + ref.Host() + "/" + ref.Owner() + "/" + ref.Name(),
		CloneURL:    "https://" + ref.Host() + "/" + ref.Owner() + "/" + ref.Name() + ".git",
		UpdatedAt:   time.Now(),
		Description: opts.Description,
		Homepage:    opts.Homepage,
		Private:     opts.Private,
		IsTemplate:  opts.IsTemplate,
	}, nil
}

func (m *mockHostingService) CreateRepositoryFromTemplate(ctx context.Context, ref repository.Reference, template repository.Reference, opts hosting.CreateRepositoryFromTemplateOptions) (*hosting.Repository, error) {
	if m.createRepositoryFromTemplateFunc != nil {
		return m.createRepositoryFromTemplateFunc(ctx, ref, template, opts)
	}
	return &hosting.Repository{
		Ref:         ref,
		URL:         "https://" + ref.Host() + "/" + ref.Owner() + "/" + ref.Name(),
		CloneURL:    "https://" + ref.Host() + "/" + ref.Owner() + "/" + ref.Name() + ".git",
		UpdatedAt:   time.Now(),
		Description: opts.Description,
		Private:     opts.Private,
	}, nil
}

func (m *mockHostingService) ForkRepository(ctx context.Context, ref repository.Reference, target repository.Reference, opts hosting.ForkRepositoryOptions) (*hosting.Repository, error) {
	if m.forkRepositoryFunc != nil {
		return m.forkRepositoryFunc(ctx, ref, target, opts)
	}
	parent := &hosting.ParentRepository{
		Ref:      ref,
		CloneURL: "https://" + ref.Host() + "/" + ref.Owner() + "/" + ref.Name() + ".git",
	}
	return &hosting.Repository{
		Ref:       target,
		URL:       "https://" + target.Host() + "/" + target.Owner() + "/" + target.Name(),
		CloneURL:  "https://" + target.Host() + "/" + target.Owner() + "/" + target.Name() + ".git",
		UpdatedAt: time.Now(),
		Fork:      true,
		Parent:    parent,
	}, nil
}

func TestHostingService_GetURLOf(t *testing.T) {
	ref := repository.NewReference("github.com", "kyoh86", "gogh")

	t.Run("success", func(t *testing.T) {
		service := &mockHostingService{}
		u, err := service.GetURLOf(ref)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "https://github.com/kyoh86/gogh"
		if u.String() != expected {
			t.Errorf("expected URL %s, got %s", expected, u.String())
		}
	})

	t.Run("custom implementation", func(t *testing.T) {
		service := &mockHostingService{
			getURLOfFunc: func(r repository.Reference) (*url.URL, error) {
				return url.Parse("ssh://git@" + r.Host() + "/" + r.Owner() + "/" + r.Name() + ".git")
			},
		}
		u, err := service.GetURLOf(ref)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "ssh://git@github.com/kyoh86/gogh.git"
		if u.String() != expected {
			t.Errorf("expected URL %s, got %s", expected, u.String())
		}
	})

	t.Run("error", func(t *testing.T) {
		service := &mockHostingService{
			getURLOfFunc: func(r repository.Reference) (*url.URL, error) {
				return nil, errors.New("failed to get URL")
			},
		}
		_, err := service.GetURLOf(ref)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestHostingService_ParseURL(t *testing.T) {
	u, _ := url.Parse("https://github.com/kyoh86/gogh")

	t.Run("error", func(t *testing.T) {
		service := &mockHostingService{}
		_, err := service.ParseURL(u)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("success", func(t *testing.T) {
		service := &mockHostingService{
			parseURLFunc: func(u *url.URL) (*repository.Reference, error) {
				// Simple implementation for test
				owner := "kyoh86"
				name := "gogh"
				ref := repository.NewReference(u.Host, owner, name)
				return &ref, nil
			},
		}
		ref, err := service.ParseURL(u)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ref.Host() != "github.com" || ref.Owner() != "kyoh86" || ref.Name() != "gogh" {
			t.Errorf("unexpected reference: %v", ref)
		}
	})
}

func TestHostingService_GetTokenFor(t *testing.T) {
	ctx := context.Background()

	t.Run("default", func(t *testing.T) {
		service := &mockHostingService{}
		user, token, err := service.GetTokenFor(ctx, "github.com", "kyoh86")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user != "user" {
			t.Errorf("expected user 'user', got %s", user)
		}
		if token.AccessToken != "token" {
			t.Errorf("expected token 'token', got %s", token.AccessToken)
		}
	})

	t.Run("custom", func(t *testing.T) {
		service := &mockHostingService{
			getTokenForFunc: func(ctx context.Context, host, owner string) (string, auth.Token, error) {
				return owner, oauth2.Token{AccessToken: "secret-" + host}, nil
			},
		}
		user, token, err := service.GetTokenFor(ctx, "github.com", "kyoh86")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user != "kyoh86" {
			t.Errorf("expected user 'kyoh86', got %s", user)
		}
		if token.AccessToken != "secret-github.com" {
			t.Errorf("expected token 'secret-github.com', got %s", token.AccessToken)
		}
	})

	t.Run("error", func(t *testing.T) {
		service := &mockHostingService{
			getTokenForFunc: func(ctx context.Context, host, owner string) (string, auth.Token, error) {
				return "", oauth2.Token{}, errors.New("authentication failed")
			},
		}
		_, _, err := service.GetTokenFor(ctx, "github.com", "kyoh86")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestHostingService_GetRepository(t *testing.T) {
	ctx := context.Background()
	ref := repository.NewReference("github.com", "kyoh86", "gogh")

	t.Run("default", func(t *testing.T) {
		service := &mockHostingService{}
		repo, err := service.GetRepository(ctx, ref)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.Ref.String() != ref.String() {
			t.Errorf("expected ref %s, got %s", ref.String(), repo.Ref.String())
		}
		if repo.URL != "https://github.com/kyoh86/gogh" {
			t.Errorf("unexpected URL: %s", repo.URL)
		}
		if repo.CloneURL != "https://github.com/kyoh86/gogh.git" {
			t.Errorf("unexpected CloneURL: %s", repo.CloneURL)
		}
	})

	t.Run("with metadata", func(t *testing.T) {
		service := &mockHostingService{
			getRepositoryFunc: func(ctx context.Context, ref repository.Reference) (*hosting.Repository, error) {
				return &hosting.Repository{
					Ref:         ref,
					URL:         "https://github.com/kyoh86/gogh",
					CloneURL:    "https://github.com/kyoh86/gogh.git",
					UpdatedAt:   time.Now(),
					Description: "A repository manager",
					Language:    "Go",
					Private:     true,
					Archived:    true,
				}, nil
			},
		}
		repo, err := service.GetRepository(ctx, ref)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.Description != "A repository manager" {
			t.Errorf("expected description 'A repository manager', got %s", repo.Description)
		}
		if repo.Language != "Go" {
			t.Errorf("expected language 'Go', got %s", repo.Language)
		}
		if !repo.Private {
			t.Error("expected private repository")
		}
		if !repo.Archived {
			t.Error("expected archived repository")
		}
	})

	t.Run("error", func(t *testing.T) {
		service := &mockHostingService{
			getRepositoryFunc: func(ctx context.Context, ref repository.Reference) (*hosting.Repository, error) {
				return nil, errors.New("repository not found")
			},
		}
		_, err := service.GetRepository(ctx, ref)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestHostingService_ListRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("empty list", func(t *testing.T) {
		service := &mockHostingService{}
		opts := hosting.ListRepositoryOptions{}
		var count int
		for _, err := range service.ListRepository(ctx, opts) {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			count++
		}
		if count != 0 {
			t.Errorf("expected 0 repositories, got %d", count)
		}
	})

	t.Run("with repositories", func(t *testing.T) {
		repos := []*hosting.Repository{
			{
				Ref:      repository.NewReference("github.com", "kyoh86", "gogh"),
				URL:      "https://github.com/kyoh86/gogh",
				CloneURL: "https://github.com/kyoh86/gogh.git",
			},
			{
				Ref:      repository.NewReference("github.com", "kyoh86", "dotfiles"),
				URL:      "https://github.com/kyoh86/dotfiles",
				CloneURL: "https://github.com/kyoh86/dotfiles.git",
			},
		}

		service := &mockHostingService{
			listRepositoryFunc: func(ctx context.Context, opts hosting.ListRepositoryOptions) iter.Seq2[*hosting.Repository, error] {
				return func(yield func(*hosting.Repository, error) bool) {
					for _, repo := range repos {
						if !yield(repo, nil) {
							return
						}
					}
				}
			},
		}

		opts := hosting.ListRepositoryOptions{
			OrderBy: hosting.RepositoryOrder{
				Direction: hosting.OrderDirectionDesc,
				Field:     hosting.RepositoryOrderFieldUpdatedAt,
			},
			Privacy:           hosting.RepositoryPrivacyPublic,
			OwnerAffiliations: []hosting.RepositoryAffiliation{hosting.RepositoryAffiliationOwner},
			Limit:             10,
			IsFork:            typ.TristateTrue,
			IsArchived:        typ.TristateFalse,
		}

		var count int
		for repo, err := range service.ListRepository(ctx, opts) {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if count < len(repos) && repo.Ref.String() != repos[count].Ref.String() {
				t.Errorf("unexpected repository at index %d: %v", count, repo.Ref)
			}
			count++
		}
		if count != len(repos) {
			t.Errorf("expected %d repositories, got %d", len(repos), count)
		}
	})

	t.Run("with error", func(t *testing.T) {
		service := &mockHostingService{
			listRepositoryFunc: func(ctx context.Context, opts hosting.ListRepositoryOptions) iter.Seq2[*hosting.Repository, error] {
				return func(yield func(*hosting.Repository, error) bool) {
					yield(nil, errors.New("list failed"))
				}
			},
		}

		opts := hosting.ListRepositoryOptions{}
		for _, err := range service.ListRepository(ctx, opts) {
			if err == nil {
				t.Error("expected error, got nil")
			}
			break
		}
	})
}

func TestHostingService_DeleteRepository(t *testing.T) {
	ctx := context.Background()
	ref := repository.NewReference("github.com", "kyoh86", "test-repo")

	t.Run("success", func(t *testing.T) {
		service := &mockHostingService{}
		err := service.DeleteRepository(ctx, ref)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		service := &mockHostingService{
			deleteRepositoryFunc: func(ctx context.Context, ref repository.Reference) error {
				return errors.New("permission denied")
			},
		}
		err := service.DeleteRepository(ctx, ref)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestHostingService_CreateRepository(t *testing.T) {
	ctx := context.Background()
	ref := repository.NewReference("github.com", "kyoh86", "new-repo")

	t.Run("default", func(t *testing.T) {
		service := &mockHostingService{}
		opts := hosting.CreateRepositoryOptions{
			Description: "A new repository",
			Private:     true,
		}
		repo, err := service.CreateRepository(ctx, ref, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.Ref.String() != ref.String() {
			t.Errorf("expected ref %s, got %s", ref.String(), repo.Ref.String())
		}
		if repo.Description != "A new repository" {
			t.Errorf("expected description 'A new repository', got %s", repo.Description)
		}
		if !repo.Private {
			t.Error("expected private repository")
		}
	})

	t.Run("with all options", func(t *testing.T) {
		service := &mockHostingService{
			createRepositoryFunc: func(ctx context.Context, ref repository.Reference, opts hosting.CreateRepositoryOptions) (*hosting.Repository, error) {
				return &hosting.Repository{
					Ref:         ref,
					URL:         "https://github.com/org/new-repo",
					CloneURL:    "https://github.com/org/new-repo.git",
					UpdatedAt:   time.Now(),
					Description: opts.Description,
					Homepage:    opts.Homepage,
					Private:     opts.Private,
					IsTemplate:  opts.IsTemplate,
				}, nil
			},
		}

		opts := hosting.CreateRepositoryOptions{
			Description:         "Org repository",
			Homepage:            "https://example.com",
			Organization:        "org",
			LicenseTemplate:     "mit",
			GitignoreTemplate:   "Go",
			TeamID:              123,
			DisableDownloads:    true,
			IsTemplate:          true,
			Private:             true,
			DisableWiki:         true,
			AutoInit:            true,
			DisableProjects:     true,
			DisableIssues:       true,
			PreventSquashMerge:  true,
			PreventMergeCommit:  true,
			PreventRebaseMerge:  true,
			DeleteBranchOnMerge: true,
		}

		repo, err := service.CreateRepository(ctx, ref, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.Homepage != "https://example.com" {
			t.Errorf("expected homepage 'https://example.com', got %s", repo.Homepage)
		}
		if !repo.IsTemplate {
			t.Error("expected template repository")
		}
	})

	t.Run("error", func(t *testing.T) {
		service := &mockHostingService{
			createRepositoryFunc: func(ctx context.Context, ref repository.Reference, opts hosting.CreateRepositoryOptions) (*hosting.Repository, error) {
				return nil, errors.New("repository already exists")
			},
		}
		opts := hosting.CreateRepositoryOptions{}
		_, err := service.CreateRepository(ctx, ref, opts)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestHostingService_CreateRepositoryFromTemplate(t *testing.T) {
	ctx := context.Background()
	ref := repository.NewReference("github.com", "kyoh86", "new-from-template")
	template := repository.NewReference("github.com", "org", "template-repo")

	t.Run("success", func(t *testing.T) {
		service := &mockHostingService{}
		opts := hosting.CreateRepositoryFromTemplateOptions{
			Description:        "Created from template",
			IncludeAllBranches: true,
			Private:            true,
		}
		repo, err := service.CreateRepositoryFromTemplate(ctx, ref, template, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.Ref.String() != ref.String() {
			t.Errorf("expected ref %s, got %s", ref.String(), repo.Ref.String())
		}
		if repo.Description != "Created from template" {
			t.Errorf("expected description 'Created from template', got %s", repo.Description)
		}
		if !repo.Private {
			t.Error("expected private repository")
		}
	})

	t.Run("error", func(t *testing.T) {
		service := &mockHostingService{
			createRepositoryFromTemplateFunc: func(ctx context.Context, ref, template repository.Reference, opts hosting.CreateRepositoryFromTemplateOptions) (*hosting.Repository, error) {
				return nil, errors.New("template not found")
			},
		}
		opts := hosting.CreateRepositoryFromTemplateOptions{}
		_, err := service.CreateRepositoryFromTemplate(ctx, ref, template, opts)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestHostingService_ForkRepository(t *testing.T) {
	ctx := context.Background()
	ref := repository.NewReference("github.com", "original", "repo")
	target := repository.NewReference("github.com", "kyoh86", "repo")

	t.Run("success", func(t *testing.T) {
		service := &mockHostingService{}
		opts := hosting.ForkRepositoryOptions{
			DefaultBranchOnly: true,
		}
		repo, err := service.ForkRepository(ctx, ref, target, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.Ref.String() != target.String() {
			t.Errorf("expected ref %s, got %s", target.String(), repo.Ref.String())
		}
		if !repo.Fork {
			t.Error("expected fork flag to be true")
		}
		if repo.Parent == nil {
			t.Fatal("expected parent repository")
		}
		if repo.Parent.Ref.String() != ref.String() {
			t.Errorf("expected parent ref %s, got %s", ref.String(), repo.Parent.Ref.String())
		}
	})

	t.Run("custom parent", func(t *testing.T) {
		service := &mockHostingService{
			forkRepositoryFunc: func(ctx context.Context, ref, target repository.Reference, opts hosting.ForkRepositoryOptions) (*hosting.Repository, error) {
				parent := &hosting.ParentRepository{
					Ref:      ref,
					CloneURL: "git@github.com:original/repo.git",
				}
				return &hosting.Repository{
					Ref:       target,
					URL:       "https://github.com/kyoh86/repo",
					CloneURL:  "git@github.com:kyoh86/repo.git",
					UpdatedAt: time.Now(),
					Fork:      true,
					Parent:    parent,
				}, nil
			},
		}
		opts := hosting.ForkRepositoryOptions{}
		repo, err := service.ForkRepository(ctx, ref, target, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.CloneURL != "git@github.com:kyoh86/repo.git" {
			t.Errorf("unexpected clone URL: %s", repo.CloneURL)
		}
		if repo.Parent.CloneURL != "git@github.com:original/repo.git" {
			t.Errorf("unexpected parent clone URL: %s", repo.Parent.CloneURL)
		}
	})

	t.Run("error", func(t *testing.T) {
		service := &mockHostingService{
			forkRepositoryFunc: func(ctx context.Context, ref, target repository.Reference, opts hosting.ForkRepositoryOptions) (*hosting.Repository, error) {
				return nil, errors.New("cannot fork")
			},
		}
		opts := hosting.ForkRepositoryOptions{}
		_, err := service.ForkRepository(ctx, ref, target, opts)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestRepositoryOrderConstants(t *testing.T) {
	// Test OrderDirection constants
	if hosting.OrderDirectionAsc != 1 {
		t.Errorf("OrderDirectionAsc = %d, want 1", hosting.OrderDirectionAsc)
	}
	if hosting.OrderDirectionDesc != 2 {
		t.Errorf("OrderDirectionDesc = %d, want 2", hosting.OrderDirectionDesc)
	}

	// Test RepositoryOrderField constants
	if hosting.RepositoryOrderFieldCreatedAt != 0 {
		t.Errorf("RepositoryOrderFieldCreatedAt = %d, want 0", hosting.RepositoryOrderFieldCreatedAt)
	}
	if hosting.RepositoryOrderFieldName != 1 {
		t.Errorf("RepositoryOrderFieldName = %d, want 1", hosting.RepositoryOrderFieldName)
	}
	if hosting.RepositoryOrderFieldPushedAt != 2 {
		t.Errorf("RepositoryOrderFieldPushedAt = %d, want 2", hosting.RepositoryOrderFieldPushedAt)
	}
	if hosting.RepositoryOrderFieldStargazers != 3 {
		t.Errorf("RepositoryOrderFieldStargazers = %d, want 3", hosting.RepositoryOrderFieldStargazers)
	}
	if hosting.RepositoryOrderFieldUpdatedAt != 4 {
		t.Errorf("RepositoryOrderFieldUpdatedAt = %d, want 4", hosting.RepositoryOrderFieldUpdatedAt)
	}
}

func TestRepositoryPrivacyConstants(t *testing.T) {
	if hosting.RepositoryPrivacyNone != 0 {
		t.Errorf("RepositoryPrivacyNone = %d, want 0", hosting.RepositoryPrivacyNone)
	}
	if hosting.RepositoryPrivacyPrivate != 1 {
		t.Errorf("RepositoryPrivacyPrivate = %d, want 1", hosting.RepositoryPrivacyPrivate)
	}
	if hosting.RepositoryPrivacyPublic != 2 {
		t.Errorf("RepositoryPrivacyPublic = %d, want 2", hosting.RepositoryPrivacyPublic)
	}
}

func TestRepositoryAffiliationConstants(t *testing.T) {
	if hosting.RepositoryAffiliationCollaborator != 0 {
		t.Errorf("RepositoryAffiliationCollaborator = %d, want 0", hosting.RepositoryAffiliationCollaborator)
	}
	if hosting.RepositoryAffiliationOrganizationMember != 1 {
		t.Errorf("RepositoryAffiliationOrganizationMember = %d, want 1", hosting.RepositoryAffiliationOrganizationMember)
	}
	if hosting.RepositoryAffiliationOwner != 2 {
		t.Errorf("RepositoryAffiliationOwner = %d, want 2", hosting.RepositoryAffiliationOwner)
	}
}
