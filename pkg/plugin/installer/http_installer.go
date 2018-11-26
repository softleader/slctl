package installer

import (
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
	localInstaller
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
	dl, err := newDownloader(source, home, filepath.Base(source))
	if err != nil {
		return nil, err
	}

	hi := httpInstaller{}
	hi.source = source
	hi.home = home
	hi.downloader = dl
	return hi, nil
}

func (i httpInstaller) retrievePlugin() error {
	saved, err := i.downloader.download()
	if err != nil {
		return err
	}
	v.Println(saved, "downloaded.")
	extractDir := filepath.Join(i.home.CachePlugins(), filepath.Base(saved))
	if err := ensureDirEmpty(extractDir); err != nil {
		return err
	}
	return extract(saved, extractDir)
}
