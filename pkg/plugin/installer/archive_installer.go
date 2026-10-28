package installer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archives"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/plugin"
)

type archiveInstaller struct {
	localInstaller
	downloader downloader
}

func newArchiveInstaller(log *logrus.Logger, source string, home paths.Home, opt *InstallOption) (ai *archiveInstaller, err error) {
	log.Debugf("downloading the archive: %s\n", source)
	ai = &archiveInstaller{}
	ai.log = log
	ai.source = source
	ai.home = home
	ai.opt = opt
	if plugin.IsLocalReference(source) {
		var r io.Reader
		if r, err = os.Open(source); err != nil {
			return nil, err
		}
		ai.downloader = newReaderDownloader(&r, home, filepath.Base(source))
	} else {
		if environment.Settings.Offline {
			return nil, errNonResolvableInOfflineMode
		}
		ai.downloader = newURLDownloader(source, home, filepath.Base(source))
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
	if err := os.WriteFile(filepath.Join(extractDir, plugin.SourceFileName), []byte(i.source), 0644); err != nil {
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

func extract(source, destination string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()

	ctx := context.Background()
	format, stream, err := archives.Identify(ctx, source, input)
	if err != nil {
		return err
	}
	extractor, ok := format.(archives.Extractor)
	if !ok {
		return fmt.Errorf("%s is not a supported archive", source)
	}
	if err := os.MkdirAll(destination, 0755); err != nil {
		return err
	}
	root, err := os.OpenRoot(destination)
	if err != nil {
		return err
	}
	defer root.Close()

	return extractor.Extract(ctx, stream, func(_ context.Context, file archives.FileInfo) error {
		// Reject Windows separators, drive names and alternate data streams on
		// every platform, so an archive has the same safety rules everywhere.
		name := file.NameInArchive
		if strings.ContainsAny(name, `\:`) || !filepath.IsLocal(name) {
			return fmt.Errorf("unsafe archive path %q", name)
		}
		if file.LinkTarget != "" || (!file.IsDir() && !file.Mode().IsRegular()) {
			return fmt.Errorf("unsupported archive entry %q: links and special files are not allowed", name)
		}
		name = filepath.FromSlash(name)
		if file.IsDir() {
			return root.MkdirAll(name, 0755)
		}
		if err := root.MkdirAll(filepath.Dir(name), 0755); err != nil {
			return err
		}
		in, err := file.Open()
		if err != nil {
			return err
		}
		defer in.Close()

		// Root confines filesystem operations even when an existing directory
		// is a symlink. O_EXCL also prevents duplicate entries from overwriting
		// files or following an existing file symlink/hard link.
		out, err := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, file.Mode().Perm())
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(out, in)
		return errors.Join(copyErr, out.Close())
	})
}
