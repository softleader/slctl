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
	// expose 一個 plugin envs 範例讓 command 可以輸出給使用者參考
	Envs = func() (m map[string]string) {
		plugName := "foo"
		e := environment.Settings
		return envsMap(
			plugName,
			filepath.Join(e.Home.Plugins(), plugName),
			filepath.Join(e.Home.Mounts(), plugName),
			"<version.of.slctl>",
			e.PluginDirs(),
			e.Home.String(),
			e.Verbose,
			e.Offline,
			"<github.token>")
	}()
)

// 載入 plugin env
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
		environment.Settings.PluginDirs(),
		environment.Settings.Home.String(),
		environment.Settings.Verbose,
		environment.Settings.Offline,
		conf.Token)
	for key, val := range m {
		os.Setenv(key, val)
	}
	return nil
}

// 抽一層 func 只是為了 Envs 可以拿到一樣的 map 而已
func envsMap(pluginName, pluginDir, pluginMount, version, plugin, home string, verbose, offline bool, token string) (e map[string]string) {
	e = map[string]string{
		"SL_BIN":          os.Args[0],
		"SL_CLI":          "slctl",
		"SL_PLUGIN_NAME":  pluginName,
		"SL_PLUGIN_DIR":   pluginDir,
		"SL_PLUGIN_MOUNT": pluginMount,
		"SL_VERSION":      version,
		"SL_PLUGIN":       plugin,
		"SL_HOME":         home,
		"SL_VERBOSE":      strconv.FormatBool(verbose),
		"SL_OFFLINE":      strconv.FormatBool(offline),
		"SL_TOKEN":        token,
	}
	return
}
