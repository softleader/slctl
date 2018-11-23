package installer

import (
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/softleader/slctl/pkg/v"
	"os"
	"path/filepath"
)

type base struct {
	// Source is the reference to a plugin
	Source string
	// HelmHome is the $HELM_HOME directory
	SlHome slpath.Home
}

func newBase(source string, home slpath.Home) base {
	return base{source, home}
}

// link creates a symlink from the plugin source to $HELM_HOME.
func (b *base) link(from string) error {
	v.Printf("symlinking %s to %s", from, b.Path())
	return os.Symlink(from, b.Path())
}

// Path is where the plugin will be symlinked to.
func (b *base) Path() string {
	if b.Source == "" {
		return ""
	}
	return filepath.Join(b.SlHome.Plugins(), filepath.Base(b.Source))
}
