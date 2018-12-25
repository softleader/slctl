package main

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/spf13/cobra"
)

const (
	pluginExtsDesc = `列出所有安裝 plugin 時支援的壓縮檔案 

	$ {{.}} plugin exts
`
)

func newPluginExtsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exts",
		Short: "list supported plugin archive extension to install",
		Long:  usage(pluginExtsDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("this command does not accept arguments")
			}
			for _, ext := range plugin.SupportedExtensions {
				logrus.Println(ext)
			}
			return nil
		},
	}
	return cmd
}
