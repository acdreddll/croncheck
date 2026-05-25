// Package audit provides conflict detection and misconfiguration analysis
// for parsed cron expressions.
package audit

import (
	"fmt"

	"github.com/user/croncheck/internal/parser"
)

// ConflictType describes the kind of conflict detected.
type ConflictType string

const (
	ConflictOverlap    ConflictType = "overlap"
	ConflictDuplicate  ConflictType = "duplicate"
	ConflictSubsumed   ConflictType = "subsumed"
)

// Conflict represents a detected scheduling conflict between two cron entries.
type Conflict struct {
	Type    ConflictType
	A       parser.CronEntry
	B       parser.CronEntry
	Message string
}

func (c Conflict) Error() string {
	return fmt.Sprintf("%s conflict between %q and %q: %s", c.Type, c.A.Raw, c.B.Raw, c.Message)
}

// DetectConflicts compares all pairs of cron entries and returns any conflicts found.
func DetectConflicts(entries []parser.CronEntry) []Conflict {
	var conflicts []Conflict
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if c, ok := checkPair(entries[i], entries[j]); ok {
				conflicts = append(conflicts, c)
			}
		}
	}
	return conflicts
}

func checkPair(a, b parser.CronEntry) (Conflict, bool) {
	if a.Raw == b.Raw {
		return Conflict{
			Type:    ConflictDuplicate,
			A:       a,
			B:       b,
			Message: "identical cron expressions",
		}, true
	}
	if fieldsOverlap(a, b) {
		return Conflict{
			Type:    ConflictOverlap,
			A:       a,
			B:       b,
			Message: "schedules may fire at the same time",
		}, true
	}
	return Conflict{}, false
}

// fieldsOverlap returns true when both entries use wildcard (*) in every field,
// meaning they will always fire together.
func fieldsOverlap(a, b parser.CronEntry) bool {
	aFields := [5]string{a.Minute, a.Hour, a.DayOfMonth, a.Month, a.DayOfWeek}
	bFields := [5]string{b.Minute, b.Hour, b.DayOfMonth, b.Month, b.DayOfWeek}
	for i := range aFields {
		if aFields[i] == "*" && bFields[i] == "*" {
			continue
		}
		if aFields[i] != bFields[i] {
			return false
		}
	}
	return true
}
