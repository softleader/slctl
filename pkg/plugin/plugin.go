package plugin

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/environment"
	"github.com/softleader/slctl/pkg/paths"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const MetadataFileName = "metadata.yaml"
const SourceFileName = ".source"

var (
	Envs = envs()
)

type Metadata struct {
	Name              string   `json:"name"`
	Version           string   `json:"version"`
	Usage             string   `json:"usage"`
	Description       string   `json:"description"`
	Exec              Commands `json:"exec"`
	Hook              Commands `json:"hook"`
	IgnoreGlobalFlags bool     `json:"ignoreGlobalFlags"`
	GitHub            GitHub   `json:"github"`
}

type Plugin struct {
	Metadata *Metadata
	Dir      string
	Source   string // 在安裝 plugin 時的 source, 只有非本機的 source 才會紀錄, 為了方便之後做 github plugin 的 upgrade
	Mount    string
}

func (p *Plugin) FromGitHub() bool {
	return strings.HasPrefix(p.Source, "github.com/")
}

// PrepareCommand takes a Plugin.Command and prepares it for execution.
//
// It merges extraArgs into any arguments supplied in the plugin. It
// returns the name of the command and an args array.
//
// The result is suitable to pass to exec.Command.
func (p *Plugin) PrepareCommand(command string, extraArgs []string) (main string, argv []string, err error) {
	parts := strings.Split(os.ExpandEnv(command), " ")
	main = parts[0]
	if len(parts) > 1 {
		argv = parts[1:]
	}
	if !p.Metadata.IgnoreGlobalFlags && extraArgs != nil {
		argv = append(argv, extraArgs...)
	}
	return
}

// SetupPluginEnv prepares os.Env for plugins. It operates on os.Env because
// the plugin subsystem itself needs access to the environment variables
// created here.
func SetupPluginEnv(
	plugName, plugDir, plugMount, cli, version string) (err error) {
	var conf *config.ConfFile
	if conf, err = config.LoadConfFile(environment.Settings.Home.ConfigFile()); err != nil && err != config.ErrTokenNotExist {
		return err
	}
	paths.EnsureDirectory(logrus.StandardLogger(), plugMount)
	for key, val := range pluginEnv(plugName, plugDir, plugMount, cli, version, conf.Token) {
		os.Setenv(key, val)
	}

	return nil
}

func pluginEnv(pluginName, pluginDir, pluginMount, cli, version, token string) map[string]string {
	return map[string]string{
		"SL_CLI":          cli,
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
}

func envs() (m map[string]string) {
	plug := "foo"
	h := environment.Settings.Home
	return pluginEnv(
		plug,
		filepath.Join(h.Plugins(), plug),
		filepath.Join(h.Mounts(), plug),
		"slctl",
		"<version.of.slctl>",
		"<github.token>")
}
