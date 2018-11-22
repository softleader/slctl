package cmd

import (
	"errors"
	"fmt"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/spf13/cobra"
	"io"
	"path/filepath"
	"strings"
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

Plugin 本身沒有撰寫的語言限制, {{.}} 也已內含了幾種語言的範本 ({{.}} 推薦並預設產生 golang 的範本)
使用 '--lang' 指定你要產生的語言範本
或使用 'plugin create langs' 列出所有內含的範本語言
	
	$ {{.}} plugin create foo --lang java
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
		Description: fmt.Sprintf("the %s plugin written in %s", pname, strings.Title(c.lang)),
		Version:     "0.1.0",
		Command:     "$SL_PLUGIN_DIR/" + pname,
	}
	_, err = plugin.Create(c.lang, pfile, filepath.Dir(c.name))
	return
}
