package slpath

import (
	"os"
	"path/filepath"
	"strings"
)

type Home string

func (h Home) String() string {
	return os.ExpandEnv(string(h))
}

func (h Home) Path(elem ...string) string {
	p := []string{h.String()}
	p = append(p, elem...)
	return filepath.Join(p...)
}

func (h Home) Plugins() string {
	return h.Path("plugins")
}

func (h Home) Config() string {
	return h.Path("config")
}

func (h Home) ConfigFile() string {
	return h.Path("config", "configs.yaml")
}

func (h Home) Cache() string {
	return h.Path("cache")
}

func (h Home) CachePlugins() string {
	return h.Path("cache", "plugins")
}

func (h Home) CacheArchives() string {
	return h.Path("cache", "archives")
}

func (h Home) ContainsAnySpace() bool {
	return strings.ContainsAny(h.String(), " ")
}
