package loader

import "github.com/croncheck/internal/audit"

// ToAuditEntries converts loader entries into the audit.Entry type
// expected by the audit package, preserving file and line metadata
// in the Label field.
func ToAuditEntries(entries []Entry) []audit.Entry {
	result := make([]audit.Entry, 0, len(entries))
	for _, e := range entries {
		result = append(result, audit.Entry{
			Label:      formatLabel(e),
			Expression: e.Expression,
		})
	}
	return result
}

// formatLabel builds a human-readable label from file path and line number.
func formatLabel(e Entry) string {
	if e.File == "" {
		return e.Raw
	}
	return fmt.Sprintf("%s:%d", e.File, e.Line)
}

import "fmt" // grouped here for clarity; move to top in real code
