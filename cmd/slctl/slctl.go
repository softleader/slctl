package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/formatter"
	"github.com/softleader/slctl/pkg/plugin"
	ver "github.com/softleader/slctl/pkg/version"
	"github.com/spf13/cobra"
	"os"
)

const (
	globalUsage = `{{.|title}} is a command line interface for running commands against SoftLeader services.

To begin working with {{.}}, run the '{{.}} init' command:

	$ slctl init

It will set up any necessary local configuration.

Common actions from this point include:

Environment:
  $SL_HOME           set an alternative location for {{.}} files. By default, these are stored in ~/.sl
  $SL_NO_PLUGINS     disable plugins. Set $SL_NO_PLUGINS=true to disable plugins.
  $SL_OFFLINE        work offline. Set $SL_OFFLINE=true to work offline.
`
)

var (
	organization = "softleader"
	version      string
	commit       string
	metadata     *ver.BuildMetadata
)

func main() {
	metadata = ver.NewBuildMetadata(version, commit)
	if cmd, err := newRootCmd(os.Args[1:]); err != nil {
		exit(err)
	} else if err = cmd.Execute(); err != nil {
		exit(err)
	}
}

func exit(err error) {
	switch e := err.(type) {
	case plugin.ExitError:
		os.Exit(e.ExitStatus)
	default:
		os.Exit(1)
	}
}

func newRootCmd(args []string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:          "slctl",
		Short:        "slctl against SoftLeader services.",
		Long:         globalUsage,
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

	environment.Settings.AddGlobalFlags(flags)

	plugCommands, err := plugin.LoadPluginCommands(environment.Settings, metadata)
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(plugCommands...)

	return cmd, nil
}

//
//func usage(tpl string) string {
//	funcMap := template.FuncMap{
//		"title": strings.Title,
//	}
//	var buf bytes.Buffer
//	parsed := template.Must(template.New("").Funcs(funcMap).Parse(tpl))
//	err := parsed.Execute(&buf, name)
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//	return buf.String()
//}
