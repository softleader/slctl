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
	)
	return cmd
}

// runHook will execute a plugin hook.
func runHook(p *plugin.Plugin, event string) error {
	hook := p.Metadata.Hooks.Get(event)
	if hook == "" {
		return nil
	}

	prog := exec.Command("sh", "-c", hook)
	// TODO make this work on windows
	// I think its ... ¯\_(ツ)_/¯
	// prog := exec.Command("cmd", "/C", p.Metadata.Hooks.Install())

	v.Printf("running %s hook: %v\n", event, prog)

	if err := plugin.SetupPluginEnv(p.Metadata.Name, p.Dir); err != nil {
		return err
	}
	prog.Stdout, prog.Stderr = os.Stdout, os.Stderr
	if err := prog.Run(); err != nil {
		if eerr, ok := err.(*exec.ExitError); ok {
			os.Stderr.Write(eerr.Stderr)
			return fmt.Errorf("plugin %s hook for %q exited with error", event, p.Metadata.Name)
		}
		return err
	}
	return nil
}
