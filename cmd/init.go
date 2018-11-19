package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"
	"io"
	"os"
	"strings"
	"syscall"
)

const (
	note         = "slctl token (https://github.com/softleader/slctl)"
	organization = "softleader"
	initDesc     = `
This command grants Github access and sets up local configuration in $SL_HOME (default ~/.sl/).

	$ {{.}} init -t <github-token>
`
)

type initCmd struct {
	out      io.Writer
	home     slpath.Home
	dryRun   bool
	username string
	password string
	token    string
}

func newInitCmd(out io.Writer) *cobra.Command {
	i := &initCmd{out: out}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize " + Name,
		Long:  usage(initDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("this command does not accept arguments")
			}
			i.home = settings.Home
			return i.run()
		},
	}

	f := cmd.Flags()
	f.BoolVar(&i.dryRun, "dry-run", false, "do not login github")
	f.StringVarP(&i.token, "token", "t", "", "github access token")
	f.StringVarP(&i.username, "username", "u", "", "github username")
	f.StringVarP(&i.password, "password", "p", "", "github password")
	return cmd
}

func (i *initCmd) run() (err error) {
	if i.dryRun {
		return
	}

	if err = ensureDirectories(i.home, i.out); err != nil {
		return err
	}
	fmt.Fprintf(i.out, "$SL_HOME has been configured at %s.\n", settings.Home)

	if err = ensureConfigFile(i.home, i.out); err != nil {
		return err
	}

	if i.token == "" {
		if i.token, err = grantToken(i.username, i.password, i.out); err != nil {
			return err
		}
	}

	var username string
	if username, err = confirmToken(i.token, i.out); err != nil {
		return err
	}
	if err = refreshConfig(i.home, i.token, i.out); err != nil {
		return err
	}

	fmt.Fprintf(i.out, "Welcome aboard %s!\n", username)
	return
}
func grantToken(username, password string, out io.Writer) (token string, err error) {
	r := bufio.NewReader(os.Stdin)
	if username == "" {
		fmt.Fprint(out, "GitHub Username: ")
		username, _ = r.ReadString('\n')
	}
	if password == "" {
		fmt.Fprint(out, "GitHub Password: ")
		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
		password = string(bytePassword)
	}

	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	client := github.NewClient(tp.Client())
	ctx := context.Background()

	var auths []*github.Authorization
	auths, _, err = client.Authorizations.List(ctx, &github.ListOptions{})
	if _, ok := err.(*github.TwoFactorAuthError); ok {
		fmt.Fprint(out, "\nGitHub OTP: ")
		otp, _ := r.ReadString('\n')
		tp.OTP = strings.TrimSpace(otp)
		auths, _, err = client.Authorizations.List(ctx, &github.ListOptions{})
	}

	for _, auth := range auths {
		if auth.GetNote() == note {
			if settings.Verbose {
				fmt.Fprint(out, "\nRemoving exist token")
			}
			if _, err = client.Authorizations.Delete(ctx, auth.GetID()); err != nil {
				return "", err
			}
			break
		}
	}

	var auth *github.Authorization
	auth, _, err = client.Authorizations.Create(ctx, authorizationRequest())
	if err != nil {
		return "", err
	}

	return auth.GetToken(), nil
}

func authorizationRequest() *github.AuthorizationRequest {
	n := note
	return &github.AuthorizationRequest{
		Scopes: []github.Scope{github.ScopeReadOrg, github.ScopeUser},
		Note:   &n,
	}
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
	if settings.Verbose {
		fmt.Fprintf(out, "%s", github.Stringify(mem))
	}
	if mem.GetState() != "active" {
		return "", fmt.Errorf("you are not a active member of %s", organization)
	}
	var user *github.User
	if user, _, err = client.Users.Get(ctx, ""); err != nil {
		return "", err
	}
	if settings.Verbose {
		fmt.Fprintf(out, "%s", github.Stringify(user))
	}
	return user.GetName(), err
}

func ensureDirectories(home slpath.Home, out io.Writer) (err error) {
	configDirectories := []string{
		home.String(),
		home.Config(),
		home.Plugins(),
	}
	for _, p := range configDirectories {
		if fi, err := os.Stat(p); err != nil {
			fmt.Fprintf(out, "Creating %s \n", p)
			if err = os.MkdirAll(p, 0755); err != nil {
				return fmt.Errorf("could not create %s: %s", p, err)
			}
		} else if !fi.IsDir() {
			return fmt.Errorf("%s must be a directory", p)
		}
	}

	return
}

func ensureConfigFile(home slpath.Home, out io.Writer) (err error) {
	conf := home.ConfigFile()
	if fi, err := os.Stat(conf); err != nil {
		fmt.Fprintf(out, "Creating %s \n", conf)
		f := config.NewConfFile()
		if err := f.WriteFile(conf, config.ReadWrite); err != nil {
			return err
		}
	} else if fi.IsDir() {
		return fmt.Errorf("%s must be a file, not a directory", conf)
	}
	return
}

func refreshConfig(home slpath.Home, token string, out io.Writer) (err error) {
	conf, err := config.LoadConfFile(home.ConfigFile())
	if err != nil && err != config.ErrTokenNotExist {
		return fmt.Errorf("failed to load file (%v)", err)
	}
	conf.Token = token
	return conf.WriteFile(home.ConfigFile(), config.ReadWrite)
}
