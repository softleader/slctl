package main

import (
	"fmt"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/plugin/installer"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/softleader/slctl/pkg/v"
	"io"

	"github.com/spf13/cobra"
)

type pluginInstallCmd struct {
	source  string
	version string
	home    slpath.Home
	out     io.Writer
}

const pluginInstallDesc = `
To install a plugin from a url or a local path.

Example usage:
    $ helm plugin install https://github.com/softleader/slctl-whereis
`

func newPluginInstallCmd(out io.Writer) *cobra.Command {
	pcmd := &pluginInstallCmd{out: out}
	cmd := &cobra.Command{
		Use:   "install [options] <path|url>...",
		Short: "install one or more plugins",
		Long:  pluginInstallDesc,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return pcmd.complete(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return pcmd.run()
		},
	}
	cmd.Flags().StringVar(&pcmd.version, "tag", "", "specify a tag constraint. If this is not specified, the latest version is installed")
	return cmd
}

func (pcmd *pluginInstallCmd) complete(args []string) error {
	if err := checkArgsLength(len(args), "plugin"); err != nil {
		return err
	}
	pcmd.source = args[0]
	pcmd.home = environment.Settings.Home
	return nil
}

func (pcmd *pluginInstallCmd) run() error {
	i, err := installer.NewInstaller(pcmd.out, pcmd.source, pcmd.version, pcmd.home)
	if err != nil {
		return err
	}
	var p *plugin.Plugin

	if p, err = i.Install(); err != nil {
		return err
	}

	v.Printf("loading plugin from %s\n", p.Dir)

	if err := runHook(p, plugin.Install); err != nil {
		return err
	}

	fmt.Fprintf(pcmd.out, "Installed plugin: %s\n", p.Metadata.Name)
	return nil
}