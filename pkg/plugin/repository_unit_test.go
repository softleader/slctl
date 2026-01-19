package plugin

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-github/v69/github"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
)

func TestLoadRepository_Unit(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := github.NewClient(nil)
	u, _ := url.Parse(server.URL + "/")
	client.BaseURL = u
	client.UploadURL = u

	org := "softleader"
	mux.HandleFunc("/search/repositories", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"total_count": 1, "items": [{"full_name": "softleader/p1", "description": "d1"}]}`)
	})

	tempHome, _ := os.MkdirTemp("", "sl-home-repo-unit")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)
	os.MkdirAll(hh.CachePlugins(), 0755)

	r, err := LoadRepository(context.Background(), logrus.New(), hh, org, true, client)
	if err != nil {
		t.Fatalf("LoadRepository failed: %v", err)
	}

	if len(r.Repos) != 1 || r.Repos[0].Source != "github.com/softleader/p1" {
		t.Errorf("expected p1, got %v", r.Repos)
	}

	// Test load from local
	r2, err := LoadRepository(context.Background(), logrus.New(), hh, org, false, client)
	if err != nil {
		t.Fatal(err)
	}
	if len(r2.Repos) != 1 {
		t.Error("expected 1 repo from cache")
	}
}

func TestRepo_Contains(t *testing.T) {
	r := Repo{Source: "github.com/softleader/slctl", Description: "softleader slctl"}
	if !r.Contains("slctl") {
		t.Error("expected to contain slctl")
	}
	if !r.Contains("softleader") {
		t.Error("expected to contain softleader")
	}
	if r.Contains("not-found") {
		t.Error("expected NOT to contain not-found")
	}
}

func TestExpired(t *testing.T) {
	if !expired(nil) {
		t.Error("expected nil repo to be expired")
	}

	r := &Repository{Expires: time.Now().Add(-1 * time.Hour)}
	if !expired(r) {
		t.Error("expected past expiration to be expired")
	}

	r.Expires = time.Now().Add(1 * time.Hour)
	if expired(r) {
		t.Error("expected future expiration NOT to be expired")
	}
}

func TestFetchOnline_OfflineError(t *testing.T) {
	environment.Settings.Offline = true
	defer func() { environment.Settings.Offline = false }()

	log := logrus.New()
	hh := paths.Home("/tmp")
	_, err := fetchOnline(context.Background(), log, hh, "org", nil)
	if err == nil || !strings.Contains(err.Error(), "offline mode") {
		t.Errorf("expected offline mode error, got %v", err)
	}
}

func TestFetchOnline_NoClient_ConfigError(t *testing.T) {
	log := logrus.New()
	hh := paths.Home("/non/existent/home")
	_, err := fetchOnline(context.Background(), log, hh, "org", nil)
	if err == nil {
		t.Error("expected error for missing config")
	}
}

func TestFetchOnline_SearchError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := github.NewClient(nil)
	u, _ := url.Parse(server.URL + "/")
	client.BaseURL = u
	client.UploadURL = u

	mux.HandleFunc("/search/repositories", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	log := logrus.New()
	hh := paths.Home("/tmp")
	r, err := fetchOnline(context.Background(), log, hh, "org", client)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Repos) != 0 {
		t.Error("expected 0 repos on error")
	}
}
