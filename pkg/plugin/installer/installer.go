package installer

import (
	"fmt"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"io"
	"os"
	"strings"
)

var ErrNonResolvableInOfflineMode = fmt.Errorf("non-resolvable plugin source in offline mode")

type Installer interface {
	Install() (*plugin.Plugin, error)
}

func NewInstaller(out io.Writer, source string, tag string, asset int, home slpath.Home, force bool) (Installer, error) {
	if isLocalDirReference(source) {
		return newLocalInstaller(out, source, home, force)
	}
	if isArchive(source) {
		return newArchiveInstaller(out, source, home, force)
	}
	if isGitHubRepo(source) {
		return newGitHubInstaller(out, source, tag, asset, home, force)
	}

	return nil, fmt.Errorf("unsupported plugin source: %s", source)
}

func isLocalDirReference(source string) bool {
	f, err := os.Stat(source)
	return err == nil && f.IsDir()
}

func isLocalReference(source string) bool {
	_, err := os.Stat(source)
	return err == nil
}

func isArchive(source string) bool {
	for _, suffix := range supportedExtensions {
		if strings.HasSuffix(source, suffix) {
			return true
		}
	}
	return false
}

func isGitHubRepo(source string) bool {
	return gitHubRepo.MatchString(source)
}
