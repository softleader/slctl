package paths

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// EnsureDirectories 確保傳入資料夾路徑存在, 不存在會主動建立之
func EnsureDirectories(log *logrus.Logger, dirs ...string) (err error) {
	for _, dir := range dirs {
		if err = EnsureDirectory(log, dir); err != nil {
			return
		}
	}
	return
}

// EnsureDirectory 確保傳入資料夾路徑存在, 不存在會主動建立之
func EnsureDirectory(log *logrus.Logger, dir string) error {
	if fi, err := os.Stat(dir); err != nil {
		log.Printf("Creating %s \n", dir)
		if err = os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("could not create %s: %s", dir, err)
		}
	} else if !fi.IsDir() {
		return fmt.Errorf("%s must be a directory", dir)
	}
	return nil
}

// IsExistDirectory 判斷傳入的資料夾路徑存在, 且是目錄
func IsExistDirectory(dir string) bool {
	if fi, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	} else if !fi.IsDir() {
		return false
	}
	return true
}

// IsExistFile 判斷傳入的資料夾路徑存在, 且是檔案
func IsExistFile(dir string) bool {
	if fi, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	} else if fi.IsDir() {
		return false
	}
	return true
}
