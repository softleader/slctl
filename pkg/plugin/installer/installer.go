package installer

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/paths"
)

var ErrNonResolvableInOfflineMode = fmt.Errorf("non-resolvable plugin source in offline mode")

type Installer interface {
	Install() (*plugin.Plugin, error)
}

func NewInstaller(log *logrus.Logger, source string, tag string, asset int, home paths.Home, force, soft bool) (Installer, error) {
	if plugin.IsLocalDirReference(source) {
		return newLocalInstaller(log, source, home, force, soft)
	}
	if plugin.IsSupportedArchive(source) {
		return newArchiveInstaller(log, source, home, force, soft)
	}
	if plugin.IsGitHubRepo(source) {
		return newGitHubInstaller(log, source, tag, asset, home, force, soft)
	}
	return nil, fmt.Errorf("unsupported plugin source: %s", source)
}
