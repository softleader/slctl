package environment

import (
	"github.com/mitchellh/go-homedir"
	"github.com/softleader/slctl/pkg/slpath"
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
	"regexp"
)

var (
	Settings = new(EnvSettings)
	envMap   = map[string]string{
		"home":    "SL_HOME",
		"offline": "SL_OFFLINE",
		"verbose": "SL_VERBOSE",
	}
	Flags       = flags()
	leadingDash = regexp.MustCompile(`^[-]{1,2}(.+)`)
)

type EnvSettings struct {
	Home    slpath.Home
	Verbose bool
	Offline bool
}

func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) error {
	var found bool
	var defaultHome string
	if defaultHome, found = os.LookupEnv("SL_HOME"); found {
		if expanded, err := homedir.Expand(defaultHome); err != nil {
			defaultHome = expanded
		}
	} else {
		if h, err := homedir.Dir(); err != nil {
			return err
		} else {
			defaultHome = DefaultHome(h)
		}
	}
	fs.StringVar((*string)(&s.Home), "home", defaultHome, "location of your config. Overrides $SL_HOME")
	fs.BoolVarP(&s.Verbose, "verbose", "v", false, "enable verbose output")
	fs.BoolVar(&s.Offline, "offline", false, "work offline")
	return nil
}

func DefaultHome(base string) string {
	return filepath.Join(base, ".sl")
}

func (s *EnvSettings) Init(fs *pflag.FlagSet) {
	for name, envar := range envMap {
		setFlagFromEnv(name, envar, fs)
	}
}

func flags() (flags []string) {
	for env := range envMap {
		flags = append(flags, "--"+env)
	}
	flags = append(flags, "-v")
	return
}

func IsGlobalFlag(flag string) (global bool) {
	for _, f := range Flags {
		if f == flag {
			return true
		}
	}
	return false
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
