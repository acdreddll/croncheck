package audit

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/croncheck/internal/parser"
)

// Misconfiguration describes a suspicious or invalid cron configuration.
type Misconfiguration struct {
	Entry   parser.CronEntry
	Field   string
	Message string
}

func (m Misconfiguration) Error() string {
	return fmt.Sprintf("misconfiguration in field %q of %q: %s", m.Field, m.Entry.Raw, m.Message)
}

// DetectMisconfigurations inspects each entry for common issues.
func DetectMisconfigurations(entries []parser.CronEntry) []Misconfiguration {
	var issues []Misconfiguration
	for _, e := range entries {
		issues = append(issues, checkEntry(e)...)
	}
	return issues
}

func checkEntry(e parser.CronEntry) []Misconfiguration {
	var issues []Misconfiguration
	checks := []struct {
		field string
		value string
		min   int
		max   int
	}{
		{"minute", e.Minute, 0, 59},
		{"hour", e.Hour, 0, 23},
		{"day-of-month", e.DayOfMonth, 1, 31},
		{"month", e.Month, 1, 12},
		{"day-of-week", e.DayOfWeek, 0, 7},
	}
	for _, c := range checks {
		if c.value == "*" {
			continue
		}
		if msg := validateRange(c.value, c.min, c.max); msg != "" {
			issues = append(issues, Misconfiguration{
				Entry:   e,
				Field:   c.field,
				Message: msg,
			})
		}
	}
	return issues
}

func validateRange(value string, min, max int) string {
	parts := strings.Split(value, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if strings.Contains(p, "/") {
			sub := strings.SplitN(p, "/", 2)
			if step, err := strconv.Atoi(sub[1]); err != nil || step <= 0 {
				return fmt.Sprintf("invalid step value %q", sub[1])
			}
			continue
		}
		if strings.Contains(p, "-") {
			sub := strings.SplitN(p, "-", 2)
			lo, err1 := strconv.Atoi(sub[0])
			hi, err2 := strconv.Atoi(sub[1])
			if err1 != nil || err2 != nil {
				return fmt.Sprintf("non-numeric range %q", p)
			}
			if lo > hi {
				return fmt.Sprintf("range start %d > end %d", lo, hi)
			}
			if lo < min || hi > max {
				return fmt.Sprintf("range %d-%d out of bounds [%d,%d]", lo, hi, min, max)
			}
			continue
		}
		v, err := strconv.Atoi(p)
		if err != nil {
			return fmt.Sprintf("non-numeric value %q", p)
		}
		if v < min || v > max {
			return fmt.Sprintf("value %d out of bounds [%d,%d]", v, min, max)
		}
	}
	return ""
}
