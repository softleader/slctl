package cmd

import (
	"fmt"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/softleader/slctl/pkg/v"
	"io"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

type pluginListCmd struct {
	home slpath.Home
	out  io.Writer
}

func newPluginListCmd(out io.Writer) *cobra.Command {
	pcmd := &pluginListCmd{out: out}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list installed plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			pcmd.home = settings.Home
			return pcmd.run()
		},
	}
	return cmd
}

func (c *pluginListCmd) run() error {
	v.Println("search in plugin dirs: %s", settings.PluginDirs())
	plugins, err := findPlugins(settings.PluginDirs())
	if err != nil {
		return err
	}

	table := uitable.New()
	table.AddRow("NAME", "VERSION", "DESCRIPTION")
	for _, p := range plugins {
		table.AddRow(p.Metadata.Name, p.Metadata.Version, p.Metadata.Description)
	}
	fmt.Fprintln(c.out, table)
	return nil
}
