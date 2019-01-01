package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/spf13/cobra"
)

type pluginEnvsCmd struct {
	home paths.Home
}

func newPluginEnvsCmd() *cobra.Command {
	c := &pluginEnvsCmd{}
	cmd := &cobra.Command{
		Use:   "envs",
		Short: "list all environment variables a plugin can get",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c.home = environment.Settings.Home
			return c.run()
		},
	}
	return cmd
}

func (c *pluginEnvsCmd) run() error {
	for k, v := range plugin.Envs {
		logrus.Printf("%s=%s", k, v)
	}
	return nil
}
