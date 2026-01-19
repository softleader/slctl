package plugin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUnmount(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "sl-mount-test")
	defer os.RemoveAll(tempDir)

	mountDir := filepath.Join(tempDir, "mount")
	os.MkdirAll(mountDir, 0755)

	p := &Plugin{
		Metadata: &Metadata{Name: "test"},
		Mount:    mountDir,
	}

	if err := p.Unmount(); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(mountDir); !os.IsNotExist(err) {
		t.Error("expected mount dir to be removed")
	}
}
