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
	date    string
)

type Version struct {
	GitVersion string
	GitCommit  string
	BuildDate  string
}

func newVersionCmd() *cobra.Command {
	var short bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: usage(versionHelp),
		Long:  usage(versionHelp),
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if short {
				logrus.Println(ver().Short())
			} else {
				logrus.Printf("%#v\n", ver())
			}
		},
	}

	f := cmd.Flags()
	f.BoolVar(&short, "short", false, "print only the version number plus first 7 digits of the commit hash")

	return cmd
}

func (v *Version) Short() string {
	return fmt.Sprintf("%s+%s", v.GitVersion, v.GitCommit[:7])
}

func ver() *Version {
	if version = strings.TrimSpace(version); version == "" {
		version = unreleased
	}
	if commit = strings.TrimSpace(commit); commit == "" {
		commit = unknown
	}
	if date = strings.TrimSpace(date); date == "" {
		date = unknown
	}
	return &Version{
		GitVersion: version,
		GitCommit:  commit,
		BuildDate:  date,
	}
}
