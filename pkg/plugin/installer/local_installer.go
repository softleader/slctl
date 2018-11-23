package installer

import (
	"fmt"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"os"
	"path/filepath"
)

type LocalInstaller struct {
	home   slpath.Home
	source string
}

func (i LocalInstaller) supports(source string) bool {
	_, err := os.Stat(source)
	return err == nil
}

func (i LocalInstaller) new(source, _ string, home slpath.Home) (Installer, error) {
	src, err := filepath.Abs(source)
	if err != nil {
		return nil, fmt.Errorf("unable to get absolute path to plugin: %v", err)
	}
	return LocalInstaller{
		source: src,
		home:   home,
	}, nil
}

func (i LocalInstaller) Install() (*plugin.Plugin, error) {
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
