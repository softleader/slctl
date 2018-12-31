package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/config/token"
	"github.com/softleader/slctl/pkg/dir"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/spf13/cobra"
	"os"
)

const (
	initDesc = `This command grants Github access token and sets up local configuration in $SL_HOME (default ~/.sl/).

執行 '{{.}} init' 透過互動式的問答產生並儲存 GitHub Personal Access Token (https://github.com/settings/tokens)
也可以傳入 '--username' 或 '--password' 來整合非互動式的情境 (e.g. DevOps pipeline):

	$ {{.}} init
	$ {{.}} init -u GITHUB_USERNAME -p GITHUB-PASSWORD

使用 '--force' 在發現有重複的 Token 時, 會強制刪除並產生一個全新的 Access Token

	$ {{.}} init -f

若你想自己維護 Access Token (請務必確保有足夠的權限), 可以使用 '--token' 讓 slctl 驗證後直接儲存起來
執行 'scopes'' 可以列出所有 slctl 需要的 Access Token 權限

	$ {{.}} init --token GITHUB_TOKEN
	$ {{.}} init scopes

使用 '--offline' 則 {{.}} 不會跟 GitHub API 有任何互動, 只會配置 $SL_HOME 環境目錄.

同時使用 '--offline' 及 '--token' 可跳過 Token 驗證直接儲存起來 (e.g. 沒網路環境下)
`
)

type initCmd struct {
	home     slpath.Home
	username string
	password string
	token    string
	force    bool
}

func newInitCmd() *cobra.Command {
	i := &initCmd{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize " + name,
		Long:  usage(initDesc),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			i.home = environment.Settings.Home
			return i.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&i.force, "force", "f", false, "force to re-generate a new one if token already exists")
	f.StringVar(&i.token, "token", "", "github access token")
	f.StringVarP(&i.username, "username", "u", "", "github username")
	f.StringVarP(&i.password, "password", "p", "", "github password")

	cmd.AddCommand(
		newInitScopesCmd(),
	)

	return cmd
}

func (c *initCmd) run() (err error) {
	if c.home.ContainsAnySpace() {
		return fmt.Errorf(`default home path contains space which is not allowed (%s).
You might need to specify another SL_HOME without space and set to system variable.
For more details: https://github.com/softleader/slctl/wiki/Home-Path
`, c.home.String())
	}

	if err = ensureDirectories(c.home, logrus.StandardLogger()); err != nil {
		return err
	}
	logrus.Printf("Slctl home has been configured at %s.\n", environment.Settings.Home)

	if err = ensureConfigFile(c.home, logrus.StandardLogger()); err != nil {
		return err
	}
	var username string
	if !environment.Settings.Offline {
		if c.token == "" {
			if c.token, err = token.Grant(c.username, c.password, logrus.StandardLogger(), c.force); err != nil {
				return err
			}
		}
		if username, err = token.Confirm(organization, c.token, logrus.StandardLogger()); err != nil {
			return err
		}
	}
	if err = config.Refresh(c.home, c.token, logrus.StandardLogger()); err != nil {
		return err
	}
	logrus.Printf("Welcome aboard %s!\n", username)
	return
}

func ensureDirectories(home slpath.Home, log *logrus.Logger) (err error) {
	configDirectories := []string{
		home.String(),
		home.Config(),
		home.Plugins(),
		home.Cache(),
		home.CachePlugins(),
		home.CacheArchives(),
		home.Mounts(),
	}
	return dir.EnsureDirectories(log, configDirectories...)
}

func ensureConfigFile(home slpath.Home, log *logrus.Logger) (err error) {
	conf := home.ConfigFile()
	if fi, err := os.Stat(conf); err != nil {
		log.Printf("Creating %s \n", conf)
		f := config.NewConfFile()
		if err := f.WriteFile(conf, config.ReadWrite); err != nil {
			return err
		}
	} else if fi.IsDir() {
		return fmt.Errorf("%s must be a file, not a directory", conf)
	}
	return
}
