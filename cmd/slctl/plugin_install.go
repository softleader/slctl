package main

import (
	"fmt"
	"github.com/softleader/slctl/pkg/config/token"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/plugin/installer"
	"github.com/softleader/slctl/pkg/slpath"
	"io"

	"github.com/spf13/cobra"
)

type pluginInstallCmd struct {
	source string
	tag    string
	asset  int
	home   slpath.Home
	out    io.Writer
	force  bool
	soft   bool
}

const pluginInstallDesc = `To install a plugin from a local path, a archive, or a GitHub repo

Plugin 可以是本機上的任何目錄, 透過給予絕對或相對路徑來安裝

	$ slctl plugin install /path/to/plugin-dir/

Plugin 也可以是來自於網路上或在本機中壓縮檔, 透過給予網址或路徑來安裝

	$ slctl plugin install /path/to/plugin-archive.zip
	$ slctl plugin install http://host/plugin-archive.zip

Plugin 也可以是一個 GitHub repo, 傳入 'github.com/OWNER/REPO', {{.}} 會自動收尋最新一版的 release
並從該 release 的所有下載檔中, 嘗試找出含有當前 OS 名稱的壓縮檔來安裝, 當找不到時會改下載第一個壓縮檔來安裝

	$ slctl plugin install github.com/softleader/slctl-whereis

傳入 '--tag' 可以指定 release 版本

	$ slctl plugin install github.com/softleader/slctl-whereis --tag 1.0.0

傳入 '--tag' 及 '--asset' 可以指定 release 版本以及要下載第幾個 asset 檔案 (從 0 開始) 來安裝

	$ slctl plugin install github.com/softleader/slctl-whereis --tag 1.0.0 --asset 2

傳入 '--force' 在 install 時強制刪除已存在的 plugin

	$ slctl plugin install github.com/softleader/slctl-whereis -f
`

func newPluginInstallCmd(out io.Writer) *cobra.Command {
	c := &pluginInstallCmd{out: out}
	cmd := &cobra.Command{
		Use:   "install [options] <SOURCE>...",
		Short: "install one or more plugins",
		Long:  usage(pluginInstallDesc),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.complete(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.run()
		},
	}

	f := cmd.Flags()
	f.StringVar(&c.tag, "tag", "", "specify a tag constraint. If this is not specified, the latest release tag is installed")
	f.IntVar(&c.asset, "asset", -1, "specify a asset number, start from zero, to download")
	f.BoolVarP(&c.force, "force", "f", false, "force to re-install if plugin already exists")
	f.BoolVarP(&c.soft, "soft", "s", false, "force to remove exist plugin only if version is different")

	return cmd
}

func (c *pluginInstallCmd) complete(args []string) error {
	if err := checkArgsLength(len(args), "plugin"); err != nil {
		return err
	}
	c.source = args[0]
	c.home = environment.Settings.Home
	return nil
}

func (c *pluginInstallCmd) run() error {
	return install(c.out, c.source, c.tag, c.asset, c.home, c.force, c.soft)
}

func install(out io.Writer, source string, tag string, asset int, home slpath.Home, force, soft bool) error {
	i, err := installer.NewInstaller(out, source, tag, asset, home, force, soft)
	if err != nil {
		return err
	}
	var p *plugin.Plugin

	if p, err = i.Install(); err != nil {
		if err == installer.ErrAlreadyUpToDate {
			fmt.Fprintf(out, "Plugin \"%s@%s\" already up-to-date\n", p.Metadata.Name, p.Metadata.Version)
			return nil
		}
		return err
	}

	if err = token.EnsureScopes(out, p.Metadata.GitHub.Scopes); err != nil {
		return err
	}

	if err := runHook(p); err != nil {
		if _, ok := err.(*plugin.ErrNoCommandFound); !ok {
			return err
		}
	}

	fmt.Fprintf(out, "Installed plugin: %s@%s\n", p.Metadata.Name, p.Metadata.Version)
	return nil
}
