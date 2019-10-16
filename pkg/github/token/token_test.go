package token

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-github/v21/github"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	gh "github.com/softleader/slctl/pkg/github"
)

// need read:org & read:user permission
func TestConfirm(t *testing.T) {
	environment.Settings.Verbose = true
	token, found := os.LookupEnv("GITHUB_TOKEN_TEST")
	if !found {
		t.Skipf("provide $GITHUB_TOKEN_TEST to run the test")
	}
	var name string
	client, err := gh.NewTokenClient(context.Background(), token)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if name, err = Confirm(context.Background(), client, "text", logrus.StandardLogger()); err != nil {
		t.Error(err)
	}
	if name == "" {
		t.Errorf("name should not be empty")
	}
	fmt.Printf("Hello, %s!\n", name)
}

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
