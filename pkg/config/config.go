package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const ReadWrite = 0644

var ErrTokenNotExist = errors.New("token is no exist")

type ConfigFile struct {
	Token string `json:"token"`
}

func NewConfigFile() *ConfigFile {
	return &ConfigFile{}
}

func LoadConfigFile(path string) (*ConfigFile, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(
				"Couldn't load config file (%s).\n"+
					"You might need to run `slctl init --help`", path)
		}
		return nil, err
	}

	conf := &ConfigFile{}
	err = yaml.Unmarshal(b, conf)
	if err != nil {
		return nil, err
	}

	if conf.Token == "" {
		return conf, ErrTokenNotExist
	}

	return conf, nil
}

func (r *ConfigFile) WriteFile(path string, perm os.FileMode) error {
	data, err := yaml.Marshal(r)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, perm)
}
