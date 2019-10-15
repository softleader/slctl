package environment

import (
	"github.com/softleader/slctl/pkg/paths"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func moveHome(from, to string) error {
	if err := os.Rename(from, to); err != nil {
		return err
	}
	h := paths.Home(to)
	if p := h.Plugins(); paths.IsExistDirectory(p) {
		if plugins, err := ioutil.ReadDir(p); err == nil {
			for _, p := range plugins {
				relink(from, h, p) // 如果 link 失敗也不回傳 err, 就當成 plugin 不存在了讓使用者重新安裝就好
			}
		}
	}
	return nil
}

func relink(from string, to paths.Home, plugin os.FileInfo) error {
	path := filepath.Join(to.Plugins(), plugin.Name())
	if _, err := filepath.EvalSymlinks(path); err != nil {
		target, err := os.Readlink(path)
		if err != nil {
			return err
		}
		base := strings.ReplaceAll(target, from, "")
		newTarget := filepath.Join(to.String(), base)
		if err := os.Remove(path); err != nil {
			return err
		}
		if err := os.Symlink(newTarget, path); err != nil {
			return err
		}
	}
	return nil
}
