package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"io/ioutil"
	"os"
	"testing"
)

func TestRefreshConfig(t *testing.T) {
	home, err := ioutil.TempDir("", "sl_home")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(home)

	log := logrus.New()
	hh := paths.Home(home)
	environment.Settings.Home = hh

	environment.Settings.Verbose = true

	if err = ensureDirectories(hh, log); err != nil {
		t.Error(err)
	}
	if err := ensureConfigFile(hh, log); err != nil {
		t.Error(err)
	}

	token := "this.is.a.fake.token"
	if err = config.Refresh(hh, token, log); err != nil {
		t.Error(err)
	}

	var conf *config.ConfFile
	if conf, err = config.LoadConfFile(hh.ConfigFile()); err != nil {
		t.Error(err)
	}
	if conf.Token != token {
		t.Errorf("expected token to be %s", conf.Token)
	}
}
