package plugin

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v69/github"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"golang.org/x/oauth2"
)

func TestLoadRepository_Mock(t *testing.T) {
	// 1. Setup Mock Server
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/search/repositories", func(w http.ResponseWriter, r *http.Request) {
		// Calculate total_count, etc.
		// Verify Token?
		auth := r.Header.Get("Authorization")
		if auth != "Bearer mock-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"total_count": 1, "items": [{"full_name": "softleader/slctl-mock-plugin", "description": "mock plugin"}]}`)
	})

	// 2. Setup Temp Home & Config
	tmpHome := t.TempDir()
	home := paths.Home(tmpHome)

	// Create config file with token
	cfg := config.NewConfFile()
	cfg.Token = "mock-token"
	if err := os.MkdirAll(filepath.Dir(home.ConfigFile()), 0755); err != nil {
		t.Fatal(err)
	}
	if err := cfg.WriteFile(home.ConfigFile(), 0644); err != nil {
		t.Fatal(err)
	}

	environment.Settings.Home = home
	environment.Settings.Offline = false

	// 3. Call LoadRepository
	// PROBLEM: We cannot easily inject the Mock Server URL into LoadRepository -> fetchOnline
	// because it hardcodes github.NewClient which defaults to api.github.com.
	// This test expects to fail or we need a way to inject it.
	// For the purpose of "Writing a failing test", assuming we want to enable this injection.

	// To make this test actually meaningful as a "Failing Test" that we intend to pass:
	// We will attempt to run it. use 'force=true' to trigger fetchOnline.

	// We might expect a connection error or a "No search results" if it hits real github (and the token is invalid for real github)
	// But we WANT it to hit our mock.

	// Since we haven't refactored the code yet, this test serves to demonstrate the inability to mock.

	// Note: In a real "Failing Test" scenario for a feature, we usually implement the call but assert the result.
	// Here, we can't even tell it to use the mock.

	// Create client with BaseURL pointing to mock
	u, _ := url.Parse(server.URL + "/")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "mock-token"},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	client.BaseURL = u
	client.UploadURL = u

	t.Log("Attempting to load repository from mock...")
	r, err := LoadRepository(context.Background(), logrus.StandardLogger(), home, "softleader", true, client)
	if err != nil {
		// Ensure it's not a real network error, but we expect it to be successful if it hit our mock.
		t.Logf("Got error: %v", err)
	}

	// 4. Assert
	if r == nil || len(r.Repos) != 1 {
		t.Errorf("Expected 1 repo from mock, got %d. (This test is expected to fail before refactoring)", len(r.Repos))
	}
	if r != nil && len(r.Repos) > 0 {
		if r.Repos[0].Source != "github.com/softleader/slctl-mock-plugin" {
			t.Errorf("Expected mock plugin, got %s", r.Repos[0].Source)
		}
	}
}
