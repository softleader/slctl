package paths

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestEnsureDirectory(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "sl-paths-test")
	defer os.RemoveAll(tempDir)

	log := logrus.New()
	target := filepath.Join(tempDir, "a/b/c")

	if err := EnsureDirectory(log, target); err != nil {
		t.Fatalf("EnsureDirectory failed: %v", err)
	}

	if !IsExistDirectory(target) {
		t.Errorf("expected %s to exist as directory", target)
	}

	// Case: existing file instead of directory
	filePath := filepath.Join(tempDir, "file.txt")
	os.WriteFile(filePath, []byte("test"), 0644)
	if err := EnsureDirectory(log, filePath); err == nil {
		t.Error("expected error when path is an existing file")
	}
}

func TestEnsureDirectories(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "sl-dirs-test")
	defer os.RemoveAll(tempDir)

	d1 := filepath.Join(tempDir, "d1")
	d2 := filepath.Join(tempDir, "d2")

	if err := EnsureDirectories(logrus.New(), d1, d2); err != nil {
		t.Fatal(err)
	}
	if !IsExistDirectory(d1) || !IsExistDirectory(d2) {
		t.Error("expected d1 and d2 to exist")
	}
}

func TestIsExistFile(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "sl-file-test")
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "test.txt")
	if IsExistFile(filePath) {
		t.Error("expected false for non-existent file")
	}

	os.WriteFile(filePath, []byte("test"), 0644)
	if !IsExistFile(filePath) {
		t.Error("expected true for existing file")
	}

	if IsExistFile(tempDir) {
		t.Error("expected false for directory")
	}
}

func TestHome(t *testing.T) {
	h := Home("/tmp/sl")
	if h.Plugins() != "/tmp/sl/plugins" {
		t.Errorf("expected /tmp/sl/plugins, got %s", h.Plugins())
	}
	if h.Config() != "/tmp/sl/config" {
		t.Errorf("expected /tmp/sl/config, got %s", h.Config())
	}
	if h.ConfigFile() != "/tmp/sl/config/configs.yaml" {
		t.Errorf("expected /tmp/sl/config/configs.yaml, got %s", h.ConfigFile())
	}
	if h.Cache() != "/tmp/sl/cache" {
		t.Errorf("expected /tmp/sl/cache, got %s", h.Cache())
	}
	if h.CachePlugins() != "/tmp/sl/cache/plugins" {
		t.Errorf("expected /tmp/sl/cache/plugins, got %s", h.CachePlugins())
	}
	if h.CacheRepositoryFile() != "/tmp/sl/cache/plugins/repository.yaml" {
		t.Errorf("expected /tmp/sl/cache/plugins/repository.yaml, got %s", h.CacheRepositoryFile())
	}
	if h.CacheArchives() != "/tmp/sl/cache/archives" {
		t.Errorf("expected /tmp/sl/cache/archives, got %s", h.CacheArchives())
	}
	if h.Mounts() != "/tmp/sl/mounts" {
		t.Errorf("expected /tmp/sl/mounts, got %s", h.Mounts())
	}
	if h.ContainsAnySpace() {
		t.Error("expected false for /tmp/sl")
	}

	h2 := Home("/tmp/sl space")
	if !h2.ContainsAnySpace() {
		t.Error("expected true for /tmp/sl space")
	}
}
