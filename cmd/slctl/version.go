package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/spf13/cobra"
)

const (
	versionHelp = `print slctl version.
`
)

type versionCmd struct {
	full  bool
	check bool
}

func newVersionCmd() *cobra.Command {
	c := versionCmd{}
	cmd := &cobra.Command{
		Use:   "version",
		Short: versionHelp,
		Long:  versionHelp,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVar(&c.full, "full", false, "print full version number and commit hash")
	f.BoolVar(&c.check, "check", false, "check for slctl updates")

	return cmd
}

func (c *versionCmd) run() error {
	if c.full {
		logrus.Printf(metadata.FullString())
	} else {
		logrus.Println(metadata.String())
	}
	if c.check {
		if environment.Settings.Offline {
			return fmt.Errorf("can not check for updates in offline mode")
		}
		return environment.CheckForUpdates(logrus.StandardLogger(), environment.Settings.Home, version, true)
	}
	return nil
}
