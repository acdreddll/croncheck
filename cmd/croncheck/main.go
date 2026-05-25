package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/croncheck/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: croncheck <file>")
		os.Exit(1)
	}

	path := os.Args[1]
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	var (
		line    int
		errors  int
		parsed  int
	)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line++
		text := strings.TrimSpace(scanner.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		expr, err := parser.Parse(text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "line %d: %v\n", line, err)
			errors++
			continue
		}
		fmt.Printf("line %d: OK  %s\n", line, expr.Raw)
		parsed++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nSummary: %d valid, %d invalid\n", parsed, errors)
	if errors > 0 {
		os.Exit(2)
	}
}
