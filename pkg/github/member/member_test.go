package member

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/v69/github"
)

func TestPublicizeOrganization(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := github.NewClient(nil)
	u, _ := url.Parse(server.URL + "/")
	client.BaseURL = u
	client.UploadURL = u

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

	client := github.NewClient(nil)
	u, _ := url.Parse(server.URL + "/")
	client.BaseURL = u
	client.UploadURL = u

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	err := PublicizeOrganization(context.Background(), client, "org")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
