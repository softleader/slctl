package cmd

import (
	"bytes"
	"fmt"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/spf13/cobra"
	"strings"
	"text/template"
)

var (
	settings environment.EnvSettings
)

const (
	Name        = "slctl"
	globalUsage = `{{.|title}} is a command line interface for running commands against SoftLeader services.

To begin working with {{.}}, run the '{{.}} init' command:

	$ {{.}} init

It will set up any necessary local configuration.

Common actions from this point include:

Environment:
  $SL_HOME           set an alternative location for {{.}} files. By default, these are stored in ~/.sl
  $SL_NO_PLUGINS     disable plugins. Set $SL_NO_PLUGINS=true to disable plugins.
  $SL_OFFLINE   	 work offline. Set $SL_OFFLINE=true to work offline.
`
)

func usage(tpl string) string {
	funcMap := template.FuncMap{
		"title": strings.Title,
	}
	var buf bytes.Buffer
	parsed := template.Must(template.New("").Funcs(funcMap).Parse(tpl))
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
		newHomeCmd(out),
		newInitCmd(out),
		newPluginCmd(out),
		newVersionCmd(out),
	)

	flags.Parse(args)

	// set defaults from environment
	settings.Init(flags)

	// Find and add plugins
	loadPlugins(cmd, out)

	return cmd
}

func checkArgsLength(argsReceived int, requiredArgs ...string) error {
	expectedNum := len(requiredArgs)
	if argsReceived != expectedNum {
		arg := "arguments"
		if expectedNum == 1 {
			arg = "argument"
		}
		return fmt.Errorf("this command needs %v %s: %s", expectedNum, arg, strings.Join(requiredArgs, ", "))
	}
	return nil
}
