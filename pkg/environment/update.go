package environment

import (
	"context"
	"time"

	"github.com/softleader/slctl/pkg/release"

	"github.com/google/go-github/v28/github"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/ver"
)

const (
	owner = "softleader"
	repo  = "slctl"
)

// CheckForUpdates 檢查是否有更新的 slctl 版本
func CheckForUpdates(log *logrus.Logger, home paths.Home, metadata *release.Metadata, force bool) error {
	if !metadata.IsReleased() {
		return nil
	}
	conf, err := config.LoadConfFile(home.ConfigFile())
	if err != nil && err != config.ErrTokenNotExist {
		return err
	}
	if !force {
		if !needsToCheckOnline(conf.CheckUpdates) {
			return nil
		}
	}
	if updateAvailable, err := checkOnline(log, metadata.GitVersion); err != nil {
		return err
	} else if updateAvailable {
		conf.UpdateCheckUpdatesTimeInDays(1) // 如果已經發現可以更新, 一天後需要檢查是否更新了沒
	} else {
		conf.UpdateCheckUpdatesTime() // 如果沒有任何更新檔, 就可以等久一點再檢查吧
	}
	return conf.WriteFile(home.ConfigFile(), 0644)
}

func needsToCheckOnline(dueDate time.Time) bool {
	return dueDate.Before(time.Now())
}

func checkOnline(log *logrus.Logger, currentVersion string) (bool, error) {
	log.Println("Checking for latest slctl version...")
	client := github.NewClient(nil)
	ctx := context.Background()
	log.Debugf("fetching the latest published release from github.com/%s/%s", owner, repo)
	rr, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return false, err
	}
	latest := rr.GetTagName()
	log.Debugf("found latest version: %s+%s", latest, rr.GetTargetCommitish())
	updateAvailable, err := ver.Revision(latest).IsGreaterThan(currentVersion)
	if err != nil {
		return false, err
	}
	if updateAvailable {
		log.Printf(`A new version %q is available
It is recommended using package managers to update
Read more: https://github.com/softleader/slctl#upgrade`, latest)
	} else {
		log.Println("No update available")
	}
	return updateAvailable, nil
}
