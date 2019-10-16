package plugin

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/paths"
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
		log.Printf(`'slctl cleanup' has not been run in %v days, running now...`, config.CleanupDueDays)
	}

	installedPlugins := make(map[string]interface{})
	plugins, err := LoadPaths(home.Plugins())
	if err != nil {
		return err
	}
	for _, p := range plugins {
		if link, err := os.Readlink(p.Dir); err == nil {
			installedPlugins[filepath.Base(link)] = nil
		}
	}

	log.Debugf("cleaning up %s", home.CachePlugins())
	if err := remove(log, home, home.CachePlugins(), installedPlugins, dryRun); err != nil {
		return nil
	}

	log.Debugf("cleaning up %s", home.CacheArchives())
	if err := remove(log, home, home.CacheArchives(), installedPlugins, dryRun); err != nil {
		return nil
	}

	conf.UpdateCleanupTime()
	return conf.WriteFile(home.ConfigFile(), 0644)
}

func remove(log *logrus.Logger, home paths.Home, root string, installedPlugins map[string]interface{}, dryRun bool) error {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}
	for _, f := range files {
		if _, installed := installedPlugins[f.Name()]; !installed {
			wd := filepath.Join(root, f.Name())
			if needsToRemove(log, home, wd) {
				log.Printf(`removing: %s...`, wd)
				if !dryRun {
					os.RemoveAll(wd)
				}
			}
		}
	}
	return nil
}

func needsToCleanup(dueDate time.Time) bool {
	return dueDate.Before(time.Now())
}

// 這邊處理特殊的 clean up 檔案, 例如要比對 due date 之類的, 而一般(非特殊)的處理方式就是把檔案刪掉
func needsToRemove(log *logrus.Logger, home paths.Home, file string) bool {
	crf := home.CacheRepositoryFile()
	if crf == file && paths.IsExistFile(file) {
		if r, err := loadLocal(log, crf); err == nil && !expired(r) {
			log.Debugf("skip '%s' which is still up to date", file)
			return false
		}
	}
	return true
}
