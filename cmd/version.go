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

func newVersionCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: usage(versionHelp),
		Long:  usage(versionHelp),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(out, ver())
		},
	}
	return cmd
}

func ver() string {
	if version = strings.TrimSpace(version); version == "" {
		version = unreleased
	}
	if commit = strings.TrimSpace(commit); commit == "" {
		commit = none
	}
	if date = strings.TrimSpace(date); date == "" {
		date = unknown
	}
	return fmt.Sprintf("%v, commit %v, built at %v", version, commit, date)
}
