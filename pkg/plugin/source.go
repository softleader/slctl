package plugin

import (
	"os"
	"regexp"
	"strings"
)

var (
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
	GitHubRepo = regexp.MustCompile(`^(http[s]?://)?github.com/([^/]+)/([^/]+)[/]?$`)
)

func IsLocalDirReference(source string) bool {
	f, err := os.Stat(source)
	return err == nil && f.IsDir()
}

func IsLocalReference(source string) bool {
	_, err := os.Stat(source)
	return err == nil
}

func IsSupportedArchive(source string) bool {
	for _, suffix := range SupportedExtensions {
		if strings.HasSuffix(source, suffix) {
			return true
		}
	}
	return false
}

func IsGitHubRepo(source string) bool {
	return GitHubRepo.MatchString(source)
}
