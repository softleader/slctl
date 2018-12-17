package main

import (
	"fmt"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"io"

	"github.com/spf13/cobra"
)

type pluginEnvsCmd struct {
	home slpath.Home
	out  io.Writer
}

func newPluginEnvsCmd(out io.Writer) *cobra.Command {
	c := &pluginEnvsCmd{out: out}
	cmd := &cobra.Command{
		Use:   "envs",
		Short: "list all environment variables a plugin can get",
		RunE: func(cmd *cobra.Command, args []string) error {
			c.home = environment.Settings.Home
			return c.run()
		},
	}
	return cmd
}

func (c *pluginEnvsCmd) run() error {
	for _, env := range plugin.Envs {
		fmt.Fprintln(c.out, env)
	}
	return nil
}
