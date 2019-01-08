package plugin

import "github.com/google/go-github/v21/github"

type GitHub struct {
	Scopes []github.Scope `json:"scopes"`
}
