package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/spf13/cobra"
	"path/filepath"
)

const pluginCreateDesc = `產生 Plugin 範本, 如: '{{.}} plugin create foo' 將會產生 golang plugin 範本, 目錄結構大致如下:

	foo/
	  |
	  |- plugin.xml
	  |
	  |- main.go
	  |
	  |- Makefile

Plugin 本身沒有撰寫的語言限制, {{.}} 推薦並預設產生 golang 的範本
選擇不同撰寫語言時, 需注意該語言本身的限制: 如執行 java plugin 的 runtime 必須有 JVM
{{.}} 已內含了幾種語言的範本, 使用 '--lang' 來指定產生語言範本
	
	$ {{.}} plugin create foo --lang java

使用 'plugin create langs' 來列出所有內含的範本語言

	$ {{.}} plugin create langs

{{.|title}} 預設會在當前目錄下, 建立一個名為 Plugin 名稱的目錄, 並將範本產生在該目錄中
可以傳入 '--output' 來指定 Plugin 的產生目錄

	$ {{.}} plugin create foo -o /path/to/plugin-dir
`

type pluginCreateCmd struct {
	home   paths.Home
	name   string
	lang   string
	output string
}

func newPluginCreateCmd() *cobra.Command {
	pcc := &pluginCreateCmd{}
	cmd := &cobra.Command{
		Use:   "create NAME",
		Short: "create a new plugin with the given name",
		Long:  usage(pluginCreateDesc),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pcc.home = environment.Settings.Home
			pcc.name = args[0]
			return pcc.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&pcc.lang, "lang", "", "golang", "language of template to create")
	f.StringVarP(&pcc.output, "output", "o", "", "output directory name, uses plugin name if leave blank")

	cmd.AddCommand(
		newPluginCreateLangsCmd(),
	)

	return cmd
}

func (c *pluginCreateCmd) run() (err error) {
	pname := filepath.Base(c.name)
	logrus.Printf("Creating %s plugin %q\n", c.lang, c.name)
	pfile := &plugin.Metadata{
		Name:        pname,
		Usage:       fmt.Sprintf("the %s plugin", pname),
		Description: fmt.Sprintf("The %s plugin", pname),
		Version:     "0.1.0",
	}
	path, err := plugin.Create(c.lang, pfile, c.output)
	logrus.Printf("Successfully created plugin and saved it to: %s\n", path)
	return
}
