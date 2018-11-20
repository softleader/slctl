package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var longHomeHelp = `
This command displays the location of SL_HOME.
`

func newHomeCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "home",
		Short: "displays the location of SL_HOME",
		Long:  usage(longHomeHelp),
		Run: func(cmd *cobra.Command, args []string) {
			h := settings.Home
			fmt.Fprintln(out, h)
			if settings.Verbose {
				fmt.Fprintf(out, "Config: %s\n", h.Config())
				fmt.Fprintf(out, "ConfigFile: %s\n", h.ConfigFile())
				fmt.Fprintf(out, "Plugins: %s\n", h.Plugins())
			}
		},
	}
	return cmd
}
