package scheduler

import (
	"fmt"
	"sort"
	"time"
)

// ScheduleSummary holds a human-readable breakdown of when a cron expression fires.
type ScheduleSummary struct {
	Expression  string
	NextTimes   []time.Time
	Frequency   string
	DailyCount  int
	WeeklyCount int
}

// Summarize returns a ScheduleSummary for the given cron expression,
// computing the next n occurrences starting from base.
func Summarize(expr string, base time.Time, n int) (*ScheduleSummary, error) {
	if n <= 0 {
		n = 5
	}

	next, err := NextN(expr, base, n)
	if err != nil {
		return nil, fmt.Errorf("summarize: %w", err)
	}

	daily := countInWindow(expr, base, 24*time.Hour)
	weekly := countInWindow(expr, base, 7*24*time.Hour)

	freq := describeFrequency(daily)

	return &ScheduleSummary{
		Expression:  expr,
		NextTimes:   next,
		Frequency:   freq,
		DailyCount:  daily,
		WeeklyCount: weekly,
	}, nil
}

// countInWindow counts how many times expr fires within duration d starting from base.
func countInWindow(expr string, base time.Time, d time.Duration) int {
	end := base.Add(d)
	times, err := NextN(expr, base, 1500)
	if err != nil {
		return 0
	}
	count := 0
	for _, t := range times {
		if t.Before(end) {
			count++
		}
	}
	return count
}

// describeFrequency returns a human-readable frequency label based on daily fire count.
func describeFrequency(dailyCount int) string {
	switch {
	case dailyCount >= 1440:
		return "every minute"
	case dailyCount >= 60:
		return "multiple times per hour"
	case dailyCount >= 24:
		return "multiple times per day"
	case dailyCount >= 2:
		return fmt.Sprintf("%d times per day", dailyCount)
	case dailyCount == 1:
		return "once per day"
	default:
		return "less than once per day"
	}
}

// SortSummaries sorts a slice of ScheduleSummary by DailyCount descending.
func SortSummaries(summaries []*ScheduleSummary) {
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].DailyCount > summaries[j].DailyCount
	})
}
