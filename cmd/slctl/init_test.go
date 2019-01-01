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

// need read:org & read:user permission
func TestConfirmToken(t *testing.T) {
	//environment.SettingsVerbose = true
	//b := bytes.NewBuffer(nil)
	//token := "997f19253fccc351bfcf4cf1622f494f7708522a"
	//var err error
	//var name string
	//
	//if name, err = confirmToken(token, b); err != nil {
	//	t.Error(err)
	//}
	//if name == "" {
	//	t.Errorf("name should not be empty")
	//}
	//
	//fmt.Printf("Hello, %s!\n", name)
}

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
