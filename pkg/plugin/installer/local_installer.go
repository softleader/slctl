package installer

import (
	"errors"
	"fmt"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/softleader/slctl/pkg/v"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

var (
	ErrMissingMetadata   = errors.New("plugin metadata (" + plugin.MetadataFileName + ") missing")
	legalPluginName      = regexp.MustCompile(`^[\w\d_-]+$`)
	ErrIllegalPluginName = errors.New("plugin name must match " + legalPluginName.String())
)

type localInstaller struct {
	out    io.Writer
	home   slpath.Home
	source string
}

func newLocalInstaller(out io.Writer, source string, home slpath.Home) (*localInstaller, error) {
	src, err := filepath.Abs(source)
	if err != nil {
		return nil, fmt.Errorf("unable to get absolute path to plugin: %v", err)
	}
	return &localInstaller{
		out:    out,
		source: src,
		home:   home,
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

	link := filepath.Join(i.home.Plugins(), plug.Metadata.Name)
	v.Printf("symlinking %s to %s\n", plug.Dir, link)

	if _, err := os.Stat(link); !os.IsNotExist(err) {
		return nil, errors.New(`plugin '` + plug.Metadata.Name + `' already exists.`)
	}
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
