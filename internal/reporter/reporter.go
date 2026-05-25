// Package reporter formats and outputs audit results.
package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/croncheck/internal/audit"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds the results of an audit run.
type Report struct {
	Conflicts        []audit.Conflict
	Misconfigurations []audit.Misconfiguration
}

// Writer writes a Report to an io.Writer in the specified format.
type Writer struct {
	format Format
	out    io.Writer
}

// NewWriter creates a new Writer with the given format and output destination.
func NewWriter(format Format, out io.Writer) *Writer {
	return &Writer{format: format, out: out}
}

// Write outputs the report to the writer's destination.
func (w *Writer) Write(r Report) error {
	switch w.format {
	case FormatJSON:
		return writeJSON(w.out, r)
	default:
		return writeText(w.out, r)
	}
}

func writeText(out io.Writer, r Report) error {
	if len(r.Conflicts) == 0 && len(r.Misconfigurations) == 0 {
		_, err := fmt.Fprintln(out, "No issues found.")
		return err
	}

	var sb strings.Builder
	if len(r.Conflicts) > 0 {
		sb.WriteString(fmt.Sprintf("Conflicts (%d):\n", len(r.Conflicts)))
		for _, c := range r.Conflicts {
			sb.WriteString(fmt.Sprintf("  [CONFLICT] %s\n", c.Message))
		}
	}
	if len(r.Misconfigurations) > 0 {
		sb.WriteString(fmt.Sprintf("Misconfigurations (%d):\n", len(r.Misconfigurations)))
		for _, m := range r.Misconfigurations {
			sb.WriteString(fmt.Sprintf("  [MISCONFIG] %s\n", m.Message))
		}
	}
	_, err := fmt.Fprint(out, sb.String())
	return err
}
