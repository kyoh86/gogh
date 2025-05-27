package repository_test

import (
	"errors"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/core/repository"
)

func TestValidateHost(t *testing.T) {
	testCases := []struct {
		name        string
		host        string
		expectError bool
		errorIs     error
	}{
		{
			name:        "valid domain",
			host:        "github.com",
			expectError: false,
		},
		{
			name:        "valid domain with subdomain",
			host:        "gist.github.com",
			expectError: false,
		},
		{
			name:        "valid IP address",
			host:        "192.168.1.1",
			expectError: false,
		},
		{
			name:        "empty host",
			host:        "",
			expectError: true,
			errorIs:     testtarget.ErrEmptyHost,
		},
		{
			name:        "host with path",
			host:        "github.com/path",
			expectError: true,
		},
		{
			name:        "host with scheme",
			host:        "https://github.com",
			expectError: true,
		},
		{
			name:        "host with user info",
			host:        "user@github.com",
			expectError: true,
		},
		{
			name:        "host with port",
			host:        "github.com:443",
			expectError: false,
		},
		{
			name:        "host with invalid characters",
			host:        "github com",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := testtarget.ValidateHost(tc.host)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for host '%s', but got nil", tc.host)
					return
				}

				if tc.errorIs != nil && !errors.Is(err, tc.errorIs) {
					t.Errorf("Expected error to be '%v', but got '%v'", tc.errorIs, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for host '%s', but got: %v", tc.host, err)
				}
			}
		})
	}
}

func TestValidateOwner(t *testing.T) {
	testCases := []struct {
		name        string
		owner       string
		expectError bool
		errorIs     error
	}{
		{
			name:        "valid lowercase",
			owner:       "user",
			expectError: false,
		},
		{
			name:        "valid with numbers",
			owner:       "user123",
			expectError: false,
		},
		{
			name:        "valid with hyphen",
			owner:       "user-name",
			expectError: false,
		},
		{
			name:        "valid uppercase",
			owner:       "UserName",
			expectError: false,
		},
		{
			name:        "valid mixed case",
			owner:       "User-Name123",
			expectError: false,
		},
		{
			name:        "empty owner",
			owner:       "",
			expectError: true,
			errorIs:     testtarget.ErrEmptyOwner,
		},
		{
			name:        "invalid with underscore",
			owner:       "user_name",
			expectError: true,
		},
		{
			name:        "invalid starts with hyphen",
			owner:       "-username",
			expectError: true,
		},
		{
			name:        "invalid ends with hyphen",
			owner:       "username-",
			expectError: true,
		},
		{
			name:        "invalid double hyphen",
			owner:       "user--name",
			expectError: true,
		},
		{
			name:        "invalid with special characters",
			owner:       "user.name",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := testtarget.ValidateOwner(tc.owner)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for owner '%s', but got nil", tc.owner)
					return
				}

				if tc.errorIs != nil && !errors.Is(err, tc.errorIs) {
					t.Errorf("Expected error to be '%v', but got '%v'", tc.errorIs, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for owner '%s', but got: %v", tc.owner, err)
				}
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	testCases := []struct {
		name        string
		repoName    string
		expectError bool
		errorIs     error
	}{
		{
			name:        "valid simple name",
			repoName:    "project",
			expectError: false,
		},
		{
			name:        "valid with numbers",
			repoName:    "project123",
			expectError: false,
		},
		{
			name:        "valid with hyphen",
			repoName:    "my-project",
			expectError: false,
		},
		{
			name:        "valid with dot",
			repoName:    "project.js",
			expectError: false,
		},
		{
			name:        "valid with underscore",
			repoName:    "my_project",
			expectError: false,
		},
		{
			name:        "empty name",
			repoName:    "",
			expectError: true,
			errorIs:     testtarget.ErrEmptyName,
		},
		{
			name:        "single dot (reserved)",
			repoName:    ".",
			expectError: true,
		},
		{
			name:        "double dot (reserved)",
			repoName:    "..",
			expectError: true,
		},
		{
			name:        "invalid with space",
			repoName:    "my project",
			expectError: true,
		},
		{
			name:        "invalid with special characters",
			repoName:    "project/name",
			expectError: true,
		},
		{
			name:        "invalid with colon",
			repoName:    "project:name",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := testtarget.ValidateName(tc.repoName)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for name '%s', but got nil", tc.repoName)
					return
				}

				if tc.errorIs != nil && !errors.Is(err, tc.errorIs) {
					t.Errorf("Expected error to be '%v', but got '%v'", tc.errorIs, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for name '%s', but got: %v", tc.repoName, err)
				}
			}
		})
	}
}
