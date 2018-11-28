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
To install a plugin from a local path, a archive url, or a GitHub repo

Plugin 可以是本機上的任何目錄, 透過給予絕對或相對路徑來安裝

	$ slctl plugin install /path/to/plugin-dir/

Plugin 也可以是來自於網路上或在本機中壓縮檔, 透過給予網址或路徑來安裝

	$ slctl plugin install /path/to/plugin-archive.zip
	$ slctl plugin install http://host/plugin-archive.zip

Plugin 也可以是一個 GitHub repo, 傳入 'github.com/OWNER/REPO', {{.}} 會自動收尋最新一版的 release
並從該 release 的所有下載檔中, 嘗試找出含有當前 OS 名稱的壓縮檔來安裝, 當找不到時會改下載第一個壓縮檔來安裝

	$ slctl plugin install github.com/softleader/slctl-whereis

傳入 '--tag' 可以指定 release 版本

	$ slctl plugin install github.com/softleader/slctl-whereis --tag v1.0.0
`

func newPluginInstallCmd(out io.Writer) *cobra.Command {
	pcmd := &pluginInstallCmd{out: out}
	cmd := &cobra.Command{
		Use:   "install [options] <SOURCE>...",
		Short: "install one or more plugins",
		Long:  pluginInstallDesc,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return pcmd.complete(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return pcmd.run()
		},
	}
	cmd.Flags().StringVar(&pcmd.version, "tag", "", "specify a tag constraint. If this is not specified, the latest release version is installed")
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

	if err := runHook(p); err != nil {
		return err
	}

	fmt.Fprintf(pcmd.out, "Installed plugin: %s\n", p.Metadata.Name)
	return nil
}
