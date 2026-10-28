package member

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v93/github"
)

func TestPublicizeOrganization(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client, clientErr := github.NewClient(github.WithURLs(github.Ptr(server.URL+"/"), nil))
	if clientErr != nil {
		t.Fatal(clientErr)
	}

	org := "test-org"
	username := "test-user"

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fmt.Sprintf(`{"login": "%s"}`, username))
	})
	mux.HandleFunc("/orgs/"+org+"/public_members/"+username, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := PublicizeOrganization(context.Background(), client, org)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestPublicizeOrganization_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client, clientErr := github.NewClient(github.WithURLs(github.Ptr(server.URL+"/"), nil))
	if clientErr != nil {
		t.Fatal(clientErr)
	}

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	err := PublicizeOrganization(context.Background(), client, "org")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
