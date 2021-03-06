package installer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/plugin"
)

var (
	errMissingMetadata      = errors.New("plugin metadata (" + plugin.MetadataFileName + ") missing")
	legalPluginName         = regexp.MustCompile(`^[\w\d_-]+$`)
	errIllegalPluginName    = errors.New("plugin name must match " + legalPluginName.String())
	errIllegalPluginVersion = errors.New("require a Semantic 2 version: https://semver.org/")
	// ErrAlreadyUpToDate 表示 plugin 版本已經最新
	ErrAlreadyUpToDate = errors.New("already up-to-date")
)

type localInstaller struct {
	log    *logrus.Logger
	home   paths.Home
	source string
	opt    *InstallOption
}

func newLocalInstaller(log *logrus.Logger, source string, home paths.Home, opt *InstallOption) (*localInstaller, error) {
	if expanded, err := homedir.Expand(source); err != nil {
		source = expanded
	}
	src, err := filepath.Abs(source)
	if err != nil {
		return nil, fmt.Errorf("unable to get absolute path to plugin: %v", err)
	}
	return &localInstaller{
		log:    log,
		source: src,
		home:   home,
		opt:    opt,
	}, nil
}

func (i *localInstaller) Install() (*plugin.Plugin, error) {
	if !isPlugin(i.source) {
		return nil, errMissingMetadata
	}
	plug, err := plugin.LoadDir(i.source)
	if err != nil {
		return nil, err
	}

	if !isPluginNameLegal(plug.Metadata.Name) {
		return nil, errIllegalPluginName
	}

	if !plug.Metadata.IsVersionLegal() {
		return nil, errIllegalPluginVersion
	}

	link := filepath.Join(i.home.Plugins(), plug.Metadata.Name)

	if _, err := os.Stat(link); !os.IsNotExist(err) {
		if !i.opt.Force {
			return nil, fmt.Errorf("plugin %q already exists", plug.Metadata.Name)
		}
		if i.opt.Soft {
			exist, err := plugin.LoadDir(link)
			if err != nil {
				return nil, err
			}
			if outdated, err := plug.Metadata.IsVersionGreaterThan(exist.Metadata); err != nil {
				return nil, err
			} else if !outdated {
				return exist, ErrAlreadyUpToDate
			}
		}
		i.log.Debugf("plugin %q already exists, force to remove it\n", plug.Metadata.Name)
		if !i.opt.DryRun {
			os.RemoveAll(link)
		}
	}

	i.log.Debugf("symbolic linking %s to %s\n", plug.Dir, link)
	if !i.opt.DryRun {
		if err := os.Symlink(plug.Dir, link); err != nil {
			return nil, err
		}
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
