package scheduler

import (
	"fmt"
	"time"

	"github.com/croncheck/internal/audit"
)

// Window represents a time range with a label.
type Window struct {
	Start time.Time
	End   time.Time
	Label string
}

// Duration returns the duration of the window.
func (w Window) Duration() time.Duration {
	return w.End.Sub(w.Start)
}

// String returns a human-readable representation of the window.
func (w Window) String() string {
	return fmt.Sprintf("%s [%s – %s] (%s)",
		w.Label,
		w.Start.Format("15:04"),
		w.End.Format("15:04"),
		w.Duration().Round(time.Minute),
	)
}

// ActiveWindows returns the time windows within [from, to) during which
// the given cron entry is scheduled to run at least once.
// The granularity is per-minute.
func ActiveWindows(entry audit.Entry, from, to time.Time) ([]Window, error) {
	times, err := NextN(entry, 1440, from)
	if err != nil {
		return nil, fmt.Errorf("expanding schedule for %q: %w", entry.Label, err)
	}

	var windows []Window
	var current *Window

	for _, t := range times {
		if t.Before(from) || !t.Before(to) {
			continue
		}
		if current == nil {
			current = &Window{Start: t, End: t.Add(time.Minute), Label: entry.Label}
		} else if t.Sub(current.End) <= time.Minute {
			current.End = t.Add(time.Minute)
		} else {
			windows = append(windows, *current)
			current = &Window{Start: t, End: t.Add(time.Minute), Label: entry.Label}
		}
	}
	if current != nil {
		windows = append(windows, *current)
	}
	return windows, nil
}
