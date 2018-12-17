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

const pluginUpgradeDesc = `Upgrade plugin which installed from GitHub Repo

	$ slctl plugin upgrade NAME...

NAME 可傳入指定要更新的 Plugin 完整名稱 (一或多個, 以空白區隔); 反之更新全部

	$ slctl plugin upgrade
	$ slctl plugin upgrade slctl-whereis

傳入 '--tag' 可以指定要更新的 release 版本

	$ slctl plugin upgrade slctl-whereis --tag 1.0.0

傳入 '--tag' 及 '--asset' 可以指定要更新的 release 版本以及要下載第幾個 asset 檔案 (從 0 開始)

	$ slctl plugin upgrade slctl-whereis --tag 1.0.0 --asset 2
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
		Use:   "upgrade NAME...",
		Short: "upgrade plugin  which installed from GitHub",
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

	f := cmd.Flags()
	f.StringVar(&c.tag, "tag", "", "specify a tag constraint. If this is not specified, the latest release tag is installed")
	f.IntVar(&c.asset, "asset", -1, "specify a asset number, start from zero, to download")

	return cmd
}

func (c *pluginUpgradeCmd) run() error {
	plugins, err := findPlugins(environment.Settings.PluginDirs())
	if err != nil {
		return err
	}
	var errors []string
	for _, p := range plugins {
		if !plugin.IsGitHubRepo(p.Source) {
			continue
		}
		if len(c.names) == 0 || match(p, c.names) {
			fmt.Fprintf(c.out, "Upgrading %q plugin\n", p.Metadata.Name)
			if err := install(c.out, p.Source, c.tag, c.asset, c.home, true, true); err != nil {
				errors = append(errors, err.Error())
			}
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil
}

func match(p *plugin.Plugin, names []string) bool {
	for _, n := range names {
		if strings.ToLower(p.Metadata.Name) == strings.ToLower(n) {
			return true
		}
	}
	return false
}
