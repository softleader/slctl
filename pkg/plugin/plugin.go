package plugin

import (
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

// Metadata describes a plugin.
//
// This is the plugin equivalent of a chart.Metadata.
type Metadata struct {
	// Name is the name of the plugin
	Name string `json:"name"`

	// Version is a SemVer 2 version of the plugin.
	Version string `json:"version"`

	// Usage is the single-line usage text shown in help
	Usage string `json:"usage"`

	// Description is a long description shown in places like `helm help`
	Description string `json:"description"`

	// Command is the command, as a single string.
	//
	// The command will be passed through environment expansion, so env vars can
	// be present in this command. Unless IgnoreFlags is set, this will
	// also merge the flags passed from Helm.
	//
	// Note that command is not executed in a shell. To do so, we suggest
	// pointing the command to a shell script.
	Command string `json:"command"`

	// IgnoreFlags ignores any flags passed in from Helm
	//
	// For example, if the plugin is invoked as `helm --debug myplugin`, if this
	// is false, `--debug` will be appended to `--command`. If this is true,
	// the `--debug` flag will be discarded.
	IgnoreGlobalFlags bool `json:"ignoreGlobalFlags"`
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
func (p *Plugin) PrepareCommand(extraArgs []string) (string, []string) {
	parts := strings.Split(os.ExpandEnv(p.Metadata.Command), " ")
	main := parts[0]
	baseArgs := []string{}
	if len(parts) > 1 {
		baseArgs = parts[1:]
	}
	if !p.Metadata.IgnoreGlobalFlags {
		baseArgs = append(baseArgs, extraArgs...)
	}
	return main, baseArgs
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
	plugins := []*Plugin{}
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

// FindPlugins returns a list of YAML files that describe plugins.
//func FindPlugins(plugdirs string) ([]*Plugin, error) {
//	found := []*Plugin{}
//	// Let's get all UNIXy and allow path separators
//	for _, p := range filepath.SplitList(plugdirs) {
//		matches, err := LoadAll(p)
//		if err != nil {
//			return matches, err
//		}
//		found = append(found, matches...)
//	}
//	return found, nil
//}

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
