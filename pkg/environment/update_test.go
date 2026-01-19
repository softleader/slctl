package environment

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/go-github/v69/github"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/release"
)

func TestCheckForUpdates_NotReleased(t *testing.T) {
	log := logrus.New()
	metadata := release.NewMetadata("v1.2.3", "") // commit will be "unknown"
	err := CheckForUpdates(log, "some-home", metadata, false)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestCheckForUpdates(t *testing.T) {
	log := logrus.New()
	tempHome, _ := os.MkdirTemp("", "sl-home-update")
	defer os.RemoveAll(tempHome)
	home := paths.Home(tempHome)
	os.MkdirAll(home.Config(), 0755)

	// Create a dummy config file
	configFile := home.ConfigFile()
	os.WriteFile(configFile, []byte("token: some-token\ncheckUpdates: 2020-01-01T00:00:00Z"), 0644)

	metadata := &release.Metadata{GitVersion: "1.0.0", GitCommit: "abcdef123456"}

	// Mock latestRelease to return a newer version
	oldLatestRelease := latestRelease
	defer func() { latestRelease = oldLatestRelease }()

	latestRelease = func(ctx context.Context, log *logrus.Logger) (*github.RepositoryRelease, error) {
		tag := "1.1.0"
		commit := "abcdef"
		return &github.RepositoryRelease{
			TagName:         &tag,
			TargetCommitish: &commit,
		}, nil
	}

	err := CheckForUpdates(log, home, metadata, true)
	if err != nil {
		t.Fatalf("CheckForUpdates failed: %v", err)
	}

	// Mock latestRelease to return same version
	latestRelease = func(ctx context.Context, log *logrus.Logger) (*github.RepositoryRelease, error) {
		tag := "1.0.0"
		commit := "abcdef"
		return &github.RepositoryRelease{
			TagName:         &tag,
			TargetCommitish: &commit,
		}, nil
	}

	err = CheckForUpdates(log, home, metadata, true)
	if err != nil {
		t.Fatalf("CheckForUpdates failed: %v", err)
	}

	// Test force=false, no need to update
	future := time.Now().Add(24 * time.Hour)
	os.WriteFile(configFile, []byte(fmt.Sprintf("token: secret\ncheckUpdates: %s", future.Format(time.RFC3339))), 0644)
	err = CheckForUpdates(log, home, metadata, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNeedsToCheckOnline(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	future := time.Now().Add(1 * time.Hour)

	if !needsToCheckOnline(past) {
		t.Error("expected true for past time")
	}
	if needsToCheckOnline(future) {
		t.Error("expected false for future time")
	}
}
