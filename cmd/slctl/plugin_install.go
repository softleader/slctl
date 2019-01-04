package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/config/token"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/plugin/installer"
	"io"

	"github.com/spf13/cobra"
)

type pluginInstallCmd struct {
	source string
	tag    string
	asset  int
	home   paths.Home
	out    io.Writer
	opt    *installer.InstallOption
}

const pluginInstallDesc = `To install a plugin from a local path, a archive, or a GitHub repo

Plugin 可以是本機上的任何目錄, 透過給予絕對或相對路徑來安裝

	$ {{.}} plugin install /path/to/plugin-dir/

Plugin 也可以是來自於網路上或在本機中壓縮檔, 透過給予網址或路徑來安裝

	$ {{.}} plugin install /path/to/plugin-archive.zip
	$ {{.}} plugin install http://host/plugin-archive.zip

Plugin 也可以是一個 GitHub repo, 傳入 'github.com/OWNER/REPO', {{.}} 會自動收尋最新一版的 release
並從該 release 的所有下載檔中, 嘗試找出含有當前 OS 名稱的壓縮檔來安裝, 當找不到時會改下載第一個壓縮檔來安裝

	$ {{.}} plugin install github.com/softleader/slctl-whereis

傳入 '--tag' 可以指定 release 版本

	$ {{.}} plugin install github.com/softleader/slctl-whereis --tag 1.0.0

傳入 '--tag' 及 '--asset' 可以指定 release 版本以及要下載第幾個 asset 檔案 (從 0 開始) 來安裝

	$ {{.}} plugin install github.com/softleader/slctl-whereis --tag 1.0.0 --asset 2

傳入 '--force' 在 install 時強制刪除已存在的 plugin

	$ {{.}} plugin install github.com/softleader/slctl-whereis -f

傳入 '--dry-run' 可以模擬真實的 install, 但不會真的影響當前的配置 

	$ {{.}} plugin install github.com/softleader/slctl-whereis -f --dry-run
`

func newPluginInstallCmd() *cobra.Command {
	c := &pluginInstallCmd{
		opt: &installer.InstallOption{},
	}
	cmd := &cobra.Command{
		Use:   "install [options] <SOURCE>...",
		Short: "install one or more plugins",
		Long:  usage(pluginInstallDesc),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c.source = args[0]
			c.home = environment.Settings.Home
			return c.run()
		},
	}

	f := cmd.Flags()
	f.StringVar(&c.tag, "tag", "", "specify a tag constraint. If this is not specified, the latest release tag is installed")
	f.IntVar(&c.asset, "asset", -1, "specify a asset number, start from zero, to download")
	f.BoolVar(&c.opt.DryRun, "dry-run", false, `simulate an install "for real"`)
	f.BoolVarP(&c.opt.Force, "force", "f", false, "force to re-install if plugin already exists")
	f.BoolVarP(&c.opt.Soft, "soft", "s", false, "force to remove exist plugin only if version is different")

	return cmd
}

//func (c *pluginInstallCmd) complete(args []string) error {
//	if err := checkArgsLength(len(args), "plugin"); err != nil {
//		return err
//	}
//	c.source = args[0]
//	c.home = environment.Settings.Home
//	return nil
//}

func (c *pluginInstallCmd) run() error {
	if c.opt.DryRun {
		logrus.Warnln("running in dry-run mode, specify the '-v' flag if you want to turn on verbose output")
	}
	return install(c.source, c.tag, c.asset, c.home, c.opt)
}

func install(source string, tag string, asset int, home paths.Home, opt *installer.InstallOption) error {
	i, err := installer.NewInstaller(logrus.StandardLogger(), source, tag, asset, home, opt)
	if err != nil {
		return err
	}
	var p *plugin.Plugin

	if p, err = i.Install(); err != nil {
		if err == installer.ErrAlreadyUpToDate {
			logrus.Printf("Plugin \"%s@%s\" already up-to-date\n", p.Metadata.Name, p.Metadata.Version)
			return nil
		}
		return err
	}

	if err = token.EnsureScopes(logrus.StandardLogger(), p.Metadata.GitHub.Scopes); err != nil {
		return err
	}

	if err := runHook(p); err != nil {
		if _, ok := err.(*plugin.ErrNoCommandFound); !ok {
			return err
		}
	}

	logrus.Printf("Installed plugin: %s@%s\n", p.Metadata.Name, p.Metadata.Version)
	return nil
}
