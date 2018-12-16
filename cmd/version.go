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
	none        = "none"
	unknown     = "unknown"
)

var (
	version string
	commit  string
	date    string
)

type BuildMetadata struct {
	GitVersion string
	GitCommit  string
	BuildDate  string
}

func newVersionCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: usage(versionHelp),
		Long:  usage(versionHelp),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(out, "%#v\n", buildMetadata())
		},
	}
	return cmd
}

func buildMetadata() BuildMetadata {
	if version = strings.TrimSpace(version); version == "" {
		version = unreleased
	}
	if commit = strings.TrimSpace(commit); commit == "" {
		commit = none
	}
	if date = strings.TrimSpace(date); date == "" {
		date = unknown
	}
	return BuildMetadata{
		GitVersion: version,
		GitCommit:  commit,
		BuildDate:  date,
	}
}
