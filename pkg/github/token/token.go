package token

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/v21/github"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	gh "github.com/softleader/slctl/pkg/github"
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

	ctx := context.Background()
	client, err := gh.NewBasicAuthClient(ctx, log, "", "")
	if err != nil {
		return err
	}
	auths, _, err := client.Authorizations.List(ctx, &github.ListOptions{})
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
func Grant(ctx context.Context, client *github.Client, log *logrus.Logger, force bool) (token string, err error) {
	auths, _, err := client.Authorizations.List(ctx, &github.ListOptions{})
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
func Confirm(ctx context.Context, client *github.Client, org string, _ *logrus.Logger) (name string, err error) {
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
