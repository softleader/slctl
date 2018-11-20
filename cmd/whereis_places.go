package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
)

const (
	whereisPlacesDesc = `
可以列出所有 {{.}} whereis 支援的地點

	$ {{.}} whereis places
`
)

func newWhereisPlacesCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "places",
		Short: "list all places",
		Long:  usage(whereisPlacesDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("this command does not accept arguments")
			}
			// TODO 要向後端查
			fmt.Fprintln(out, "1:本部")
			fmt.Fprintln(out, "2:仁愛")
			fmt.Fprintln(out, "3:建北")
			fmt.Fprintln(out, "4:內湖")
			return nil
		},
	}
	return cmd
}
