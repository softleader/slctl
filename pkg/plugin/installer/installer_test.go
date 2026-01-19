package installer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/plugin"
)

func TestNewInstaller(t *testing.T) {
	log := logrus.New()
	tempHome, _ := os.MkdirTemp("", "sl-home-new-installer")
	defer os.RemoveAll(tempHome)
	home := paths.Home(tempHome)
	os.MkdirAll(home.Config(), 0755)
	os.WriteFile(home.ConfigFile(), []byte("token: secret"), 0644)

	opt := &InstallOption{}

	tempDir, _ := os.MkdirTemp("", "sl-installer-test")
	defer os.RemoveAll(tempDir)

	// Local
	i, err := NewInstaller(log, tempDir, "", 0, home, opt)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := i.(*localInstaller); !ok {
		t.Error("expected localInstaller")
	}

	// Archive
	i, err = NewInstaller(log, "plugin.zip", "", 0, home, opt)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := i.(*archiveInstaller); !ok {
		t.Error("expected archiveInstaller")
	}

	// Unsupported
	_, err = NewInstaller(log, "unsupported://source", "", 0, home, opt)
	if err == nil {
		t.Error("expected error for unsupported source")
	}
}

func TestNewLocalInstaller_Tilde(t *testing.T) {
	log := logrus.New()
	home := paths.Home("/tmp")
	opt := &InstallOption{}

	i, err := newLocalInstaller(log, "~/test-plugin", home, opt)
	if err != nil {
		t.Fatal(err)
	}
	if !filepath.IsAbs(i.source) {
		t.Errorf("expected absolute path, got %s", i.source)
	}
}

func TestLocalInstaller_Install(t *testing.T) {
	log := logrus.New()
	tempHome, _ := os.MkdirTemp("", "sl-home-local")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)
	os.MkdirAll(hh.Plugins(), 0755)

	pluginDir, _ := os.MkdirTemp("", "my-plugin")
	defer os.RemoveAll(pluginDir)
	metadataContent := `name: my-plugin
version: 1.0.0
`
	os.WriteFile(filepath.Join(pluginDir, plugin.MetadataFileName), []byte(metadataContent), 0644)

	i := &localInstaller{
		log:    log,
		home:   hh,
		source: pluginDir,
		opt:    &InstallOption{Force: true},
	}

	plug, err := i.Install()
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}
	if plug.Metadata.Name != "my-plugin" {
		t.Errorf("expected my-plugin, got %s", plug.Metadata.Name)
	}

	// Test already exists
	i.opt.Force = false
	_, err = i.Install()
	if err == nil {
		t.Error("expected error when plugin already exists and force is false")
	}

	// Test soft install (up-to-date)
	i.opt.Force = true
	i.opt.Soft = true
	_, err = i.Install()
	if err != ErrAlreadyUpToDate {
		t.Errorf("expected ErrAlreadyUpToDate, got %v", err)
	}

	// Test soft install (outdated)
	pluginDir2, _ := os.MkdirTemp("", "my-plugin-v2")
	defer os.RemoveAll(pluginDir2)
	metadataContent2 := `name: my-plugin
version: 1.1.0
`
	os.WriteFile(filepath.Join(pluginDir2, plugin.MetadataFileName), []byte(metadataContent2), 0644)
	i.source = pluginDir2
	plug3, err := i.Install()
	if err != nil {
		t.Fatalf("Soft install (outdated) failed: %v", err)
	}
	if plug3.Metadata.Version != "1.1.0" {
		t.Errorf("expected 1.1.0, got %s", plug3.Metadata.Version)
	}

	// Test DryRun
	i.opt.DryRun = true
	i.opt.Soft = false
	_, err = i.Install()
	if err != nil {
		t.Errorf("DryRun failed: %v", err)
	}

	// Test soft mode - error loading existing
	os.WriteFile(filepath.Join(tempHome, "plugins", "my-plugin", plugin.MetadataFileName), []byte("invalid yaml : : "), 0644)
	i.opt.DryRun = false
	_, err = i.Install()
	if err == nil {
		t.Error("expected error when existing plugin has invalid metadata")
	}
}

func TestLocalInstaller_Errors(t *testing.T) {
	log := logrus.New()
	tempHome, _ := os.MkdirTemp("", "sl-home-errors")
	defer os.RemoveAll(tempHome)
	hh := paths.Home(tempHome)

	// Illegal name
	pluginDir, _ := os.MkdirTemp("", "illegal")
	defer os.RemoveAll(pluginDir)
	os.WriteFile(filepath.Join(pluginDir, plugin.MetadataFileName), []byte("name: '!!'"), 0644)

	i := &localInstaller{log: log, home: hh, source: pluginDir, opt: &InstallOption{}}
	_, err := i.Install()
	if err != errIllegalPluginName {
		t.Errorf("expected %v, got %v", errIllegalPluginName, err)
	}

	// Metadata missing
	os.Remove(filepath.Join(pluginDir, plugin.MetadataFileName))
	_, err = i.Install()
	if err != errMissingMetadata {
		t.Errorf("expected %v, got %v", errMissingMetadata, err)
	}

	// LoadDir error (invalid YAML)
	os.WriteFile(filepath.Join(pluginDir, plugin.MetadataFileName), []byte("invalid: yaml: : "), 0644)
	_, err = i.Install()
	if err == nil {
		t.Error("expected error for invalid yaml metadata")
	}

	// Illegal version
	os.WriteFile(filepath.Join(pluginDir, plugin.MetadataFileName), []byte("name: legal\nversion: invalid"), 0644)
	_, err = i.Install()
	if err != errIllegalPluginVersion {
		t.Errorf("expected %v, got %v", errIllegalPluginVersion, err)
	}
}
