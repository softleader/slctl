package main

import (
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/spf13/cobra"
)

const (
	pluginCreateLangsDesc = `列出所有 plugin 範本的語言 

	$ slctl plugin create langs
`
)

func newPluginCreateLangsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "langs",
		Short: "list languages of plugin template",
		Long:  pluginCreateLangsDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, c := range plugin.Creators {
				logrus.Println(reflect.TypeOf(c).Name())
			}
			return nil
		},
	}
	return cmd
}
