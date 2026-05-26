package audit

import (
	"fmt"
	"time"

	"github.com/user/croncheck/internal/scheduler"
)

const (
	// overlapLookAhead is the number of future occurrences used for conflict detection.
	overlapLookAhead = 1440 // ~24 hours at per-minute granularity
)

// DetectConflicts finds pairs of cron entries whose schedules overlap.
// It returns an Issue for each conflicting pair.
func DetectConflicts(entries []Entry) []Issue {
	var issues []Issue
	base := time.Now().Truncate(time.Minute)

	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if issue := checkPair(entries[i], entries[j], base); issue != nil {
				issues = append(issues, *issue)
			}
		}
	}
	return issues
}

func checkPair(a, b Entry, base time.Time) *Issue {
	// Fast path: identical expressions are always conflicting.
	if a.Expression == b.Expression {
		return &Issue{
			Severity: SeverityWarning,
			Type:     IssueConflict,
			Message: fmt.Sprintf(
				"duplicate expression %q: %q and %q fire at identical times",
				a.Expression, a.Label, b.Label,
			),
			Entries: []string{a.Label, b.Label},
		}
	}

	// Use window-based overlap detection for richer diagnostics.
	windows, err := scheduler.FindOverlapWindows(a.Expression, b.Expression, base, overlapLookAhead)
	if err != nil || len(windows) == 0 {
		return nil
	}

	totalOverlaps := 0
	for _, w := range windows {
		totalOverlaps += w.Count
	}

	return &Issue{
		Severity: SeverityWarning,
		Type:     IssueConflict,
		Message: fmt.Sprintf(
			"%q and %q overlap %d time(s) in the next 24h (first window: %s)",
			a.Label, b.Label, totalOverlaps, windows[0].String(),
		),
		Entries: []string{a.Label, b.Label},
	}
}

// fieldsOverlap is retained for unit-testable field-level checks.
func fieldsOverlap(a, b string) bool {
	if a == "*" || b == "*" {
		return true
	}
	return a == b
}
