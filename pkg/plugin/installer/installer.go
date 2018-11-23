package installer

import (
	"errors"
	"fmt"
	"github.com/mholt/archiver"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"os"
	"path/filepath"
)

var ErrMissingMetadata = errors.New("plugin metadata (" + plugin.MetadataFileName + ") missing")

type Installer interface {
	Install() (*plugin.Plugin, error)
	supports(source string) bool
	new(source, version string, home slpath.Home) (Installer, error)
}

var installers = []Installer{
	LocalInstaller{},
	GitHubInstaller{},
}

func NewInstaller(source string, version string, home slpath.Home) (Installer, error) {
	for _, i := range installers {
		if i.supports(source) {
			return i.new(source, version, home)
		}
	}
	return nil, fmt.Errorf("unsupported plugin source: %s", source)
}

func isPlugin(dirname string) bool {
	_, err := os.Stat(filepath.Join(dirname, plugin.MetadataFileName))
	return err == nil
}

func extract(source, destination string) (err error) {
	if err = archiver.Unarchive(source, destination); err != nil { // find Unarchiver by header
		var arc interface{}
		if arc, err = archiver.ByExtension(source); err != nil { // try again to find by extension
			return err
		}
		if err = arc.(archiver.Unarchiver).Unarchive(source, destination); err != nil {
			return err
		}
	}
	return
}

func ensureDirEmpty(path string) error {
	if fi, err := os.Stat(path); err != nil {
		if err = os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("could not create %s: %s", path, err)
		}
		return nil
	} else if !fi.IsDir() {
		return fmt.Errorf("%s must be a directory", path)
	}
	// if goes here, dir already exist, so let's delete it
	return os.RemoveAll(path)
}
