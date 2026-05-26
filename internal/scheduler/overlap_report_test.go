package scheduler

import (
	"testing"
	"time"
)

func baseTime() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}

func TestFindOverlapWindows_SameExpression(t *testing.T) {
	base := baseTime()
	windows, err := FindOverlapWindows("* * * * *", "* * * * *", base, 60)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(windows) == 0 {
		t.Fatal("expected overlap windows for identical expressions")
	}
	// All 60 minutes should merge into a single contiguous window
	if len(windows) != 1 {
		t.Errorf("expected 1 merged window, got %d", len(windows))
	}
	if windows[0].Count != 60 {
		t.Errorf("expected count=60, got %d", windows[0].Count)
	}
}

func TestFindOverlapWindows_NoOverlap(t *testing.T) {
	base := baseTime()
	// One fires at :00, other at :30
	windows, err := FindOverlapWindows("0 * * * *", "30 * * * *", base, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(windows) != 0 {
		t.Errorf("expected no overlap, got %d windows", len(windows))
	}
}

func TestFindOverlapWindows_PartialOverlap(t *testing.T) {
	base := baseTime()
	// Every 15 min vs every 30 min — overlap at :00 and :30
	windows, err := FindOverlapWindows("*/15 * * * *", "*/30 * * * *", base, 120)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(windows) == 0 {
		t.Fatal("expected some overlap windows")
	}
}

func TestFindOverlapWindows_InvalidExprA(t *testing.T) {
	_, err := FindOverlapWindows("invalid", "* * * * *", baseTime(), 10)
	if err == nil {
		t.Fatal("expected error for invalid exprA")
	}
}

func TestFindOverlapWindows_InvalidExprB(t *testing.T) {
	_, err := FindOverlapWindows("* * * * *", "bad expr", baseTime(), 10)
	if err == nil {
		t.Fatal("expected error for invalid exprB")
	}
}

func TestOverlapWindow_String(t *testing.T) {
	w := OverlapWindow{
		Start: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		End:   time.Date(2024, 1, 1, 12, 5, 0, 0, time.UTC),
		Count: 6,
	}
	s := w.String()
	if s == "" {
		t.Error("expected non-empty string from OverlapWindow.String()")
	}
}
