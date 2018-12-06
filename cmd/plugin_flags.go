package main

import (
	"fmt"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"io"

	"github.com/spf13/cobra"
)

const pluginFlagsDesc = "Global flags will passed to plugin only if 'ignoreGlobalFlags' in " + plugin.MetadataFileName + " of the plugin is true"

type pluginFlagsCmd struct {
	home slpath.Home
	out  io.Writer
}

func newPluginFlagsCmd(out io.Writer) *cobra.Command {
	c := &pluginFlagsCmd{out: out}
	cmd := &cobra.Command{
		Use:   "flags",
		Short: "list all global flags",
		Long:  usage(pluginFlagsDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			c.home = environment.Settings.Home
			return c.run()
		},
	}
	return cmd
}

func (c *pluginFlagsCmd) run() error {
	for _, flag := range environment.Flags {
		fmt.Fprintln(c.out, flag)
	}
	return nil
}
