package token

import "github.com/google/go-github/github"

var Scopes = []github.Scope{github.ScopeReadOrg, github.ScopeUser}