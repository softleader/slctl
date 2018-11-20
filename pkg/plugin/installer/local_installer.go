package installer

import (
	"fmt"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/softleader/slctl/pkg/v"
	"path/filepath"
)

// LocalInstaller installs plugins from the filesystem.
type LocalInstaller struct {
	base
}

// NewLocalInstaller creates a new LocalInstaller.
func NewLocalInstaller(source string, home slpath.Home) (*LocalInstaller, error) {
	src, err := filepath.Abs(source)
	if err != nil {
		return nil, fmt.Errorf("unable to get absolute path to plugin: %v", err)
	}
	i := &LocalInstaller{
		base: newBase(src, home),
	}
	return i, nil
}

// Install creates a symlink to the plugin directory in $HELM_HOME.
//
// Implements Installer.
func (i *LocalInstaller) Install() error {
	if !isPlugin(i.Source) {
		return ErrMissingMetadata
	}
	return i.link(i.Source)
}

// Update updates a local repository
func (i *LocalInstaller) Update() error {
	v.Println("local repository is auto-updated")
	return nil
}
