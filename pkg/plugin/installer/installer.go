package installer

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/plugin"
)

var ErrNonResolvableInOfflineMode = fmt.Errorf("non-resolvable plugin source in offline mode")

type Installer interface {
	Install() (*plugin.Plugin, error)
}

func NewInstaller(log *logrus.Logger, source string, tag string, asset int, home paths.Home, dryRun, force, soft bool) (Installer, error) {
	if plugin.IsLocalDirReference(source) {
		return newLocalInstaller(log, source, home, dryRun, force, soft)
	}
	if plugin.IsSupportedArchive(source) {
		return newArchiveInstaller(log, source, home, dryRun, force, soft)
	}
	if plugin.IsGitHubRepo(source) {
		return newGitHubInstaller(log, source, tag, asset, home, dryRun, force, soft)
	}
	return nil, fmt.Errorf("unsupported plugin source: %s", source)
}
