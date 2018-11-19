package cmd

import (
	"bytes"
	"fmt"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/slpath"
	"io/ioutil"
	"os"
	"testing"
)

// need read:org & read:user permission
func TestConfirmToken(t *testing.T) {
	settings.Debug = true
	b := bytes.NewBuffer(nil)
	token := "997f19253fccc351bfcf4cf1622f494f7708522a"
	var err error
	var name string

	if name, err = confirmToken(token, b); err != nil {
		t.Error(err)
	}
	if name == "" {
		t.Errorf("name should not be empty")
	}

	fmt.Printf("Hello, %s!\n", name)
}

func TestRefreshConfig(t *testing.T) {
	home, err := ioutil.TempDir("", "sl_home")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(home)

	b := bytes.NewBuffer(nil)
	hh := slpath.Home(home)
	settings.Home = hh

	settings.Debug = true

	if err = ensureDirectories(hh, b); err != nil {
		t.Error(err)
	}
	if err := ensureConfigFile(hh, b); err != nil {
		t.Error(err)
	}

	token := "this.is.a.fake.token"
	if err = refreshConfig(hh, token, b); err != nil {
		t.Error(err)
	}

	var conf *config.ConfigFile
	if conf, err = config.LoadConfigFile(hh.ConfigFile()); err != nil {
		t.Error(err)
	}
	if conf.Token != token {
		t.Errorf("expected token to be %s", conf.Token)
	}
}
