package ver

import (
	"fmt"

	"github.com/blang/semver"
)

// Revision 版號
type Revision string

func (r Revision) String() string {
	return string(r)
}

// IsGreaterThan 判斷此 revision 是否 > other revision
func (r Revision) IsGreaterThan(other string) (bool, error) {
	rr, err := semver.ParseRange(fmt.Sprintf(">%s", other))
	if err != nil {
		return false, err
	}
	v, err := semver.Parse(r.String())
	if err != nil {
		return false, err
	}
	return rr(v), nil
}
