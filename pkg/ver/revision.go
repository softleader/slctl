package ver

import (
	"fmt"

	"github.com/Masterminds/semver"
)

// Revision 版號
type Revision string

func (r Revision) String() string {
	return string(r)
}

// IsGreaterThan 判斷此 revision 是否 > other revision
func (r Revision) IsGreaterThan(other string) (bool, error) {
	c, err := semver.NewConstraint(fmt.Sprintf(">%s", other))
	if err != nil {
		return false, err
	}
	v, err := semver.NewVersion(r.String())
	if err != nil {
		return false, err
	}
	return c.Check(v), nil
}
