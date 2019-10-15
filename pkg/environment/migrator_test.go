package environment

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/slctl/pkg/paths"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestMoveHome(t *testing.T) {
	root, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(root)

	oldHome, err := createTempHome(root)
	if err != nil {
		t.Fatal(err)
	}

	home := filepath.Join(root, "home")
	if err := MoveHome(oldHome, home); err != nil {
		t.Fatal(err)
	}

	// 舊的目錄應該是空的
	if paths.IsExistDirectory(oldHome) {
		t.Fatal("舊目錄應該不存在")
	}

	// 新目錄應該有舊目錄的資料
	files, err := ioutil.ReadDir(home)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(files); l != 2 {
		t.Fatal(fmt.Printf(" 新目錄應該有舊目錄的資料, 也就是2個資料夾, 但目前只有%v個\n", l))
	}

	// 所有 link 應該都要正常
	hh := paths.Home(home)
	plugins, err := ioutil.ReadDir(hh.Plugins())
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range plugins {
		link := filepath.Join(hh.Plugins(), p.Name())
		if _, err := filepath.EvalSymlinks(link); err != nil {
			t.Fatal(err)
		}
	}
}

func createTempHome(root string) (string, error) {
	oldHome, err := ioutil.TempDir(root, "old-home-")
	if err != nil {
		return "", err
	}
	home := paths.Home(oldHome)
	configDirectories := []string{
		home.String(),
		home.Config(),
		home.Plugins(),
	}
	if err := paths.EnsureDirectories(logrus.StandardLogger(), configDirectories...); err != nil {
		return "", err
	}
	lh := paths.Home(oldHome)
	if err := os.Mkdir(filepath.Join(lh.Plugins(), "not-a-link"), 0755); err != nil {
		return "", err
	}
	if err := os.Symlink(lh.Config(), filepath.Join(lh.Plugins(), "should-relink")); err != nil {
		return "", err
	}
	static := filepath.Join(root, "should-not-relink")
	if err := os.Mkdir(static, 0755); err != nil {
		return "", err
	}
	if err := os.Symlink(static, filepath.Join(lh.Plugins(), "should-not-relink")); err != nil {
		return "", err
	}
	return lh.String(), nil
}
