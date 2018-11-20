package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const ReadWrite = 0644

var ErrTokenNotExist = errors.New(`token not exist.
You might need to run 'slctl init'
`)

type ConfFile struct {
	Token string `json:"token"`
}

func NewConfFile() *ConfFile {
	return &ConfFile{}
}

func LoadConfFile(path string) (*ConfFile, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(
				"Couldn't load config file (%s).\n"+
					"You might need to run `slctl init`", path)
		}
		return nil, err
	}

	conf := &ConfFile{}
	err = yaml.Unmarshal(b, conf)
	if err != nil {
		return nil, err
	}

	if conf.Token == "" {
		return conf, ErrTokenNotExist
	}

	return conf, nil
}

func (c *ConfFile) WriteFile(path string, perm os.FileMode) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, perm)
}
