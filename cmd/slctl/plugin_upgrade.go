package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/plugin/installer"
	"github.com/spf13/cobra"
	"strings"
)

const pluginUpgradeDesc = `Upgrade plugin which installed from GitHub Repo

	$ {{.}} plugin upgrade NAME...

NAME 可傳入指定要更新的 Plugin 完整名稱 (一或多個, 以空白區隔); 反之更新全部

	$ {{.}} plugin upgrade
	$ {{.}} plugin upgrade whereis

傳入 '--tag' 可以指定要更新的 release 版本

	$ {{.}} plugin upgrade whereis --tag 1.0.0

傳入 '--tag' 及 '--asset' 可以指定要更新的 release 版本以及要下載第幾個 asset 檔案 (從 0 開始)

	$ {{.}} plugin upgrade whereis --tag 1.0.0 --asset 2

傳入 '--dry-run' 可以模擬真實的 upgrade, 但不會真的影響當前的配置
通常可以用來檢查 plugin 是否有新版的再決定是否要更新

	$ {{.}} plugin upgrade --dry-run
`

type pluginUpgradeCmd struct {
	home  paths.Home
	names []string
	opt   *installer.InstallOption
	tag   string
	asset int
}

func newPluginUpgradeCmd() *cobra.Command {
	c := &pluginUpgradeCmd{
		opt: &installer.InstallOption{
			Force: true,
			Soft:  true,
		},
	}
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
	f.BoolVar(&c.opt.DryRun, "dry-run", false, `simulate an upgrade "for real"`)
	f.StringVar(&c.tag, "tag", "", "specify a tag constraint. If this is not specified, the latest release tag is installed")
	f.IntVar(&c.asset, "asset", -1, "specify a asset number, start from zero, to download")

	return cmd
}

func (c *pluginUpgradeCmd) run() error {
	if c.opt.DryRun {
		logrus.Warnln("running in dry-run mode, specify the '-v' flag if you want to turn on verbose output")
	}
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
			logrus.Printf("Upgrading %q plugin\n", p.Metadata.Name)
			if err := install(p.Source, c.tag, c.asset, c.home, c.opt); err != nil {
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
