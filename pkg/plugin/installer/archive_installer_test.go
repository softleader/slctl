package installer

import (
	"compress/flate"
	"os"
	"path/filepath"
	"testing"

	"github.com/mholt/archiver"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/plugin"
)

func TestArchiveInstaller_Install(t *testing.T) {
	home, err := os.MkdirTemp("", "sl_home")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(home)
	hh := paths.Home(home)
	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      false,
		ImplicitTopLevelFolder: false,
	}
	arcName := "test.zip"
	arcPath := filepath.Join(hh.String(), arcName)
	arcSrc := filepath.Join(hh.String(), "file.txt")
	if err := os.WriteFile(arcSrc, []byte("test"), 0744); err != nil {
		t.Error(err)
		return
	}
	if err := z.Archive([]string{arcSrc}, arcPath); err != nil {
		t.Error(err)
		return
	}

	log := logrus.New()
	i, err := newArchiveInstaller(log, arcPath, hh, &InstallOption{
		DryRun: false,
		Force:  true,
		Soft:   false,
	})
	if err != nil {
		t.Error(err)
		return
	}

	downloaded, err := i.downloader.download()
	if err != nil {
		t.Error(err)
		return
	}

	dst := filepath.Join(hh.CacheArchives(), arcName)
	if downloaded != dst {
		t.Errorf("expected downloaded dir to be %s", dst)
	}

	if !plugin.IsLocalReference(arcPath) {
		t.Errorf("expected downloaded dir to be a legal local reference: %s", dst)
	}
}

func TestArchiveInstaller_Install_Full(t *testing.T) {
	home, _ := os.MkdirTemp("", "sl_home_full")
	defer os.RemoveAll(home)
	hh := paths.Home(home)
	os.MkdirAll(hh.Plugins(), 0755)
	os.MkdirAll(hh.CacheArchives(), 0755)

	// Create a plugin directory to zip
	plugDir, _ := os.MkdirTemp("", "plug_src")
	defer os.RemoveAll(plugDir)
	metadataContent := `name: arch-plugin
version: 1.0.0
`
	os.WriteFile(filepath.Join(plugDir, plugin.MetadataFileName), []byte(metadataContent), 0644)

	arcPath := filepath.Join(home, "arch-plugin.zip")
	z := archiver.Zip{
		ImplicitTopLevelFolder: false,
	}
	if err := z.Archive([]string{filepath.Join(plugDir, plugin.MetadataFileName)}, arcPath); err != nil {
		t.Fatal(err)
	}

	log := logrus.New()
	i, err := newArchiveInstaller(log, arcPath, hh, &InstallOption{Force: true})
	if err != nil {
		t.Fatal(err)
	}

	plug, err := i.Install()
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	if plug.Metadata.Name != "arch-plugin" {
		t.Errorf("expected arch-plugin, got %s", plug.Metadata.Name)
	}
}

func TestArchiveInstaller_Errors(t *testing.T) {
	hh := paths.Home("/non/existent/home")
	log := logrus.New()
	
	// Non-existent source
	i, err := newArchiveInstaller(log, "/non/existent/file.zip", hh, &InstallOption{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = i.Install()
	if err == nil {
		t.Error("expected error for non-existent source")
	}

	// Download error
	i, _ = newArchiveInstaller(log, "http://invalid-url-123.com/file.zip", hh, &InstallOption{})
	_, err = i.Install()
	if err == nil {
		t.Error("expected error for download failure")
	}
}



