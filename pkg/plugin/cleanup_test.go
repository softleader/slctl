package plugin

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/paths"
)

func TestCleanup(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "sl-home-cleanup")
	defer os.RemoveAll(tempHome)
	home := paths.Home(tempHome)
	os.MkdirAll(home.Config(), 0755)
	os.MkdirAll(home.Plugins(), 0755)
	os.MkdirAll(home.CachePlugins(), 0755)
	os.MkdirAll(home.CacheArchives(), 0755)

	// Create a dummy config file
	configFile := home.ConfigFile()
	os.WriteFile(configFile, []byte("cleanup: 2020-01-01T00:00:00Z"), 0644)

	// Create some dummy cache files
	os.MkdirAll(filepath.Join(home.CachePlugins(), "outdated-plugin"), 0755)
	os.WriteFile(filepath.Join(home.CacheArchives(), "outdated.zip"), []byte("test"), 0644)

	log := logrus.New()
	err := Cleanup(log, home, true, false)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Check if outdated files are removed
	if _, err := os.Stat(filepath.Join(home.CachePlugins(), "outdated-plugin")); !os.IsNotExist(err) {
		t.Error("expected outdated-plugin to be removed")
	}
	if _, err := os.Stat(filepath.Join(home.CacheArchives(), "outdated.zip")); !os.IsNotExist(err) {
		t.Error("expected outdated.zip to be removed")
	}
}

func TestNeedsToCleanup(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	future := time.Now().Add(1 * time.Hour)

	if !needsToCleanup(past) {
		t.Error("expected true for past time")
	}
	if needsToCleanup(future) {
		t.Error("expected false for future time")
	}
}

func TestNeedsToRemove(t *testing.T) {
	log := logrus.New()
	tempHome, _ := os.MkdirTemp("", "sl-home-rem")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)
	
	// Case 1: non-special file
	if !needsToRemove(log, hh, "/tmp/some-file") {
		t.Error("expected true for random file")
	}

	// Case 2: Repository file
	repoFile := hh.CacheRepositoryFile()
	os.MkdirAll(filepath.Dir(repoFile), 0755)
	
	r := &Repository{Expires: time.Now().Add(1 * time.Hour)}
	r.save(repoFile)
	
	if needsToRemove(log, hh, repoFile) {
		t.Error("expected false for fresh repository file")
	}
	
	r.Expires = time.Now().Add(-1 * time.Hour)
	r.save(repoFile)
	if !needsToRemove(log, hh, repoFile) {
		t.Error("expected true for expired repository file")
	}
}

func TestCleanup_ConfigError(t *testing.T) {
	log := logrus.New()
	tempHome, _ := os.MkdirTemp("", "sl-home-cleanup-err")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)
	os.MkdirAll(hh.Config(), 0755)
	os.WriteFile(hh.ConfigFile(), []byte("invalid yaml : :"), 0644)

	err := Cleanup(log, hh, false, false)
	if err == nil {
		t.Error("expected error for invalid config yaml")
	}
}


