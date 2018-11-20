package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"io"
)

const whereisHelp = `
Lookup where the SoftLeader member is.
`

type whereisCmd struct {
	out   io.Writer
	name  string
	limit string
	place string
	date  string
}

func newWhereisCmd(out io.Writer) *cobra.Command {
	w := &whereisCmd{out: out}

	cmd := &cobra.Command{
		Use:   "whereis <name>",
		Short: "lookup where the SoftLeader member is",
		Long:  usage(whereisHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("please provide name to lookup")
			}
			return w.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&w.limit, "limit", "l", "20/1", "limit output <size>/[page]")
	f.StringVarP(&w.date, "date", "d", "today", "filter the specified date <from>..[to]")
	f.StringVarP(&w.place, "place", "p", "", "filter the specified place")

	cmd.AddCommand(
		newWhereisPlacesCmd(out),
	)

	return cmd
}

func (i *whereisCmd) run() (err error) {
	// TODO
	return nil
}
