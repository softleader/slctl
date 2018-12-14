package main

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/spf13/cobra"
	"io"
	"strings"
)

const pluginSearchDesc = `Search SoftLeader official plugin

	$ {{.}} plugin search NAME

NAME 可傳入指定的 Plugin 名稱, 會視為模糊條件來過濾; 反之列出全部

	$ {{.}} plugin search
	$ {{.}} plugin search whereis

查詢的結果將會被 cache 並留存一天
傳入 '--force' 可以先強制更新 cache

	$ {{.}} plugin search -f
`

type pluginSearchCmd struct {
	home  slpath.Home
	out   io.Writer
	name  string
	force bool
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

	f := cmd.Flags()
	f.BoolVarP(&c.force, "force", "f", false, "force to update cache before searching plugins")

	return cmd
}

func (c *pluginSearchCmd) run() (err error) {
	r, err := plugin.LoadRepository(c.out, c.home, organization, c.force)
	if err != nil {
		return err
	}
	if len(r.Repos) == 0 {
		fmt.Fprintln(c.out, "No search results")
		return
	}
	plugins, err := findPlugins(environment.Settings.PluginDirs())
	if err != nil {
		return err
	}
	table := uitable.New()
	table.AddRow("INSTALLED", "NAME", "SOURCE", "DESCRIPTION")
	for _, repo := range r.Repos {
		table.AddRow(installed(plugins, repo.Source), repo.Name, repo.Source, repo.Description)
	}
	fmt.Fprintln(c.out, table)
	return
}

func installed(plugins []*plugin.Plugin, source string) string {
	for _, plugin := range plugins {
		if strings.Contains(plugin.Source, source) {
			return "V"
		}
	}
	return ""
}
