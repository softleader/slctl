package environment

import (
	"github.com/softleader/slctl/pkg/homedir"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

var DefaultHome = filepath.Join(homedir.HomeDir(), ".sl")

type EnvSettings struct {
	Home    slpath.Home
	Verbose bool
}

// AddFlags binds flags to the given flagset.
func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar((*string)(&s.Home), "home", DefaultHome, "location of your config. Overrides $SL_HOME")
	fs.BoolVarP(&s.Verbose, "verbose", "v", false, "enable verbose output")
}

// Init sets values from the environment.
func (s *EnvSettings) Init(fs *pflag.FlagSet) {
	for name, envar := range envMap {
		setFlagFromEnv(name, envar, fs)
	}
}

// envMap maps flag names to envvars
var envMap = map[string]string{
	"debug": "SL_DEBUG",
	"home":  "SL_HOME",
}

func (s EnvSettings) PluginDirs() string {
	if d, ok := os.LookupEnv("SL_PLUGIN"); ok {
		return d
	}
	return s.Home.Plugins()
}

// setFlagFromEnv looks up and sets a flag if the corresponding environment variable changed.
// if the flag with the corresponding name was set during fs.Parse(), then the environment
// variable is ignored.
func setFlagFromEnv(name, envar string, fs *pflag.FlagSet) {
	if fs.Changed(name) {
		return
	}
	if v, ok := os.LookupEnv(envar); ok {
		fs.Set(name, v)
	}
}
