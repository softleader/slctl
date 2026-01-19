package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequestDeviceCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"device_code": "dc123", "user_code": "uc456", "verification_uri": "https://gh.com/verify", "expires_in": 900, "interval": 5}`)
	}))
	defer server.Close()

	oldURL := deviceCodeURL
	oldClient := httpClient
	deviceCodeURL = server.URL
	httpClient = server.Client()
	defer func() {
		deviceCodeURL = oldURL
		httpClient = oldClient
	}()

	resp, err := RequestDeviceCode(context.Background(), "", Scopes)
	if err != nil {
		t.Fatalf("RequestDeviceCode failed: %v", err)
	}
	if resp.DeviceCode != "dc123" {
		t.Errorf("expected dc123, got %s", resp.DeviceCode)
	}
}

func TestPollAccessToken(t *testing.T) {
	count := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if count == 0 {
			fmt.Fprint(w, `{"error": "authorization_pending"}`)
		} else {
			fmt.Fprint(w, `{"access_token": "at123", "token_type": "bearer"}`)
		}
		count++
	}))
	defer server.Close()

	oldURL := accessTokenURL
	oldClient := httpClient
	accessTokenURL = server.URL
	httpClient = server.Client()
	defer func() {
		accessTokenURL = oldURL
		httpClient = oldClient
	}()

	// Use a small interval for testing
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token, err := PollAccessToken(ctx, "", "dc123", 1)
	if err != nil {
		t.Fatalf("PollAccessToken failed: %v", err)
	}
	if token != "at123" {
		t.Errorf("expected at123, got %s", token)
	}
}

func TestPollAccessToken_Errors(t *testing.T) {
	tests := []struct {
		githubError string
		expectedErr string
	}{
		{"expired_token", "the device code has expired"},
		{"access_denied", "access denied by user"},
		{"other_error", "oauth error: other_error"},
	}

	for _, tt := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"error": "%s"}`, tt.githubError)
		}))

		oldURL := accessTokenURL
		oldClient := httpClient
		accessTokenURL = server.URL
		httpClient = server.Client()

		_, err := PollAccessToken(context.Background(), "", "dc123", 1)
		if err == nil || err.Error() != tt.expectedErr {
			t.Errorf("expected error %q, got %v", tt.expectedErr, err)
		}

		server.Close()
		accessTokenURL = oldURL
		httpClient = oldClient
	}
}
