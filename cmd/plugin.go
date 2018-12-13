package main

import (
	"fmt"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/v"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
)

const pluginHelp = `
Manage {{.}} plugins.

'plugin install' command is not supported for now! 
Please manually drop plugin folder into $SL_PLUGIN (default $SL_HOME/plugins).
`

func newPluginCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "add, list, remove, or create plugins",
		Long:  usage(pluginHelp),
	}
	cmd.AddCommand(
		newPluginListCmd(out),
		newPluginInstallCmd(out),
		newPluginRemoveCmd(out),
		newPluginCreateCmd(out),
		newPluginEnvsCmd(out),
		newPluginFlagsCmd(out),
	)
	return cmd
}

func runHook(p *plugin.Plugin) error {
	if err := plugin.SetupPluginEnv(p.Metadata.Name, p.Dir, name, version); err != nil {
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
	v.Printf("running hook: %v\n", command)
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
