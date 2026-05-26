package scheduler

import (
	"time"

	"github.com/croncheck/internal/audit"
)

// NextN returns the next n scheduled times for a given cron entry
// starting from the provided base time.
func NextN(entry audit.Entry, base time.Time, n int) ([]time.Time, error) {
	if n <= 0 {
		return nil, nil
	}

	results := make([]time.Time, 0, n)
	current := base

	for len(results) < n {
		next, err := nextAfter(entry, current)
		if err != nil {
			return nil, err
		}
		results = append(results, next)
		current = next.Add(time.Minute)
	}

	return results, nil
}

// Overlaps reports whether two entries share at least one scheduled
// minute within the given time window starting from base.
func Overlaps(a, b audit.Entry, base time.Time, window time.Duration) (bool, error) {
	end := base.Add(window)
	timesA, err := collectMinutes(a, base, end)
	if err != nil {
		return false, err
	}
	timesB, err := collectMinutes(b, base, end)
	if err != nil {
		return false, err
	}

	for _, ta := range timesA {
		for _, tb := range timesB {
			if ta.Equal(tb) {
				return true, nil
			}
		}
	}
	return false, nil
}

func collectMinutes(entry audit.Entry, from, to time.Time) ([]time.Time, error) {
	var times []time.Time
	current := from
	for !current.After(to) {
		next, err := nextAfter(entry, current)
		if err != nil {
			return nil, err
		}
		if next.After(to) {
			break
		}
		times = append(times, next)
		current = next.Add(time.Minute)
	}
	return times, nil
}
