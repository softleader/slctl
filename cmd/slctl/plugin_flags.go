package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/spf13/cobra"
)

const pluginFlagsDesc = "Global flags will passed to plugin only if 'ignoreGlobalFlags' in " + plugin.MetadataFileName + " of the plugin is true"

type pluginFlagsCmd struct {
	home paths.Home
}

func newPluginFlagsCmd() *cobra.Command {
	c := &pluginFlagsCmd{}
	cmd := &cobra.Command{
		Use:   "flags",
		Short: "list all global flags",
		Long:  usage(pluginFlagsDesc),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c.home = environment.Settings.Home
			return c.run()
		},
	}
	return cmd
}

func (c *pluginFlagsCmd) run() error {
	for _, flag := range environment.Flags {
		logrus.Println(flag)
	}
	return nil
}
