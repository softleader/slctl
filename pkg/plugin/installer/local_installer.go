package installer

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"os"
	"path/filepath"
	"regexp"
)

var (
	ErrMissingMetadata      = errors.New("plugin metadata (" + plugin.MetadataFileName + ") missing")
	legalPluginName         = regexp.MustCompile(`^[\w\d_-]+$`)
	ErrIllegalPluginName    = errors.New("plugin name must match " + legalPluginName.String())
	ErrIllegalPluginVersion = errors.New("require a Semantic 2 version: https://semver.org/")
	ErrAlreadyUpToDate      = errors.New("already up-to-date")
)

type localInstaller struct {
	log    *logrus.Logger
	home   slpath.Home
	source string
	force  bool
	soft   bool // soft means remove exist plugin only if version is different
}

func newLocalInstaller(log *logrus.Logger, source string, home slpath.Home, force, soft bool) (*localInstaller, error) {
	src, err := filepath.Abs(source)
	if err != nil {
		return nil, fmt.Errorf("unable to get absolute path to plugin: %v", err)
	}
	return &localInstaller{
		log:    log,
		source: src,
		home:   home,
		force:  force,
		soft:   soft,
	}, nil
}

func (i *localInstaller) Install() (*plugin.Plugin, error) {
	if !isPlugin(i.source) {
		return nil, ErrMissingMetadata
	}
	plug, err := plugin.LoadDir(i.source)
	if err != nil {
		return nil, err
	}

	if !isPluginNameLegal(plug.Metadata.Name) {
		return nil, ErrIllegalPluginName
	}

	if !isVersionLegel(plug.Metadata.Version) {
		return nil, ErrIllegalPluginVersion
	}

	link := filepath.Join(i.home.Plugins(), plug.Metadata.Name)

	if _, err := os.Stat(link); !os.IsNotExist(err) {
		if !i.force {
			return nil, fmt.Errorf("plugin %q already exists", plug.Metadata.Name)
		}
		if i.soft {
			exist, err := plugin.LoadDir(link)
			if err != nil {
				return nil, err
			}
			if exist.Metadata.Version == plug.Metadata.Version {
				return exist, ErrAlreadyUpToDate
			}
		}
		i.log.Debugf("plugin %q already exists, force to remove it\n", plug.Metadata.Name)
		os.RemoveAll(link)
	}

	i.log.Printf("symlinking %s to %s\n", plug.Dir, link)
	if err := os.Symlink(plug.Dir, link); err != nil {
		return nil, err
	}

	return plugin.LoadDir(link)
}

func isPlugin(dirname string) bool {
	_, err := os.Stat(filepath.Join(dirname, plugin.MetadataFileName))
	return err == nil
}

func isPluginNameLegal(name string) bool {
	return legalPluginName.MatchString(name)
}

func isVersionLegel(version string) bool {
	if _, err := semver.NewVersion(version); err != nil {
		return false
	}
	return true
}
