package repotab_test

import (
	"testing"
	"time"

	testtarget "github.com/kyoh86/gogh/v4/app/repoprint/repotab"
)

func TestFormatTimeAgo(t *testing.T) {
	now := time.Date(2023, 5, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		at       time.Time
		expected string
	}{
		{
			name:     "just now",
			at:       now,
			expected: "now",
		},
		{
			name:     "10 minutes ago",
			at:       now.Add(-10 * time.Minute),
			expected: "10m",
		},
		{
			name:     "2 hours ago",
			at:       now.Add(-2 * time.Hour),
			expected: "2h",
		},
		{
			name:     "5 days ago",
			at:       now.Add(-5 * 24 * time.Hour),
			expected: "5d",
		},
		{
			name:     "45 days ago",
			at:       now.Add(-45 * 24 * time.Hour),
			expected: "2023-03-17",
		},
		{
			name:     "10 minutes in future",
			at:       now.Add(10 * time.Minute),
			expected: "now",
		},
		{
			name:     "exactly 1 hour ago",
			at:       now.Add(-1 * time.Hour),
			expected: "1h",
		},
		{
			name:     "exactly 1 day ago",
			at:       now.Add(-24 * time.Hour),
			expected: "1d",
		},
		{
			name:     "abount 30 days ago",
			at:       now.Add(-30*24*time.Hour + 1),
			expected: "30d",
		},
		{
			name:     "exactly 30 days ago",
			at:       now.Add(-30 * 24 * time.Hour),
			expected: "2023-04-01",
		},
		{
			name:     "zero time",
			at:       time.Time{},
			expected: "0001-01-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testtarget.FuzzyAgoAbbr(now, tt.at)
			if result != tt.expected {
				t.Errorf("FuzzyAgoAbbr() = %v, want %v", result, tt.expected)
			}
		})
	}
}
