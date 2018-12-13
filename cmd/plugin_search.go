package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/gosuri/uitable"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"io"
	"strings"
)

const pluginSearchDesc = `Search SoftLeader official plugin

	$ slctl plugin search NAME

NAME 可傳入指定的 Plugin 名稱, 會視為模糊條件來過濾; 反之列出全部

	$ slctl plugin search
	$ slctl plugin search whereis
`

type pluginSearchCmd struct {
	home slpath.Home
	out  io.Writer
	name string
}

func newPluginSearchCmd(out io.Writer) *cobra.Command {
	c := &pluginSearchCmd{out: out}
	cmd := &cobra.Command{
		Use:   "search NAME",
		Short: "search SoftLeader official plugin",
		Long:  usage(pluginSearchDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			if environment.Settings.Offline {
				return fmt.Errorf("can not run the command in offline mode")
			}
			if len(args) > 0 {
				c.name = strings.TrimSpace(args[0])
			}
			c.home = environment.Settings.Home
			return c.run()
		},
	}
	return cmd
}

func (c *pluginSearchCmd) run() (err error) {
	cfg, err := config.LoadConfFile(c.home.ConfigFile())
	if err != nil {
		return
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	repos, _, err := client.Repositories.ListByOrg(ctx, organization, &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 999999},
	})
	if err != nil {
		return
	}
	if len(repos) == 0 {
		fmt.Fprintln(c.out, "No search results")
		return
	}
	plugins, err := findPlugins(environment.Settings.PluginDirs())
	if err != nil {
		return err
	}
	table := uitable.New()
	table.AddRow("INSTALLED", "NAME", "SOURCE", "DESCRIPTION")
	for _, repo := range repos {
		if name := repo.GetName(); strings.HasPrefix(name, "slctl-") {
			if c.name != "" && !strings.Contains(name, c.name) {
				continue
			}
			source := fmt.Sprintf("github.com/%s", repo.GetFullName())
			table.AddRow(installed(plugins, source), name, source, repo.GetDescription())
		}

	}
	fmt.Fprintln(c.out, table)
	return
}

func installed(plugins []*plugin.Plugin, source string) string {
	for _, plugin := range plugins {
		if plugin.Source == source {
			return "V"
		}
	}
	return ""
}
