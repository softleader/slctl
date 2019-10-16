package main

import (
	"fmt"
	"strings"

	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/spf13/cobra"
)

const pluginOpenDesc = `Open plugins

當 Plugin Source 來自於 GitHub, 則以預設瀏覽器開啟 GitHub Repo 網址; 反之開啟 Plugin 所在目錄

	$ slctl plugin open PLUGIN

傳入 '--app' 使用指定 app 名稱來開啟非 GitHub Source 的 Plugin

	$ slctl plugin open PLUGIN -a "Sublime Text"
`

type pluginOpenCmd struct {
	home   paths.Home
	plugin string
	app    string
}

func newPluginOpenCmd() *cobra.Command {
	c := &pluginOpenCmd{}
	cmd := &cobra.Command{
		Use:   "open <PLUGIN>",
		Short: "open plugin",
		Long:  pluginOpenDesc,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if environment.Settings.Offline {
				return fmt.Errorf("can not run the command in offline mode")
			}
			c.plugin = args[0]
			c.home = environment.Settings.Home
			return c.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&c.app, "app", "a", "", "opens plugin with the specified application")
	return cmd
}

func (c *pluginOpenCmd) run() (err error) {
	plugins, err := plugin.LoadPaths(environment.Settings.PluginDirs())
	if err != nil {
		return err
	}
	for _, p := range plugins {
		if strings.EqualFold(p.Metadata.Name, c.plugin) {
			return p.Open(c.app)
		}
	}
	return fmt.Errorf("no plugin named %q found", c.plugin)
}
