package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

const pluginHelp = `Manage {{.}} plugins.
`

func newPluginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "add, list, remove, or create plugins",
		Long:  usage(pluginHelp),
	}
	cmd.AddCommand(
		newPluginListCmd(),
		newPluginInstallCmd(),
		newPluginRemoveCmd(),
		newPluginCreateCmd(),
		newPluginEnvsCmd(),
		newPluginFlagsCmd(),
		newPluginSearchCmd(),
		newPluginUpgradeCmd(),
		newPluginExtsCmd(),
	)
	return cmd
}

func runHook(p *plugin.Plugin) error {
	if err := plugin.SetupPluginEnv(p.Metadata.Name, p.Dir, name, ver().Short()); err != nil {
		return err
	}
	command, err := p.Metadata.Hook.GetCommand()
	if err != nil {
		return err
	}
	main, argv, err := p.PrepareCommand(command, nil)
	if err != nil {
		return err
	}
	prog := exec.Command(main, argv...)
	logrus.Debugf("running hook: %v\n", command)
	prog.Stdout, prog.Stderr = os.Stdout, os.Stderr
	if err := prog.Run(); err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			os.Stderr.Write(e.Stderr)
			return fmt.Errorf("plugin hook for %q exited with error", p.Metadata.Name)
		}
		return err
	}
	return nil
}
