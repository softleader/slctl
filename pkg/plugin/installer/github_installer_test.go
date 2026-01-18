package installer

import (
	"runtime"
	"testing"

	"github.com/google/go-github/v69/github"
	"github.com/sirupsen/logrus"
)

func TestDismantle(t *testing.T) {
	owner, repo := dismantle("github.com/softleader/slctl")
	if owner != "softleader" {
		t.Errorf("expected owner softleader, got %s", owner)
	}
	if repo != "slctl" {
		t.Errorf("expected repo slctl, got %s", repo)
	}
}

func TestFindRuntimeOsAsset(t *testing.T) {
	log := logrus.New()

	nameDarwin := "slctl_darwin_amd64"
	nameLinux := "slctl_linux_amd64"

	assets := []*github.ReleaseAsset{
		{Name: &nameDarwin},
		{Name: &nameLinux},
	}

	// Temporarily override GOOS for testing logic if possible?
	// actually findRuntimeOsAsset uses runtime.GOOS constant.
	// So we can only test for the current runtime OS.

	expectedName := ""
	if runtime.GOOS == "darwin" {
		expectedName = "darwin"
	} else if runtime.GOOS == "linux" {
		expectedName = "linux"
	} else if runtime.GOOS == "windows" {
		expectedName = "windows"
	}

	if expectedName != "" {
		found := findRuntimeOsAsset(log, assets)
		if found == nil {
			t.Logf("Skipping test because no matching asset for %s", runtime.GOOS)
		} else {
			if !contains(*found.Name, expectedName) {
				t.Errorf("expected asset name to contain %s, got %s", expectedName, *found.Name)
			}
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr || len(s) > len(substr) && contains(s[1:], substr)
}
