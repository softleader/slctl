package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	ghapi "github.com/google/go-github/v93/github"
)

func TestNewTokenClientAuthenticatesRequests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/user" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if got := r.Header.Get("Authorization"); got != "Bearer dummy-token" {
			t.Errorf("Authorization = %q, want Bearer dummy-token", got)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"login":"test-user","id":123}`)
	}))
	t.Cleanup(server.Close)

	ctx := context.Background()
	client, err := NewTokenClient(ctx, "dummy-token")
	if err != nil {
		t.Fatal(err)
	}
	client, err = client.Clone(ghapi.WithURLs(ghapi.Ptr(server.URL+"/"), nil))
	if err != nil {
		t.Fatal(err)
	}
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if user.GetLogin() != "test-user" || user.GetID() != 123 {
		t.Fatalf("unexpected authenticated user: %+v", user)
	}
}
