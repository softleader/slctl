package main

import (
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/spf13/cobra"
	"path/filepath"
)

const longHomeHelp = `home displays the location of SL_HOME.

傳入 '--move' 可以將當前的 home 搬移到指定目錄

	$ slctl home --move /some/where/else

可以配合 '--home' 搬移指定 home 目錄

	$ slctl home --home /another/home/path --move /some/where/else
`

type homeCmd struct {
	home paths.Home
	move string
}

func newHomeCmd() *cobra.Command {
	c := &homeCmd{}
	cmd := &cobra.Command{
		Use:   "home",
		Short: "displays the location of SL_HOME",
		Long:  longHomeHelp,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c.home = environment.Settings.Home
			return c.run()
		},
	}

	f := cmd.Flags()
	f.StringVar(&c.move, "move", "", "moves home to specified path")

	return cmd
}

func (c *homeCmd) run() (err error) {
	if len(c.move) == 0 {
		logrus.Println(c.home)
		logrus.Debugf("Config: %s\n", c.home.Config())
		logrus.Debugf("ConfigFile: %s\n", c.home.ConfigFile())
		logrus.Debugf("Plugins: %s\n", c.home.Plugins())
	} else {
		if c.move, err = homedir.Expand(c.move); err != nil {
			return
		}
		c.move, err = filepath.Abs(c.move)
		if err != nil {
			return
		}
		if err = environment.MoveHome(c.home.String(), c.move); err != nil {
			return
		}
		logrus.Printf("Successfully moved home from %q to %q", c.home.String(), c.move)
	}
	return nil
}
