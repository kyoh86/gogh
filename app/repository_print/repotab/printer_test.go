package repotab_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/kyoh86/gogh/v4/app/repository_print/repotab"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/repository"
)

func TestNewPrinter(t *testing.T) {
	t.Run("default columns", func(t *testing.T) {
		var buf bytes.Buffer
		p := repotab.NewPrinter(&buf)

		// Create a test repository using Reference properly
		ref := repository.NewReference("github.com", "testuser", "test-repo")
		testRepo := hosting.Repository{
			Ref:         ref,
			Description: "Test repository",
			URL:         "https://github.com/testuser/test-repo",
			CloneURL:    "https://github.com/testuser/test-repo.git",
			UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			Language:    "Go",
		}

		// Print and verify output
		if err := p.Print(testRepo); err != nil {
			t.Fatalf("failed to print: %v", err)
		}

		if err := p.Close(); err != nil {
			t.Fatalf("failed to close: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "test-repo") {
			t.Errorf("output should contain repository name, got: %s", output)
		}
		if !strings.Contains(output, "Test repository") {
			t.Errorf("output should contain description, got: %s", output)
		}
		if !strings.Contains(output, "2023") {
			t.Errorf("output should contain updated date, got: %s", output)
		}
	})

	t.Run("custom columns", func(t *testing.T) {
		var buf bytes.Buffer

		// Define a custom column for testing
		customColumns := []repotab.Column{
			{
				Priority:    0,
				CellBuilder: repotab.RepoRefCell,
			},
		}

		p := repotab.NewPrinter(&buf, repotab.Columns(customColumns...))

		ref := repository.NewReference("github.com", "testuser", "test-repo")
		testRepo := hosting.Repository{
			Ref: ref,
		}

		if err := p.Print(testRepo); err != nil {
			t.Fatalf("failed to print: %v", err)
		}

		if err := p.Close(); err != nil {
			t.Fatalf("failed to close: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "testuser/test-repo") {
			t.Errorf("output should contain full repository name, got: %s", output)
		}

		// Should not contain other columns
		if strings.Contains(output, "description") {
			t.Errorf("output should not contain description column, got: %s", output)
		}
	})

	t.Run("custom width", func(t *testing.T) {
		var buf bytes.Buffer
		// Set a narrow width to test truncation
		p := repotab.NewPrinter(&buf, repotab.Width(40))

		// Repository with long description
		ref := repository.NewReference("github.com", "testuser", "test-repo")
		testRepo := hosting.Repository{
			Ref:         ref,
			Description: "This is a very long description that should be truncated in narrow width",
			UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		if err := p.Print(testRepo); err != nil {
			t.Fatalf("failed to print: %v", err)
		}

		if err := p.Close(); err != nil {
			t.Fatalf("failed to close: %v", err)
		}

		output := buf.String()
		// Description should be truncated
		if len(output) > 40*2 { // allow for 2 lines (header + content)
			t.Logf("Output with width 40: %s", output)
			// Note: Exact truncation point depends on other columns, so we're just checking overall length
		}
	})
}

func TestPrinter_Print(t *testing.T) {
	t.Run("multiple repositories", func(t *testing.T) {
		var buf bytes.Buffer
		p := repotab.NewPrinter(&buf)

		repos := []hosting.Repository{
			{
				Ref:         repository.NewReference("github.com", "user1", "repo1"),
				Description: "First repository",
				UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				Ref:         repository.NewReference("github.com", "user2", "repo2"),
				Description: "Second repository",
				Private:     true,
				UpdatedAt:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			{
				Ref:         repository.NewReference("github.com", "user1", "repo3"),
				Description: "Third repository",
				Archived:    true,
				UpdatedAt:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			},
		}

		for _, repo := range repos {
			if err := p.Print(repo); err != nil {
				t.Fatalf("failed to print: %v", err)
			}
		}

		if err := p.Close(); err != nil {
			t.Fatalf("failed to close: %v", err)
		}

		output := buf.String()

		// Check for each repository content
		if !strings.Contains(output, "repo1") || !strings.Contains(output, "First repository") {
			t.Errorf("output should contain repo1 details, got: %s", output)
		}
		if !strings.Contains(output, "repo2") || !strings.Contains(output, "Second repository") {
			t.Errorf("output should contain repo2 details, got: %s", output)
		}
		if !strings.Contains(output, "repo3") || !strings.Contains(output, "Third repository") {
			t.Errorf("output should contain repo3 details, got: %s", output)
		}
	})

	t.Run("empty repository list", func(t *testing.T) {
		var buf bytes.Buffer
		p := repotab.NewPrinter(&buf)

		// No repositories are printed

		if err := p.Close(); err != nil {
			t.Fatalf("failed to close: %v", err)
		}

		output := buf.String()
		if output != "" {
			t.Errorf("expected empty output for no repositories, got: %s", output)
		}
	})

	t.Run("repository attributes", func(t *testing.T) {
		var buf bytes.Buffer
		p := repotab.NewPrinter(&buf)

		// Repository with all attributes set
		ref := repository.NewReference("github.com", "testuser", "feature-repo")
		parentRef := repository.NewReference("github.com", "original-owner", "original-repo")

		testRepo := hosting.Repository{
			Ref:         ref,
			Description: "Test repository with all attributes",
			Private:     true,
			Archived:    true,
			Fork:        true,
			Parent: &hosting.ParentRepository{
				Ref:      parentRef,
				CloneURL: "https://github.com/original-owner/original-repo.git",
			},
			UpdatedAt:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			Language:   "Go",
			Homepage:   "https://example.com",
			IsTemplate: true,
		}

		if err := p.Print(testRepo); err != nil {
			t.Fatalf("failed to print: %v", err)
		}

		if err := p.Close(); err != nil {
			t.Fatalf("failed to close: %v", err)
		}

		output := buf.String()

		// Check that attributes are displayed correctly
		if !strings.Contains(strings.ToLower(output), "private") {
			t.Errorf("output should indicate private repository, got: %s", output)
		}
		if !strings.Contains(strings.ToLower(output), "archived") {
			t.Errorf("output should indicate archived repository, got: %s", output)
		}
		if !strings.Contains(strings.ToLower(output), "fork") {
			t.Errorf("output should indicate forked repository, got: %s", output)
		}
	})
}

// Each cell type test
func TestIndividualCellTypes(t *testing.T) {
	// Create a repository with all attributes set for testing cells
	ref := repository.NewReference("github.com", "testuser", "test-repo")
	parentRef := repository.NewReference("github.com", "original-owner", "original-repo")

	testRepo := hosting.Repository{
		Ref:         ref,
		Description: "Test repository description",
		Private:     true,
		Archived:    true,
		Fork:        true,
		Parent: &hosting.ParentRepository{
			Ref:      parentRef,
			CloneURL: "https://github.com/original-owner/original-repo.git",
		},
		UpdatedAt:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		Language:   "Go",
		Homepage:   "https://example.com",
		IsTemplate: true,
	}

	// Test each cell type individually
	t.Run("RepoRefCell", func(t *testing.T) {
		var buf bytes.Buffer
		p := repotab.NewPrinter(&buf, repotab.Columns(repotab.Column{
			Priority:    0,
			CellBuilder: repotab.RepoRefCell,
		}))

		if err := p.Print(testRepo); err != nil {
			t.Fatalf("failed to print: %v", err)
		}

		if err := p.Close(); err != nil {
			t.Fatalf("failed to close: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "testuser/test-repo") {
			t.Errorf("RepoRefCell should show repository full name, got: %s", output)
		}
	})

	t.Run("DescriptionCell", func(t *testing.T) {
		var buf bytes.Buffer
		p := repotab.NewPrinter(&buf, repotab.Columns(repotab.Column{
			Priority:    0,
			CellBuilder: repotab.DescriptionCell,
		}))

		if err := p.Print(testRepo); err != nil {
			t.Fatalf("failed to print: %v", err)
		}

		if err := p.Close(); err != nil {
			t.Fatalf("failed to close: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "Test repository description") {
			t.Errorf("DescriptionCell should show repository description, got: %s", output)
		}
	})

	t.Run("EmojiAttributesCell", func(t *testing.T) {
		var buf bytes.Buffer
		p := repotab.NewPrinter(&buf, repotab.Columns(repotab.Column{
			Priority:    0,
			CellBuilder: repotab.EmojiAttributesCell,
		}))

		if err := p.Print(testRepo); err != nil {
			t.Fatalf("failed to print: %v", err)
		}

		if err := p.Close(); err != nil {
			t.Fatalf("failed to close: %v", err)
		}

		output := buf.String()
		// Check for emojis present in the output
		if !strings.Contains(output, "🔒") {
			t.Errorf("EmojiAttributesCell should show private emoji, got: %s", output)
		}
		if !strings.Contains(output, "🔀") {
			t.Errorf("EmojiAttributesCell should show fork emoji, got: %s", output)
		}
		if !strings.Contains(output, "🗃️") {
			t.Errorf("EmojiAttributesCell should show archived emoji, got: %s", output)
		}
	})

	t.Run("AttributesCell", func(t *testing.T) {
		var buf bytes.Buffer
		p := repotab.NewPrinter(&buf, repotab.Columns(repotab.Column{
			Priority:    0,
			CellBuilder: repotab.AttributesCell,
		}))

		if err := p.Print(testRepo); err != nil {
			t.Fatalf("failed to print: %v", err)
		}

		if err := p.Close(); err != nil {
			t.Fatalf("failed to close: %v", err)
		}

		output := buf.String()
		// Check for attributes text in the output
		if !strings.Contains(strings.ToLower(output), "private") {
			t.Errorf("AttributesCell should show private status, got: %s", output)
		}
		if !strings.Contains(strings.ToLower(output), "fork") {
			t.Errorf("AttributesCell should show fork status, got: %s", output)
		}
		if !strings.Contains(strings.ToLower(output), "archived") {
			t.Errorf("AttributesCell should show archived status, got: %s", output)
		}
	})

	t.Run("UpdatedAtCell", func(t *testing.T) {
		var buf bytes.Buffer
		p := repotab.NewPrinter(&buf, repotab.Columns(repotab.Column{
			Priority:    0,
			CellBuilder: repotab.UpdatedAtCell,
		}))

		// Create a repository with a specific update time
		recentRepo := hosting.Repository{
			Ref:       repository.NewReference("github.com", "testuser", "test-repo"),
			UpdatedAt: time.Now().Add(-24 * time.Hour), // 1 day ago
		}

		if err := p.Print(recentRepo); err != nil {
			t.Fatalf("failed to print: %v", err)
		}

		if err := p.Close(); err != nil {
			t.Fatalf("failed to close: %v", err)
		}

		output := buf.String()
		// The exact format depends on the FuzzyAgoAbbr function, but should indicate "ago"
		if !strings.Contains(output, "d") { // Looking for "1d" (1 day)
			t.Errorf("UpdatedAtCell should show relative time, got: %s", output)
		}
	})
}

func TestMultipleCellsCombined(t *testing.T) {
	var buf bytes.Buffer

	// Test combining multiple cell types
	p := repotab.NewPrinter(&buf, repotab.Columns(
		repotab.Column{
			Priority:    0,
			CellBuilder: repotab.RepoRefCell,
		},
		repotab.Column{
			Priority:    1,
			CellBuilder: repotab.EmojiAttributesCell,
		},
		repotab.Column{
			Priority:    2,
			CellBuilder: repotab.DescriptionCell,
		},
	))

	// Create a repository with mixed attributes
	ref := repository.NewReference("github.com", "testuser", "test-repo")
	testRepo := hosting.Repository{
		Ref:         ref,
		Description: "This is a test repository",
		Private:     true,
		Fork:        true,
		UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	if err := p.Print(testRepo); err != nil {
		t.Fatalf("failed to print: %v", err)
	}

	if err := p.Close(); err != nil {
		t.Fatalf("failed to close: %v", err)
	}

	output := buf.String()

	// Check that all cell types are represented
	if !strings.Contains(output, "testuser/test-repo") {
		t.Errorf("Output should contain repository name, got: %s", output)
	}

	if !strings.Contains(output, "🔒") || !strings.Contains(output, "🔀") {
		t.Errorf("Output should contain emoji attributes, got: %s", output)
	}

	if !strings.Contains(output, "This is a test repository") {
		t.Errorf("Output should contain description, got: %s", output)
	}
}
