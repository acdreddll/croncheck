package scheduler

import (
	"testing"

	"github.com/croncheck/internal/audit"
)

func heatEntry(label, expr string) audit.Entry {
	return audit.Entry{Label: label, Expression: expr}
}

func TestBuildHeatmap_EveryMinute(t *testing.T) {
	entries := []audit.Entry{heatEntry("all", "* * * * *")}
	hm, err := BuildHeatmap(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, c := range hm.Cells {
		if c.Count != 1 {
			t.Errorf("cell %d: want count 1, got %d", i, c.Count)
		}
	}
	if hm.MaxLoad != 1 {
		t.Errorf("MaxLoad: want 1, got %d", hm.MaxLoad)
	}
}

func TestBuildHeatmap_HourlyOverlap(t *testing.T) {
	entries := []audit.Entry{
		heatEntry("a", "0 * * * *"),
		heatEntry("b", "0 * * * *"),
	}
	hm, err := BuildHeatmap(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hm.MaxLoad != 2 {
		t.Errorf("MaxLoad: want 2, got %d", hm.MaxLoad)
	}
	// Minute 0 of every hour should have count 2.
	for h := 0; h < 24; h++ {
		c := hm.Cells[h*60]
		if c.Count != 2 {
			t.Errorf("hour %d minute 0: want 2, got %d", h, c.Count)
		}
	}
}

func TestBuildHeatmap_InvalidExpr(t *testing.T) {
	entries := []audit.Entry{heatEntry("bad", "99 * * * *")}
	_, err := BuildHeatmap(entries)
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}

func TestHotSpots_Threshold(t *testing.T) {
	entries := []audit.Entry{
		heatEntry("a", "0 0 * * *"), // once daily at midnight
		heatEntry("b", "0 0 * * *"),
		heatEntry("c", "0 0 * * *"),
	}
	hm, err := BuildHeatmap(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hots := hm.HotSpots(3)
	if len(hots) != 1 {
		t.Errorf("want 1 hot spot, got %d", len(hots))
	}
	if len(hots) > 0 && hots[0].Count != 3 {
		t.Errorf("hot spot count: want 3, got %d", hots[0].Count)
	}
}

func TestASCIIRow_Length(t *testing.T) {
	entries := []audit.Entry{heatEntry("x", "*/5 * * * *")}
	hm, err := BuildHeatmap(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	row := hm.ASCIIRow(0)
	if len(row) != 60 {
		t.Errorf("ASCIIRow length: want 60, got %d", len(row))
	}
}
