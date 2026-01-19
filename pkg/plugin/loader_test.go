package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/release"
	"github.com/spf13/cobra"
)

func TestTransformToCommand(t *testing.T) {
	p := &Plugin{
		Metadata: &Metadata{
			Name:  "test-cmd",
			Usage: "test usage",
		},
		Dir: "/tmp",
	}

	metadata := &release.Metadata{GitVersion: "1.0.0", GitCommit: "abcdef"}
	cmd := p.transformToCommand(metadata, func(cmd *cobra.Command, args []string) ([]string, error) {
		return args, nil
	})

	if cmd.Use != "test-cmd" {
		t.Errorf("expected test-cmd, got %s", cmd.Use)
	}
	if cmd.Short != "test usage" {
		t.Errorf("expected test usage, got %s", cmd.Short)
	}

	// Test RunE
	tempHome, _ := os.MkdirTemp("", "sl-home-cmd")
	defer os.RemoveAll(tempHome)

	oldHome := environment.Settings.Home
	environment.Settings.Home = paths.Home(tempHome)
	defer func() { environment.Settings.Home = oldHome }()

	os.MkdirAll(environment.Settings.Home.Config(), 0755)
	os.WriteFile(environment.Settings.Home.ConfigFile(), []byte("token: secret"), 0644)

	p.Metadata.Exec.Command = "echo"
	err := cmd.RunE(cmd, []string{"hello"})
	if err != nil {
		t.Errorf("RunE failed: %v", err)
	}

	// Test RunE error - SetupEnv error (missing config)
	os.RemoveAll(tempHome)
	err = cmd.RunE(cmd, []string{"hello"})
	if err == nil {
		t.Error("expected error when config is missing")
	}
}

func TestLoadDir(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "sl-plugin-loader")
	defer os.RemoveAll(tempDir)

	metadataContent := `name: test-plugin
version: 1.0.0
usage: a test plugin
exec:
  command: echo hello
`
	os.WriteFile(filepath.Join(tempDir, MetadataFileName), []byte(metadataContent), 0644)
	os.WriteFile(filepath.Join(tempDir, SourceFileName), []byte("github.com/softleader/test-plugin"), 0644)

	p, err := LoadDir(tempDir)
	if err != nil {
		t.Fatalf("LoadDir failed: %v", err)
	}

	if p.Metadata.Name != "test-plugin" {
		t.Errorf("expected test-plugin, got %s", p.Metadata.Name)
	}
	if p.Source != "github.com/softleader/test-plugin" {
		t.Errorf("expected github.com/softleader/test-plugin, got %s", p.Source)
	}

	// Source file missing
	os.Remove(filepath.Join(tempDir, SourceFileName))
	p, err = LoadDir(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	if p.Source != "" {
		t.Errorf("expected empty source, got %s", p.Source)
	}
}

func TestExitError(t *testing.T) {
	err := ExitError{
		error:      fmt.Errorf("some error"),
		ExitStatus: 1,
	}
	if err.Error() != "some error" {
		t.Errorf("expected some error, got %s", err.Error())
	}
	if err.ExitStatus != 1 {
		t.Errorf("expected 1, got %d", err.ExitStatus)
	}
}

func TestLoadPaths(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "sl-paths")
	defer os.RemoveAll(tempDir)

	p1Dir := filepath.Join(tempDir, "p1")
	os.MkdirAll(p1Dir, 0755)
	os.WriteFile(filepath.Join(p1Dir, MetadataFileName), []byte("name: p1"), 0644)

	// Test with multiple paths
	pathsStr := tempDir + string(os.PathListSeparator) + "/non/existent"
	plugins, err := LoadPaths(pathsStr)
	if err != nil {
		t.Fatal(err)
	}
	if len(plugins) != 1 {
		t.Errorf("expected 1 plugin, got %d", len(plugins))
	}
}

func TestLoadAll(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "sl-plugin-all")
	defer os.RemoveAll(tempDir)

	p1Dir := filepath.Join(tempDir, "p1")
	os.MkdirAll(p1Dir, 0755)
	os.WriteFile(filepath.Join(p1Dir, MetadataFileName), []byte("name: p1"), 0644)

	p2Dir := filepath.Join(tempDir, "p2")
	os.MkdirAll(p2Dir, 0755)
	os.WriteFile(filepath.Join(p2Dir, MetadataFileName), []byte("name: p2"), 0644)

	plugins, err := LoadAll(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(plugins) != 2 {
		t.Errorf("expected 2 plugins, got %d", len(plugins))
	}
}

func TestProcessFlags(t *testing.T) {
	args := []string{"--home", "/tmp", "local-arg", "-v"}
	global, local := processFlags(args)

	if len(global) != 2 {
		t.Errorf("expected 2 global flags, got %d", len(global))
	}
	if len(local) != 2 {
		t.Errorf("expected 2 local flags, got %d", len(local))
	}
}

func TestLoadPluginCommands(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "sl-plugins-dir")
	defer os.RemoveAll(tempDir)

	p1Dir := filepath.Join(tempDir, "p1")
	os.MkdirAll(p1Dir, 0755)
	os.WriteFile(filepath.Join(p1Dir, MetadataFileName), []byte("name: p1\nusage: usage1"), 0644)

	// Set environment for PluginDirs
	oldHome := environment.Settings.Home
	environment.Settings.Home = paths.Home(tempDir)
	defer func() { environment.Settings.Home = oldHome }()

	os.Setenv("SL_PLUGIN", tempDir)
	defer os.Unsetenv("SL_PLUGIN")

	metadata := &release.Metadata{GitVersion: "1.0.0"}
	cmds, err := LoadPluginCommands(metadata)
	if err != nil {
		t.Fatal(err)
	}
	if len(cmds) != 1 {
		t.Errorf("expected 1 command, got %d", len(cmds))
	}
}

func TestLoadPluginCommands_Off(t *testing.T) {
	os.Setenv("SL_PLUGINS_OFF", "true")
	defer os.Unsetenv("SL_PLUGINS_OFF")

	metadata := &release.Metadata{GitVersion: "1.0.0"}
	cmds, err := LoadPluginCommands(metadata)
	if err != nil {
		t.Fatal(err)
	}
	if len(cmds) != 0 {
		t.Error("expected 0 commands when SL_PLUGINS_OFF is true")
	}
}
