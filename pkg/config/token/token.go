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

	fmt.Fprintf(out, "Checking authorization scopes %q for the GitHub access token\n", scopes)

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

	fmt.Fprintf(out, "granting scopes: %q\n", addScopes)

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
