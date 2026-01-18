package environment

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/softleader/slctl/pkg/paths"
)

// MoveHome moves home from "oldHome" to "newHome"
func MoveHome(oldHome, newHome string) error {
	if err := os.Rename(oldHome, newHome); err != nil {
		return err
	}
	h := paths.Home(newHome)
	if p := h.Plugins(); paths.IsExistDirectory(p) {
		if plugins, err := os.ReadDir(p); err == nil {
			for _, p := range plugins {
				relink(oldHome, h, p) // 如果 link 失敗也不回傳 err, 就當成 plugin 不存在了讓使用者重新安裝就好
			}
		}
	}
	return nil
}

func relink(from string, to paths.Home, plugin os.DirEntry) error {
	path := filepath.Join(to.Plugins(), plugin.Name())
	if _, err := filepath.EvalSymlinks(path); err == nil { // 代表 link 是正常的
		return nil
	}
	target, err := os.Readlink(path)
	if err != nil {
		return err
	}
	base := strings.ReplaceAll(target, from, "")
	newTarget := filepath.Join(to.String(), base)
	if err := os.Remove(path); err != nil {
		return err
	}
	return os.Symlink(newTarget, path)
}
