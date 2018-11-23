package main

import (
	"errors"
	"fmt"
	"github.com/softleader/slctl/pkg/plugin"
	"github.com/spf13/cobra"
	"io"
	"reflect"
)

const (
	pluginCreateLangsDesc = `
列出所有 plugin 範本的語言 

	$ {{.}} plugin create langs
`
)

func newPluginCreateLangsCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "langs",
		Short: "list languages of plugin template",
		Long:  usage(pluginCreateLangsDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("this command does not accept arguments")
			}
			for _, c := range plugin.Creators {
				fmt.Fprintln(out, reflect.TypeOf(c).Name())
			}
			return nil
		},
	}
	return cmd
}
