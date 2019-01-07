package main

import (
	"github.com/gosuri/uitable"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/spf13/cobra"
)

type pluginListCmd struct {
	home paths.Home
}

func newPluginListCmd() *cobra.Command {
	pcmd := &pluginListCmd{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list installed plugins",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			pcmd.home = environment.Settings.Home
			return pcmd.run()
		},
	}
	return cmd
}

func (c *pluginListCmd) run() error {
	logrus.Debugf("search in plugin dirs: %s", environment.Settings.PluginDirs())
	plugins, err := plugin.LoadPaths(environment.Settings.PluginDirs())
	if err != nil {
		return err
	}

	table := uitable.New()
	table.AddRow("NAME", "VERSION", "DESCRIPTION", "SOURCE")
	for _, p := range plugins {
		table.AddRow(p.Metadata.Name, p.Metadata.Version, p.Metadata.Description, p.Source)
	}
	logrus.Println(table)
	return nil
}
