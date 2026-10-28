package plugin

import "github.com/google/go-github/v93/github"

// GitHub 封裝了 GitHub-Repo 的相關資訊
type GitHub struct {
	Scopes []github.Scope `json:"scopes"`
}
