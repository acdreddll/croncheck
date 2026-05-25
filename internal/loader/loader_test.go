package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/croncheck/internal/loader"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "crontab")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return path
}

func TestLoadFile_Valid(t *testing.T) {
	path := writeTempFile(t, "# comment\n\n0 * * * * /usr/bin/backup\n*/5 * * * * /usr/bin/check\n")

	entries, err := loader.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Line != 3 {
		t.Errorf("expected line 3, got %d", entries[0].Line)
	}
	if entries[1].Line != 4 {
		t.Errorf("expected line 4, got %d", entries[1].Line)
	}
}

func TestLoadFile_InvalidExpression(t *testing.T) {
	path := writeTempFile(t, "99 * * * * /usr/bin/bad\n")

	_, err := loader.LoadFile(path)
	if err == nil {
		t.Fatal("expected error for invalid expression, got nil")
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := loader.LoadFile("/nonexistent/path/crontab")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadFiles_Multiple(t *testing.T) {
	p1 := writeTempFile(t, "0 1 * * * /bin/job1\n")
	p2 := writeTempFile(t, "0 2 * * * /bin/job2\n0 3 * * * /bin/job3\n")

	entries, err := loader.LoadFiles([]string{p1, p2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestLoadFile_EmptyFile(t *testing.T) {
	path := writeTempFile(t, "# only comments\n\n")

	entries, err := loader.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}
