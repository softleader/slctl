package plugin

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Unmount 解除 Plugin Mount Volume
func (p *Plugin) Unmount() error {
	logrus.Debugf("unmounting plugin %q", p.Metadata.Name)
	if err := os.RemoveAll(p.Mount); err != nil {
		return err
	}
	logrus.Debugf("removed mount volume: %s", p.Mount)
	return nil
}
