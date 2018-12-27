package plugin

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/gosuri/uitable"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/slpath"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Repo struct {
	Name        string
	Source      string
	Description string
}

type Repository struct {
	Repos   []Repo
	Expires time.Time
}

func LoadRepository(log *logrus.Logger, home slpath.Home, org string, force bool) (r *Repository, err error) {
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

func fetchOnline(log *logrus.Logger, home slpath.Home, org string) (r *Repository, err error) {
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
	repos, _, err := client.Repositories.ListByOrg(ctx, org, &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 999999},
	})
	if err != nil {
		return
	}
	r = &Repository{}
	r.Expires = time.Now().AddDate(0, 0, 1)
	for _, repo := range repos {
		if name := repo.GetName(); strings.HasPrefix(name, "slctl-") {
			source := fmt.Sprintf("github.com/%s", repo.GetFullName())
			r.Repos = append(r.Repos, Repo{
				Name:        name,
				Source:      source,
				Description: repo.GetDescription(),
			})
		}
	}
	log.Debugf("retrieved %v plugins\n", len(r.Repos))
	if environment.Settings.Verbose {
		table := uitable.New()
		for _, r := range r.Repos {
			table.AddRow(r.Name, r.Source)
		}
		log.Println(table)
	}
	return
}
