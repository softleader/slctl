package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"io"
	"os"
)

const (
	organization = "softleader"
	initDesc     = `
This command grants Github access and sets up local configuration in $SL_HOME (default ~/.sl/).

	$ {{.}} init -t <github-token>
`
)

type initCmd struct {
	out    io.Writer
	home   slpath.Home
	dryRun bool
	token  string
}

func newInitCmd(out io.Writer) *cobra.Command {
	i := &initCmd{out: out}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize " + Name,
		Long:  usage(initDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("This command does not accept arguments")
			}
			i.home = settings.Home
			return i.run()
		},
	}

	f := cmd.Flags()
	f.BoolVar(&i.dryRun, "dry-run", false, "do not login github")
	f.StringVarP(&i.token, "token", "t", "", "github access token")
	return cmd
}

func (i *initCmd) run() (err error) {
	if i.dryRun {
		return nil
	}

	var name string
	if name, err = confirmToken(i.token, i.out); err != nil {
		return err
	}

	if err = ensureDirectories(i.home, i.out); err != nil {
		return err
	}
	fmt.Fprintf(i.out, "$SL_HOME has been configured at %s.\n", settings.Home)

	fmt.Fprintf(i.out, "Welcome aboard %s!\n", name)
	return nil
}

func confirmToken(token string, out io.Writer) (name string, err error) {
	if token == "" {
		return "", fmt.Errorf("required flag(s) \"token\" not set")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var mem *github.Membership
	if mem, _, err = client.Organizations.GetOrgMembership(ctx, "", organization); err != nil {
		return "", err
	}
	if settings.Debug {
		fmt.Fprintf(out, "%s", mem)
	}
	if mem.GetState() != "active" {
		return "", fmt.Errorf("you are not a active member of %s", organization)
	}
	var user *github.User
	if user, _, err = client.Users.Get(ctx, ""); err != nil {
		return "", err
	}
	if settings.Debug {
		fmt.Fprintf(out, "%s", user)
	}
	return user.GetName(), nil
}

func ensureDirectories(home slpath.Home, out io.Writer) (err error) {
	configDirectories := []string{
		home.String(),
		home.Plugins(),
	}
	for _, p := range configDirectories {
		if fi, err := os.Stat(p); err != nil {
			fmt.Fprintf(out, "Creating %s \n", p)
			if err = os.MkdirAll(p, 0755); err != nil {
				return fmt.Errorf("Could not create %s: %s", p, err)
			}
		} else if !fi.IsDir() {
			return fmt.Errorf("%s must be a directory", p)
		}
	}

	return nil
}
