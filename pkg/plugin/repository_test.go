package plugin

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
)

func TestFetchOnline(t *testing.T) {
	uh, err := homedir.Dir()
	if err != nil {
		t.SkipNow()
	}
	h := environment.DefaultHome(uh)
	if _, err := os.Stat(h); os.IsNotExist(err) {
		t.SkipNow()
	}
	hh := paths.Home(h)
	r, err := fetchOnline(context.Background(), logrus.StandardLogger(), hh, "softleader", nil)
	if err != nil {
		if strings.Contains(err.Error(), "401 Bad credentials") || strings.Contains(err.Error(), "token not exist") {
			t.Skipf("maybe just token not set")
		}
		t.Error(err)
		return
	}
	if l := len(r.Repos); l < 3 {
		t.Errorf("should be a least 3 official plugins, but got %v", l)
	}
	for _, repo := range r.Repos {
		fmt.Println(repo)
	}
}
