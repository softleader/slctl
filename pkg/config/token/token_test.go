package token

import (
	"github.com/google/go-github/github"
	"testing"
)

func TestRemove(t *testing.T) {
	drop := github.ScopeGist
	a := []github.Scope{github.ScopeReadOrg, drop, github.ScopeRepo}
	a = remove(a, drop)

	for _, v := range a {
		if v == drop {
			t.Errorf("expected %q not contains %q", a, drop)
		}
	}

}
