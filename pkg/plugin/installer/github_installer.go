package installer

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/softleader/slctl/pkg/verbose"
	"golang.org/x/oauth2"
	"io"
	"path/filepath"
	"runtime"
	"strings"
)

type gitHubInstaller struct {
	archiveInstaller
}

func newGitHubInstaller(out io.Writer, source, tag string, asset int, home slpath.Home, force, soft bool) (*gitHubInstaller, error) {
	if environment.Settings.Offline {
		return nil, ErrNonResolvableInOfflineMode
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
		verbose.Fprintf(out, "fetching the latest published release from github.com/%s/%s\n", owner, repo)
		if release, _, err = client.Repositories.GetLatestRelease(ctx, owner, repo); err != nil {
			return nil, err
		}
	} else {
		verbose.Fprintf(out, "fetching the release from github.com/%s/%s with tag %q\n", owner, repo, tag)
		if release, _, err = client.Repositories.GetReleaseByTag(ctx, owner, repo, tag); err != nil {
			return nil, err
		}
	}

	ra, err := pickAsset(out, release, asset)
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
	ghi.out = out
	ghi.force = force
	ghi.soft = soft

	binary := ra.GetBrowserDownloadURL()
	verbose.Fprintf(out, "downloading the binary content: %s\n", binary)

	if url != "" {
		ghi.downloader = newUrlDownloader(url, home, filepath.Base(binary))
	} else {
		ghi.downloader = newReadCloserDownloader(&rc, home, filepath.Base(binary))
	}

	return &ghi, nil
}

func dismantle(url string) (owner, repo string) {
	match := plugin.GitHubRepo.FindStringSubmatch(url)
	owner = match[2]
	repo = match[3]
	return
}

func pickAsset(out io.Writer, release *github.RepositoryRelease, asset int) (ra *github.ReleaseAsset, err error) {
	if len(release.Assets) < 1 {
		err = fmt.Errorf("no assets found on release %q", release.GetName())
		return
	}
	if asset > 0 {
		if len := len(release.Assets); asset >= len {
			return nil, fmt.Errorf("only %v assets found on release %q but you're asking %v (start from zero)", len, release.GetName(), asset)
		}
		ra = &release.Assets[asset]
		return
	}
	verbose.Fprintf(out, "trying to find asset name contains %q from release %q\n", runtime.GOOS, release.GetName())
	if ra = findRuntimeOsAsset(out, release.Assets); ra == nil {
		verbose.Fprintf(out, "%s asset not found, using first asset\n", runtime.GOOS)
		ra = &release.Assets[0]
	}
	return
}

func findRuntimeOsAsset(_ io.Writer, assets []github.ReleaseAsset) *github.ReleaseAsset {
	for _, asset := range assets {
		if strings.Contains(asset.GetName(), runtime.GOOS) {
			return &asset
		}
	}
	return nil
}
