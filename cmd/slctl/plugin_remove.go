package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/spf13/cobra"
)

type pluginRemoveCmd struct {
	names []string
	home  paths.Home
	force bool
}

func newPluginRemoveCmd() *cobra.Command {
	pcmd := &pluginRemoveCmd{}
	cmd := &cobra.Command{
		Use:   "remove <plugin>...",
		Short: "remove one or more plugins",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pcmd.names = args
			pcmd.home = environment.Settings.Home
			return pcmd.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&pcmd.force, "force", "f", false, "do not display a diagnostic message or modify the exit status to reflect an error if plugin does not exist")

	return cmd
}

func (c *pluginRemoveCmd) run() error {
	logrus.Debugf("loading installed plugins from %s\n", environment.Settings.PluginDirs())
	plugins, err := plugin.LoadPaths(environment.Settings.PluginDirs())
	if err != nil {
		return err
	}
	var errorPlugins []string
	for _, name := range c.names {
		if found := findPlugin(plugins, name); found != nil {
			if err := os.RemoveAll(found.Dir); err != nil {
				errorPlugins = append(errorPlugins, fmt.Sprintf("Failed to remove plugin %s, got error (%v)", name, err))
			} else {
				logrus.Printf("Removed plugin: %s\n", name)
			}
		} else if !c.force {
			errorPlugins = append(errorPlugins, fmt.Sprintf("Plugin: %s not found", name))
		}
	}
	if len(errorPlugins) > 0 {
		return fmt.Errorf(strings.Join(errorPlugins, "\n"))
	}
	return nil
}

func findPlugin(plugins []*plugin.Plugin, name string) *plugin.Plugin {
	for _, p := range plugins {
		if p.Metadata.Name == name {
			return p
		}
	}
	return nil
}
