package hosting_test

import (
	"encoding/json"
	"testing"
	"time"

	testtarget "github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/repository"
)

// TestRepositoryFormatFunc tests that the RepositoryFormatFunc correctly implements the RepositoryFormat interface
func TestRepositoryFormatFunc(t *testing.T) {
	// Create a custom formatter function
	formatFunc := func(r testtarget.Repository) (string, error) {
		return "custom:" + r.Ref.String(), nil
	}

	// Convert it to a RepositoryFormatFunc
	formatter := testtarget.RepositoryFormatFunc(formatFunc)

	// Create a test repository
	ref := repository.NewReference("github.com", "kyoh86", "gogh")
	repo := testtarget.Repository{
		Ref: ref,
		URL: "https://github.com/kyoh86/gogh",
	}

	// Format the repository using the formatter
	result, err := formatter.Format(repo)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify the result
	expected := "custom:github.com/kyoh86/gogh"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestRepositoryFormatRef tests the RepositoryFormatRef formatter
func TestRepositoryFormatRef(t *testing.T) {
	// Create a test repository
	ref := repository.NewReference("github.com", "kyoh86", "gogh")
	repo := testtarget.Repository{
		Ref: ref,
	}

	// Format the repository using RepositoryFormatRef
	result, err := testtarget.RepositoryFormatRef.Format(repo)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify the result
	expected := "github.com/kyoh86/gogh"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test with a different host
	refGitlab := repository.NewReference("gitlab.com", "user", "project")
	repoGitlab := testtarget.Repository{
		Ref: refGitlab,
	}

	// Format the repository using RepositoryFormatRef
	result, err = testtarget.RepositoryFormatRef.Format(repoGitlab)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify the result
	expected = "gitlab.com/user/project"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestRepositoryFormatURL tests the RepositoryFormatURL formatter
func TestRepositoryFormatURL(t *testing.T) {
	// Create a test repository with URL
	ref := repository.NewReference("github.com", "kyoh86", "gogh")
	repo := testtarget.Repository{
		Ref: ref,
		URL: "https://github.com/kyoh86/gogh",
	}

	// Format the repository using RepositoryFormatURL
	result, err := testtarget.RepositoryFormatURL.Format(repo)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify the result
	expected := "https://github.com/kyoh86/gogh"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test with empty URL
	repoNoURL := testtarget.Repository{
		Ref: ref,
		URL: "",
	}

	// Format the repository using RepositoryFormatURL
	result, err = testtarget.RepositoryFormatURL.Format(repoNoURL)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify the result is empty
	if result != "" {
		t.Errorf("Expected empty string, got %q", result)
	}
}

// TestRepositoryFormatJSON_Basic tests basic JSON formatting
func TestRepositoryFormatJSON_Basic(t *testing.T) {
	// Fixed time for consistent testing
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	// Create a test repository
	ref := repository.NewReference("github.com", "kyoh86", "gogh")
	repo := testtarget.Repository{
		Ref:       ref,
		URL:       "https://github.com/kyoh86/gogh",
		CloneURL:  "https://github.com/kyoh86/gogh.git",
		UpdatedAt: fixedTime,
	}

	// Format the repository using RepositoryFormatJSON
	result, err := testtarget.RepositoryFormatJSON.Format(repo)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Parse the JSON
	var data map[string]any
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify the basic fields
	rawRef, ok := data["ref"].(map[string]any)
	if !ok {
		t.Fatal("Missing 'ref' field in JSON")
	}

	if rawRef["host"] != "github.com" || rawRef["owner"] != "kyoh86" || rawRef["name"] != "gogh" {
		t.Errorf("Incorrect ref values: %v", ref)
	}

	if data["url"] != "https://github.com/kyoh86/gogh" {
		t.Errorf("Incorrect URL: %v", data["url"])
	}

	if data["cloneUrl"] != "https://github.com/kyoh86/gogh.git" {
		t.Errorf("Incorrect CloneURL: %v", data["cloneUrl"])
	}
}

// TestRepositoryFormatJSON_WithParent tests JSON formatting with parent repository
func TestRepositoryFormatJSON_WithParent(t *testing.T) {
	// Fixed time for consistent testing
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	// Create a test repository with parent
	ref := repository.NewReference("github.com", "kyoh86", "gogh")
	parentRef := repository.NewReference("github.com", "original", "gogh")

	parent := testtarget.ParentRepository{
		Ref:      parentRef,
		CloneURL: "https://github.com/original/gogh.git",
	}

	repo := testtarget.Repository{
		Ref:       ref,
		URL:       "https://github.com/kyoh86/gogh",
		UpdatedAt: fixedTime,
		Parent:    &parent,
		Fork:      true,
	}

	// Format the repository using RepositoryFormatJSON
	result, err := testtarget.RepositoryFormatJSON.Format(repo)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Parse the JSON
	var data map[string]any
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify the parent fields
	rawParent, ok := data["parent"].(map[string]any)
	if !ok {
		t.Fatal("Missing 'parent' field in JSON")
	}

	rawParentRef, ok := rawParent["ref"].(map[string]any)
	if !ok {
		t.Fatal("Missing 'ref' field in parent")
	}

	if rawParentRef["host"] != "github.com" || rawParentRef["owner"] != "original" || rawParentRef["name"] != "gogh" {
		t.Errorf("Incorrect parent ref values: %v", rawParentRef)
	}

	if rawParent["cloneUrl"] != "https://github.com/original/gogh.git" {
		t.Errorf("Incorrect parent CloneURL: %v", rawParent["cloneUrl"])
	}

	if data["fork"] != true {
		t.Errorf("Expected 'fork' to be true, got %v", data["fork"])
	}
}

// TestRepositoryFormatJSON_WithMetadata tests JSON formatting with all metadata fields
func TestRepositoryFormatJSON_WithMetadata(t *testing.T) {
	// Fixed time for consistent testing
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	// Create a test repository with metadata
	ref := repository.NewReference("github.com", "kyoh86", "gogh")
	repo := testtarget.Repository{
		Ref:         ref,
		URL:         "https://github.com/kyoh86/gogh",
		UpdatedAt:   fixedTime,
		Description: "A Git repository manager",
		Homepage:    "https://kyoh86.dev/gogh",
		Language:    "Go",
		Archived:    true,
		Private:     true,
		IsTemplate:  true,
	}

	// Format the repository using RepositoryFormatJSON
	result, err := testtarget.RepositoryFormatJSON.Format(repo)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Parse the JSON
	var data map[string]any
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify the metadata fields
	if data["description"] != "A Git repository manager" {
		t.Errorf("Incorrect description: %v", data["description"])
	}

	if data["homepage"] != "https://kyoh86.dev/gogh" {
		t.Errorf("Incorrect homepage: %v", data["homepage"])
	}

	if data["language"] != "Go" {
		t.Errorf("Incorrect language: %v", data["language"])
	}

	if data["archived"] != true {
		t.Errorf("Expected 'archived' to be true, got %v", data["archived"])
	}

	if data["private"] != true {
		t.Errorf("Expected 'private' to be true, got %v", data["private"])
	}

	if data["isTemplate"] != true {
		t.Errorf("Expected 'isTemplate' to be true, got %v", data["isTemplate"])
	}
}
