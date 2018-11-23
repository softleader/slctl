package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

const (
	versionHelp = `print {{.}} version.`
	version     = "unreleased"
)

func newVersionCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: usage(versionHelp),
		Long:  usage(versionHelp),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(out, version)
		},
	}
	return cmd
}
