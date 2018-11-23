package installer

import (
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/softleader/slctl/pkg/v"
	"path/filepath"
	"strings"
)

var supportedExtensions = []string{
	".zip",
	".tar",
	".tar.gz",
	".tgz",
	".tar.bz2",
	".tbz2",
	".tar.xz",
	".txz",
	".tar.lz4",
	".tlz4",
	".tar.sz",
	".tsz",
	".rar",
	".bz2",
	".gz",
	".lz4",
	".sz",
	".xz",
}

type httpInstaller struct {
	home       slpath.Home
	source     string
	downloader downloader
}

func (i httpInstaller) supports(source string) bool {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		for _, suffix := range supportedExtensions {
			if strings.HasSuffix(source, suffix) {
				return true
			}
		}
	}
	return false
}

func (i httpInstaller) new(source, _ string, home slpath.Home) (Installer, error) {
	dl, err := newDownloader(source)
	if err != nil {
		return nil, err
	}
	return httpInstaller{
		home:       home,
		source:     source,
		downloader: dl,
	}, nil
}

func (i httpInstaller) Install() (*plugin.Plugin, error) {
	archiveName := filepath.Base(i.source)
	archivePath := filepath.Join(i.home.CacheArchives(), archiveName)
	i.downloader.downloadTo(archivePath)

	v.Println(archivePath, "downloaded.")

	extractDir := filepath.Join(i.home.CachePlugins(), archiveName)
	ensureDirEmpty(extractDir)

	if err := extract(archivePath, extractDir); err != nil {
		return nil, err
	}

	if !isPlugin(extractDir) {
		return nil, ErrMissingMetadata
	}

	plug, err := plugin.LoadDir(extractDir)
	if err != nil {
		return nil, err
	}

	linked, err := plug.LinkTo(i.home)
	if err != nil {
		return nil, err
	}

	return plugin.LoadDir(linked)
}
