package slpath

import (
	"os"
	"path/filepath"
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
