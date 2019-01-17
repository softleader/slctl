package plugin

import (
	"fmt"
	"github.com/Masterminds/semver"
)

// MetadataFileName 定義了 metadata 的檔案名稱
const MetadataFileName = "metadata.yaml"

// Metadata 描述了 Plugin 相關的資訊
type Metadata struct {
	Name              string   `json:"name"`
	Version           string   `json:"version"`
	Usage             string   `json:"usage"`
	Description       string   `json:"description"`
	Exec              Commands `json:"exec"`
	Hook              Commands `json:"hook"`
	IgnoreGlobalFlags bool     `json:"ignoreGlobalFlags"`
	GitHub            GitHub   `json:"github"`
}

// IsVersionGreaterThan checks current version is greater than other.version
func (m *Metadata) IsVersionGreaterThan(other *Metadata) (bool, error) {
	c, err := semver.NewConstraint(fmt.Sprintf(">%s", other.Version))
	if err != nil {
		return false, err
	}
	v, err := semver.NewVersion(m.Version)
	if err != nil {
		return false, err
	}
	return c.Check(v), nil
}

// IsVersionLegal check is Version meet SemanticVersion2 spec, more details: https://semver.org/
func (m *Metadata) IsVersionLegal() bool {
	_, err := semver.NewVersion(m.Version)
	return err == nil
}
