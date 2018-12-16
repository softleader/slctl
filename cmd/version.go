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
	var v, c, d string
	if v = strings.TrimSpace(version); v == "" {
		v = unreleased
	}
	if c = strings.TrimSpace(commit); c == "" {
		c = none
	}
	if d = strings.TrimSpace(date); d == "" {
		d = unknown
	}
	return fmt.Sprintf("%v, commit %v, built at %v", version, commit, date)
}
