package plugin

import (
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/slpath"
	"testing"
)

func TestFetchOnline(t *testing.T) {
	//home, err := ioutil.TempDir("", "sl_home")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//defer os.RemoveAll(home)
	//hh := slpath.Home(home)
	userHome, err := homedir.Dir()
	if err != nil {
		t.SkipNow()
	}
	hh := slpath.Home(environment.DefaultHome(userHome))
	r, err := fetchOnline(logrus.StandardLogger(), hh, "softleader")
	if err != nil {
		t.Error(err)
	}
	if l := len(r.Repos); l < 3 {
		t.Errorf("should be a least 3 official plugins, but got %v", l)
	}
}
