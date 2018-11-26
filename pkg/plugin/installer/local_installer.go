package installer

import (
	"fmt"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/softleader/slctl/pkg/v"
	"io"
	"os"
	"path/filepath"
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
	v.Fprintf(i.out, "loading plugin from source: %s\n", i.source)
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

func isPlugin(dirname string) bool {
	_, err := os.Stat(filepath.Join(dirname, plugin.MetadataFileName))
	return err == nil
}
