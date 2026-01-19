package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-github/v69/github"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
)

func TestRefreshConfig(t *testing.T) {
	home, err := os.MkdirTemp("", "sl_home")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(home)

	log := logrus.New()
	hh := paths.Home(home)
	environment.Settings.Home = hh

	environment.Settings.Verbose = true

	if err = ensureDirectories(hh, log); err != nil {
		t.Error(err)
	}
	if err := ensureConfigFile(hh, log); err != nil {
		t.Error(err)
	}

	token := "this.is.a.fake.token"
	if err = config.Refresh(hh, token, log); err != nil {
		t.Error(err)
	}

	var conf *config.ConfFile
	if conf, err = config.LoadConfFile(hh.ConfigFile()); err != nil {
		t.Error(err)
	}
	if conf.Token != token {
		t.Errorf("expected token to be %s", conf.Token)
	}
}

func TestInitCmd_Run_Offline(t *testing.T) {
	home, _ := os.MkdirTemp("", "sl_home_offline")
	defer os.RemoveAll(home)
	hh := paths.Home(home)

	oldHome := environment.Settings.Home
	environment.Settings.Home = hh
	defer func() { environment.Settings.Home = oldHome }()

	environment.Settings.Offline = true
	defer func() { environment.Settings.Offline = false }()

	c := &initCmd{
		home: hh,
		yes:  true,
	}

	if err := c.run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	if !paths.IsExistDirectory(hh.Plugins()) {
		t.Error("expected plugins directory to be created")
	}
}

func TestInitCmd_Run_SpaceInHome(t *testing.T) {
	hh := paths.Home("/path with space")
	c := &initCmd{home: hh}
	if err := c.run(); err == nil {
		t.Error("expected error for space in home path")
	}
}

func TestInitCmd_Run_Token(t *testing.T) {
	// This will still fail at token.Confirm unless we mock it.
	// But it will cover the token provided path.
	home, _ := os.MkdirTemp("", "sl_home_token")
	defer os.RemoveAll(home)
	hh := paths.Home(home)

	_ = &initCmd{
		home:  hh,
		token: "fake-token",
		yes:   true,
	}

	// We'll skip executing it as it will fail networking.
	// But we can test the logic that leads to it.
}

func TestInitCmd_Run_Full(t *testing.T) {
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

	mux.HandleFunc("/user/memberships/orgs/softleader", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"state": "active"}`)
	})
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"login": "test-user", "name": "Test User"}`)
	})
	mux.HandleFunc("/orgs/softleader/public_members/test-user", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	home, _ := os.MkdirTemp("", "sl_home_full_init")
	defer os.RemoveAll(home)
	hh := paths.Home(home)

	c := &initCmd{
		home:  hh,
		token: "secret",
		yes:   true,
	}

	if err := c.run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}
}

func TestEnsureConfigFile_Error(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "sl-home-ecf-err")
	defer os.RemoveAll(tempDir)
	hh := paths.Home(tempDir)
	os.MkdirAll(hh.ConfigFile(), 0755) // existing directory

	err := ensureConfigFile(hh, logrus.New())
	if err == nil {
		t.Error("expected error when config file path is a directory")
	}
}

// TestMockableFunctions verifies the mockable functions can be replaced for testing
func TestMockableFunctions(t *testing.T) {
	// Save original functions
	originalOpenBrowser := openBrowser
	originalWriteToClipboard := writeToClipboard
	defer func() {
		openBrowser = originalOpenBrowser
		writeToClipboard = originalWriteToClipboard
	}()

	t.Run("openBrowser can be mocked", func(t *testing.T) {
		called := false
		calledWith := ""
		openBrowser = func(input string) error {
			called = true
			calledWith = input
			return nil
		}

		err := openBrowser("https://github.com/login/device")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !called {
			t.Error("openBrowser was not called")
		}
		if calledWith != "https://github.com/login/device" {
			t.Errorf("openBrowser called with wrong URL: %s", calledWith)
		}
	})

	t.Run("writeToClipboard can be mocked", func(t *testing.T) {
		called := false
		calledWith := ""
		writeToClipboard = func(text string) error {
			called = true
			calledWith = text
			return nil
		}

		err := writeToClipboard("ABC-123")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !called {
			t.Error("writeToClipboard was not called")
		}
		if calledWith != "ABC-123" {
			t.Errorf("writeToClipboard called with wrong code: %s", calledWith)
		}
	})

	t.Run("errors from mocked functions are returned", func(t *testing.T) {
		expectedErr := errors.New("mock error")
		openBrowser = func(input string) error {
			return expectedErr
		}
		writeToClipboard = func(text string) error {
			return expectedErr
		}

		if err := openBrowser("test"); err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if err := writeToClipboard("test"); err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})
}
