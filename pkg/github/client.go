package github

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/google/go-github/v21/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"
)

// NewBasicAuthClient 產生一個 BasicAuth 的 GitHub Client, 並且會檢查完 OTP
func NewBasicAuthClient(ctx context.Context, log *logrus.Logger, username, password string) (*github.Client, error) {
	r := bufio.NewReader(os.Stdin)
	if username == "" {
		fmt.Fprint(log.Out, "GitHub username: ")
		username, _ = r.ReadString('\n')
	}
	if password == "" {
		fmt.Fprint(log.Out, "GitHub password: ")
		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
		password = string(bytePassword)
		fmt.Fprintln(log.Out, "")
	}

	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	client := github.NewClient(tp.Client())

	// call api for checking OTP
	_, _, err := client.APIMeta(ctx)
	if err == nil {
		return client, nil
	}

	if _, ok := err.(*github.TwoFactorAuthError); ok {
		fmt.Fprint(log.Out, "GitHub two-factor authentication code: ")
		otp, _ := r.ReadString('\n')
		tp.OTP = strings.TrimSpace(otp)
		if _, _, err = client.APIMeta(ctx); err == nil {
			return client, nil
		}
	}

	return nil, err
}

// NewTokenClient 產生一個 Token 的 GitHub Client
func NewTokenClient(ctx context.Context, token string) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), nil
}
