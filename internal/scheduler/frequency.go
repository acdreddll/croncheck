package scheduler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/croncheck/internal/parser"
)

// FrequencyInfo holds computed frequency metadata for a cron expression.
type FrequencyInfo struct {
	// Expression is the original cron string.
	Expression string
	// RunsPerHour is the number of times the job fires in a 60-minute window.
	RunsPerHour int
	// RunsPerDay is the number of times the job fires in a 24-hour window.
	RunsPerDay int
	// Label is a human-readable description such as "every minute" or "once daily".
	Label string
}

// ComputeFrequency returns frequency statistics for the given cron expression.
// It returns an error if the expression cannot be parsed.
func ComputeFrequency(expr string) (FrequencyInfo, error) {
	entry, err := parser.Parse(expr)
	if err != nil {
		return FrequencyInfo{}, fmt.Errorf("ComputeFrequency: %w", err)
	}

	minutesPerHour := countUniqueMinutes(entry.Minute)
	hours := countUniqueValues(entry.Hour, 0, 23)

	runsPerHour := minutesPerHour
	runsPerDay := minutesPerHour * hours

	return FrequencyInfo{
		Expression:  expr,
		RunsPerHour: runsPerHour,
		RunsPerDay:  runsPerDay,
		Label:       buildLabel(runsPerHour, runsPerDay),
	}, nil
}

// countUniqueMinutes returns how many distinct minute values the field covers.
func countUniqueMinutes(field string) int {
	return countUniqueValues(field, 0, 59)
}

// countUniqueValues expands a cron field string and counts distinct values in [min, max].
func countUniqueValues(field string, min, max int) int {
	set := make(map[int]struct{})
	parts := strings.Split(field, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "*" {
			for v := min; v <= max; v++ {
				set[v] = struct{}{}
			}
			continue
		}
		if strings.Contains(p, "/") {
			expandStep(p, min, max, set)
			continue
		}
		if strings.Contains(p, "-") {
			expandRange(p, set)
			continue
		}
		if v, err := strconv.Atoi(p); err == nil {
			set[v] = struct{}{}
		}
	}
	return len(set)
}

func expandStep(field string, min, max int, set map[int]struct{}) {
	parts := strings.SplitN(field, "/", 2)
	step, err := strconv.Atoi(parts[1])
	if err != nil || step <= 0 {
		return
	}
	start := min
	if parts[0] != "*" {
		if v, err := strconv.Atoi(parts[0]); err == nil {
			start = v
		}
	}
	for v := start; v <= max; v += step {
		set[v] = struct{}{}
	}
}

func expandRange(field string, set map[int]struct{}) {
	parts := strings.SplitN(field, "-", 2)
	lo, err1 := strconv.Atoi(parts[0])
	hi, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return
	}
	for v := lo; v <= hi; v++ {
		set[v] = struct{}{}
	}
}

func buildLabel(perHour, perDay int) string {
	switch {
	case perHour == 60:
		return "every minute"
	case perHour > 1:
		return fmt.Sprintf("%d times per hour", perHour)
	case perDay == 24:
		return "once per hour"
	case perDay > 1:
		return fmt.Sprintf("%d times per day", perDay)
	default:
		return "once daily"
	}
}
