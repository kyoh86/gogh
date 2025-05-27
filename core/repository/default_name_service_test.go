package repository_test

import (
	"testing"

	"github.com/kyoh86/gogh/v4/core/gogh"
	testtarget "github.com/kyoh86/gogh/v4/core/repository"
)

func TestNewDefaultNameService(t *testing.T) {
	// Create a new service
	service := testtarget.NewDefaultNameService()

	// Check initial values
	if service.GetDefaultHost() != gogh.DefaultHost {
		t.Errorf("Expected default host to be %s, got %s", gogh.DefaultHost, service.GetDefaultHost())
	}

	if hosts := service.GetMap(); hosts == nil || len(hosts) != 0 {
		t.Errorf("Expected empty hosts map, got %v", hosts)
	}

	if service.HasChanges() {
		t.Error("New service should not have changes")
	}
}

func TestGetDefaultHost(t *testing.T) {
	// Create a new service
	service := testtarget.NewDefaultNameService()

	// Test with default initialization
	if host := service.GetDefaultHost(); host != gogh.DefaultHost {
		t.Errorf("Expected default host to be %s, got %s", gogh.DefaultHost, host)
	}

	// Set a new default host
	err := service.SetDefaultHost("custom.example.com")
	if err != nil {
		t.Fatalf("Failed to set default host: %v", err)
	}

	// Test with custom value
	if host := service.GetDefaultHost(); host != "custom.example.com" {
		t.Errorf("Expected default host to be custom.example.com, got %s", host)
	}

	// Test empty default host (should return gogh.DefaultHost)
	service = testtarget.NewDefaultNameService()
	if err := service.SetDefaultHost(""); err != testtarget.ErrEmptyHost {
		t.Errorf("Expected error for empty host, got %v", err)
	}
}

func TestGetMap(t *testing.T) {
	// Create a new service
	service := testtarget.NewDefaultNameService()

	// Test with empty map
	initialMap := service.GetMap()
	if initialMap == nil {
		t.Error("Expected non-nil map, got nil")
	}
	if len(initialMap) != 0 {
		t.Errorf("Expected empty map, got %v", initialMap)
	}

	// Add some host-owner pairs
	err := service.SetDefaultOwnerFor("github.com", "user1")
	if err != nil {
		t.Fatalf("Failed to set default owner: %v", err)
	}

	if err := service.SetDefaultOwnerFor("gitlab.com", "user2"); err != nil {
		t.Fatalf("Failed to set default owner: %v", err)
	}

	// Test with populated map
	hostMap := service.GetMap()
	if len(hostMap) != 2 {
		t.Errorf("Expected map with 2 entries, got %d", len(hostMap))
	}

	if owner, exists := hostMap["github.com"]; !exists || owner != "user1" {
		t.Errorf("Expected github.com to map to user1, got %s, exists: %v", owner, exists)
	}

	if owner, exists := hostMap["gitlab.com"]; !exists || owner != "user2" {
		t.Errorf("Expected gitlab.com to map to user2, got %s, exists: %v", owner, exists)
	}
}

func TestGetDefaultHostAndOwner(t *testing.T) {
	// Create a new service
	service := testtarget.NewDefaultNameService()

	// Test with default initialization (no owner set)
	host, owner := service.GetDefaultHostAndOwner()
	if host != gogh.DefaultHost {
		t.Errorf("Expected default host to be %s, got %s", gogh.DefaultHost, host)
	}
	if owner != "" {
		t.Errorf("Expected default owner to be empty, got %s", owner)
	}

	// Set default host and owner
	err := service.SetDefaultHost("github.com")
	if err != nil {
		t.Fatalf("Failed to set default host: %v", err)
	}

	if err := service.SetDefaultOwnerFor("github.com", "testuser"); err != nil {
		t.Fatalf("Failed to set default owner: %v", err)
	}

	// Test with set values
	host, owner = service.GetDefaultHostAndOwner()
	if host != "github.com" {
		t.Errorf("Expected host to be github.com, got %s", host)
	}
	if owner != "testuser" {
		t.Errorf("Expected owner to be testuser, got %s", owner)
	}

	// Test with missing owner for default host
	if err := service.SetDefaultHost("gitlab.com"); err != nil {
		t.Fatalf("Failed to set default host: %v", err)
	}

	host, owner = service.GetDefaultHostAndOwner()
	if host != "gitlab.com" {
		t.Errorf("Expected host to be gitlab.com, got %s", host)
	}
	if owner != "" {
		t.Errorf("Expected owner to be empty, got %s", owner)
	}
}

func TestGetDefaultOwnerFor(t *testing.T) {
	// Create a new service
	service := testtarget.NewDefaultNameService()

	// Test with no owner set
	owner, err := service.GetDefaultOwnerFor("github.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if owner != "" {
		t.Errorf("Expected owner to be empty, got %s", owner)
	}

	// Set owner for host
	if err := service.SetDefaultOwnerFor("github.com", "testuser"); err != nil {
		t.Fatalf("Failed to set default owner: %v", err)
	}

	// Test with owner set
	owner, err = service.GetDefaultOwnerFor("github.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if owner != "testuser" {
		t.Errorf("Expected owner to be testuser, got %s", owner)
	}

	// Test with different host
	owner, err = service.GetDefaultOwnerFor("gitlab.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if owner != "" {
		t.Errorf("Expected owner to be empty, got %s", owner)
	}
}

func TestSetDefaultHost(t *testing.T) {
	// Create a new service
	service := testtarget.NewDefaultNameService()

	// Test with valid host
	err := service.SetDefaultHost("github.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if host := service.GetDefaultHost(); host != "github.com" {
		t.Errorf("Expected host to be github.com, got %s", host)
	}
	if !service.HasChanges() {
		t.Error("Expected HasChanges to return true after setting host")
	}

	// Reset changes flag
	service.MarkSaved()

	// Test with invalid host (assuming ValidateHost returns an error for invalid hosts)
	if err := service.SetDefaultHost("invalid/host"); err == nil {
		t.Error("Expected error for invalid host, got nil")
	}
	if host := service.GetDefaultHost(); host != "github.com" {
		t.Errorf("Expected host to remain github.com, got %s", host)
	}
	if service.HasChanges() {
		t.Error("Expected HasChanges to remain false after failed set")
	}
}

func TestSetDefaultOwnerFor(t *testing.T) {
	// Create a new service
	service := testtarget.NewDefaultNameService()

	// Test with valid host and owner
	err := service.SetDefaultOwnerFor("github.com", "testuser")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	owner, err := service.GetDefaultOwnerFor("github.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if owner != "testuser" {
		t.Errorf("Expected owner to be testuser, got %s", owner)
	}
	if !service.HasChanges() {
		t.Error("Expected HasChanges to return true after setting owner")
	}

	// Reset changes flag
	service.MarkSaved()

	// Test with invalid host
	if err := service.SetDefaultOwnerFor("invalid/host", "testuser"); err == nil {
		t.Error("Expected error for invalid host, got nil")
	}
	if service.HasChanges() {
		t.Error("Expected HasChanges to remain false after failed set")
	}

	// Test with invalid owner
	if err := service.SetDefaultOwnerFor("github.com", "invalid/user"); err == nil {
		t.Error("Expected error for invalid owner, got nil")
	}
	if service.HasChanges() {
		t.Error("Expected HasChanges to remain false after failed set")
	}

	// Verify original owner is unchanged
	owner, err = service.GetDefaultOwnerFor("github.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if owner != "testuser" {
		t.Errorf("Expected owner to remain testuser, got %s", owner)
	}
}

func TestHasChangesAndMarkSaved(t *testing.T) {
	// Create a new service
	service := testtarget.NewDefaultNameService()

	// Initially should have no changes
	if service.HasChanges() {
		t.Error("New service should not have changes")
	}

	// Make a change
	err := service.SetDefaultHost("github.com")
	if err != nil {
		t.Fatalf("Failed to set default host: %v", err)
	}

	// Should now have changes
	if !service.HasChanges() {
		t.Error("Expected HasChanges to return true after change")
	}

	// Mark as saved
	service.MarkSaved()

	// Should no longer have changes
	if service.HasChanges() {
		t.Error("Expected HasChanges to return false after MarkSaved")
	}

	// Another change
	if err := service.SetDefaultOwnerFor("github.com", "testuser"); err != nil {
		t.Fatalf("Failed to set default owner: %v", err)
	}

	// Should have changes again
	if !service.HasChanges() {
		t.Error("Expected HasChanges to return true after another change")
	}
}
