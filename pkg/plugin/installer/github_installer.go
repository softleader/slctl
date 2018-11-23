package installer

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/softleader/slctl/pkg/v"
	"golang.org/x/oauth2"
	"io"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var gitHubRepo = regexp.MustCompile(`^(http[s]?://)?github.com/([^/]+)/([^/]+)[/]?$`)

type GitHubInstaller struct {
	home   slpath.Home
	source string
	owner  string
	repo   string
	tag    string
}

func (i GitHubInstaller) supports(source string) bool {
	return gitHubRepo.MatchString(source)
}

func (i GitHubInstaller) new(source, tag string, home slpath.Home) (Installer, error) {
	match := gitHubRepo.FindStringSubmatch(source)
	return GitHubInstaller{
		home:   home,
		source: source,
		owner:  match[2],
		repo:   match[3],
		tag:    tag,
	}, nil
}

func findAssert(assets []github.ReleaseAsset) (*github.ReleaseAsset, error) {
	if len(assets) < 1 {
		return nil, fmt.Errorf("no assets to find")
	}
	for _, asset := range assets {
		if strings.Contains(asset.GetName(), runtime.GOOS) {
			return &asset, nil
		}
	}
	return &assets[0], nil
}

func (i GitHubInstaller) Install() (*plugin.Plugin, error) {
	conf, err := config.LoadConfFile(i.home.ConfigFile())
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var release *github.RepositoryRelease
	if i.tag == "" {
		if release, _, err = client.Repositories.GetLatestRelease(ctx, i.owner, i.repo); err != nil {
			return nil, err
		}
	} else {
		if release, _, err = client.Repositories.GetReleaseByTag(ctx, i.owner, i.repo, i.tag); err != nil {
			return nil, err
		}
	}

	assert, err := findAssert(release.Assets)
	if err != nil {
		return nil, err
	}

	var dl downloader
	var url string
	var rc io.ReadCloser
	if rc, url, err = client.Repositories.DownloadReleaseAsset(ctx, i.owner, i.repo, assert.GetID()); err != nil {
		return nil, err
	}
	if url != "" {
		if dl, err = newDownloader(url); err != nil {
			return nil, err
		}
	} else {
		if dl, err = newDownloader(rc); err != nil {
			return nil, err
		}
	}

	archiveName := filepath.Base(assert.GetBrowserDownloadURL())
	archivePath := filepath.Join(i.home.CacheArchives(), archiveName)
	dl.downloadTo(archivePath)

	v.Println(archivePath, "downloaded.")

	extractDir := filepath.Join(i.home.CachePlugins(), archiveName)
	ensureDirEmpty(extractDir)

	if err = extract(archivePath, extractDir); err != nil {
		return nil, err
	}

	if !isPlugin(extractDir) {
		return nil, ErrMissingMetadata
	}

	plug, err := plugin.LoadDir(extractDir)
	if err != nil {
		return nil, err
	}

	linked, err := plug.LinkTo(i.home)
	if err != nil {
		return nil, err
	}

	return plugin.LoadDir(linked)
}
