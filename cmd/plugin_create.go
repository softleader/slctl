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
This command creates a plugin directory.

For example, '{{.}} plugin create foo' will create a directory structure like this:

	foo/
	  |
	  |- plugin.xml
	  |
	  |- main.go
	  |
	  |- Makefile
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
		Long:  pluginCreateDesc,
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
	pluginName := filepath.Base(c.name)
	fmt.Fprintf(c.out, "Creating %s\n", c.name)
	pfile := &plugin.Metadata{
		Name:        pluginName,
		Usage:       pluginName,
		Description: fmt.Sprintf("the %s plugin", pluginName),
		Version:     "0.1.0",
		Command:     "$SL_PLUGIN_DIR/" + pluginName,
	}
	_, err = plugin.Create(c.lang, pfile, filepath.Dir(c.name))
	return
}
