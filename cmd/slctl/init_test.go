package main

import (
	"errors"
	"os"
	"testing"

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
