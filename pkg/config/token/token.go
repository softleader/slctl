package token

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v21/github"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"
	"os"
	"strings"
	"syscall"
)

const (
	note = "slctl token (https://github.com/softleader/slctl)"
)

var (
	// Scopes 表示此 app 預設需要的 scopes
	Scopes = []github.Scope{github.ScopeReadOrg, github.ScopeUser, github.ScopeRepo}
	// ErrOauthAccessAlreadyExists 表示 token 已存在
	ErrOauthAccessAlreadyExists = errors.New(`access token already exists
To store a token on https://github.com/settings/tokens, use '--token' flag  
To re-generate a new token, use '--force' flag
Use 'init --help' for more information about the command`)
)

// EnsureScopes 確保當前的 token 有傳入的所有 scopes
func EnsureScopes(log *logrus.Logger, scopes []github.Scope) (err error) {
	if environment.Settings.Offline {
		return
	}
	var addScopes []github.Scope
	for _, scope := range scopes {
		if !contains(Scopes, scope) {
			addScopes = append(addScopes, scope)
		}
	}
	if len(addScopes) == 0 {
		// if goes here, plugin's scopes are the same as 'slctl init scopes', so we don't have to against GitHub api
		return
	}

	fmt.Fprintf(log.Out, "Checking authorization scopes %q for the GitHub access token\n", scopes)

	r := bufio.NewReader(os.Stdin)

	fmt.Fprint(log.Out, "GitHub username: ")
	username, _ := r.ReadString('\n')

	fmt.Fprint(log.Out, "GitHub password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)
	fmt.Fprintln(log.Out, "")

	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	client := github.NewClient(tp.Client())
	ctx := context.Background()

	auths, _, err := client.Authorizations.List(ctx, &github.ListOptions{})
	if _, ok := err.(*github.TwoFactorAuthError); ok {
		fmt.Fprint(log.Out, "GitHub two-factor authentication code: ")
		otp, _ := r.ReadString('\n')
		tp.OTP = strings.TrimSpace(otp)
		if auths, _, err = client.Authorizations.List(ctx, &github.ListOptions{}); err != nil {
			return
		}
	}
	auth := findAuth(auths, note)
	if auth == nil {
		return fmt.Errorf(
			"Couldn't find access token with note: %q.\n"+
				"You might need to run `slctl init`", note)
	}
	for _, scope := range auth.Scopes {
		addScopes = remove(addScopes, scope)
	}
	if len(addScopes) == 0 {
		// if goes here, plugin's scopes are already granted for the token
		return
	}

	fmt.Fprintf(log.Out, "granting scopes: %q\n", addScopes)

	_, _, err = client.Authorizations.Edit(ctx, auth.GetID(), newAuthorizationUpdateRequest(addScopes))
	return
}

func newAuthorizationUpdateRequest(scopes []github.Scope) (r *github.AuthorizationUpdateRequest) {
	r = &github.AuthorizationUpdateRequest{
		AddScopes: make([]string, len(scopes)),
	}
	for i, s := range scopes {
		r.AddScopes[i] = string(s)
	}
	return
}

func findAuth(auths []*github.Authorization, note string) *github.Authorization {
	for _, auth := range auths {
		if auth.GetNote() == note {
			return auth
		}
	}
	return nil
}

func remove(a []github.Scope, s github.Scope) []github.Scope {
	var i *int
	for idx, v := range a {
		if v == s {
			i = &idx
			break
		}
	}
	if i == nil {
		return a
	}
	return append(a[:*i], a[*i+1:]...)
}

func contains(base []github.Scope, target github.Scope) bool {
	for _, v := range base {
		if v == target {
			return true
		}
	}
	return false
}

// Grant 產生傳入的 username/password 的 token
func Grant(username, password string, log *logrus.Logger, force bool) (token string, err error) {
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
	ctx := context.Background()

	auths, _, err := client.Authorizations.List(ctx, &github.ListOptions{})
	if _, ok := err.(*github.TwoFactorAuthError); ok {
		fmt.Fprint(log.Out, "GitHub two-factor authentication code: ")
		otp, _ := r.ReadString('\n')
		tp.OTP = strings.TrimSpace(otp)
		if auths, _, err = client.Authorizations.List(ctx, &github.ListOptions{}); err != nil {
			return
		}
	}

	for _, auth := range auths {
		if auth.GetNote() == note {
			if !force {
				return "", ErrOauthAccessAlreadyExists
			}
			log.Debugf("Removing exist token\n")
			if _, err = client.Authorizations.Delete(ctx, auth.GetID()); err != nil {
				return "", err
			}
			break
		}
	}

	n := note
	auth, _, err := client.Authorizations.Create(ctx, &github.AuthorizationRequest{
		Scopes: Scopes,
		Note:   &n,
	})
	if err != nil {
		return "", err
	}

	return auth.GetToken(), nil
}

// Confirm 確保 token 的使用者存在於傳入的 org 中
func Confirm(org, token string, _ *logrus.Logger) (name string, err error) {
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
	if mem, _, err = client.Organizations.GetOrgMembership(ctx, "", org); err != nil {
		return "", err
	}
	// log.Debugf("%s\n", github.Stringify(mem))
	if mem.GetState() != "active" {
		return "", fmt.Errorf("you are not a active member of %s", org)
	}
	var user *github.User
	if user, _, err = client.Users.Get(ctx, ""); err != nil {
		return "", err
	}
	// log.Debugf("%s\n", github.Stringify(user))
	return user.GetName(), err
}
