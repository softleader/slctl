package installer

import (
	"fmt"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"io"
)

var ErrNonResolvableInOfflineMode = fmt.Errorf("non-resolvable plugin source in offline mode")

type Installer interface {
	Install() (*plugin.Plugin, error)
}

func NewInstaller(out io.Writer, source string, tag string, asset int, home slpath.Home, force, soft bool) (Installer, error) {
	if plugin.IsLocalDirReference(source) {
		return newLocalInstaller(out, source, home, force, soft)
	}
	if plugin.IsArchive(source) {
		return newArchiveInstaller(out, source, home, force, soft)
	}
	if plugin.IsGitHubRepo(source) {
		return newGitHubInstaller(out, source, tag, asset, home, force, soft)
	}
	return nil, fmt.Errorf("unsupported plugin source: %s", source)
}
