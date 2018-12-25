package main

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/config/token"
	"github.com/spf13/cobra"
)

const (
	initScopesDesc = `列出所有 {{.}} 需要的 GitHub Personal Access Token 權限 (https://github.com/settings/tokens)

	$ {{.}} init scopes
`
)

func newInitScopesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scopes",
		Short: "list scopes of token that " + name + " required",
		Long:  usage(initScopesDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("this command does not accept arguments")
			}
			for _, scope := range token.Scopes {
				logrus.Println(scope)
			}
			return nil
		},
	}
	return cmd
}
