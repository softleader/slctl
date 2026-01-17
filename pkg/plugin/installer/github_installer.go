package installer

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/go-github/v28/github"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/plugin"
	"golang.org/x/oauth2"
)

type gitHubInstaller struct {
	archiveInstaller
}

func newGitHubInstaller(log *logrus.Logger, source, tag string, asset int, home paths.Home, opt *InstallOption) (*gitHubInstaller, error) {
	if environment.Settings.Offline {
		return nil, errNonResolvableInOfflineMode
	}
	conf, err := config.LoadConfFile(home.ConfigFile())
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	owner, repo := dismantle(source)

	var release *github.RepositoryRelease
	if tag == "" {
		log.Debugf("fetching the latest published release from github.com/%s/%s\n", owner, repo)
		if release, _, err = client.Repositories.GetLatestRelease(ctx, owner, repo); err != nil {
			return nil, err
		}
		log.Printf("fetching latest published release %q", release.GetTagName())
	} else {
		log.Debugf("fetching the release from github.com/%s/%s with tag %q\n", owner, repo, tag)
		if release, _, err = client.Repositories.GetReleaseByTag(ctx, owner, repo, tag); err != nil {
			return nil, err
		}
	}

	if body := release.GetBody(); len(body) > 0 {
		for _, line := range strings.Split(body, "\n") {
			log.Println(line)
		}
	}

	ra, err := pickAsset(log, release, asset)
	if err != nil {
		return nil, err
	}

	rc, url, err := client.Repositories.DownloadReleaseAsset(ctx, owner, repo, ra.GetID())
	if err != nil {
		return nil, err
	}

	ghi := gitHubInstaller{}
	ghi.source = source
	ghi.home = home
	ghi.log = log
	ghi.opt = opt

	binary := ra.GetBrowserDownloadURL()
	log.Debugf("downloading %s\n", binary)

	if url != "" {
		log.Debugf("downloading from %s\n", url)
		ghi.downloader = newURLDownloader(url, home, filepath.Base(binary))
	} else {
		ghi.downloader = newReadCloserDownloader(&rc, ra.GetSize(), home, filepath.Base(binary))
	}

	return &ghi, nil
}

func dismantle(url string) (owner, repo string) {
	match := plugin.GitHubRepo.FindStringSubmatch(url)
	owner = match[2]
	repo = match[3]
	return
}

func pickAsset(log *logrus.Logger, release *github.RepositoryRelease, asset int) (ra *github.ReleaseAsset, err error) {
	if len(release.Assets) < 1 {
		err = fmt.Errorf("no assets found on release %q", release.GetTagName())
		return
	}
	if asset > 0 {
		if len := len(release.Assets); asset >= len {
			return nil, fmt.Errorf("only %v assets found on release %q but you're asking %v (start from zero)", len, release.GetName(), asset)
		}
		ra = &release.Assets[asset]
		return
	}
	log.Debugf("trying to find asset for %q\n", runtime.GOOS)
	if ra = findRuntimeOsAsset(log, release.Assets); ra == nil {
		log.Debugf("%s asset not found, using first asset\n", runtime.GOOS)
		ra = &release.Assets[0]
	}
	return
}

func findRuntimeOsAsset(_ *logrus.Logger, assets []github.ReleaseAsset) *github.ReleaseAsset {
	for _, asset := range assets {
		if strings.Contains(asset.GetName(), runtime.GOOS) {
			return &asset
		}
	}
	return nil
}
