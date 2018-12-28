package main

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/spf13/cobra"
	"strings"
)

const pluginSearchDesc = `Search SoftLeader official plugin

	$ {{.}} plugin search <FILTER...>

NAME 可傳入多個 FILTERS, 會以 Or 且模糊條件來過濾 SOURCE; 反之列出全部

	$ {{.}} plugin search
	$ {{.}} plugin search whereis

傳入 '--installed' 只列出已安裝的 Plugin

	$ {{.}} plugin search -i

查詢的結果將會被 cache 並留存一天, 傳入 '--force' 在查詢前強制更新 cache

	$ {{.}} plugin search -f
`

type pluginSearchCmd struct {
	home              slpath.Home
	filters           []string
	onlyShowInstalled bool
	force             bool
}

func newPluginSearchCmd() *cobra.Command {
	c := &pluginSearchCmd{}
	cmd := &cobra.Command{
		Use:   "search <FILTER...>",
		Short: "search SoftLeader official plugin",
		Long:  usage(pluginSearchDesc),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if environment.Settings.Offline {
				return fmt.Errorf("can not run the command in offline mode")
			}
			if len(args) > 0 {
				c.filters = args
			}
			c.home = environment.Settings.Home
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&c.force, "force", "f", false, "force to update cache before searching plugins")
	f.BoolVarP(&c.onlyShowInstalled, "installed", "i", false, "only shows installed plugins")

	return cmd
}

func (c *pluginSearchCmd) run() (err error) {
	r, err := plugin.LoadRepository(logrus.StandardLogger(), c.home, organization, c.force)
	if err != nil {
		return err
	}
	if len(r.Repos) == 0 {
		logrus.Println("No search results")
		return
	}
	plugins, err := findPlugins(environment.Settings.PluginDirs())
	if err != nil {
		return err
	}
	table := uitable.New()
	table.AddRow("INSTALLED", "SOURCE", "DESCRIPTION")
	for _, repo := range r.Repos {
		hasInstalled := installed(plugins, repo.Source)
		if c.onlyShowInstalled && !hasInstalled { // 要求只顯示安裝過的
			continue
		}
		if len(c.filters) > 0 && !repo.Contains(c.filters...) {
			continue
		}
		var installed string
		if hasInstalled {
			installed = "V"
		}
		table.AddRow(installed, repo.Source, repo.Description)
	}
	logrus.Println(table)
	return
}

func installed(plugins []*plugin.Plugin, source string) bool {
	for _, plugin := range plugins {
		if strings.Contains(plugin.Source, source) {
			return true
		}
	}
	return false
}
