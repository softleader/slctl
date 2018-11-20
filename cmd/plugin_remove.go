package cmd

import (
	"errors"
	"fmt"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/softleader/slctl/pkg/v"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type pluginRemoveCmd struct {
	names []string
	home  slpath.Home
	out   io.Writer
}

func newPluginRemoveCmd(out io.Writer) *cobra.Command {
	pcmd := &pluginRemoveCmd{out: out}
	cmd := &cobra.Command{
		Use:   "remove <plugin>...",
		Short: "remove one or more plugins",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return pcmd.complete(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return pcmd.run()
		},
	}
	return cmd
}

func (pcmd *pluginRemoveCmd) complete(args []string) error {
	if len(args) == 0 {
		return errors.New("please provide plugin name to remove")
	}
	pcmd.names = args
	pcmd.home = settings.Home
	return nil
}

func (c *pluginRemoveCmd) run() error {
	v.Println("loading installed plugins from %s", settings.PluginDirs())
	plugins, err := findPlugins(settings.PluginDirs())
	if err != nil {
		return err
	}
	var errorPlugins []string
	for _, name := range c.names {
		if found := findPlugin(plugins, name); found != nil {
			if err := removePlugin(found); err != nil {
				errorPlugins = append(errorPlugins, fmt.Sprintf("Failed to remove plugin %s, got error (%v)", name, err))
			} else {
				fmt.Fprintf(c.out, "Removed plugin: %s\n", name)
			}
		} else {
			errorPlugins = append(errorPlugins, fmt.Sprintf("Plugin: %s not found", name))
		}
	}
	if len(errorPlugins) > 0 {
		return fmt.Errorf(strings.Join(errorPlugins, "\n"))
	}
	return nil
}

func removePlugin(p *plugin.Plugin) error {
	if err := os.RemoveAll(p.Dir); err != nil {
		return err
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
