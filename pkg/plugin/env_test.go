package plugin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/release"
)

func TestSetupEnv(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "sl-home-env")
	defer os.RemoveAll(tempHome)

	// Set global settings
	oldHome := environment.Settings.Home
	environment.Settings.Home = paths.Home(tempHome)
	defer func() { environment.Settings.Home = oldHome }()

	os.MkdirAll(environment.Settings.Home.Config(), 0755)
	os.WriteFile(environment.Settings.Home.ConfigFile(), []byte("token: secret-token"), 0644)

	p := &Plugin{
		Metadata: &Metadata{Name: "test-plugin"},
		Dir:      "/plugins/test-plugin",
		Mount:    filepath.Join(tempHome, "mounts", "test-plugin"),
	}

	metadata := &release.Metadata{GitVersion: "1.0.0", GitCommit: "abcdef"}
	err := p.SetupEnv(metadata)
	if err != nil {
		t.Fatalf("SetupEnv failed: %v", err)
	}

	if os.Getenv("SL_TOKEN") != "secret-token" {
		t.Errorf("expected SL_TOKEN to be secret-token, got %s", os.Getenv("SL_TOKEN"))
	}
	if os.Getenv("SL_PLUGIN_NAME") != "test-plugin" {
		t.Errorf("expected SL_PLUGIN_NAME to be test-plugin, got %s", os.Getenv("SL_PLUGIN_NAME"))
	}
}

func TestSetupEnv_ConfigError(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "sl-home-env-err")
	defer os.RemoveAll(tempHome)

	oldHome := environment.Settings.Home
	environment.Settings.Home = paths.Home(tempHome)
	defer func() { environment.Settings.Home = oldHome }()

	os.MkdirAll(environment.Settings.Home.Config(), 0755)
	// Invalid YAML
	os.WriteFile(environment.Settings.Home.ConfigFile(), []byte("invalid: yaml: :"), 0644)

	p := &Plugin{
		Metadata: &Metadata{Name: "test-plugin"},
		Dir:      "/tmp",
		Mount:    filepath.Join(tempHome, "mounts", "test-plugin"),
	}

	metadata := &release.Metadata{GitVersion: "1.0.0"}
	err := p.SetupEnv(metadata)
	if err == nil {
		t.Error("expected error for invalid config yaml")
	}
}

func TestEnvsMap(t *testing.T) {
	m := envsMap("p", "d", "m", "v", "t")
	if m["SL_PLUGIN_NAME"] != "p" {
		t.Errorf("expected p, got %s", m["SL_PLUGIN_NAME"])
	}
}
