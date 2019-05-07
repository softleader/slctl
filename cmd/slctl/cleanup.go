package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/spf13/cobra"
)

const (
	cleanupDesc = `Remove outdated downloads of plugin
`
)

type cleanupCmd struct {
	home   paths.Home
	dryRun bool
}

func newCleanupCmd() *cobra.Command {
	c := &cleanupCmd{}

	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "remove outdated downloads of plugin",
		Long:  cleanupDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c.home = environment.Settings.Home
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVar(&c.dryRun, "dry-run", false, "show what would be removed, but do not actually remove anything.")

	return cmd
}

func (c *cleanupCmd) run() (err error) {
	return plugin.Cleanup(logrus.StandardLogger(), c.home, true, c.dryRun)
}
