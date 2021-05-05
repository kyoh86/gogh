package view

import (
	"fmt"
	"time"
)

func FuzzyAgoAbbr(now time.Time, at time.Time) string {
	ago := now.Sub(at)
	if ago < time.Hour {
		return fmt.Sprintf("%d%s", int(ago.Minutes()), "m")
	}
	if ago < 24*time.Hour {
		return fmt.Sprintf("%d%s", int(ago.Hours()), "h")
	}
	if ago < 30*24*time.Hour {
		return fmt.Sprintf("%d%s", int(ago.Hours())/24, "d")
	}
	return at.Format("2006-01-02")
}
