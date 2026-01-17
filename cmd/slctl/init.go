package main

import (
	"context"
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/google/go-github/v69/github"
	"github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/environment"
	gh "github.com/softleader/slctl/pkg/github"
	"github.com/softleader/slctl/pkg/github/member"
	"github.com/softleader/slctl/pkg/github/token"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/prompt"
	"github.com/spf13/cobra"
)

// Mockable functions for testing
var (
	openBrowser      = open.Run
	writeToClipboard = clipboard.WriteAll
)

const (
	initDesc = `This command authorizes slctl with GitHub and sets up local configuration in $SL_HOME (default ~/.config/slctl/).

To authenticate, slctl will use the GitHub Device Flow. You will be prompted to enter a code in your browser.

Alternatively, you can provide a Personal Access Token (PAT) directly using the '--token' flag.
The PAT must have the required scopes. To see them, run 'slctl init scopes'.

	$ slctl init
	$ slctl init --token GITHUB_TOKEN

Using '--offline' skips all GitHub API interactions and only sets up local directories.
`
	askForPublicizeOrg = `Do you want to publicize the membership in SoftLeader?`
)

type initCmd struct {
	home  paths.Home
	token string
	force bool
	yes   bool
}

func newInitCmd() *cobra.Command {
	c := &initCmd{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize slctl",
		Long:  initDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c.home = environment.Settings.Home
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&c.force, "force", "f", false, "force to re-generate a new one if token already exists")
	f.StringVar(&c.token, "token", "", "github access token")
	f.BoolVarP(&c.yes, "yes", "y", false, "automatic 'yes' to prompts. Assume 'yes' as answer to all prompts and run non-interactively")

	cmd.AddCommand(
		newInitScopesCmd(),
	)

	return cmd
}

func (c *initCmd) run() (err error) {
	if c.home.ContainsAnySpace() {
		return fmt.Errorf(`default home path contains space which is not allowed (%s).
You might need to specify another SL_HOME without space and set to system variable.
For more details: https://github.com/softleader/slctl/wiki/Home-Path`, c.home.String())
	}
	if err = ensureDirectories(c.home, logrus.StandardLogger()); err != nil {
		return err
	}
	logrus.Printf("Slctl home has been configured at %s.\n", environment.Settings.Home)

	if err = ensureConfigFile(c.home, logrus.StandardLogger()); err != nil {
		return err
	}
	var username string
	var client *github.Client
	ctx := context.Background()
	if !environment.Settings.Offline {
		if c.token == "" {
			dcr, err := gh.RequestDeviceCode(ctx, "", gh.Scopes)
			if err != nil {
				return fmt.Errorf("failed to request device code: %w", err)
			}

			logrus.Printf("Please go to %s and enter the code: %s", dcr.VerificationURI, dcr.UserCode)

			// Attempt to copy user code to clipboard (non-fatal if fails)
			if err := writeToClipboard(dcr.UserCode); err != nil {
				logrus.Debugf("Failed to copy code to clipboard: %v", err)
			} else {
				logrus.Debug("Code copied to clipboard")
			}

			// Attempt to open browser (non-fatal if fails)
			if err := openBrowser(dcr.VerificationURI); err != nil {
				logrus.Debugf("Failed to open browser: %v", err)
			} else {
				logrus.Debug("Browser opened")
			}

			c.token, err = gh.PollAccessToken(ctx, "", dcr.DeviceCode, dcr.Interval)
			if err != nil {
				return fmt.Errorf("failed to poll for access token: %w", err)
			}
			logrus.Println("Successfully authenticated.")
		}

		if client, err = gh.NewTokenClient(ctx, c.token); err != nil {
			return err
		}

		if username, err = token.Confirm(ctx, client, organization, logrus.StandardLogger()); err != nil {
			return err
		}
	}
	if err = config.Refresh(c.home, c.token, logrus.StandardLogger()); err != nil {
		return err
	}
	if c.yes || prompt.YesNoQuestion(logrus.StandardLogger().Out, askForPublicizeOrg) {
		if err = member.PublicizeOrganization(ctx, client, organization); err != nil {
			return err
		}
	}
	logrus.Printf("Welcome aboard %s!\n", username)
	return
}

func ensureDirectories(home paths.Home, log *logrus.Logger) (err error) {
	configDirectories := []string{
		home.String(),
		home.Config(),
		home.Plugins(),
		home.Cache(),
		home.CachePlugins(),
		home.CacheArchives(),
		home.Mounts(),
	}
	return paths.EnsureDirectories(log, configDirectories...)
}

func ensureConfigFile(home paths.Home, log *logrus.Logger) (err error) {
	conf := home.ConfigFile()
	if fi, err := os.Stat(conf); err != nil {
		log.Printf("Creating %s \n", conf)
		f := config.NewConfFile()
		if err := f.WriteFile(conf, 0644); err != nil {
			return err
		}
	} else if fi.IsDir() {
		return fmt.Errorf("%s must be a file, not a directory", conf)
	}
	return
}
