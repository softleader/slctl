package token

import (
	"fmt"
	"github.com/google/go-github/v21/github"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"os"
	"testing"
)

// need read:org & read:user permission
func TestConfirm(t *testing.T) {
	environment.Settings.Verbose = true
	token, found := os.LookupEnv("GITHUB_TOKEN_TEST")
	if !found {
		t.Skipf("provide $GITHUB_TOKEN_TEST to run the test")
	}
	var err error
	var name string
	if name, err = Confirm("test", token, logrus.StandardLogger()); err != nil {
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
