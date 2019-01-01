package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/spf13/cobra"
)

const (
	versionHelp = `print {{.}} version.`
	unreleased  = "unreleased"
	unknown     = "unknown"
)

var (
	version string
	commit  string
)

type Version struct {
	GitVersion string
	GitCommit  string
}

func newVersionCmd() *cobra.Command {
	var full bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: usage(versionHelp),
		Long:  usage(versionHelp),
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if full {
				logrus.Printf(ver().FullString())
			} else {
				logrus.Println(ver().String())
			}
		},
	}

	f := cmd.Flags()
	f.BoolVar(&full, "full", false, "print full version number and commit hash")

	return cmd
}

func (v *Version) FullString() string {
	return fmt.Sprintf("%#v", v)
}

func (v *Version) String() string {
	return fmt.Sprintf("%s+%s", v.GitVersion, v.GitCommit[:7])
}

func ver() *Version {
	if version = strings.TrimSpace(version); version == "" {
		version = unreleased
	}
	if commit = strings.TrimSpace(commit); commit == "" {
		commit = unknown
	}
	return &Version{
		GitVersion: version,
		GitCommit:  commit,
	}
}
