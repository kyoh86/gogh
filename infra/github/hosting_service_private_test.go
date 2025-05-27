package github

import (
	"testing"
	"time"

	"github.com/google/go-github/v69/github"
	"github.com/kyoh86/gogh/v4/core/repository"
)

func TestConvertRepository(t *testing.T) {
	// Create a sample GitHub repository
	now := time.Now()
	login := "owner"
	name := "repo"
	description := "Test repository"
	homepage := "https://example.com"
	language := "Go"
	htmlURL := "https://github.com/owner/repo"
	cloneURL := "https://github.com/owner/repo.git"

	// Create owner
	owner := &github.User{
		Login: &login,
	}

	// Create the repository
	ghRepo := &github.Repository{
		Name:        &name,
		Description: &description,
		Homepage:    &homepage,
		Language:    &language,
		HTMLURL:     &htmlURL,
		CloneURL:    &cloneURL,
		Owner:       owner,
		Archived:    github.Ptr(true),
		Private:     github.Ptr(true),
		IsTemplate:  github.Ptr(true),
		Fork:        github.Ptr(true),
		UpdatedAt:   &github.Timestamp{Time: now},
	}

	// Create parent repository
	parentLogin := "parentOwner"
	parentName := "parentRepo"
	parentHTMLURL := "https://github.com/parentOwner/parentRepo"
	parentCloneURL := "https://github.com/parentOwner/parentRepo"

	parentOwner := &github.User{
		Login: &parentLogin,
	}

	parent := &github.Repository{
		Name:     &parentName,
		Owner:    parentOwner,
		HTMLURL:  &parentHTMLURL,
		CloneURL: &parentCloneURL,
	}

	ghRepo.Parent = parent

	// Convert using the function being tested
	ref := repository.NewReference("github.com", "owner", "repo")
	repo, err := convertRepository(ref, ghRepo)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Validate the conversion
	if repo.Ref != ref {
		t.Errorf("Expected ref %v, got %v", ref, repo.Ref)
	}

	if repo.Description != description {
		t.Errorf("Expected description %q, got %q", description, repo.Description)
	}

	if repo.Homepage != homepage {
		t.Errorf("Expected homepage %q, got %q", homepage, repo.Homepage)
	}

	if repo.Language != language {
		t.Errorf("Expected language %q, got %q", language, repo.Language)
	}

	if repo.URL != htmlURL {
		t.Errorf("Expected URL %q, got %q", htmlURL, repo.URL)
	}

	if repo.CloneURL != cloneURL {
		t.Errorf("Expected CloneURL %q, got %q", cloneURL, repo.CloneURL)
	}

	if !repo.Archived {
		t.Error("Expected Archived to be true")
	}

	if !repo.Private {
		t.Error("Expected Private to be true")
	}

	if !repo.IsTemplate {
		t.Error("Expected IsTemplate to be true")
	}

	if !repo.Fork {
		t.Error("Expected Fork to be true")
	}

	if repo.UpdatedAt != now {
		t.Errorf("Expected UpdatedAt %v, got %v", now, repo.UpdatedAt)
	}

	// Validate parent repository
	if repo.Parent == nil {
		t.Fatal("Expected Parent to be non-nil")
	}

	if repo.Parent.Ref.Host() != "github.com" {
		t.Errorf("Expected parent host github.com, got %s", repo.Parent.Ref.Host())
	}

	if repo.Parent.Ref.Owner() != parentLogin {
		t.Errorf("Expected parent owner %q, got %q", parentLogin, repo.Parent.Ref.Owner())
	}

	if repo.Parent.Ref.Name() != parentName {
		t.Errorf("Expected parent name %q, got %q", parentName, repo.Parent.Ref.Name())
	}

	if repo.Parent.CloneURL != parentCloneURL {
		t.Errorf("Expected parent CloneURL %q, got %q", parentCloneURL, repo.Parent.CloneURL)
	}
}

func TestConvertSSHToHTTPS(t *testing.T) {
	testCases := []struct {
		name          string
		sshURL        string
		expectedHTTPS string
	}{
		{
			name:          "valid SSH URL",
			sshURL:        "git@github.com:kyoh86/gogh.git",
			expectedHTTPS: "https://github.com/kyoh86/gogh",
		},
		{
			name:          "empty URL",
			sshURL:        "",
			expectedHTTPS: "",
		},
		{
			name:          "malformed SSH URL without colon",
			sshURL:        "git@github.com/kyoh86/gogh.git",
			expectedHTTPS: "git@github.com/kyoh86/gogh.git", // Returns original
		},
		{
			name:          "malformed SSH URL without @",
			sshURL:        "gitgithub.com:kyoh86/gogh.git",
			expectedHTTPS: "gitgithub.com:kyoh86/gogh.git", // Returns original
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convertSSHToHTTPS(tc.sshURL)
			if result != tc.expectedHTTPS {
				t.Errorf("Expected %q, got %q", tc.expectedHTTPS, result)
			}
		})
	}
}
