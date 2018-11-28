package plugin

import (
	"github.com/google/go-github/github"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/environment"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const MetadataFileName = "metadata.yaml"

type Metadata struct {
	Name              string         `json:"name"`
	Version           string         `json:"version"`
	Usage             string         `json:"usage"`
	Description       string         `json:"description"`
	Command           Commands       `json:"command"`
	Hook              Commands       `json:"hook"`
	IgnoreGlobalFlags bool           `json:"ignoreGlobalFlags"`
	Scopes            []github.Scope `json:"scopes"`
}

type Plugin struct {
	Metadata *Metadata
	Dir      string
}

// PrepareCommand takes a Plugin.Command and prepares it for execution.
//
// It merges extraArgs into any arguments supplied in the plugin. It
// returns the name of the command and an args array.
//
// The result is suitable to pass to exec.Command.
func (p *Plugin) PrepareCommand(extraArgs []string) (main string, argv []string, err error) {
	command, err := p.Metadata.Command.GetCommand()
	if err != nil {
		return
	}
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

// LoadDir loads a plugin from the given directory.
func LoadDir(dirname string) (*Plugin, error) {
	data, err := ioutil.ReadFile(filepath.Join(dirname, MetadataFileName))
	if err != nil {
		return nil, err
	}

	plug := &Plugin{Dir: dirname}
	if err := yaml.Unmarshal(data, &plug.Metadata); err != nil {
		return nil, err
	}
	return plug, nil
}

// LoadAll loads all plugins found beneath the base directory.
//
// This scans only one directory level.
func LoadAll(basedir string) ([]*Plugin, error) {
	var plugins []*Plugin
	// We want basedir/*/plugin.yaml
	scanpath := filepath.Join(basedir, "*", MetadataFileName)
	matches, err := filepath.Glob(scanpath)
	if err != nil {
		return plugins, err
	}

	if matches == nil {
		return plugins, nil
	}

	for _, yaml := range matches {
		dir := filepath.Dir(yaml)
		p, err := LoadDir(dir)
		if err != nil {
			return plugins, err
		}
		plugins = append(plugins, p)
	}
	return plugins, nil
}

// SetupPluginEnv prepares os.Env for plugins. It operates on os.Env because
// the plugin subsystem itself needs access to the environment variables
// created here.
func SetupPluginEnv(
	shortName, base string) (err error) {
	var conf *config.ConfFile
	if conf, err = config.LoadConfFile(environment.Settings.Home.ConfigFile()); err != nil && err != config.ErrTokenNotExist {
		return err
	}

	for key, val := range map[string]string{
		"SL_PLUGIN_NAME": shortName,
		"SL_PLUGIN_DIR":  base,
		"SL_BIN":         os.Args[0],

		// Set vars that may not have been set, and save client the
		// trouble of re-parsing.
		"SL_PLUGIN":  environment.Settings.PluginDirs(),
		"SL_HOME":    environment.Settings.Home.String(),
		"SL_VERBOSE": strconv.FormatBool(environment.Settings.Verbose),
		"SL_OFFLINE": strconv.FormatBool(environment.Settings.Offline),
		"SL_TOKEN":   conf.Token,
	} {
		os.Setenv(key, val)
	}

	return nil
}
