package loader

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/croncheck/internal/parser"
)

// Entry represents a single cron entry read from a file.
type Entry struct {
	File       string
	Line       int
	Raw        string
	Expression *parser.CronExpression
}

// LoadFile reads a crontab-style file and returns parsed cron entries.
// Lines beginning with '#' or that are blank are skipped.
func LoadFile(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("loader: open %q: %w", path, err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		expr, err := parser.Parse(line)
		if err != nil {
			return nil, fmt.Errorf("loader: %q line %d: %w", path, lineNum, err)
		}

		entries = append(entries, Entry{
			File:       path,
			Line:       lineNum,
			Raw:        line,
			Expression: expr,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("loader: scanning %q: %w", path, err)
	}

	return entries, nil
}

// LoadFiles loads cron entries from multiple files, merging results.
func LoadFiles(paths []string) ([]Entry, error) {
	var all []Entry
	for _, p := range paths {
		entries, err := LoadFile(p)
		if err != nil {
			return nil, err
		}
		all = append(all, entries...)
	}
	return all, nil
}
