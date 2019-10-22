package release

import (
	"fmt"
	"strings"
)

const (
	unreleased = "unreleased"
	unknown    = "unknown"
)

// Metadata 代表此 app 的 release 資訊
type Metadata struct {
	GitVersion string
	GitCommit  string
}

// NewMetadata 產生一個 app 的 release 資訊
func NewMetadata(version, commit string) (b *Metadata) {
	b = &Metadata{
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

func (b *Metadata) String() string {
	trunc := 7
	if len := len(b.GitCommit); len < 7 {
		trunc = len
	}
	return fmt.Sprintf("%s+%s", b.GitVersion, b.GitCommit[:trunc])
}

// FullString 回傳完整的 release 資訊
func (b *Metadata) FullString() string {
	return fmt.Sprintf("%#v", b)
}

// IsReleased 回傳 Metadata 是否已經 released
func (b *Metadata) IsReleased() bool {
	return b.GitVersion != unreleased && b.GitCommit != unknown
}
