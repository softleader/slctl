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
)

var version string


func newVersionCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: usage(versionHelp),
		Long:  usage(versionHelp),
		Run: func(cmd *cobra.Command, args []string) {
			if v := strings.TrimSpace(version); v != "" {
				fmt.Fprintln(out, v)
			} else {
				fmt.Fprintln(out, unreleased)
			}
		},
	}
	return cmd
}