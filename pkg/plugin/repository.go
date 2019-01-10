package plugin

import (
	"context"
	"fmt"
	"github.com/google/go-github/v21/github"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	officialPluginTopic = "slctl-plugin"
)

// Repo 代表 Plugin 的 GitHub-Repo 的來源跟描述
type Repo struct {
	Source      string
	Description string
}

// Contains 回傳此 Repo 是否包含任何傳入的 filters 文字
func (r *Repo) Contains(filters ...string) bool {
	for _, filter := range filters {
		if strings.Contains(r.Source, filter) || strings.Contains(r.Description, filter) {
			return true
		}
	}
	return false
}

// Repository 代表 GitHub-Repo plugin 的匯總
type Repository struct {
	Repos   []Repo
	Expires time.Time
}

// LoadRepository 載入 Repository
func LoadRepository(log *logrus.Logger, home paths.Home, org string, force bool) (r *Repository, err error) {
	cached := home.CacheRepositoryFile()
	if force {
		if r, err = fetchOnline(log, home, org); err == nil {
			r.save(cached)
		}
		return
	}
	r, err = loadLocal(log, cached)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if expired(r) {
		log.Debugln("cache is out of date")
		if r, err = fetchOnline(log, home, org); err == nil {
			r.save(cached)
		}
	}
	return
}

func expired(r *Repository) bool {
	return r == nil || r.Expires.Before(time.Now())
}

func (r *Repository) save(path string) error {
	data, err := yaml.Marshal(r)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

func loadLocal(log *logrus.Logger, path string) (r *Repository, err error) {
	log.Debugf("loading cached plugin repositories from: %s\n", path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	r = &Repository{}
	err = yaml.Unmarshal(data, r)
	return
}

func fetchOnline(log *logrus.Logger, home paths.Home, org string) (r *Repository, err error) {
	if environment.Settings.Offline {
		return nil, fmt.Errorf("can not fetch plugin repository in offline mode")
	}
	log.Debugf("fetching the plugin repositories\n")
	cfg, err := config.LoadConfFile(home.ConfigFile())
	if err != nil {
		return
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	query := fmt.Sprintf("org:%s+topic:%s", org, officialPluginTopic)
	log.Debugf("specifying searching qualifiers: %s\n", query)
	var allRepos []github.Repository
	opt := &github.SearchOptions{}
	for {
		result, resp, err := client.Search.Repositories(ctx, query, opt)
		if err != nil {
			break
		}
		allRepos = append(allRepos, result.Repositories...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	if err != nil {
		return
	}
	r = &Repository{}
	r.Expires = time.Now().AddDate(0, 0, 1)
	for _, repo := range allRepos {
		source := fmt.Sprintf("github.com/%s", repo.GetFullName())
		r.Repos = append(r.Repos, Repo{
			Source:      source,
			Description: repo.GetDescription(),
		})
	}
	log.Debugf("retrieved %v plugins\n", len(r.Repos))
	return
}
