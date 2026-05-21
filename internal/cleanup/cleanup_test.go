package cleanup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRemoverDeletesPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	result := Remover{}.Remove(path)
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected file to be deleted, stat err was %v", err)
	}
}

func TestRemoverDryRunDoesNotDelete(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	result := Remover{DryRun: true}.Remove(path)
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if !result.DryRun {
		t.Fatal("expected dry-run result")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to remain: %v", err)
	}
}
