package config

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestWriteFile(t *testing.T) {
	cf := NewConfFile()
	cf.Token = "this.is.a.fake.token"

	repoFile, err := ioutil.TempFile("", "sl-config")
	if err != nil {
		t.Errorf("failed to create test-file (%v)", err)
	}
	defer os.Remove(repoFile.Name())

	if err := cf.WriteFile(repoFile.Name(), 0644); err != nil {
		t.Errorf("failed to write file (%v)", err)
	}

	loaded, err := LoadConfFile(repoFile.Name())
	if err != nil {
		t.Errorf("failed to load file (%v)", err)
	}
	if loaded.Token != cf.Token {
		t.Errorf("expected token to be %s", cf.Token)
	}
}

func TestConfigNotExists(t *testing.T) {
	_, err := LoadConfFile("/this/path/does/not/exist.yaml")
	if err == nil {
		t.Errorf("expected err to be non-nil when path does not exist")
	} else if !strings.Contains(err.Error(), "You might need to run `slctl init`") {
		t.Errorf("expected prompt to run `slctl init` when config file does not exist")
	}
}
