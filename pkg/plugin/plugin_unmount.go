package plugin

import (
	"github.com/sirupsen/logrus"
	"os"
)

func (p *Plugin) Unmount() error {
	logrus.Debugf("unmounting plugin %q", p.Metadata.Name)
	if err := os.RemoveAll(p.Mount); err != nil {
		return err
	}
	logrus.Debugf("removed mount volume: %s", p.Mount)
	return nil
}
