package token

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-github/v69/github"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	gh "github.com/softleader/slctl/pkg/github"
)

func TestConfirmUnit(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := github.NewClient(nil)
	u, _ := url.Parse(server.URL + "/")
	client.BaseURL = u
	client.UploadURL = u

	org := "test-org"

	mux.HandleFunc("/user/memberships/orgs/"+org, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"state": "active"}`)
	})
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"login": "test-user", "name": "Test User"}`)
	})

	name, err := Confirm(context.Background(), client, org, logrus.StandardLogger())
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if name != "Test User" {
		t.Errorf("expected Test User, got %s", name)
	}
}

func TestConfirmUnit_Inactive(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := github.NewClient(nil)
	u, _ := url.Parse(server.URL + "/")
	client.BaseURL = u
	client.UploadURL = u

	org := "test-org"

	mux.HandleFunc("/user/memberships/orgs/"+org, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"state": "pending"}`)
	})

	_, err := Confirm(context.Background(), client, org, logrus.StandardLogger())
	if err == nil {
		t.Errorf("expected error for inactive membership")
	}
}

// need read:org & read:user permission
func TestConfirm(t *testing.T) {
	environment.Settings.Verbose = true
	token, found := os.LookupEnv("GITHUB_TOKEN_TEST")
	if !found {
		t.Skipf("provide $GITHUB_TOKEN_TEST to run the test")
	}
	var name string
	client, err := gh.NewTokenClient(context.Background(), token)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if name, err = Confirm(context.Background(), client, "text", logrus.StandardLogger()); err != nil {
		t.Error(err)
	}
	if name == "" {
		t.Errorf("name should not be empty")
	}
	fmt.Printf("Hello, %s!\n", name)
}

func TestRemove(t *testing.T) {
	drop := github.ScopeGist
	a := []github.Scope{github.ScopeReadOrg, drop, github.ScopeRepo}
	a = remove(a, drop)

	for _, v := range a {
		if v == drop {
			t.Errorf("expected %q not contains %q", a, drop)
		}
	}

	// Not found
	a = remove(a, github.ScopeAdminOrg)
	if len(a) != 2 {
		t.Errorf("expected length 2, got %d", len(a))
	}
}
