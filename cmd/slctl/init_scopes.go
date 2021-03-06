package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/github/token"
	"github.com/spf13/cobra"
)

const (
	initScopesDesc = `列出所有 slctl 需要的 GitHub Personal Access Token 權限 (https://github.com/settings/tokens)

	$ slctl init scopes
`
)

func newInitScopesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scopes",
		Short: "list scopes of token that slctl required",
		Long:  initScopesDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, scope := range token.Scopes {
				logrus.Println(scope)
			}
			return nil
		},
	}
	return cmd
}
