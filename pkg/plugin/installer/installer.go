package installer

import (
	"errors"
	"fmt"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"io"
	"os"
	"strings"
)

var (
	ErrMissingMetadata = errors.New("plugin metadata (" + plugin.MetadataFileName + ") missing")
)

type Installer interface {
	Install() (*plugin.Plugin, error)
}

func NewInstaller(out io.Writer, source string, version string, home slpath.Home) (Installer, error) {
	if isLocalReference(source) {
		return newLocalInstaller(out, source, home)
	}
	if environment.Settings.Offline {
		return nil, fmt.Errorf("non-resolvable plugin source (%s) in offline mode", source)
	}
	if isRemoteHTTPArchive(source) {
		return newHttpInstaller(out, source, home)
	} else if isGitHubRepo(source) {
		return newGitHubInstaller(out, source, version, home)
	}

	return nil, fmt.Errorf("unsupported plugin source: %s", source)
}

func isLocalReference(source string) bool {
	_, err := os.Stat(source)
	return err == nil
}

func isRemoteHTTPArchive(source string) bool {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		for _, suffix := range supportedExtensions {
			if strings.HasSuffix(source, suffix) {
				return true
			}
		}
	}
	return false
}

func isGitHubRepo(source string) bool {
	return gitHubRepo.MatchString(source)
}
