package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/spf13/cobra"
	"os"
)

type pluginUnmountCmd struct {
	name []string
}

const pluginUnmountDesc = `To unmount a plugin volume

將 Plugin 的 Mount Volume 完整移除 
For more details: https://github.com/softleader/slctl/wiki/Plugins-Guide#mount-volume
`

func newPluginUnmountCmd() *cobra.Command {
	c := &pluginUnmountCmd{}
	cmd := &cobra.Command{
		Use:   "unmount PLUGIN_NAME....",
		Short: "unmount one or more plugins",
		Long:  usage(pluginUnmountDesc),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c.name = args
			return c.run()
		},
	}

	return cmd
}

func (c *pluginUnmountCmd) run() error {
	plugs, err := findPlugins(environment.Settings.PluginDirs())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load plugins: %s", err)
		return err
	}
	for _, name := range c.name {
		if p, found := pick(plugs, name); found {
			if err := p.Unmount(); err != nil {
				return err
			}
		} else {
			logrus.Debugf("Skip unmounting %q, it is not a installed plugin", name)
		}
	}
	return nil
}

func pick(plugs []*plugin.Plugin, name string) (*plugin.Plugin, bool) {
	for _, p := range plugs {
		if p.Metadata.Name == name {
			return p, true
		}
	}
	return nil, false
}
