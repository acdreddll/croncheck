package scheduler

import (
	"testing"
)

func TestComputeFrequency_EveryMinute(t *testing.T) {
	info, err := ComputeFrequency("* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.RunsPerHour != 60 {
		t.Errorf("RunsPerHour = %d, want 60", info.RunsPerHour)
	}
	if info.RunsPerDay != 1440 {
		t.Errorf("RunsPerDay = %d, want 1440", info.RunsPerDay)
	}
	if info.Label != "every minute" {
		t.Errorf("Label = %q, want \"every minute\"", info.Label)
	}
}

func TestComputeFrequency_Hourly(t *testing.T) {
	info, err := ComputeFrequency("0 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.RunsPerHour != 1 {
		t.Errorf("RunsPerHour = %d, want 1", info.RunsPerHour)
	}
	if info.RunsPerDay != 24 {
		t.Errorf("RunsPerDay = %d, want 24", info.RunsPerDay)
	}
	if info.Label != "once per hour" {
		t.Errorf("Label = %q, want \"once per hour\"", info.Label)
	}
}

func TestComputeFrequency_OnceDaily(t *testing.T) {
	info, err := ComputeFrequency("30 2 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.RunsPerHour != 1 {
		t.Errorf("RunsPerHour = %d, want 1", info.RunsPerHour)
	}
	if info.RunsPerDay != 24 {
		t.Errorf("RunsPerDay = %d, want 24", info.RunsPerDay)
	}
}

func TestComputeFrequency_EveryFiveMinutes(t *testing.T) {
	info, err := ComputeFrequency("*/5 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.RunsPerHour != 12 {
		t.Errorf("RunsPerHour = %d, want 12", info.RunsPerHour)
	}
	if info.Label != "12 times per hour" {
		t.Errorf("Label = %q, want \"12 times per hour\"", info.Label)
	}
}

func TestComputeFrequency_TwiceDaily(t *testing.T) {
	info, err := ComputeFrequency("0 0,12 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.RunsPerDay != 2 {
		t.Errorf("RunsPerDay = %d, want 2", info.RunsPerDay)
	}
	if info.Label != "2 times per day" {
		t.Errorf("Label = %q, want \"2 times per day\"", info.Label)
	}
}

func TestComputeFrequency_InvalidExpr(t *testing.T) {
	_, err := ComputeFrequency("not a cron")
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}

func TestComputeFrequency_ExpressionPreserved(t *testing.T) {
	expr := "15 4 * * *"
	info, err := ComputeFrequency(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Expression != expr {
		t.Errorf("Expression = %q, want %q", info.Expression, expr)
	}
}
