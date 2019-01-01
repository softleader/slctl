package main

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/formatter"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/template"
)

const (
	name        = "slctl"
	globalUsage = `{{.|title}} is a command line interface for running commands against SoftLeader services.

To begin working with {{.}}, run the '{{.}} init' command:

	$ {{.}} init

It will set up any necessary local configuration.

Common actions from this point include:

Environment:
  $SL_HOME           set an alternative location for {{.}} files. By default, these are stored in ~/.sl
  $SL_NO_PLUGINS     disable plugins. Set $SL_NO_PLUGINS=true to disable plugins.
  $SL_OFFLINE        work offline. Set $SL_OFFLINE=true to work offline.
`
)

var organization = "softleader"

func main() {
	if cmd, err := newRootCmd(os.Args[1:]); err != nil {
		exit(err)
	} else if err = cmd.Execute(); err != nil {
		exit(err)
	}
}

func exit(err error) {
	switch e := err.(type) {
	case PluginError:
		os.Exit(e.Code)
	default:
		os.Exit(1)
	}
}

func newRootCmd(args []string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:          name,
		Short:        name + " against SoftLeader services.",
		Long:         usage(globalUsage),
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logrus.SetOutput(cmd.OutOrStdout())
			logrus.SetFormatter(&formatter.PlainFormatter{})
			if environment.Settings.Verbose {
				logrus.SetLevel(logrus.DebugLevel)
			}
		},
	}
	flags := cmd.PersistentFlags()

	if err := environment.Settings.AddFlags(flags); err != nil {
		return nil, err
	}

	cmd.AddCommand(
		newHomeCmd(),
		newInitCmd(),
		newPluginCmd(),
		newVersionCmd(),
	)

	flags.Parse(args)

	environment.Settings.Init(flags)

	loadPlugins(cmd)

	return cmd, nil
}

func usage(tpl string) string {
	funcMap := template.FuncMap{
		"title": strings.Title,
	}
	var buf bytes.Buffer
	parsed := template.Must(template.New("").Funcs(funcMap).Parse(tpl))
	err := parsed.Execute(&buf, name)
	if err != nil {
		panic(err)
	}
	return buf.String()
}
