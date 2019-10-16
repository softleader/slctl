package installer

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/paths"
	"github.com/softleader/slctl/pkg/plugin"
)

var errNonResolvableInOfflineMode = fmt.Errorf("non-resolvable plugin source in offline mode")

// Installer 抽象化了 plugin 的安裝實作
type Installer interface {
	Install() (*plugin.Plugin, error)
}

// NewInstaller 依照傳入的資訊取得 plugin 的安裝實作
func NewInstaller(log *logrus.Logger, source string, tag string, asset int, home paths.Home, opt *InstallOption) (Installer, error) {
	if plugin.IsLocalDirReference(source) {
		return newLocalInstaller(log, source, home, opt)
	}
	if plugin.IsSupportedArchive(source) {
		return newArchiveInstaller(log, source, home, opt)
	}
	if plugin.IsGitHubRepo(source) {
		return newGitHubInstaller(log, source, tag, asset, home, opt)
	}
	return nil, fmt.Errorf("unsupported plugin source: %s", source)
}

// InstallOption 封裝了跟安裝有關的控制
type InstallOption struct {
	DryRun bool // 模擬 install, 只會印出相關訊息, 但所有的 install 指令都不會執行
	Force  bool // 表示如果當前已經安裝過, 會強制移除重新安裝
	Soft   bool // soft means remove exist plugin only if version is different
}
