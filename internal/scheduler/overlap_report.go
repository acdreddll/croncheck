package scheduler

import (
	"fmt"
	"time"
)

// OverlapWindow describes a time window where two cron expressions coincide.
type OverlapWindow struct {
	Start time.Time
	End   time.Time
	Count int
}

// String returns a human-readable representation of the overlap window.
func (w OverlapWindow) String() string {
	return fmt.Sprintf("%s to %s (%d occurrence(s))",
		w.Start.Format(time.RFC3339),
		w.End.Format(time.RFC3339),
		w.Count,
	)
}

// FindOverlapWindows returns contiguous time windows where exprA and exprB
// both fire within the next n occurrences starting from base.
// Windows are merged if overlapping minutes are adjacent.
func FindOverlapWindows(exprA, exprB string, base time.Time, n int) ([]OverlapWindow, error) {
	timesA, err := NextN(exprA, base, n)
	if err != nil {
		return nil, fmt.Errorf("exprA: %w", err)
	}
	timesB, err := NextN(exprB, base, n)
	if err != nil {
		return nil, fmt.Errorf("exprB: %w", err)
	}

	setA := toMinuteSet(timesA)
	setB := toMinuteSet(timesB)

	var shared []time.Time
	for t := range setA {
		if setB[t] {
			shared = append(shared, t)
		}
	}
	if len(shared) == 0 {
		return nil, nil
	}

	// Sort shared times
	sortTimes(shared)
	return mergeIntoWindows(shared), nil
}

func toMinuteSet(times []time.Time) map[time.Time]bool {
	m := make(map[time.Time]bool, len(times))
	for _, t := range times {
		// Truncate to minute precision
		m[t.Truncate(time.Minute)] = true
	}
	return m
}

func sortTimes(ts []time.Time) {
	for i := 1; i < len(ts); i++ {
		for j := i; j > 0 && ts[j].Before(ts[j-1]); j-- {
			ts[j], ts[j-1] = ts[j-1], ts[j]
		}
	}
}

func mergeIntoWindows(sorted []time.Time) []OverlapWindow {
	if len(sorted) == 0 {
		return nil
	}
	var windows []OverlapWindow
	cur := OverlapWindow{Start: sorted[0], End: sorted[0], Count: 1}
	for _, t := range sorted[1:] {
		if t.Sub(cur.End) <= time.Minute {
			cur.End = t
			cur.Count++
		} else {
			windows = append(windows, cur)
			cur = OverlapWindow{Start: t, End: t, Count: 1}
		}
	}
	windows = append(windows, cur)
	return windows
}
