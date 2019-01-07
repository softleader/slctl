package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/spf13/cobra"
)

const (
	pluginExtsDesc = `列出所有安裝 plugin 時支援的壓縮檔案 

	$ slctl plugin exts
`
)

func newPluginExtsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exts",
		Short: "list supported plugin archive extension to install",
		Long:  pluginExtsDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, ext := range plugin.SupportedExtensions {
				logrus.Println(ext)
			}
			return nil
		},
	}
	return cmd
}
