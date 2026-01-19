package plugin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsLocalDirReference(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "sl-source-test")
	defer os.RemoveAll(tempDir)

	if !IsLocalDirReference(tempDir) {
		t.Errorf("expected %s to be local dir reference", tempDir)
	}

	tempFile := filepath.Join(tempDir, "file.txt")
	os.WriteFile(tempFile, []byte("test"), 0644)
	if IsLocalDirReference(tempFile) {
		t.Errorf("expected %s NOT to be local dir reference", tempFile)
	}

	if IsLocalDirReference("https://github.com/softleader/slctl") {
		t.Error("expected URL NOT to be local dir reference")
	}
}

func TestIsLocalReference(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "sl-source-test-ref")
	defer os.RemoveAll(tempDir)

	if !IsLocalReference(tempDir) {
		t.Errorf("expected %s to be local reference", tempDir)
	}

	tempFile := filepath.Join(tempDir, "file.txt")
	os.WriteFile(tempFile, []byte("test"), 0644)
	if !IsLocalReference(tempFile) {
		t.Errorf("expected %s to be local reference", tempFile)
	}

	if IsLocalReference("https://github.com/softleader/slctl") {
		t.Error("expected URL NOT to be local reference")
	}
}

func TestIsSupportedArchive(t *testing.T) {
	tests := []struct {
		source string
		want   bool
	}{
		{"plugin.zip", true},
		{"plugin.tar.gz", true},
		{"plugin.tgz", true},
		{"plugin.txt", false},
		{"plugin", false},
	}

	for _, tt := range tests {
		if got := IsSupportedArchive(tt.source); got != tt.want {
			t.Errorf("IsSupportedArchive(%q) = %v, want %v", tt.source, got, tt.want)
		}
	}
}

func TestIsGitHubRepo(t *testing.T) {
	tests := []struct {
		source string
		want   bool
	}{
		{"github.com/softleader/slctl", true},
		{"https://github.com/softleader/slctl", true},
		{"http://github.com/softleader/slctl/", true},
		{"gitlab.com/softleader/slctl", false},
		{"not-a-repo", false},
	}

	for _, tt := range tests {
		if got := IsGitHubRepo(tt.source); got != tt.want {
			t.Errorf("IsGitHubRepo(%q) = %v, want %v", tt.source, got, tt.want)
		}
	}
}
