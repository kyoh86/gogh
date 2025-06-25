package hosting_test

import (
	"testing"
	"time"

	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/repository"
)

func TestRepository(t *testing.T) {
	// Create test data
	ref := repository.NewReference("github.com", "kyoh86", "gogh")
	parentRef := repository.NewReference("github.com", "original", "gogh")
	updatedAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	t.Run("basic repository", func(t *testing.T) {
		repo := hosting.Repository{
			Ref:       ref,
			URL:       "https://github.com/kyoh86/gogh",
			CloneURL:  "https://github.com/kyoh86/gogh.git",
			UpdatedAt: updatedAt,
		}

		if repo.Ref.String() != "github.com/kyoh86/gogh" {
			t.Errorf("expected ref 'github.com/kyoh86/gogh', got %s", repo.Ref.String())
		}
		if repo.URL != "https://github.com/kyoh86/gogh" {
			t.Errorf("expected URL 'https://github.com/kyoh86/gogh', got %s", repo.URL)
		}
		if repo.CloneURL != "https://github.com/kyoh86/gogh.git" {
			t.Errorf("expected CloneURL 'https://github.com/kyoh86/gogh.git', got %s", repo.CloneURL)
		}
		if !repo.UpdatedAt.Equal(updatedAt) {
			t.Errorf("expected UpdatedAt %v, got %v", updatedAt, repo.UpdatedAt)
		}
		if repo.Parent != nil {
			t.Error("expected nil parent")
		}
		if repo.Fork {
			t.Error("expected Fork to be false")
		}
	})

	t.Run("repository with metadata", func(t *testing.T) {
		repo := hosting.Repository{
			Ref:         ref,
			URL:         "https://github.com/kyoh86/gogh",
			CloneURL:    "https://github.com/kyoh86/gogh.git",
			UpdatedAt:   updatedAt,
			Description: "A Git repository manager",
			Homepage:    "https://kyoh86.dev/gogh",
			Language:    "Go",
			Archived:    true,
			Private:     true,
			IsTemplate:  true,
		}

		if repo.Description != "A Git repository manager" {
			t.Errorf("expected description 'A Git repository manager', got %s", repo.Description)
		}
		if repo.Homepage != "https://kyoh86.dev/gogh" {
			t.Errorf("expected homepage 'https://kyoh86.dev/gogh', got %s", repo.Homepage)
		}
		if repo.Language != "Go" {
			t.Errorf("expected language 'Go', got %s", repo.Language)
		}
		if !repo.Archived {
			t.Error("expected Archived to be true")
		}
		if !repo.Private {
			t.Error("expected Private to be true")
		}
		if !repo.IsTemplate {
			t.Error("expected IsTemplate to be true")
		}
	})

	t.Run("fork repository", func(t *testing.T) {
		parent := &hosting.ParentRepository{
			Ref:      parentRef,
			CloneURL: "https://github.com/original/gogh.git",
		}

		repo := hosting.Repository{
			Ref:       ref,
			URL:       "https://github.com/kyoh86/gogh",
			CloneURL:  "https://github.com/kyoh86/gogh.git",
			UpdatedAt: updatedAt,
			Parent:    parent,
			Fork:      true,
		}

		if !repo.Fork {
			t.Error("expected Fork to be true")
		}
		if repo.Parent == nil {
			t.Fatal("expected non-nil parent")
		}
		if repo.Parent.Ref.String() != "github.com/original/gogh" {
			t.Errorf("expected parent ref 'github.com/original/gogh', got %s", repo.Parent.Ref.String())
		}
		if repo.Parent.CloneURL != "https://github.com/original/gogh.git" {
			t.Errorf("expected parent CloneURL 'https://github.com/original/gogh.git', got %s", repo.Parent.CloneURL)
		}
	})

	t.Run("zero values", func(t *testing.T) {
		var repo hosting.Repository

		// Check that zero values behave correctly
		if repo.Ref.String() != "" {
			t.Errorf("expected empty ref string, got %s", repo.Ref.String())
		}
		if repo.URL != "" {
			t.Error("expected empty URL")
		}
		if repo.CloneURL != "" {
			t.Error("expected empty CloneURL")
		}
		if !repo.UpdatedAt.IsZero() {
			t.Error("expected zero UpdatedAt")
		}
		if repo.Parent != nil {
			t.Error("expected nil parent")
		}
		if repo.Description != "" {
			t.Error("expected empty description")
		}
		if repo.Homepage != "" {
			t.Error("expected empty homepage")
		}
		if repo.Language != "" {
			t.Error("expected empty language")
		}
		if repo.Archived {
			t.Error("expected Archived to be false")
		}
		if repo.Private {
			t.Error("expected Private to be false")
		}
		if repo.IsTemplate {
			t.Error("expected IsTemplate to be false")
		}
		if repo.Fork {
			t.Error("expected Fork to be false")
		}
	})
}

func TestParentRepository(t *testing.T) {
	ref := repository.NewReference("github.com", "original", "gogh")

	t.Run("basic parent", func(t *testing.T) {
		parent := hosting.ParentRepository{
			Ref:      ref,
			CloneURL: "https://github.com/original/gogh.git",
		}

		if parent.Ref.String() != "github.com/original/gogh" {
			t.Errorf("expected ref 'github.com/original/gogh', got %s", parent.Ref.String())
		}
		if parent.CloneURL != "https://github.com/original/gogh.git" {
			t.Errorf("expected CloneURL 'https://github.com/original/gogh.git', got %s", parent.CloneURL)
		}
	})

	t.Run("ssh clone URL", func(t *testing.T) {
		parent := hosting.ParentRepository{
			Ref:      ref,
			CloneURL: "git@github.com:original/gogh.git",
		}

		if parent.CloneURL != "git@github.com:original/gogh.git" {
			t.Errorf("expected CloneURL 'git@github.com:original/gogh.git', got %s", parent.CloneURL)
		}
	})

	t.Run("zero values", func(t *testing.T) {
		var parent hosting.ParentRepository

		if parent.Ref.String() != "" {
			t.Errorf("expected empty ref string, got %s", parent.Ref.String())
		}
		if parent.CloneURL != "" {
			t.Error("expected empty CloneURL")
		}
	})
}
