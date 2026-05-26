package scheduler

import (
	"testing"
	"time"

	"github.com/croncheck/internal/audit"
)

func makeEntry(label, expr string) audit.Entry {
	return audit.Entry{Label: label, Raw: expr}
}

func TestActiveWindows_EveryMinute(t *testing.T) {
	e := makeEntry("every-min", "* * * * *")
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(10 * time.Minute)

	windows, err := ActiveWindows(e, from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(windows) != 1 {
		t.Fatalf("expected 1 merged window, got %d", len(windows))
	}
	if windows[0].Duration() < 9*time.Minute {
		t.Errorf("expected window duration >= 9m, got %s", windows[0].Duration())
	}
}

func TestActiveWindows_HourlyOutsideRange(t *testing.T) {
	e := makeEntry("hourly", "0 * * * *")
	// window of only 30 minutes starting at :15 — no :00 tick inside
	from := time.Date(2024, 1, 1, 0, 15, 0, 0, time.UTC)
	to := from.Add(30 * time.Minute)

	windows, err := ActiveWindows(e, from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(windows) != 0 {
		t.Errorf("expected 0 windows, got %d", len(windows))
	}
}

func TestActiveWindows_TwicePerHour(t *testing.T) {
	e := makeEntry("twice", "0,30 * * * *")
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(60 * time.Minute)

	windows, err := ActiveWindows(e, from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(windows) != 2 {
		t.Errorf("expected 2 windows, got %d", len(windows))
	}
}

func TestActiveWindows_InvalidExpr(t *testing.T) {
	e := makeEntry("bad", "99 99 99 99 99")
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(time.Hour)

	_, err := ActiveWindows(e, from, to)
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}

func TestWindow_String(t *testing.T) {
	w := Window{
		Label: "test-job",
		Start: time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		End:   time.Date(2024, 1, 1, 8, 5, 0, 0, time.UTC),
	}
	s := w.String()
	if s == "" {
		t.Error("expected non-empty string from Window.String()")
	}
}
