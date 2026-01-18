package token

import (
	"context"
	"errors"
	"fmt"

	"strings"

	"github.com/google/go-github/v69/github"
	"github.com/sirupsen/logrus"
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
func EnsureScopes(ctx context.Context, client *github.Client, log *logrus.Logger, scopes []github.Scope) (err error) {
	if len(scopes) == 0 {
		return nil
	}
	_, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		return err
	}
	granted := resp.Header.Get("X-OAuth-Scopes")
	log.Debugf("granted scopes: %s", granted)

	missing := []string{}
	for _, s := range scopes {
		if !strings.Contains(granted, string(s)) {
			missing = append(missing, string(s))
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing scopes: %v", missing)
	}
	return nil
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
