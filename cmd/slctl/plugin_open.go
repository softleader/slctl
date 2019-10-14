package main

import (
	"fmt"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/spf13/cobra"
	"strings"
)

const pluginOpenDesc = `Open plugin

當 Plugin Source 來自於 GitHub, 則以預設瀏覽器開啟 GitHub Repo 網址; 反之開啟 plugin 所在目錄

	$ slctl plugin open PLUGIN

傳入 '--app' 使用指定 app 名稱來開啟 plugin source

	$ slctl plugin open PLUGIN --app firefox

傳入 '--wait' 等待 open command 執行完畢才結束

	$ slctl plugin open PLUGIN --app firefox -w
`

type pluginOpenCmd struct {
	home   paths.Home
	plugin string
	app    string
	wait   bool
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
	f.StringVar(&c.app, "app", "", "open the plugin using the specified application")
	f.BoolVarP(&c.wait, "wait", "w", false, "wait for the open command to complete")
	return cmd
}

func (c *pluginOpenCmd) run() (err error) {
	plugins, err := plugin.LoadPaths(environment.Settings.PluginDirs())
	if err != nil {
		return err
	}
	for _, p := range plugins {
		if strings.EqualFold(p.Metadata.Name, c.plugin) {
			if c.wait {
				return p.OpenAndWait(c.app)
			}
			return p.Open(c.app)
		}
	}
	return fmt.Errorf("no plugin named %q found", c.plugin)
}
