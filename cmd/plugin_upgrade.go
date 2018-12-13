package main

import (
	"fmt"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/spf13/cobra"
	"io"
	"strings"
)

const pluginUpgradeDesc = `更新 Plugin, 只會更新 SOURCE 為 GitHub 的 Plugin

	$ slctl plugin upgrade NAME...
`

type pluginUpgradeCmd struct {
	home  slpath.Home
	out   io.Writer
	names []string
	tag   string
	asset int
}

func newPluginUpgradeCmd(out io.Writer) *cobra.Command {
	c := &pluginUpgradeCmd{out: out}
	cmd := &cobra.Command{
		Use:   "upgrade NAME",
		Short: "upgrade plugin from Source Github",
		Long:  usage(pluginUpgradeDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			if environment.Settings.Offline {
				return fmt.Errorf("can not run the command in offline mode")
			}
			if len(args) > 0 {
				c.names = args
			}
			c.home = environment.Settings.Home
			return c.run()
		},
	}

	cmd.Flags().StringVar(&c.tag, "tag", "", "specify a tag constraint. If this is not specified, the latest release tag is installed")
	cmd.Flags().IntVar(&c.asset, "asset", -1, "specify a asset number, start from zero, to download")

	return cmd
}

func (c *pluginUpgradeCmd) run() error {
	plugins, err := findPlugins(environment.Settings.PluginDirs())
	if err != nil {
		return err
	}
	var upgrades []*plugin.Plugin
	for _, p := range plugins {
		if !p.FromGitHub() {
			continue
		}
		if match(p, c.names) {
			upgrades = append(upgrades, p)
		}
	}
	var errors []string
	for _, plug := range upgrades {
		fmt.Fprintf(c.out, "Upgrading %q plugin\n", plug.Metadata.Name)
		if err := install(c.out, plug.Source, c.tag, c.asset, c.home, true, true); err != nil {
			errors = append(errors, err.Error())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil
}

func match(p *plugin.Plugin, names []string) bool {
	if len(names) == 0 {
		return true
	}
	for _, n := range names {
		if strings.ToLower(p.Metadata.Name) == strings.ToLower(n) {
			return true
		}
	}
	return false
}
