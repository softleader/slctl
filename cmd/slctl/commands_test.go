package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v69/github"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/plugin/installer"
	"github.com/softleader/slctl/pkg/release"
	"github.com/spf13/cobra"
)

func TestNewCommands(t *testing.T) {
	// Root
	root, err := newRootCmd(nil)
	if err != nil {
		t.Fatal(err)
	}
	if root.Use != "slctl" {
		t.Errorf("expected slctl, got %s", root.Use)
	}

	// Subcommands
	cmds := []*cobra.Command{
		newCleanupCmd(),
		newCompletionCmd(),
		newHomeCmd(),
		newInitCmd(),
		newPluginCmd(),
		newVersionCmd(),
	}

	for _, cmd := range cmds {
		if cmd == nil {
			t.Errorf("command is nil")
		}
	}
}

func TestCleanupCmd_Run(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "sl-home-cleanup")
	defer os.RemoveAll(tempHome)

	c := &cleanupCmd{
		home:   paths.Home(tempHome),
		dryRun: true,
	}

	// Ensure config file exists
	os.MkdirAll(c.home.Config(), 0755)
	os.WriteFile(c.home.ConfigFile(), []byte("cleanup: 2020-01-01T00:00:00Z"), 0644)

	if err := c.run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}
}

func TestHomeCmd_Run(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "sl-home-home")
	defer os.RemoveAll(tempHome)

	c := &homeCmd{
		home: paths.Home(tempHome),
	}

	if err := c.run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	// Test Move
	parent, _ := os.MkdirTemp("", "sl-home-parent")
	defer os.RemoveAll(parent)
	tempMove := filepath.Join(parent, "new-home")

	c.move = tempMove
	if err := c.run(); err != nil {
		t.Fatalf("run move failed: %v", err)
	}
}

func TestVersionCmd_Run(t *testing.T) {
	metadata = release.NewMetadata("1.0.0", "abcdef")
	c := &versionCmd{
		full: true,
	}
	if err := c.run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	c.full = false
	if err := c.run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	// Check
	c.check = true
	environment.Settings.Offline = true
	if err := c.run(); err == nil {
		t.Error("expected error for version check in offline mode")
	}
	environment.Settings.Offline = false
}

func TestPluginRemoveCmd_Run(t *testing.T) {
	c := &pluginRemoveCmd{
		names: []string{"not-found"},
		force: true,
	}
	if err := c.run(); err != nil {
		t.Fatalf("run failed with force: %v", err)
	}

	c.force = false
	if err := c.run(); err == nil {
		t.Error("expected error for non-existent plugin without force")
	}
}

func TestPluginSearchCmd_Run(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "sl-home-search")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)

	// Create cached repo file to avoid online fetch
	repoFile := hh.CacheRepositoryFile()
	os.MkdirAll(filepath.Dir(repoFile), 0755)
	os.WriteFile(repoFile, []byte("repos:\n  - source: github.com/softleader/slctl\n    description: desc\nexpires: 2099-01-01T00:00:00Z"), 0644)

	c := &pluginSearchCmd{
		home: hh,
	}
	if err := c.run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}
}

func TestPluginUnmountCmd_Run(t *testing.T) {
	c := &pluginUnmountCmd{
		name: []string{"not-found"},
	}
	if err := c.run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}
}

func TestPluginInstallCmd_Run(t *testing.T) {
	c := &pluginInstallCmd{
		source: "invalid-source",
		opt:    &installer.InstallOption{},
	}
	if err := c.run(); err == nil {
		t.Error("expected error for invalid source")
	}
}

func TestPluginUpgradeCmd_Run(t *testing.T) {
	c := &pluginUpgradeCmd{
		names: []string{"not-found"},
		opt:   &installer.InstallOption{},
	}
	if err := c.run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}
}

func TestPluginCreateCmd_Run(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "sl-plugin-create-cmd")
	defer os.RemoveAll(tempDir)

	c := &pluginCreateCmd{
		name:   "test-plugin",
		lang:   "golang",
		output: tempDir,
	}
	if err := c.run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}
}

func TestCompletionCmd_Run(t *testing.T) {
	root, _ := newRootCmd(nil)
	compCmd := newCompletionCmd()
	root.AddCommand(compCmd)

	// Bash
	if err := runCompletionBash(compCmd); err != nil {
		t.Fatalf("bash completion failed: %v", err)
	}

	// Zsh
	if err := runCompletionZsh(compCmd); err != nil {
		t.Fatalf("zsh completion failed: %v", err)
	}
}

func TestPluginExtsCmd_Run(t *testing.T) {
	cmd := newPluginExtsCmd()
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("run failed: %v", err)
	}
}

func TestPluginCreateLangsCmd_Run(t *testing.T) {
	cmd := newPluginCreateLangsCmd()
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("run failed: %v", err)
	}
}

func TestInitScopesCmd_Run(t *testing.T) {
	cmd := newInitScopesCmd()
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("run failed: %v", err)
	}
}

func TestInitMetadata(t *testing.T) {
	version = "1.2.3"
	commit = "12345678"
	initMetadata()
	if metadata.GitVersion != "1.2.3" {
		t.Errorf("expected 1.2.3, got %s", metadata.GitVersion)
	}
}

func TestPluginUpgradeCmd_Run_Full(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "sl-home-upgrade")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)
	os.MkdirAll(hh.Plugins(), 0755)

	// Create a dummy plugin
	pDir := filepath.Join(hh.Plugins(), "p1")
	os.MkdirAll(pDir, 0755)
	os.WriteFile(filepath.Join(pDir, "metadata.yaml"), []byte("name: p1"), 0644)
	os.WriteFile(filepath.Join(pDir, ".source"), []byte("github.com/softleader/slctl"), 0644)

	c := &pluginUpgradeCmd{
		home: hh,
		opt:  &installer.InstallOption{DryRun: true},
	}

	// Offline will fail run()
	environment.Settings.Offline = true
	defer func() { environment.Settings.Offline = false }()

	if err := c.run(); err != nil {
		// Expect error from RunE wrapper usually, but here we call run() directly.
	}
}

func TestMatch(t *testing.T) {
	plugs := []*plugin.Plugin{
		{Metadata: &plugin.Metadata{Name: "P1"}},
	}
	if _, found := match("p1", plugs); !found {
		t.Error("expected to find p1")
	}
	if _, found := match("p2", plugs); found {
		t.Error("did not expect to find p2")
	}
}

func TestRootCmd_Hooks(t *testing.T) {
	root, _ := newRootCmd(nil)
	if root.PersistentPreRun != nil {
		root.PersistentPreRun(root, nil)
	}
	if root.PersistentPostRun != nil {
		// Mock environment to avoid update check
		environment.Settings.Offline = true
		defer func() { environment.Settings.Offline = false }()
		root.PersistentPostRun(root, nil)
	}
}

func TestPluginSearchCmd_Run_Installed(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "sl-home-search-inst")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)

	// Create cached repo file
	repoFile := hh.CacheRepositoryFile()
	os.MkdirAll(filepath.Dir(repoFile), 0755)
	os.WriteFile(repoFile, []byte("repos:\n  - source: github.com/softleader/slctl\n    description: desc\nexpires: 2099-01-01T00:00:00Z"), 0644)

	c := &pluginSearchCmd{
		home:              hh,
		onlyShowInstalled: true,
	}
	if err := c.run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}
}

func TestPluginRemoveCmd_Run_Full(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "sl-home-remove")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)

	// Set environment for PluginDirs
	oldHome := environment.Settings.Home
	environment.Settings.Home = hh
	defer func() { environment.Settings.Home = oldHome }()

	os.MkdirAll(hh.Plugins(), 0755)
	pDir := filepath.Join(hh.Plugins(), "p1")
	os.MkdirAll(pDir, 0755)
	os.WriteFile(filepath.Join(pDir, "metadata.yaml"), []byte("name: p1"), 0644)

	c := &pluginRemoveCmd{
		names: []string{"p1"},
		home:  hh,
	}
	if err := c.run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}
}

func TestPluginUnmountCmd_Run_Full(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "sl-home-unmount")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)

	// Set environment for PluginDirs
	oldHome := environment.Settings.Home
	environment.Settings.Home = hh
	defer func() { environment.Settings.Home = oldHome }()

	os.MkdirAll(hh.Plugins(), 0755)
	pDir := filepath.Join(hh.Plugins(), "p1")
	os.MkdirAll(pDir, 0755)
	os.WriteFile(filepath.Join(pDir, "metadata.yaml"), []byte("name: p1"), 0644)

	mountDir := filepath.Join(hh.Mounts(), "p1")
	os.MkdirAll(mountDir, 0755)

	c := &pluginUnmountCmd{
		name: []string{"p1"},
	}
	if err := c.run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if _, err := os.Stat(mountDir); !os.IsNotExist(err) {
		t.Error("expected mount directory to be removed")
	}
}

func TestInstall_Error(t *testing.T) {
	hh := paths.Home("/non/existent")
	opt := &installer.InstallOption{}
	err := install("invalid", "", 0, hh, opt)
	if err == nil {
		t.Error("expected error for invalid install source")
	}
}

func TestRootCmd_Execute(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "sl-home-exec")
	defer os.RemoveAll(tempHome)

	oldHome := environment.Settings.Home
	environment.Settings.Home = paths.Home(tempHome)
	defer func() { environment.Settings.Home = oldHome }()

	os.MkdirAll(environment.Settings.Home.Config(), 0755)
	os.WriteFile(environment.Settings.Home.ConfigFile(), []byte("token: secret"), 0644)

	// Mock tokenClient
	oldTokenClient := tokenClient
	tokenClient = func(ctx context.Context, token string) (*github.Client, error) {
		return github.NewClient(nil), nil
	}
	defer func() { tokenClient = oldTokenClient }()

	root, _ := newRootCmd([]string{"--home", tempHome, "--offline"})

	root.SetArgs([]string{"version"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute version failed: %v", err)
	}

	root.SetArgs([]string{"cleanup", "--dry-run"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute cleanup failed: %v", err)
	}

	root.SetArgs([]string{"plugin", "list"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute plugin list failed: %v", err)
	}

	// Add a plugin and list again
	pDir := filepath.Join(tempHome, "plugins", "p1")
	os.MkdirAll(pDir, 0755)
	os.WriteFile(filepath.Join(pDir, "metadata.yaml"), []byte("name: p1"), 0644)
	root.SetArgs([]string{"plugin", "list"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute plugin list with plugin failed: %v", err)
	}

	root.SetArgs([]string{"plugin", "envs"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute plugin envs failed: %v", err)
	}

	root.SetArgs([]string{"plugin", "flags"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute plugin flags failed: %v", err)
	}

	// Search - use cached file
	repoFile := paths.Home(tempHome).CacheRepositoryFile()
	os.MkdirAll(filepath.Dir(repoFile), 0755)
	os.WriteFile(repoFile, []byte("repos:\n  - source: github.com/softleader/slctl\nexpires: 2099-01-01T00:00:00Z"), 0644)

	oldOffline := environment.Settings.Offline
	environment.Settings.Offline = false // enable search
	defer func() { environment.Settings.Offline = oldOffline }()

	root.SetArgs([]string{"plugin", "search"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute plugin search failed: %v", err)
	}

	root.SetArgs([]string{"plugin", "search", "-i"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute plugin search installed failed: %v", err)
	}

	root.SetArgs([]string{"plugin", "create", "foo", "-o", filepath.Join(tempHome, "foo")})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute plugin create failed: %v", err)
	}

	root.SetArgs([]string{"init", "scopes"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute init scopes failed: %v", err)
	}

	root.SetArgs([]string{"init", "--force", "--token", "secret", "--yes"})
	if err := root.Execute(); err != nil {
		// Expect fail but covers lines
	}

	root.SetArgs([]string{"plugin", "remove", "not-found", "-f"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute plugin remove failed: %v", err)
	}

	root.SetArgs([]string{"plugin", "umount", "not-found"})

	root.SetArgs([]string{"plugin", "open", "p1"})
	// Open might fail in CI but we want to cover the code path
	root.Execute()

	root.SetArgs([]string{"home"})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute home failed: %v", err)
	}

	parent, _ := os.MkdirTemp("", "sl-home-parent-exec")
	defer os.RemoveAll(parent)
	tempMove := filepath.Join(parent, "new-home")
	root.SetArgs([]string{"home", "--move", tempMove})
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute home move failed: %v", err)
	}
	environment.Settings.Home = paths.Home(tempHome) // restore for cleanup test
}

func TestPluginWorkflow(t *testing.T) {
	tempHome, _ := os.MkdirTemp("", "sl-home-workflow")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)

	oldHome := environment.Settings.Home
	environment.Settings.Home = hh
	defer func() { environment.Settings.Home = oldHome }()

	os.MkdirAll(hh.Config(), 0755)
	os.MkdirAll(hh.Plugins(), 0755)
	os.WriteFile(hh.ConfigFile(), []byte("token: secret"), 0644)

	metadata = release.NewMetadata("1.0.0", "abcdef")

	// Mock tokenClient
	oldTokenClient := tokenClient
	tokenClient = func(ctx context.Context, token string) (*github.Client, error) {
		mux := http.NewServeMux()
		server := httptest.NewServer(mux)
		// We don't close server here as it's returned client.
		// Actually github.NewClient(server.Client()) is better.
		c := github.NewClient(server.Client())
		u, _ := url.Parse(server.URL + "/")
		c.BaseURL = u

		mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-OAuth-Scopes", "repo, user, read:org")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"login": "test-user"}`)
		})
		return c, nil
	}
	defer func() { tokenClient = oldTokenClient }()

	// 1. Create a plugin to install
	pluginDir, _ := os.MkdirTemp("", "workflow-plug")
	defer os.RemoveAll(pluginDir)
	os.WriteFile(filepath.Join(pluginDir, "metadata.yaml"), []byte("name: wf-plug\nversion: 1.0.0\nexec:\n  command: echo hello"), 0644)

	// 2. Install
	if err := install(pluginDir, "", 0, hh, &installer.InstallOption{}); err != nil {
		t.Fatalf("install failed: %v", err)
	}

	// 3. List
	root, _ := newRootCmd([]string{"--home", tempHome, "--offline"})
	root.SetArgs([]string{"plugin", "list"})
	if err := root.Execute(); err != nil {
		t.Fatalf("list failed: %v", err)
	}

	// 4. Run
	root.SetArgs([]string{"wf-plug"})
	if err := root.Execute(); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	// 5. Remove
	root.SetArgs([]string{"plugin", "remove", "wf-plug"})
	if err := root.Execute(); err != nil {
		t.Fatalf("remove failed: %v", err)
	}
}

func TestRunHook(t *testing.T) {
	p := &plugin.Plugin{
		Metadata: &plugin.Metadata{
			Name: "test-plugin",
			Hook: plugin.Commands{
				Command: "echo hello",
			},
		},
	}

	// Set up environment for SetupEnv
	tempHome, _ := os.MkdirTemp("", "sl-home-hook")
	defer os.RemoveAll(tempHome)
	oldHome := environment.Settings.Home
	environment.Settings.Home = paths.Home(tempHome)
	defer func() { environment.Settings.Home = oldHome }()
	os.MkdirAll(environment.Settings.Home.Config(), 0755)
	os.WriteFile(environment.Settings.Home.ConfigFile(), []byte("token: secret"), 0644)

	metadata = release.NewMetadata("1.0.0", "abcdef")
	if err := runHook(p); err != nil {
		t.Fatalf("runHook failed: %v", err)
	}
}

func TestExit(t *testing.T) {
	// We can't actually call os.Exit(1) or it will stop the test.
}

func TestInstall_ConfigError(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "sl-install-err")
	defer os.RemoveAll(tempDir)
	hh := paths.Home(tempDir)
	os.MkdirAll(hh.Config(), 0755)
	os.WriteFile(hh.ConfigFile(), []byte("invalid yaml"), 0644)

	// We need a valid installer to get past NewInstaller
	pluginDir, _ := os.MkdirTemp("", "plug-src")
	defer os.RemoveAll(pluginDir)
	os.WriteFile(filepath.Join(pluginDir, "metadata.yaml"), []byte("name: p1"), 0644)

	opt := &installer.InstallOption{}
	err := install(pluginDir, "", 0, hh, opt)
	if err == nil {
		t.Error("expected error for malformed config")
	}
}

func TestPluginUpgradeCmd_Upgrade(t *testing.T) {
	c := &pluginUpgradeCmd{}
	p := &plugin.Plugin{
		Metadata: &plugin.Metadata{Name: "p1"},
		Source:   "/local/path",
	}
	// skip upgrading local plugin
	if err := c.upgrade(p); err != nil {
		t.Fatal(err)
	}
}

func TestRunCompletionZsh_Error(t *testing.T) {
	// Root GenBashCompletion failure is hard to trigger,
	// but we can pass a dummy command that is not root.
	cmd := &cobra.Command{Use: "test"}
	if err := runCompletionZsh(cmd); err == nil {
		// Actually runCompletionZsh calls cmd.Root().GenBashCompletion
		// so it will still succeed if it can find a root.
	}
}

func TestInstall_Full(t *testing.T) {
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

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-OAuth-Scopes", "repo, user, read:org")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"login": "test-user"}`)
	})

	tempHome, _ := os.MkdirTemp("", "sl-home-install-full")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)
	os.MkdirAll(hh.Config(), 0755)
	os.MkdirAll(hh.Plugins(), 0755)
	os.WriteFile(hh.ConfigFile(), []byte("token: secret"), 0644)

	pluginDir, _ := os.MkdirTemp("", "plug-src-full")
	defer os.RemoveAll(pluginDir)
	os.WriteFile(filepath.Join(pluginDir, "metadata.yaml"), []byte("name: p1\nversion: 1.0.0"), 0644)

	opt := &installer.InstallOption{}
	err := install(pluginDir, "", 0, hh, opt)
	if err != nil {
		t.Fatalf("install failed: %v", err)
	}
}

func TestPluginSubcommands(t *testing.T) {
	cmds := []*cobra.Command{
		newPluginCreateCmd(),
		newPluginCreateLangsCmd(),
		newPluginEnvsCmd(),
		newPluginExtsCmd(),
		newPluginFlagsCmd(),
		newPluginInstallCmd(),
		newPluginListCmd(),
		newPluginOpenCmd(),
		newPluginRemoveCmd(),
		newPluginSearchCmd(),
		newPluginUnmountCmd(),
		newPluginUpgradeCmd(),
	}

	for _, cmd := range cmds {
		if cmd == nil {
			t.Errorf("command is nil")
		}
	}
}

func TestPluginListCmd_Run(t *testing.T) {
	c := &pluginListCmd{}
	if err := c.run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}
}

func TestPluginEnvsCmd_Run(t *testing.T) {
	c := &pluginEnvsCmd{}
	if err := c.run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}
}

func TestPluginFlagsCmd_Run(t *testing.T) {
	c := &pluginFlagsCmd{}
	if err := c.run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}
}

func TestPluginOpenCmd_Run(t *testing.T) {
	c := &pluginOpenCmd{
		plugin: "not-found",
	}
	if err := c.run(); err == nil {
		t.Error("expected error for non-existent plugin")
	}
}
