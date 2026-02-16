package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWatcher_DetectsFileChange(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.json")

	if err := os.WriteFile(file, []byte(`{"a":1}`), 0644); err != nil {
		t.Fatal(err)
	}

	w, err := New([]string{file, dir})
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	// Modify the file
	if err := os.WriteFile(file, []byte(`{"a":2}`), 0644); err != nil {
		t.Fatal(err)
	}

	// Debounce (300ms) + margin
	select {
	case msg := <-w.ch:
		if _, ok := msg.(FileChangedMsg); !ok {
			t.Fatalf("expected FileChangedMsg, got %T", msg)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout: FileChangedMsg not received")
	}
}

func TestWatcher_SkipsMissingPaths(t *testing.T) {
	existing := t.TempDir()
	missing := filepath.Join(existing, "no_such_dir")

	w, err := New([]string{missing, existing})
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	// Success if created without error
	if w.fsw == nil {
		t.Fatal("fsw should not be nil")
	}
}
