package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/paths"
	"gopkg.in/yaml.v2"
)

const (
	// CleanupDueDays 代表每次自動執行 cleanup 的間隔天數
	CleanupDueDays = 30
	// CheckUpdatesDueDays 代表每次自動執行檢查版本的間隔天數
	CheckUpdatesDueDays = 15
)

// ErrTokenNotExist 代表了 GitHub Token 在不存在於 config 中
var ErrTokenNotExist = errors.New(`token not exist.
You might need to run 'slctl init'`)

// ConfFile 內容等於 $SL_HOME/config/configs.yaml
type ConfFile struct {
	Token        string    `json:"token"`        // github token
	Cleanup      time.Time `json:"cleanup"`      // 下次要執行 cleanup 的時間
	CheckUpdates time.Time `json:"checkUpdates"` // 下次要執行檢查版本的時間
}

// NewConfFile return a pointer of a blank ConfFile
func NewConfFile() (c *ConfFile) {
	c = &ConfFile{}
	c.UpdateCheckUpdatesTime()
	c.UpdateCleanupTime()
	return
}

// UpdateCleanupTime updates cleanup time
func (c *ConfFile) UpdateCleanupTime() {
	c.Cleanup = time.Now().AddDate(0, 0, CleanupDueDays)
}

// UpdateCheckUpdatesTime updates check updates time
func (c *ConfFile) UpdateCheckUpdatesTime() {
	c.UpdateCheckUpdatesTimeInDays(CheckUpdatesDueDays)
}

// UpdateCheckUpdatesTimeInDays updates check updates time
func (c *ConfFile) UpdateCheckUpdatesTimeInDays(days int) {
	c.CheckUpdates = time.Now().AddDate(0, 0, days)
}

// LoadConfFile return a pointer of a ConfFile which read from path
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

// WriteFile 將當前的 ConfFile 內容寫入 path
func (c *ConfFile) WriteFile(path string, perm os.FileMode) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, perm)
}

// Refresh 將傳入的 token 寫回檔案中
func Refresh(home paths.Home, token string, _ *logrus.Logger) (err error) {
	conf, err := LoadConfFile(home.ConfigFile())
	if err != nil && err != ErrTokenNotExist {
		return fmt.Errorf("failed to load file (%v)", err)
	}
	conf.Token = token

	return conf.WriteFile(home.ConfigFile(), 0644)
}
