// Package reporter provides formatting and output functionality for croncheck
// audit results. It supports multiple output formats including plain text and
// JSON, making it easy to integrate croncheck output into both human-readable
// reports and automated pipelines.
//
// Usage:
//
//	w := reporter.NewWriter(reporter.FormatText, os.Stdout)
//	err := w.Write(reporter.Report{
//		Conflicts:        conflicts,
//		Misconfigurations: misconfigs,
//	})
package reporter
