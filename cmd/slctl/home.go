package main

import (
	"fmt"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/verbose"
	"io"

	"github.com/spf13/cobra"
)

const longHomeHelp = `This command displays the location of SL_HOME.
`

func newHomeCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "home",
		Short: "displays the location of SL_HOME",
		Long:  usage(longHomeHelp),
		Run: func(cmd *cobra.Command, args []string) {
			h := environment.Settings.Home
			fmt.Fprintln(out, h)
			verbose.Fprintf(out, "Config: %s\n", h.Config())
			verbose.Fprintf(out, "ConfigFile: %s\n", h.ConfigFile())
			verbose.Fprintf(out, "Plugins: %s\n", h.Plugins())
		},
	}
	return cmd
}
