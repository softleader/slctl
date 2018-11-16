package main

import (
	"bytes"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/spf13/cobra"
	"html/template"
	"os"
)

var (
	settings environment.EnvSettings
)

const (
	Name        = "slctl"
	globalUsage = `{{.}} controls SoftLeader services.

To begin working with {{.}}, run the '{{.}} init' command:

	$ {{.}} init

It will set up any necessary local configuration.

Common actions from this point include:

Environment:
  $SL_HOME           set an alternative location for {{.}} files.By default, these are stored in ~/.sl
`
)

func main() {
	cmd := newRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		switch e := err.(type) {
		case pluginError:
			os.Exit(e.code)
		default:
			os.Exit(1)
		}
	}
}

func usage() string {
	var buf bytes.Buffer
	t := template.Must(template.New("").Parse(globalUsage))
	err := t.Execute(&buf, Name)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func newRootCmd(args []string) *cobra.Command {

	cmd := &cobra.Command{
		Use:          Name,
		Short:        Name + " controls SoftLeader services.",
		Long:         usage(),
		SilenceUsage: true,
	}
	flags := cmd.PersistentFlags()

	settings.AddFlags(flags)

	out := cmd.OutOrStdout()

	cmd.AddCommand()

	flags.Parse(args)

	// set defaults from environment
	settings.Init(flags)

	// Find and add plugins
	loadPlugins(cmd, out)

	return cmd
}
