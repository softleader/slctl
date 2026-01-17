package github

import (
	"context"

	"github.com/google/go-github/v69/github"
	"golang.org/x/oauth2"
)

// NewTokenClient 產生一個 Token 的 GitHub Client
func NewTokenClient(ctx context.Context, token string) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), nil
}
