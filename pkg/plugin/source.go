package plugin

import (
	"os"
	"regexp"
	"strings"
)

var (
	// SupportedExtensions 代表所有支援壓縮檔的副檔名
	SupportedExtensions = []string{
		".zip",
		".tar",
		".tar.gz",
		".tgz",
		".tar.bz2",
		".tbz2",
		".tar.xz",
		".txz",
		".tar.lz4",
		".tlz4",
		".tar.sz",
		".tsz",
		".rar",
		".bz2",
		".gz",
		".lz4",
		".sz",
		".xz",
	}
	// GitHubRepo 代表檢查 github repo 文字的 regular expression
	GitHubRepo = regexp.MustCompile(`^(http[s]?://)?github.com/([^/]+)/([^/]+)[/]?$`)
)

// IsLocalDirReference 回傳 source 是否為 local 路徑並且為一個資料夾, 如 C:\some-dir, 而不是 https://...
func IsLocalDirReference(source string) bool {
	f, err := os.Stat(source)
	return err == nil && f.IsDir()
}

// IsLocalReference 回傳 source 是否為 local 路徑 如 C:\some-path, 而不是 https://...
func IsLocalReference(source string) bool {
	_, err := os.Stat(source)
	return err == nil
}

// IsSupportedArchive 回傳 source 是否支援傳入的檔案路徑的副檔名
func IsSupportedArchive(source string) bool {
	for _, suffix := range SupportedExtensions {
		if strings.HasSuffix(source, suffix) {
			return true
		}
	}
	return false
}

// IsGitHubRepo 回傳 source 是否為 GitHub Repo
func IsGitHubRepo(source string) bool {
	return GitHubRepo.MatchString(source)
}
