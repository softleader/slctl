package version

import (
	"fmt"
	"strings"
)

const (
	unreleased = "unreleased"
	unknown    = "unknown"
)

// BuildMetadata 代表此 app 的 release 資訊
type BuildMetadata struct {
	GitVersion string
	GitCommit  string
}

// NewBuildMetadata 產生一個 app 的 release 資訊
func NewBuildMetadata(version, commit string) (b *BuildMetadata) {
	b = &BuildMetadata{
		GitVersion: unreleased,
		GitCommit:  unknown,
	}
	if version = strings.TrimSpace(version); version != "" {
		b.GitVersion = version
	}
	if commit = strings.TrimSpace(commit); commit != "" {
		b.GitCommit = commit
	}
	return
}

func (b *BuildMetadata) String() string {
	trunc := 7
	if len := len(b.GitCommit); len < 7 {
		trunc = len
	}
	return fmt.Sprintf("%s+%s", b.GitVersion, b.GitCommit[:trunc])
}

// FullString 回傳完整的 release 資訊
func (b *BuildMetadata) FullString() string {
	return fmt.Sprintf("%#v", b)
}
