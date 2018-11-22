package cmd

import (
	"errors"
	"fmt"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/spf13/cobra"
	"io"
	"path/filepath"
)

const pluginCreateDesc = `
產生 Plugin 範本, 如: '{{.}} plugin create foo' 將會產生 golang plugin 範本, 目錄結構大致如下:

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

使用 'plugin create langs' 列出所有內含的範本語言

	$ {{.}} plugin create langs
`

type pluginCreateCmd struct {
	home slpath.Home
	name string
	out  io.Writer
	lang string
}

func newPluginCreateCmd(out io.Writer) *cobra.Command {
	pcc := &pluginCreateCmd{out: out}
	cmd := &cobra.Command{
		Use:   "create NAME",
		Short: "create a new plugin with the given name",
		Long:  usage(pluginCreateDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			pcc.home = settings.Home
			if len(args) == 0 {
				return errors.New("the name of the new plugin is required")
			}
			pcc.name = args[0]
			return pcc.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&pcc.lang, "lang", "", "golang", "language of template to create")

	cmd.AddCommand(
		newPluginCreateLangsCmd(out),
	)

	return cmd
}

func (c *pluginCreateCmd) run() (err error) {
	pname := filepath.Base(c.name)
	fmt.Fprintf(c.out, "Creating %s\n", c.name)
	pfile := &plugin.Metadata{
		Name:        pname,
		Usage:       pname,
		Description: fmt.Sprintf("The %s plugin", pname),
		Version:     "0.1.0",
	}
	_, err = plugin.Create(c.lang, pfile, filepath.Dir(c.name))
	return
}
