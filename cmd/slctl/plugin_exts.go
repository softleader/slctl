package main

import (
	"errors"
	"fmt"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/spf13/cobra"
	"io"
)

const (
	pluginExtsDesc = `列出所有安裝 plugin 時支援的壓縮檔案 

	$ {{.}} plugin exts
`
)

func newPluginExtsCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exts",
		Short: "list supported plugin archive extension to install",
		Long:  usage(pluginExtsDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("this command does not accept arguments")
			}
			for _, ext := range plugin.SupportedExtensions {
				fmt.Fprintln(out, ext)
			}
			return nil
		},
	}
	return cmd
}
