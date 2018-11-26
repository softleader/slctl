package installer

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/slpath"
	"golang.org/x/oauth2"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var gitHubRepo = regexp.MustCompile(`^(http[s]?://)?github.com/([^/]+)/([^/]+)[/]?$`)

type gitHubInstaller struct {
	httpInstaller
}

func newGitHubInstaller(source, tag string, home slpath.Home) (*gitHubInstaller, error) {
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
		if release, _, err = client.Repositories.GetLatestRelease(ctx, owner, repo); err != nil {
			return nil, err
		}
	} else {
		if release, _, err = client.Repositories.GetReleaseByTag(ctx, owner, repo, tag); err != nil {
			return nil, err
		}
	}

	asset, err := findAsset(release.Assets)
	if err != nil {
		return nil, err
	}

	rc, url, err := client.Repositories.DownloadReleaseAsset(ctx, owner, repo, asset.GetID())
	if err != nil {
		return nil, err
	}

	ghi := gitHubInstaller{}
	ghi.source = source
	ghi.home = home

	if url != "" {
		if ghi.downloader, err = newDownloader(url, home, filepath.Base(asset.GetBrowserDownloadURL())); err != nil {
			return nil, err
		}
	} else {
		if ghi.downloader, err = newDownloader(rc, home, filepath.Base(asset.GetBrowserDownloadURL())); err != nil {
			return nil, err
		}
	}

	return &ghi, nil
}

func dismantle(url string) (owner, repo string) {
	match := gitHubRepo.FindStringSubmatch(url)
	owner = match[2]
	repo = match[3]
	return
}

func findAsset(assets []github.ReleaseAsset) (*github.ReleaseAsset, error) {
	if len(assets) < 1 {
		return nil, fmt.Errorf("no assets found")
	}
	for _, asset := range assets {
		if strings.Contains(asset.GetName(), runtime.GOOS) {
			return &asset, nil
		}
	}
	return &assets[0], nil
}
