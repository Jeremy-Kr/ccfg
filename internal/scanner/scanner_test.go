package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeremy-kr/ccfg/internal/model"
)

func TestScanWithProject(t *testing.T) {
	// Create temporary project structure
	tmp := t.TempDir()
	projectDir := filepath.Join(tmp, "project")
	if err := os.MkdirAll(filepath.Join(projectDir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(projectDir, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Create test configuration files
	if err := os.WriteFile(filepath.Join(projectDir, "CLAUDE.md"), []byte("# Test"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, ".claude", "settings.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	s := New(projectDir)
	result, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	// Project root should be detected
	if result.RootDir != projectDir {
		t.Errorf("RootDir = %q, want %q", result.RootDir, projectDir)
	}

	// Project scope should contain files
	if len(result.Project) == 0 {
		t.Fatal("Project scope is empty")
	}

	// Verify Exists flag for existing files
	found := false
	for _, cf := range result.Project {
		if cf.Path == filepath.Join(projectDir, "CLAUDE.md") {
			found = true
			if !cf.Exists {
				t.Error("CLAUDE.md exists but Exists=false")
			}
			if cf.FileType != model.FileTypeMarkdown {
				t.Errorf("CLAUDE.md FileType = %d, want FileTypeMarkdown", cf.FileType)
			}
			if cf.Category != model.CategoryInstructions {
				t.Errorf("CLAUDE.md Category = %d, want CategoryInstructions", cf.Category)
			}
		}
	}
	if !found {
		t.Error("CLAUDE.md not found in Project scope")
	}
}

func TestScanWithoutProject(t *testing.T) {
	// Directory without .git
	tmp := t.TempDir()

	s := New(tmp)
	result, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	if result.RootDir != "" {
		t.Errorf("RootDir = %q, want empty", result.RootDir)
	}
	if len(result.Project) != 0 {
		t.Errorf("Project scope length = %d, want 0", len(result.Project))
	}
}

func TestDetectFileType(t *testing.T) {
	tests := []struct {
		path string
		want model.FileType
	}{
		{"settings.json", model.FileTypeJSON},
		{"settings.jsonc", model.FileTypeJSONC},
		{"CLAUDE.md", model.FileTypeMarkdown},
		{".mcp.json", model.FileTypeJSON},
		{"unknown.txt", model.FileTypeJSON}, // default is JSON
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := detectFileType(tt.path)
			if got != tt.want {
				t.Errorf("detectFileType(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
