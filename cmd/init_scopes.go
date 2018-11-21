package cmd

import (
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"io"
)

const (
	initScopesDesc = `
列出所有 {{.}} 需要的 GitHub Personal Access Token 權限 (https://github.com/settings/tokens)

	$ {{.}} init scopes
`
)

var tokenScopes = []github.Scope{github.ScopeReadOrg, github.ScopeUser}

func newInitScopesCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scopes",
		Short: "list scopes of token that " + Name + " required",
		Long:  usage(initScopesDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("this command does not accept arguments")
			}
			for _, scope := range tokenScopes {
				fmt.Fprintln(out, scope)
			}
			return nil
		},
	}
	return cmd
}
