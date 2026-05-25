// Package loader reads crontab-style files and produces parsed cron entries
// ready for auditing.
//
// Supported file format:
//   - One cron expression per line (5 fields + optional command).
//   - Lines starting with '#' are treated as comments and ignored.
//   - Blank lines are ignored.
//
// Example usage:
//
//	entries, err := loader.LoadFile("/etc/cron.d/myjobs")
//	if err != nil {
//		log.Fatal(err)
//	}
package loader
