package installer

import (
	"fmt"
	"github.com/mholt/archiver"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/softleader/slctl/pkg/v"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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

type archiveInstaller struct {
	localInstaller
	downloader downloader
}

func newArchiveInstaller(out io.Writer, source string, home slpath.Home, force, soft bool) (ai *archiveInstaller, err error) {
	v.Fprintf(out, "downloading the archive: %s\n", source)
	ai = &archiveInstaller{}
	ai.out = out
	ai.source = source
	ai.home = home
	ai.force = force
	ai.soft = soft
	if isLocalReference(source) {
		var r io.Reader
		if r, err = os.Open(source); err != nil {
			return nil, err
		}
		ai.downloader = newReaderDownloader(&r, home, filepath.Base(source))
	} else {
		if environment.Settings.Offline {
			return nil, ErrNonResolvableInOfflineMode
		}
		ai.downloader = newUrlDownloader(source, home, filepath.Base(source))
	}
	return ai, nil
}

func (i *archiveInstaller) Install() (*plugin.Plugin, error) {
	if err := i.retrievePlugin(); err != nil {
		return nil, err
	}
	return i.localInstaller.Install()
}

func (i *archiveInstaller) retrievePlugin() error {
	saved, err := i.downloader.download()
	if err != nil {
		return err
	}
	v.Fprintf(i.out, "successfully downloaded and saved it to: %s\n", saved)
	extractDir := filepath.Join(i.home.CachePlugins(), filepath.Base(saved))
	if err := ensureDirEmpty(extractDir); err != nil {
		return err
	}
	v.Fprintln(i.out, "extracting archive to", extractDir)
	if err = extract(saved, extractDir); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(extractDir, plugin.SourceFileName), []byte(i.source), 0644); err != nil {
		return err
	}
	i.source = extractDir
	return nil
}

func ensureDirEmpty(path string) error {
	if fi, err := os.Stat(path); err != nil {
		if err = os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("could not create %s: %s", path, err)
		}
		return nil
	} else if !fi.IsDir() {
		return fmt.Errorf("%s must be a directory", path)
	}
	// if goes here, dir already exist, so let's delete it
	return os.RemoveAll(path)
}

func extract(source, destination string) (err error) {
	if err = archiver.Unarchive(source, destination); err != nil { // find Unarchiver by header
		var arc interface{}
		if arc, err = archiver.ByExtension(source); err != nil { // try again to find by extension
			return err
		}
		if err = arc.(archiver.Unarchiver).Unarchive(source, destination); err != nil {
			return err
		}
	}
	return
}
