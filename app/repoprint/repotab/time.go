package repotab

import (
	"fmt"
	"math"
	"time"
)

func FuzzyAgoAbbr(now time.Time, at time.Time) string {
	// Handle future dates
	if at.After(now) {
		return "now"
	}

	ago := now.Sub(at)
	if ago < time.Minute {
		return "now"
	}
	if ago < time.Hour {
		return fmt.Sprintf("%dm", int(math.Round(ago.Minutes())))
	}
	if ago < 24*time.Hour {
		return fmt.Sprintf("%dh", int(math.Round(ago.Hours())))
	}
	if ago < 30*24*time.Hour {
		return fmt.Sprintf("%dd", int(math.Round(ago.Hours()/24)))
	}
	return at.Format("2006-01-02")
}
