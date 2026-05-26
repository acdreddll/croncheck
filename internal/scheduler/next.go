package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/croncheck/internal/audit"
)

// nextAfter returns the next scheduled time strictly after `from`.
func nextAfter(entry audit.Entry, from time.Time) (time.Time, error) {
	fields := strings.Fields(entry.Expression)
	if len(fields) != 5 {
		return time.Time{}, fmt.Errorf("invalid expression: %q", entry.Expression)
	}

	minuteSet, err := expandField(fields[0], 0, 59)
	if err != nil {
		return time.Time{}, fmt.Errorf("minute: %w", err)
	}
	hourSet, err := expandField(fields[1], 0, 23)
	if err != nil {
		return time.Time{}, fmt.Errorf("hour: %w", err)
	}
	daySet, err := expandField(fields[2], 1, 31)
	if err != nil {
		return time.Time{}, fmt.Errorf("day: %w", err)
	}
	monthSet, err := expandField(fields[3], 1, 12)
	if err != nil {
		return time.Time{}, fmt.Errorf("month: %w", err)
	}
	weekdaySet, err := expandField(fields[4], 0, 6)
	if err != nil {
		return time.Time{}, fmt.Errorf("weekday: %w", err)
	}

	// Start searching from the next minute.
	t := from.Truncate(time.Minute).Add(time.Minute)

	for i := 0; i < 366*24*60; i++ {
		if monthSet[int(t.Month())] &&
			daySet[t.Day()] &&
			weekdaySet[int(t.Weekday())] &&
			hourSet[t.Hour()] &&
			minuteSet[t.Minute()] {
			return t, nil
		}
		t = t.Add(time.Minute)
	}

	return time.Time{}, fmt.Errorf("no next time found for expression %q", entry.Expression)
}

// expandField converts a cron field string into a boolean set indexed by value.
func expandField(field string, min, max int) (map[int]bool, error) {
	set := make(map[int]bool)

	for _, part := range strings.Split(field, ",") {
		if part == "*" {
			for v := min; v <= max; v++ {
				set[v] = true
			}
			continue
		}
		if strings.Contains(part, "/") {
			sub := strings.SplitN(part, "/", 2)
			step, err := strconv.Atoi(sub[1])
			if err != nil || step <= 0 {
				return nil, fmt.Errorf("invalid step in %q", part)
			}
			start := min
			if sub[0] != "*" {
				start, err = strconv.Atoi(sub[0])
				if err != nil {
					return nil, fmt.Errorf("invalid start in %q", part)
				}
			}
			for v := start; v <= max; v += step {
				set[v] = true
			}
			continue
		}
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			lo, err1 := strconv.Atoi(bounds[0])
			hi, err2 := strconv.Atoi(bounds[1])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid range %q", part)
			}
			for v := lo; v <= hi; v++ {
				set[v] = true
			}
			continue
		}
		v, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid value %q", part)
		}
		set[v] = true
	}

	return set, nil
}
