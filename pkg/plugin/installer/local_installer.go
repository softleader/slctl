package installer

import (
	"fmt"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"os"
	"path/filepath"
)

type localInstaller struct {
	home   slpath.Home
	source string
}

func (i localInstaller) supports(source string) bool {
	_, err := os.Stat(source)
	return err == nil
}

func (i localInstaller) new(source, _ string, home slpath.Home) (Installer, error) {
	src, err := filepath.Abs(source)
	if err != nil {
		return nil, fmt.Errorf("unable to get absolute path to plugin: %v", err)
	}
	return localInstaller{
		source: src,
		home:   home,
	}, nil
}

func (i localInstaller) install() (*plugin.Plugin, error) {
	if !isPlugin(i.source) {
		return nil, ErrMissingMetadata
	}

	plug, err := plugin.LoadDir(i.source)
	if err != nil {
		return nil, err
	}

	linked, err := plug.LinkTo(i.home)
	if err != nil {
		return nil, err
	}

	return plugin.LoadDir(linked)
}

func (i localInstaller) retrievePlugin() error {
	// local plugin is already on the host
	return nil
}
