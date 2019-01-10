package paths

import (
	"os"
	"path/filepath"
	"strings"
)

// Home 代表此 app 寫 filesystem 的根目錄路徑
type Home string

func (h Home) String() string {
	return os.ExpandEnv(string(h))
}

func (h Home) path(elem ...string) string {
	p := []string{h.String()}
	p = append(p, elem...)
	return filepath.Join(p...)
}

// Plugins 回傳 plugins 的目錄路徑
func (h Home) Plugins() string {
	return h.path("plugins")
}

// Config 回傳 config 的目錄路徑
func (h Home) Config() string {
	return h.path("config")
}

// ConfigFile 回傳 config/configs.yaml 的檔案路徑
func (h Home) ConfigFile() string {
	return h.path("config", "configs.yaml")
}

// Cache 回傳 cache 的目錄路徑
func (h Home) Cache() string {
	return h.path("cache")
}

// CachePlugins 回傳 cache/plugins 的目錄路徑
func (h Home) CachePlugins() string {
	return h.path("cache", "plugins")
}

// CacheRepositoryFile 回傳 cache/plugins/repository.yaml 的檔案路徑
func (h Home) CacheRepositoryFile() string {
	return h.path("cache", "plugins", "repository.yaml")
}

// CacheArchives 回傳 cache/archives 的目錄路徑
func (h Home) CacheArchives() string {
	return h.path("cache", "archives")
}

// Mounts 回傳 mounts 的目錄路徑
func (h Home) Mounts() string {
	return h.path("mounts")
}

// ContainsAnySpace 回傳 Home 根目錄是否包含空白
func (h Home) ContainsAnySpace() bool {
	return strings.ContainsAny(h.String(), " ")
}
