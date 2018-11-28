package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/softleader/slctl/pkg/v"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"
	"io"
	"os"
	"strings"
	"syscall"
)

const (
	note         = name + " token (https://github.com/softleader/slctl)"
	organization = "softleader"
	initDesc     = `
This command grants Github access token and sets up local configuration in $SL_HOME (default ~/.sl/).

執行 '{{.}} init' 透過互動式的問答產生並儲存 GitHub Personal Access Token (https://github.com/settings/tokens)
也可以傳入 '--username' 或 '--password' 來整合非互動式的情境 (e.g. DevOps pipeline):

	$ {{.}} init
	$ {{.}} init -u GITHUB_USERNAME -p GITHUB-PASSWORD

執行 'scopes' 可以列出所有 {{.}} 需要的 Access Token 權限

	$ {{.}} init scopes

使用 '--refresh' 讓 {{.}} 發現有重複的 Token 時, 自動的刪除既有的並產生一個全新的 Access Token
若你想自己維護 Access Token (請務必確保有足夠的權限), 可以使用 '--token' 讓 {{.}} 驗證後直接儲存起來

	$ {{.}} init --refresh
	$ {{.}} init --token GITHUB_TOKEN

使用 '--offline' 則 {{.}} 不會跟 GitHub API 有任何互動, 只會配置 $SL_HOME 環境目錄.
同時使用 '--offline' 及 '--token' 可跳過 Token 驗證直接儲存起來 (e.g. 沒網路環境下)
`
)

var ErrOauthAccessAlreadyExists = errors.New(`access token already exists.
To store a token already on https://github.com/settings/tokens, use '--token' flag.  
To automatically re-generate a new one, use '--refresh' flag.
Use 'init --help' for more information about the command.`)

type initCmd struct {
	out      io.Writer
	home     slpath.Home
	username string
	password string
	token    string
	refresh  bool
}

func newInitCmd(out io.Writer) *cobra.Command {
	i := &initCmd{out: out}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize " + name,
		Long:  usage(initDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			i.home = environment.Settings.Home
			return i.run()
		},
	}

	f := cmd.Flags()
	f.BoolVar(&i.refresh, "refresh", false, "automatically re-generate a new one if token already exists")
	f.StringVar(&i.token, "token", "", "github access token")
	f.StringVarP(&i.username, "username", "u", "", "github username")
	f.StringVarP(&i.password, "password", "p", "", "github password")

	cmd.AddCommand(
		newInitScopesCmd(out),
	)

	return cmd
}

func (c *initCmd) run() (err error) {
	if err = ensureDirectories(c.home, c.out); err != nil {
		return err
	}
	fmt.Fprintf(c.out, "$SL_HOME has been configured at %s.\n", environment.Settings.Home)

	if err = ensureConfigFile(c.home, c.out); err != nil {
		return err
	}
	var username string
	if !environment.Settings.Offline {
		if c.token == "" {
			if c.token, err = grantToken(c.username, c.password, c.out, c.refresh); err != nil {
				return err
			}
		}
		if username, err = confirmToken(c.token, c.out); err != nil {
			return err
		}
	}
	if err = refreshConfig(c.home, c.token, c.out); err != nil {
		return err
	}
	fmt.Fprintf(c.out, "Welcome aboard %s!\n", username)
	return
}
func grantToken(username, password string, out io.Writer, refresh bool) (token string, err error) {
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
			if !refresh {
				return "", ErrOauthAccessAlreadyExists
			}
			v.Fprintln(out, "\nRemoving exist token")
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
		Scopes: tokenScopes,
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
	// v.Fprintf(out, "%s\n", github.Stringify(mem))
	if mem.GetState() != "active" {
		return "", fmt.Errorf("you are not a active member of %s", organization)
	}
	var user *github.User
	if user, _, err = client.Users.Get(ctx, ""); err != nil {
		return "", err
	}
	// v.Fprintf(out, "%s\n", github.Stringify(user))
	return user.GetName(), err
}

func ensureDirectories(home slpath.Home, out io.Writer) (err error) {
	configDirectories := []string{
		home.String(),
		home.Config(),
		home.Plugins(),
		home.Cache(),
		home.CachePlugins(),
		home.CacheArchives(),
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
