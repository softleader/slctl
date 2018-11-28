package plugin

import "github.com/google/go-github/github"

type GitHub struct {
	Scopes []github.Scope `json:"scopes"`
}
