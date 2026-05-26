package scheduler_test

import (
	"testing"
	"time"

	"github.com/croncheck/internal/audit"
	"github.com/croncheck/internal/scheduler"
)

func entry(expr string) audit.Entry {
	return audit.Entry{Label: "test", Expression: expr}
}

func TestNextN_EveryMinute(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	times, err := scheduler.NextN(entry("* * * * *"), base, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) != 3 {
		t.Fatalf("expected 3 times, got %d", len(times))
	}
	expected := []time.Time{
		time.Date(2024, 1, 1, 0, 1, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 0, 2, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 0, 3, 0, 0, time.UTC),
	}
	for i, got := range times {
		if !got.Equal(expected[i]) {
			t.Errorf("time[%d]: got %v, want %v", i, got, expected[i])
		}
	}
}

func TestNextN_HourlyAt30(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	times, err := scheduler.NextN(entry("30 * * * *"), base, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) != 2 {
		t.Fatalf("expected 2 times, got %d", len(times))
	}
	if times[0].Minute() != 30 || times[1].Minute() != 30 {
		t.Errorf("expected minute 30, got %v and %v", times[0], times[1])
	}
}

func TestNextN_ZeroCount(t *testing.T) {
	base := time.Now()
	times, err := scheduler.NextN(entry("* * * * *"), base, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) != 0 {
		t.Errorf("expected empty slice, got %d items", len(times))
	}
}

func TestOverlaps_SameExpression(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	a := entry("0 * * * *")
	b := entry("0 * * * *")
	overlap, err := scheduler.Overlaps(a, b, base, 2*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !overlap {
		t.Error("expected overlap for identical expressions")
	}
}

func TestOverlaps_NoOverlap(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	a := entry("0 0 * * *") // midnight only
	b := entry("0 12 * * *") // noon only
	overlap, err := scheduler.Overlaps(a, b, base, 23*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if overlap {
		t.Error("expected no overlap for midnight vs noon")
	}
}

func TestNextN_InvalidExpression(t *testing.T) {
	base := time.Now()
	_, err := scheduler.NextN(entry("bad expression"), base, 1)
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}
