package plugin

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/version"
	"os"
	"path/filepath"
	"strconv"
)

var (
	// Envs expose plugin envs 範例讓 command 可以輸出給使用者參考
	Envs = func() (m map[string]string) {
		plugName := "foo"
		return envsMap(
			plugName,
			filepath.Join(environment.Settings.Home.Plugins(), plugName),
			filepath.Join(environment.Settings.Home.Mounts(), plugName),
			"<version.of.slctl>",
			"<github.token>")
	}()
)

// SetupEnv 載入 plugin env
func (p *Plugin) SetupEnv(metadata *version.BuildMetadata) (err error) {
	var conf *config.ConfFile
	if conf, err = config.LoadConfFile(environment.Settings.Home.ConfigFile()); err != nil && err != config.ErrTokenNotExist {
		return err
	}
	paths.EnsureDirectory(logrus.StandardLogger(), p.Mount)
	m := envsMap(
		p.Metadata.Name,
		p.Dir,
		p.Mount,
		metadata.String(),
		conf.Token)
	for key, val := range m {
		os.Setenv(key, val)
	}
	return nil
}

// 抽一層 func 只是為了 Envs 可以拿到一樣的 map 而已
func envsMap(pluginName, pluginDir, pluginMount, version, token string) (e map[string]string) {
	e = map[string]string{
		"SL_BIN":          os.Args[0],
		"SL_CLI":          "slctl",
		"SL_PLUGIN_NAME":  pluginName,
		"SL_PLUGIN_DIR":   pluginDir,
		"SL_PLUGIN_MOUNT": pluginMount,
		"SL_VERSION":      version,
		"SL_PLUGIN":       environment.Settings.PluginDirs(),
		"SL_HOME":         environment.Settings.Home.String(),
		"SL_VERBOSE":      strconv.FormatBool(environment.Settings.Verbose),
		"SL_OFFLINE":      strconv.FormatBool(environment.Settings.Offline),
		"SL_TOKEN":        token,
	}
	return
}
