package environment

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/ver"
	"time"
)

const (
	owner = "softleader"
	repo  = "slctl"
)

// CheckForUpdates 檢查是否有更新的 slctl 版本
func CheckForUpdates(log *logrus.Logger, home paths.Home, currentVersion string, force bool) error {
	conf, err := config.LoadConfFile(home.ConfigFile())
	if err != nil && err != config.ErrTokenNotExist {
		return err
	}
	if !force {
		if !needsToCheckOnline(conf.CheckUpdates) {
			return nil
		}
	}
	if err := checkOnline(log, currentVersion); err != nil {
		return err
	}
	conf.UpdateCheckUpdatesTime()
	return conf.WriteFile(home.ConfigFile(), 0644)
}

func needsToCheckOnline(dueDate time.Time) bool {
	return dueDate.Before(time.Now())
}

func checkOnline(log *logrus.Logger, currentVersion string) error {
	log.Println("Checking for latest slctl version...")
	client := github.NewClient(nil)
	ctx := context.Background()
	log.Debugf("fetching the latest published release from github.com/%s/%s", owner, repo)
	rr, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return err
	}
	latest := rr.GetTagName()
	log.Debugf("found latest version: %s+%s", latest, rr.GetTargetCommitish())
	if updateAvailable, err := ver.Revision(latest).IsGreaterThan(currentVersion); err != nil {
		return err
	} else if updateAvailable {
		log.Printf(`A new version %q is available
It is recommended using package managers to update
Read more: https://github.com/softleader/slctl#upgrade`, latest)
	} else {
		log.Println("No update available")
	}
	return nil
}
