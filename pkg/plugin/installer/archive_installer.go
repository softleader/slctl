package installer

import (
	"fmt"
	"github.com/mholt/archiver"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/plugin"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type archiveInstaller struct {
	localInstaller
	downloader downloader
}

func newArchiveInstaller(log *logrus.Logger, source string, home paths.Home, dryRun, force, soft bool) (ai *archiveInstaller, err error) {
	log.Debugf("downloading the archive: %s\n", source)
	ai = &archiveInstaller{}
	ai.log = log
	ai.source = source
	ai.home = home
	ai.dryRun = dryRun
	ai.force = force
	ai.soft = soft
	if plugin.IsLocalReference(source) {
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
	i.log.Debugf("successfully downloaded and saved it to: %s\n", saved)
	extractDir := filepath.Join(i.home.CachePlugins(), filepath.Base(saved))
	if err := ensureDirEmpty(extractDir); err != nil {
		return err
	}
	i.log.Debugln("extracting archive to", extractDir)
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
