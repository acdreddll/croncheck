package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// CronExpression holds the parsed fields of a cron expression.
type CronExpression struct {
	Raw     string
	Minute  string
	Hour    string
	Day     string
	Month   string
	Weekday string
}

// Parse parses a standard 5-field cron expression string.
func Parse(raw string) (*CronExpression, error) {
	fields := strings.Fields(raw)
	if len(fields) != 5 {
		return nil, fmt.Errorf("invalid cron expression %q: expected 5 fields, got %d", raw, len(fields))
	}
	expr := &CronExpression{
		Raw:     raw,
		Minute:  fields[0],
		Hour:    fields[1],
		Day:     fields[2],
		Month:   fields[3],
		Weekday: fields[4],
	}
	if err := validate(expr); err != nil {
		return nil, err
	}
	return expr, nil
}

// validate performs basic range validation on each cron field.
func validate(e *CronExpression) error {
	type fieldSpec struct {
		name  string
		value string
		min   int
		max   int
	}
	fields := []fieldSpec{
		{"minute", e.Minute, 0, 59},
		{"hour", e.Hour, 0, 23},
		{"day", e.Day, 1, 31},
		{"month", e.Month, 1, 12},
		{"weekday", e.Weekday, 0, 7},
	}
	for _, f := range fields {
		if err := validateField(f.name, f.value, f.min, f.max); err != nil {
			return err
		}
	}
	return nil
}

func validateField(name, value string, min, max int) error {
	if value == "*" {
		return nil
	}
	// Handle step values like */5 or 1-5/2
	parts := strings.SplitN(value, "/", 2)
	base := parts[0]
	if base != "*" {
		for _, segment := range strings.Split(base, ",") {
			if strings.Contains(segment, "-") {
				rangeParts := strings.SplitN(segment, "-", 2)
				lo, err1 := strconv.Atoi(rangeParts[0])
				hi, err2 := strconv.Atoi(rangeParts[1])
				if err1 != nil || err2 != nil || lo < min || hi > max || lo > hi {
					return fmt.Errorf("invalid range %q in field %s", segment, name)
				}
			} else {
				v, err := strconv.Atoi(segment)
				if err != nil || v < min || v > max {
					return fmt.Errorf("invalid value %q in field %s (allowed %d-%d)", segment, name, min, max)
				}
			}
		}
	}
	if len(parts) == 2 {
		step, err := strconv.Atoi(parts[1])
		if err != nil || step < 1 {
			return fmt.Errorf("invalid step %q in field %s", parts[1], name)
		}
	}
	return nil
}
