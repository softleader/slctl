package main

import (
	"fmt"
	"io"
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

func newVersionCmd(out io.Writer) *cobra.Command {
	var short bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: usage(versionHelp),
		Long:  usage(versionHelp),
		Run: func(cmd *cobra.Command, args []string) {
			if short {
				fmt.Fprintln(out, formatVersion(ver(), short))
			} else {
				fmt.Fprintf(out, "%#v\n", ver())
			}
		},
	}

	f := cmd.Flags()
	f.BoolVar(&short, "short", false, "print the version number plus first 7 digits of the commit hash")

	return cmd
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

func formatVersion(v *Version, short bool) string {
	if short && v.GitCommit != "" {
		return fmt.Sprintf("%s+%s", v.GitCommit, v.GitCommit[:7])
	}
	return fmt.Sprintf("%#v", v)
}
