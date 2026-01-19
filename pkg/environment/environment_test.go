package environment

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
)

func TestAddFlags(t *testing.T) {
	s := &settings{}
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	if err := s.AddFlags(fs); err != nil {
		t.Fatal(err)
	}

	if !fs.HasFlags() {
		t.Error("expected flags to be added")
	}

	if h, _ := fs.GetString("home"); h == "" {
		t.Error("expected default home to be set")
	}
}

func TestDefaultHome(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sl-test-default-home")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	oh := filepath.Join(tempDir, oldHome)
	if err := os.MkdirAll(oh, 0755); err != nil {
		t.Fatal(err)
	}

	h := DefaultHome(tempDir)
	expectedNewHome := filepath.Join(tempDir, home)
	if h != expectedNewHome {
		t.Errorf("expected %s, got %s", expectedNewHome, h)
	}

	if _, err := os.Stat(oh); !os.IsNotExist(err) {
		t.Error("expected old home to be moved")
	}
	if _, err := os.Stat(expectedNewHome); err != nil {
		t.Error("expected new home to exist")
	}
}

func TestIsGlobalFlag(t *testing.T) {
	tests := []struct {
		flag string
		want bool
	}{
		{"--home", true},
		{"--offline", true},
		{"--verbose", true},
		{"-v", true},
		{"--unknown", false},
	}

	for _, tt := range tests {
		if got := IsGlobalFlag(tt.flag); got != tt.want {
			t.Errorf("IsGlobalFlag(%q) = %v, want %v", tt.flag, got, tt.want)
		}
	}
}

func TestPluginDirs(t *testing.T) {
	s := settings{Home: "/tmp/sl"}

	// Default
	os.Unsetenv("SL_PLUGIN")
	if got := s.PluginDirs(); got != s.Home.Plugins() {
		t.Errorf("expected %s, got %s", s.Home.Plugins(), got)
	}

	// From Env
	expected := "/custom/plugin/dir"
	os.Setenv("SL_PLUGIN", expected)
	defer os.Unsetenv("SL_PLUGIN")
	if got := s.PluginDirs(); got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestExpandEnvToFlags(t *testing.T) {
	s := &settings{}
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	s.AddFlags(fs)

	os.Setenv("SL_HOME", "/env/home")
	os.Setenv("SL_OFFLINE", "true")
	os.Setenv("SL_VERBOSE", "true")
	defer func() {
		os.Unsetenv("SL_HOME")
		os.Unsetenv("SL_OFFLINE")
		os.Unsetenv("SL_VERBOSE")
	}()

	s.ExpandEnvToFlags(fs)

	if h, _ := fs.GetString("home"); h != "/env/home" {
		t.Errorf("expected home to be /env/home, got %s", h)
	}
	if o, _ := fs.GetBool("offline"); !o {
		t.Error("expected offline to be true")
	}
	if v, _ := fs.GetBool("verbose"); !v {
		t.Error("expected verbose to be true")
	}
}
