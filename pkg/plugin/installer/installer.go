package installer

import (
	"errors"
	"fmt"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"os"
	"strings"
)

var ErrMissingMetadata = errors.New("plugin metadata (" + plugin.MetadataFileName + ") missing")

type Installer interface {
	Install() (*plugin.Plugin, error)
}

func NewInstaller(source string, version string, home slpath.Home) (Installer, error) {
	if isLocalReference(source) {
		return newLocalInstaller(source, home)
	} else if isRemoteHTTPArchive(source) {
		return newHttpInstaller(source, home)
	} else if isGitHubRepo(source) {
		return newGitHubInstaller(source, version, home)
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
