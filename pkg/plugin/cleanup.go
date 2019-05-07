package plugin

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/paths"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	dueDays = 30
)

// Cleanup remove outdated downloads of plugin
func Cleanup(log *logrus.Logger, home paths.Home, force bool, dryRun bool) error {
	conf, err := config.LoadConfFile(home.ConfigFile())
	if err != nil && err != config.ErrTokenNotExist {
		return err
	}
	if !force {
		if !needsToCleanup(conf.Cleanup) {
			return nil
		}
		log.Printf(`'slctl cleanup' has not been run in %v days, running now...`, dueDays)
	}

	installs := make(map[string]interface{})
	plugins, err := LoadPaths(home.Plugins())
	if err != nil {
		return err
	}
	for _, p := range plugins {
		if link, err := os.Readlink(p.Dir); err == nil {
			installs[filepath.Base(link)] = nil
		}
	}

	if err := remove(log, home.CachePlugins(), installs, dryRun); err != nil {
		return nil
	}
	if err := remove(log, home.CacheArchives(), installs, dryRun); err != nil {
		return nil
	}

	conf.Cleanup = time.Now().AddDate(0, 0, dueDays)
	return conf.WriteFile(home.ConfigFile(), 0644)
}

func remove(log *logrus.Logger, root string, installs map[string]interface{}, dryRun bool) error {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}
	for _, f := range files {
		if _, installed := installs[f.Name()]; !installed {
			wd := filepath.Join(root, f.Name())
			log.Printf(`Removing: %s...`, wd)
			if !dryRun {
				os.RemoveAll(wd)
			}
		}
	}
	return nil
}

func needsToCleanup(dueDate time.Time) bool {
	return dueDate.Before(time.Now())
}
