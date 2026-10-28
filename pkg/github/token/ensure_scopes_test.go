package token

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v93/github"
	"github.com/sirupsen/logrus"
)

func TestEnsureScopes(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client, clientErr := github.NewClient(github.WithURLs(github.Ptr(server.URL+"/"), nil))
	if clientErr != nil {
		t.Fatal(clientErr)
	}

	// Case 1: Success
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-OAuth-Scopes", "repo, user, read:org")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"login": "test-user"}`)
	})

	err := EnsureScopes(context.Background(), client, logrus.StandardLogger(), []github.Scope{github.ScopeUser})
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	// Case 2: Missing Scope
	// We need to re-configure the handler or use a new server/mux for clean state,
	// or use a dynamic handler based on request count/header?
	// Let's just create a new setup or run subtests.
}

func TestEnsureScopes_Missing(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client, clientErr := github.NewClient(github.WithURLs(github.Ptr(server.URL+"/"), nil))
	if clientErr != nil {
		t.Fatal(clientErr)
	}

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-OAuth-Scopes", "repo") // Missing 'user'
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"login": "test-user"}`)
	})

	err := EnsureScopes(context.Background(), client, logrus.StandardLogger(), []github.Scope{github.ScopeUser})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
