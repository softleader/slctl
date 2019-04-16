package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/spf13/cobra"
)

const longHomeHelp = `
This command displays the location of SL_HOME.
`

func newHomeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "home",
		Short: "displays the location of SL_HOME",
		Long:  longHomeHelp,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			h := environment.Settings.Home
			logrus.Println(h)
			logrus.Debugf("Config: %s\n", h.Config())
			logrus.Debugf("ConfigFile: %s\n", h.ConfigFile())
			logrus.Debugf("Plugins: %s\n", h.Plugins())
		},
	}
	return cmd
}
