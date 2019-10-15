package plugin

import (
	"github.com/skratchdot/open-golang/open"
)

// Open 開啟 plugin 的目錄或遠端網址
func (p *Plugin) Open(app string) error {
	src := p.Dir
	if p.FromGitHub() {
		src = "https://" + p.Source
	}
	if len(app) != 0 {
		return open.RunWith(src, app)
	}
	return open.Run(src)
}
