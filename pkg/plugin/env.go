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
	Envs = envs()
)

// 載入 plugin env
func (p *Plugin) SetupEnv(env *environment.EnvSettings, metadata *version.BuildMetadata) (err error) {
	var conf *config.ConfFile
	if conf, err = config.LoadConfFile(env.Home.ConfigFile()); err != nil && err != config.ErrTokenNotExist {
		return err
	}
	paths.EnsureDirectory(logrus.StandardLogger(), p.Mount)
	for key, val := range envsMap(p.Metadata.Name, p.Dir, p.Mount, metadata.String(), conf.Token) {
		os.Setenv(key, val)
	}
	return nil
}

func envsMap(pluginName, pluginDir, pluginMount, version, token string) (e map[string]string) {
	e = map[string]string{
		"SL_CLI":          "slctl",
		"SL_VERSION":      version,
		"SL_PLUGIN_NAME":  pluginName,
		"SL_PLUGIN_DIR":   pluginDir,
		"SL_PLUGIN_MOUNT": pluginMount,
		"SL_BIN":          os.Args[0],
		"SL_PLUGIN":       environment.Settings.PluginDirs(),
		"SL_HOME":         environment.Settings.Home.String(),
		"SL_VERBOSE":      strconv.FormatBool(environment.Settings.Verbose),
		"SL_OFFLINE":      strconv.FormatBool(environment.Settings.Offline),
		"SL_TOKEN":        token,
	}
	return
}

func envs() (m map[string]string) {
	plug := "foo"
	h := environment.Settings.Home
	return envsMap(
		plug,
		filepath.Join(h.Plugins(), plug),
		filepath.Join(h.Mounts(), plug),
		"<version.of.slctl>",
		"<github.token>")
}
