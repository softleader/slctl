package token

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/v"
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
)

var (
	Scopes                      = []github.Scope{github.ScopeReadOrg, github.ScopeUser}
	ErrOauthAccessAlreadyExists = errors.New(`access token already exists.
To store a token already on https://github.com/settings/tokens, use '--token' flag.  
To automatically re-generate a new one, use '--refresh' flag.
Use 'init --help' for more information about the command.`)
)

func EnsureScopes(out io.Writer, scopes []github.Scope) (err error) {
	if environment.Settings.Offline {
		return nil
	}
	var addScopes []string
	for _, scope := range scopes {
		if !contains(Scopes, scope) {
			addScopes = append(addScopes, string(scope))
		}
	}
	if len(addScopes) <= 0 {
		return nil
	}

	v.Fprintf(out, "granting more scopes: %q\n", addScopes)

	r := bufio.NewReader(os.Stdin)

	fmt.Fprint(out, "GitHub Username: ")
	username, _ := r.ReadString('\n')

	fmt.Fprint(out, "GitHub Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)

	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	client := github.NewClient(tp.Client())
	ctx := context.Background()

	auths, _, err := client.Authorizations.List(ctx, &github.ListOptions{})
	if _, ok := err.(*github.TwoFactorAuthError); ok {
		fmt.Fprint(out, "\nGitHub OTP: ")
		otp, _ := r.ReadString('\n')
		tp.OTP = strings.TrimSpace(otp)
		if auths, _, err = client.Authorizations.List(ctx, &github.ListOptions{}); err != nil {
			return
		}
	}
	for _, auth := range auths {
		if auth.GetNote() == note {
			_, _, err = client.Authorizations.Edit(ctx, auth.GetID(), &github.AuthorizationUpdateRequest{
				AddScopes: addScopes,
			})
			return
		}
	}
	return fmt.Errorf(
		"Couldn't find access token with note: '%s'.\n"+
			"You might need to run `slctl init`", note)
}

func contains(base []github.Scope, target github.Scope) bool {
	for _, v := range base {
		if v == target {
			return true
		}
	}
	return false
}

func Grant(username, password string, out io.Writer, refresh bool) (token string, err error) {
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

	auths, _, err := client.Authorizations.List(ctx, &github.ListOptions{})
	if _, ok := err.(*github.TwoFactorAuthError); ok {
		fmt.Fprint(out, "\nGitHub OTP: ")
		otp, _ := r.ReadString('\n')
		tp.OTP = strings.TrimSpace(otp)
		if auths, _, err = client.Authorizations.List(ctx, &github.ListOptions{}); err != nil {
			return
		}
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

func Confirm(token string, out io.Writer) (name string, err error) {
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
