package environment

import (
	"github.com/softleader/slctl/pkg/homedir"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

var (
	DefaultHome = filepath.Join(homedir.HomeDir(), ".sl")
	Settings    = new(EnvSettings)
)

type EnvSettings struct {
	Home    slpath.Home
	Verbose bool
	Offline bool
}

func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar((*string)(&s.Home), "home", DefaultHome, "location of your config. Overrides $SL_HOME")
	fs.BoolVarP(&s.Verbose, "verbose", "v", false, "enable verbose output")
	fs.BoolVar(&s.Offline, "offline", false, "work offline")
}

func (s *EnvSettings) Init(fs *pflag.FlagSet) {
	for name, envar := range envMap {
		setFlagFromEnv(name, envar, fs)
	}
}

var envMap = map[string]string{
	"home":    "SL_HOME",
	"verbose": "SL_VERBOSE",
	"offline": "SL_OFFLINE",
}

func (s EnvSettings) PluginDirs() string {
	if d, ok := os.LookupEnv("SL_PLUGIN"); ok {
		return d
	}
	return s.Home.Plugins()
}

func setFlagFromEnv(name, envar string, fs *pflag.FlagSet) {
	if fs.Changed(name) {
		return
	}
	if v, ok := os.LookupEnv(envar); ok {
		fs.Set(name, v)
	}
}
