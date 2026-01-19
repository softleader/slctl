package installer

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"runtime"
	"testing"

	"github.com/google/go-github/v69/github"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
)

func TestNewGitHubInstaller(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := github.NewClient(nil)
	u, _ := url.Parse(server.URL + "/")
	client.BaseURL = u
	client.UploadURL = u

	// Mock tokenClient
	oldTokenClient := tokenClient
	tokenClient = func(ctx context.Context, token string) (*github.Client, error) {
		return client, nil
	}
	defer func() { tokenClient = oldTokenClient }()

	org := "softleader"
	repo := "slctl"
	
	mux.HandleFunc("/repos/"+org+"/"+repo+"/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"tag_name": "v1.0.0", "assets": [{"id": 1, "name": "asset1", "browser_download_url": "http://gh.com/asset1"}]}`)
	})
	mux.HandleFunc("/repos/"+org+"/"+repo+"/releases/assets/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "asset content")
	})

	tempHome, _ := os.MkdirTemp("", "sl-home-ghi")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)
	os.MkdirAll(hh.Config(), 0755)
	os.WriteFile(hh.ConfigFile(), []byte("token: secret"), 0644)

	log := logrus.New()
	ghi, err := newGitHubInstaller(log, "github.com/"+org+"/"+repo, "", 0, hh, &InstallOption{})
	if err != nil {
		t.Fatalf("newGitHubInstaller failed: %v", err)
	}

	if ghi.source != "github.com/"+org+"/"+repo {
		t.Errorf("expected source github.com/softleader/slctl, got %s", ghi.source)
	}
}

func TestNewGitHubInstaller_Offline(t *testing.T) {
	environment.Settings.Offline = true
	defer func() { environment.Settings.Offline = false }()

	hh := paths.Home("/tmp")
	log := logrus.New()
	_, err := newGitHubInstaller(log, "github.com/softleader/slctl", "", 0, hh, &InstallOption{})
	if err != errNonResolvableInOfflineMode {
		t.Errorf("expected %v, got %v", errNonResolvableInOfflineMode, err)
	}
}

func TestNewGitHubInstaller_TokenError(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "sl-home-ghi-err")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)
	os.MkdirAll(hh.Config(), 0755)
	// Config file exists but is empty/malformed
	os.WriteFile(hh.ConfigFile(), []byte("token: "), 0644)

	log := logrus.New()
	_, err := newGitHubInstaller(log, "github.com/softleader/slctl", "", 0, hh, &InstallOption{})
	if err == nil {
		t.Error("expected error for missing token")
	}
}

func TestNewGitHubInstaller_LatestReleaseError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := github.NewClient(nil)
	u, _ := url.Parse(server.URL + "/")
	client.BaseURL = u
	client.UploadURL = u

	oldTokenClient := tokenClient
	tokenClient = func(ctx context.Context, token string) (*github.Client, error) {
		return client, nil
	}
	defer func() { tokenClient = oldTokenClient }()

	mux.HandleFunc("/repos/softleader/slctl/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	tempHome, _ := os.MkdirTemp("", "sl-home-ghi-lre")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)
	os.MkdirAll(hh.Config(), 0755)
	os.WriteFile(hh.ConfigFile(), []byte("token: secret"), 0644)

	log := logrus.New()
	_, err := newGitHubInstaller(log, "github.com/softleader/slctl", "", 0, hh, &InstallOption{})
	if err == nil {
		t.Error("expected error for failed latest release fetch")
	}
}

func TestNewGitHubInstaller_DownloadAssetError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := github.NewClient(nil)
	u, _ := url.Parse(server.URL + "/")
	client.BaseURL = u
	client.UploadURL = u

	oldTokenClient := tokenClient
	tokenClient = func(ctx context.Context, token string) (*github.Client, error) {
		return client, nil
	}
	defer func() { tokenClient = oldTokenClient }()

	mux.HandleFunc("/repos/softleader/slctl/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"tag_name": "v1.0.0", "assets": [{"id": 1, "name": "asset1"}]}`)
	})
	mux.HandleFunc("/repos/softleader/slctl/releases/assets/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	tempHome, _ := os.MkdirTemp("", "sl-home-ghi-dae")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)
	os.MkdirAll(hh.Config(), 0755)
	os.WriteFile(hh.ConfigFile(), []byte("token: secret"), 0644)

	log := logrus.New()
	_, err := newGitHubInstaller(log, "github.com/softleader/slctl", "", 0, hh, &InstallOption{})
	if err == nil {
		t.Error("expected error for failed asset download")
	}
}

func TestNewGitHubInstaller_NoAssetsError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := github.NewClient(nil)
	u, _ := url.Parse(server.URL + "/")
	client.BaseURL = u
	client.UploadURL = u

	oldTokenClient := tokenClient
	tokenClient = func(ctx context.Context, token string) (*github.Client, error) {
		return client, nil
	}
	defer func() { tokenClient = oldTokenClient }()

	mux.HandleFunc("/repos/softleader/slctl/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"tag_name": "v1.0.0", "assets": []}`)
	})

	tempHome, _ := os.MkdirTemp("", "sl-home-ghi-nae")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)
	os.MkdirAll(hh.Config(), 0755)
	os.WriteFile(hh.ConfigFile(), []byte("token: secret"), 0644)

	log := logrus.New()
	_, err := newGitHubInstaller(log, "github.com/softleader/slctl", "", 0, hh, &InstallOption{})
	if err == nil {
		t.Error("expected error for no assets")
	}
}

func TestNewGitHubInstaller_ReleaseByTagError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := github.NewClient(nil)
	u, _ := url.Parse(server.URL + "/")
	client.BaseURL = u
	client.UploadURL = u

	oldTokenClient := tokenClient
	tokenClient = func(ctx context.Context, token string) (*github.Client, error) {
		return client, nil
	}
	defer func() { tokenClient = oldTokenClient }()

	mux.HandleFunc("/repos/softleader/slctl/releases/tags/v1.2.3", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	tempHome, _ := os.MkdirTemp("", "sl-home-ghi-rbe")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)
	os.MkdirAll(hh.Config(), 0755)
	os.WriteFile(hh.ConfigFile(), []byte("token: secret"), 0644)

	log := logrus.New()
	_, err := newGitHubInstaller(log, "github.com/softleader/slctl", "v1.2.3", 0, hh, &InstallOption{})
	if err == nil {
		t.Error("expected error for failed release fetch by tag")
	}
}

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

func TestPickAsset(t *testing.T) {
	log := logrus.New()
	
	name1 := "asset1"
	name2 := "asset2"
	assets := []*github.ReleaseAsset{
		{Name: &name1},
		{Name: &name2},
	}
	release := &github.RepositoryRelease{
		TagName: github.Ptr("v1.0.0"),
		Assets:  assets,
	}

	// Pick specific asset index
	ra, err := pickAsset(log, release, 1)
	if err != nil {
		t.Fatal(err)
	}
	if ra.GetName() != "asset2" {
		t.Errorf("expected asset2, got %s", ra.GetName())
	}

	// Pick asset index out of range
	_, err = pickAsset(log, release, 5)
	if err == nil {
		t.Error("expected error for out of range index")
	}

	// Auto pick (uses findRuntimeOsAsset)
	ra, err = pickAsset(log, release, 0)
	if err != nil {
		t.Fatal(err)
	}
	if ra == nil {
		t.Error("expected an asset to be picked")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr || len(s) > len(substr) && contains(s[1:], substr)
}
