package plugin

import (
	"github.com/ghodss/yaml"
	"github.com/softleader/slctl/pkg/config"
	"github.com/softleader/slctl/pkg/environment"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/softleader/slctl/pkg/v"
	"errors"
)

const MetadataFileName = "metadata.yaml"

// Downloaders represents the plugins capability if it can retrieve
// charts from special sources
type Downloaders struct {
	// Protocols are the list of schemes from the charts URL.
	Protocols []string `json:"protocols"`
	// Command is the executable path with which the plugin performs
	// the actual download for the corresponding Protocols
	Command string `json:"command"`
}

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
	IgnoreFlags bool `json:"ignoreFlags"`

	// UseTunnel indicates that this command needs a tunnel.
	// Setting this will cause a number of side effects, such as the
	// automatic setting of HELM_HOST.
	UseTunnel bool `json:"useTunnel"`

	// Hooks are commands that will run on events.
	Hooks Hooks

	// Downloaders field is used if the plugin supply downloader mechanism
	// for special protocols.
	Downloaders []Downloaders `json:"downloaders"`
}

// Plugin represents a plugin.
type Plugin struct {
	// Metadata is a parsed representation of a plugin.yaml
	Metadata *Metadata
	// Dir is the string path to the directory that holds the plugin.
	Dir string
}

func (p *Plugin) LinkTo(home slpath.Home) (string, error) {
	linked := filepath.Join(home.Plugins(), p.Metadata.Name)
	if _, err := os.Stat(linked); !os.IsNotExist(err) {
		return "", errors.New("plugin already exists")
	}
	v.Printf("symlinking %s to %s", p.Dir, linked)
	if err := os.Symlink(p.Dir, linked); err != nil {
		return "", err
	}
	return linked, nil
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
	if !p.Metadata.IgnoreFlags {
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
func FindPlugins(plugdirs string) ([]*Plugin, error) {
	found := []*Plugin{}
	// Let's get all UNIXy and allow path separators
	for _, p := range filepath.SplitList(plugdirs) {
		matches, err := LoadAll(p)
		if err != nil {
			return matches, err
		}
		found = append(found, matches...)
	}
	return found, nil
}

// SetupPluginEnv prepares os.Env for plugins. It operates on os.Env because
// the plugin subsystem itself needs access to the environment variables
// created here.
func SetupPluginEnv(settings environment.EnvSettings,
	shortName, base string) (err error) {
	var conf *config.ConfFile
	if conf, err = config.LoadConfFile(settings.Home.ConfigFile()); err != nil && err != config.ErrTokenNotExist {
		return err
	}

	for key, val := range map[string]string{
		"SL_PLUGIN_NAME": shortName,
		"SL_PLUGIN_DIR":  base,
		"SL_BIN":         os.Args[0],

		// Set vars that may not have been set, and save client the
		// trouble of re-parsing.
		"SL_PLUGIN":  settings.PluginDirs(),
		"SL_HOME":    settings.Home.String(),
		"SL_VERBOSE": strconv.FormatBool(settings.Verbose),
		"SL_OFFLINE": strconv.FormatBool(settings.Offline),
		"SL_TOKEN":   conf.Token,
	} {
		os.Setenv(key, val)
	}

	return nil
}
