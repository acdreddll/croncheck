package scheduler

import (
	"testing"
	"time"
)

var summaryBase = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

func TestSummarize_EveryMinute(t *testing.T) {
	s, err := Summarize("* * * * *", summaryBase, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.NextTimes) != 5 {
		t.Errorf("expected 5 next times, got %d", len(s.NextTimes))
	}
	if s.DailyCount != 1440 {
		t.Errorf("expected 1440 daily fires, got %d", s.DailyCount)
	}
	if s.Frequency != "every minute" {
		t.Errorf("expected 'every minute', got %q", s.Frequency)
	}
}

func TestSummarize_Hourly(t *testing.T) {
	s, err := Summarize("0 * * * *", summaryBase, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.DailyCount != 24 {
		t.Errorf("expected 24 daily fires, got %d", s.DailyCount)
	}
	if s.Frequency != "multiple times per day" {
		t.Errorf("expected 'multiple times per day', got %q", s.Frequency)
	}
}

func TestSummarize_OnceDaily(t *testing.T) {
	s, err := Summarize("0 9 * * *", summaryBase, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.DailyCount != 1 {
		t.Errorf("expected 1 daily fire, got %d", s.DailyCount)
	}
	if s.Frequency != "once per day" {
		t.Errorf("expected 'once per day', got %q", s.Frequency)
	}
	if s.WeeklyCount != 7 {
		t.Errorf("expected 7 weekly fires, got %d", s.WeeklyCount)
	}
}

func TestSummarize_InvalidExpr(t *testing.T) {
	_, err := Summarize("invalid", summaryBase, 5)
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestSummarize_DefaultN(t *testing.T) {
	s, err := Summarize("0 * * * *", summaryBase, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.NextTimes) != 5 {
		t.Errorf("expected default 5 next times, got %d", len(s.NextTimes))
	}
}

func TestSortSummaries(t *testing.T) {
	summaries := []*ScheduleSummary{
		{Expression: "0 9 * * *", DailyCount: 1},
		{Expression: "* * * * *", DailyCount: 1440},
		{Expression: "0 * * * *", DailyCount: 24},
	}
	SortSummaries(summaries)
	if summaries[0].DailyCount != 1440 {
		t.Errorf("expected first to have highest daily count, got %d", summaries[0].DailyCount)
	}
	if summaries[2].DailyCount != 1 {
		t.Errorf("expected last to have lowest daily count, got %d", summaries[2].DailyCount)
	}
}
