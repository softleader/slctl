package cmd

import (
	"bytes"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/spf13/cobra"
	"html/template"
	"strings"
)

var (
	settings environment.EnvSettings
)

const (
	Name        = "slctl"
	globalUsage = `{{.|title}}] is a command line interface for running commands against SoftLeader services.

To begin working with {{.}}, run the '{{.}} init' command:

	$ {{.}} init

It will set up any necessary local configuration.

Common actions from this point include:

Environment:
  $SL_HOME           set an alternative location for {{.}} files. By default, these are stored in ~/.sl
`
)

func usage(tmpl string) string {
	funcMap := template.FuncMap{
		"title": strings.Title,
	}
	var buf bytes.Buffer
	parsed := template.Must(template.New("").Funcs(funcMap).Parse(tmpl))
	err := parsed.Execute(&buf, Name)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func NewRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          Name,
		Short:        Name + " against SoftLeader services.",
		Long:         usage(globalUsage),
		SilenceUsage: true,
	}
	flags := cmd.PersistentFlags()

	settings.AddFlags(flags)

	out := cmd.OutOrStdout()

	cmd.AddCommand(
		newInitCmd(out),
	)

	flags.Parse(args)

	// set defaults from environment
	settings.Init(flags)

	// Find and add plugins
	loadPlugins(cmd, out)

	return cmd
}
