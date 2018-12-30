package plugin

import (
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/slpath"
	"os"
	"testing"
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
	hh := slpath.Home(h)
	r, err := fetchOnline(logrus.StandardLogger(), hh, "softleader")
	if err != nil {
		t.Error(err)
	}
	if l := len(r.Repos); l < 3 {
		t.Errorf("should be a least 3 official plugins, but got %v", l)
	}
}
