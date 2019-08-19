package environment

import (
	"github.com/mitchellh/go-homedir"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
	"regexp"
)

var (
	// Settings 代表此 app 的環境變數
	Settings = new(settings)
	envMap   = map[string]string{
		"home":    "SL_HOME",
		"offline": "SL_OFFLINE",
		"verbose": "SL_VERBOSE",
	}
	// Flags 代表此 app 的 global flags
	Flags       = flags()
	leadingDash = regexp.MustCompile(`^[-]{1,2}(.+)`)
	// oldHome 代表了在 3.6.x 版本之前的預設 Home 位置
	oldHome = ".sl"
	// newHome 代表了在 3.7.x 版本5之後的預設 Home 位置
	newHome = filepath.Join(".config", "slctl")
)

type settings struct {
	Home    paths.Home
	Verbose bool
	Offline bool
}

// AddFlags 設定 Settings 會用到的環境變數到 flag 中
func (s *settings) AddFlags(fs *pflag.FlagSet) error {
	var found bool
	var defaultHome string
	if defaultHome, found = os.LookupEnv("SL_HOME"); found {
		if expanded, err := homedir.Expand(defaultHome); err != nil {
			defaultHome = expanded
		}
	} else {
		h, err := homedir.Dir()
		if err != nil {
			return err
		}
		defaultHome = DefaultHome(h)
	}
	fs.StringVar((*string)(&s.Home), "home", defaultHome, "location of your config. Overrides $SL_HOME")
	fs.BoolVarP(&s.Verbose, "verbose", "v", false, "enable verbose output")
	fs.BoolVar(&s.Offline, "offline", false, "work offline")
	return nil
}

// DefaultHome 回傳此 app 的預設 home 目錄名稱
func DefaultHome(base string) string {
	oh := filepath.Join(base, oldHome)
	nh := filepath.Join(base, newHome)
	if paths.IsExistDirectory(oh) {
		if err := os.Rename(oh, nh); err != nil {
			return oh
		}
	}
	return nh
}

// ExpandEnvToFlags 將當前系統參數中已經設的值複寫掉 flags 的設定
func (s *settings) ExpandEnvToFlags(fs *pflag.FlagSet) {
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

// IsGlobalFlag 回傳此 flag 是否為 global flag
func IsGlobalFlag(flag string) (global bool) {
	for _, f := range Flags {
		if f == flag {
			return true
		}
	}
	return false
}

// PluginDirs 回傳 plugin 目錄, 若環境變數中已經有 SL_PLUGIN 設定則會以環境參數為主
func (s settings) PluginDirs() string {
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
