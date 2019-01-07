package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	versionHelp = `print slctl version.`
)

func newVersionCmd() *cobra.Command {
	var full bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: versionHelp,
		Long:  versionHelp,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if full {
				logrus.Printf(metadata.FullString())
			} else {
				logrus.Println(metadata.String())
			}
		},
	}

	f := cmd.Flags()
	f.BoolVar(&full, "full", false, "print full version number and commit hash")

	return cmd
}
