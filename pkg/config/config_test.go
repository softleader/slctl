package config

import (
	"os"
	"strings"
	"testing"

	"github.com/softleader/slctl/pkg/paths"
)

func TestWriteFile(t *testing.T) {
	cf := NewConfFile()
	cf.Token = "this.is.a.fake.token"

	repoFile, err := os.CreateTemp("", "sl-config")
	if err != nil {
		t.Errorf("failed to create test-file (%v)", err)
	}
	defer os.Remove(repoFile.Name())

	if err := cf.WriteFile(repoFile.Name(), 0644); err != nil {
		t.Errorf("failed to write file (%v)", err)
	}

	loaded, err := LoadConfFile(repoFile.Name())
	if err != nil {
		t.Errorf("failed to load file (%v)", err)
	}
	if loaded.Token != cf.Token {
		t.Errorf("expected token to be %s", cf.Token)
	}
}

func TestLoadConfFile_MalformedYAML(t *testing.T) {
	malformed := "this: is: not: valid: yaml"
	repoFile, err := os.CreateTemp("", "sl-config-malformed")
	if err != nil {
		t.Errorf("failed to create test-file (%v)", err)
	}
	defer os.Remove(repoFile.Name())

	if err := os.WriteFile(repoFile.Name(), []byte(malformed), 0644); err != nil {
		t.Errorf("failed to write file (%v)", err)
	}

	_, err = LoadConfFile(repoFile.Name())
	if err == nil {
		t.Errorf("expected err to be non-nil when loading malformed yaml")
	}
}

func TestLoadConfFile_TokenNotExist(t *testing.T) {
	noToken := "cleanup: 2026-01-19T00:00:00Z"
	repoFile, err := os.CreateTemp("", "sl-config-no-token")
	if err != nil {
		t.Errorf("failed to create test-file (%v)", err)
	}
	defer os.Remove(repoFile.Name())

	if err := os.WriteFile(repoFile.Name(), []byte(noToken), 0644); err != nil {
		t.Errorf("failed to write file (%v)", err)
	}

	conf, err := LoadConfFile(repoFile.Name())
	if err != ErrTokenNotExist {
		t.Errorf("expected ErrTokenNotExist, got %v", err)
	}
	if conf == nil {
		t.Errorf("expected conf to be non-nil even if token does not exist")
	}
}

func TestRefresh(t *testing.T) {
	tempHome, err := os.MkdirTemp("", "sl-home")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempHome)

	home := paths.Home(tempHome)
	if err := os.MkdirAll(home.Config(), 0755); err != nil {
		t.Fatal(err)
	}

	token := "new-token"
	// Refresh fails if file doesn't exist
	if err := Refresh(home, token, nil); err == nil {
		t.Error("expected error when Refreshing non-existent config file")
	}

	// Create file with no token
	if err := os.WriteFile(home.ConfigFile(), []byte("cleanup: 2026-01-19T00:00:00Z"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := Refresh(home, token, nil); err != nil {
		t.Errorf("Refresh failed: %v", err)
	}

	loaded, err := LoadConfFile(home.ConfigFile())
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Token != token {
		t.Errorf("expected token to be %s, got %s", token, loaded.Token)
	}
}

func TestConfigNotExists(t *testing.T) {
	_, err := LoadConfFile("/this/path/does/not/exist.yaml")
	if err == nil {
		t.Errorf("expected err to be non-nil when path does not exist")
	} else if !strings.Contains(err.Error(), "You might need to run `slctl init`") {
		t.Errorf("expected prompt to run `slctl init` when config file does not exist")
	}
}
