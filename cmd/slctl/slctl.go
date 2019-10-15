package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/formatter"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/release"
	"github.com/spf13/cobra"
	"os"
)

const (
	globalUsage = `Slctl is a command line interface for running commands against SoftLeader services.

To begin working with slctl, run the 'slctl init' command:

	$ slctl init

It will set up any necessary local configuration.

Common actions from this point include:

Environment:
  $SL_HOME           set an alternative location for slctl files. By default, these are stored in ~/.config/slctl
  $SL_PLUGINS_OFF    disable plugins. Set $SL_PLUGINS_OFF=true to disable plugins.
  $SL_OFFLINE        work offline. Set $SL_OFFLINE=true to work offline.
`
)

const (
	organization = "softleader"
)

var (
	version, commit string
	metadata        *release.Metadata
)

func main() {
	initMetadata()
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
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			environment.CheckForUpdates(logrus.StandardLogger(), environment.Settings.Home, version, false)
			plugin.Cleanup(logrus.StandardLogger(), environment.Settings.Home, false, false)
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
		newCompletionCmd(),
		newCleanupCmd(),
	)

	flags.Parse(args)
	environment.Settings.ExpandEnvToFlags(flags)

	plugCommands, err := plugin.LoadPluginCommands(metadata)
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(plugCommands...)

	return cmd, nil
}

// initMetadata 準備 app 的 release 資訊
func initMetadata() {
	metadata = release.NewMetadata(version, commit)
}
