package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

const pluginHelp = `
Manage {{.}} plugins.

'plugin install' command is not supported for now! please manually drop plugin into $SL_HOME/plugins.
`

func newPluginCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "list, or remove plugins",
		Long:  usage(pluginHelp),
	}
	cmd.AddCommand(
		newPluginListCmd(out),
		newPluginRemoveCmd(out),
	)
	return cmd
}
